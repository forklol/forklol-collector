package stats

import (
	"forklol-collector/bitcoin"
	"forklol-collector/db"
	"log"
	"forklol-collector/config"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"time"
	"strings"
	"sync"
)

type StatisticTimeframe struct {
	ID          string
	From, To    int64
	CompactType CompactType
}

var jobChan chan Job
var resultChan chan result

func StartDispatcher(workers int) {
	jobChan = make(chan Job)
	resultChan = make(chan result, workers)

	d := newDispatcher(workers, resultChan)
	d.Run()
}

type PackedStats struct {
	Step      int64 `json:"step"`
	From      int64 `json:"from"`
	To        int64 `json:"to"`
	Range     int64 `json:"range"`
	Compacter string `json:"compacter"`
	Index     *[]int64 `json:"index"`
	Values    map[string]*[]Value `json:"values"`
}

var lock sync.Mutex

func Package(coins []bitcoin.Coin) {

	allStart := time.Now()

	log.Printf("\u2794 Packing new statistics..\n")

	lastTime, err := db.GetLastTime()
	lastHeights := db.GetLastHeights()

	if err != nil {
		log.Println("Could not get time of the lastTime block.")
		return
	}

	day := int64(24 * 3600)
	timeframes := []StatisticTimeframe{
		StatisticTimeframe{ID: "2016", From: 2016, CompactType: COMPACT_HEIGHT},
		StatisticTimeframe{ID: "8064", From: 8064, CompactType: COMPACT_HEIGHT},
		StatisticTimeframe{ID: "genesis", CompactType: COMPACT_HEIGHT},

		StatisticTimeframe{ID: "d7", From: lastTime - 7*day, To: lastTime, CompactType: COMPACT_TIME},
		StatisticTimeframe{ID: "d30", From: lastTime - 30*day, To: lastTime, CompactType: COMPACT_TIME},
		StatisticTimeframe{ID: "fork", From: config.CHAINSPLIT_TIMESTAMP, To: lastTime, CompactType: COMPACT_TIME},
		StatisticTimeframe{ID: "d180", From: lastTime - 180*day, To: lastTime, CompactType: COMPACT_TIME},
		StatisticTimeframe{ID: "y1", From: lastTime - 365*day, To: lastTime, CompactType: COMPACT_TIME},
		StatisticTimeframe{ID: "y2", From: lastTime - 2*365*day, To: lastTime, CompactType: COMPACT_TIME},
		StatisticTimeframe{ID: "y5", From: lastTime - 5*365*day, To: lastTime, CompactType: COMPACT_TIME},
		StatisticTimeframe{ID: "genesis", From: 1231469665, To: lastTime, CompactType: COMPACT_TIME},
	}

	numSteps := int64(72)
	presets := (*GetPresets())

	packed := make(map[string]map[string]PackedStats, 0)

	for _, coin := range coins {
		packed[coin.Symbol] = make(map[string]PackedStats)
	}

	for _, tf := range timeframes {
		for _, coin := range coins {

			workerTf := StatisticTimeframe{
				ID:          tf.ID,
				From:        tf.From,
				To:          tf.To,
				CompactType: tf.CompactType,
			}

			if tf.CompactType == COMPACT_HEIGHT {
				if tf.ID == "genesis" {
					workerTf.To = lastHeights[coin.Symbol]
					workerTf.From = 1
				} else {
					workerTf.To = lastHeights[coin.Symbol]
					workerTf.From = workerTf.To - workerTf.From
				}
			}

			stepSize := GetStepSize(workerTf.From, workerTf.To, numSteps)
			index := GetCompacterIndex(workerTf.From, workerTf.To, stepSize)
			values := make(map[string]*[]Value, 0)

			from := (*index)[0]
			to := (*index)[len(*index)-1]

			p := PackedStats{
				Step:      stepSize,
				From:      from,
				To:        to,
				Range:     to - from,
				Compacter: string(workerTf.CompactType),
				Index:     index,
				Values:    values,
			}

			packed[coin.Symbol][string(workerTf.CompactType)+"_"+workerTf.ID] = p

			for n, v := range presets {
				jobChan <- Job{
					Coin:       coin,
					StatName:   n,
					StatPreset: v,
					Timeframe:  workerTf,
					Step:       stepSize,
					Done:       resultChan,
				}
			}
		}
	}

	numJobs := len(presets) * len(timeframes) * len(coins)

	for n := 0; n < numJobs; n++ {
		r := <-resultChan

		mapKey := string(r.TF.CompactType) + "_" + r.TF.ID

		lock.Lock()
		packed[r.Coin.Symbol][mapKey].Values[r.StatName] = r.Values
		lock.Unlock()

		if len(packed[r.Coin.Symbol][mapKey].Values) == len(presets) {

			j, err := json.Marshal(packed[r.Coin.Symbol][mapKey])
			if err != nil {
				// ?
			}

			log.Printf("\u2714 Finished packing %s %s %s\n", r.Coin.Symbol, string(r.TF.CompactType), r.TF.ID)

			ioutil.WriteFile(
				fmt.Sprintf("/tmp/forklol_%s_%s_%s.json", string(r.TF.CompactType), strings.ToLower(r.Coin.Symbol), r.TF.ID),
				j,
				0644,
			)
		}
	}

	allEnd := time.Now()
	log.Printf("\u2714\u2714 Finished packing results after %s\n", allEnd.Sub(allStart))
}

type result struct {
	Coin     bitcoin.Coin
	TF       StatisticTimeframe
	StatName string
	Values   *[]Value
}

type Job struct {
	Coin       bitcoin.Coin
	StatName   string
	StatPreset StatPreset
	Timeframe  StatisticTimeframe
	Step       int64
	Done       chan result
}

type worker struct {
	Pool    chan chan Job
	JobChan chan Job
	Done    chan result
}

type dispatcher struct {
	NumWorkers int
	Pool       chan chan Job
	Done       chan result
}

func newDispatcher(workers int, done chan result) *dispatcher {
	pool := make(chan chan Job, workers)
	return &dispatcher{Pool: pool, Done: done, NumWorkers: workers}
}

func (d dispatcher) Run() {
	for n := 0; n < d.NumWorkers; n++ {
		worker := NewWorker(d.Pool, d.Done)
		worker.Start()
	}

	go d.dispatch()
}

func (d *dispatcher) dispatch() {
	for {
		select {
		case job := <-jobChan:
			go func(job Job) {
				jobChannel := <-d.Pool
				jobChannel <- job
			}(job)
		}
	}
}

func NewWorker(pool chan chan Job, done chan result) worker {
	return worker{
		Pool:    pool,
		JobChan: make(chan Job),
		Done:    done,
	}
}

func (w worker) Start() {
	go func() {
		for {
			w.Pool <- w.JobChan

			select {
			case job := <-w.JobChan:
				builder := NewStatBuilder(job.Coin)
				s, _ := builder.GetStatByPreset(
					job.StatPreset,
					job.Timeframe.CompactType,
					job.Timeframe.From,
					job.Timeframe.To,
					job.Step,
				)

				w.Done <- result{
					StatName: job.StatName,
					Values:   s,
					TF:       job.Timeframe,
					Coin:     job.Coin,
				}
			}
		}
	}()
}
