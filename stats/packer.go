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
)

type StatisticTimeframe struct {
	ID       string
	From, To uint64
}

var jobChan chan Job
var resultChan chan result

func StartDispatcher(workers int) {
	jobChan = make(chan Job)
	resultChan = make(chan result)

	d := newDispatcher(workers, resultChan)
	d.Run()
}

func Package(coins []bitcoin.Coin) {

	allStart := time.Now()

	log.Printf("\u2794 Packing new statistics..\n")

	last, err := db.GetLastTime()

	if err != nil {
		log.Println("Could not get time of the last block.")
		return
	}

	day := uint64(24 * 3600)
	timeframes := []StatisticTimeframe{
		StatisticTimeframe{ID: "d7", From: last - 7*day, To: last},
		StatisticTimeframe{ID: "d30", From: last - 30*day, To: last},
		StatisticTimeframe{ID: "d180", From: last - 180*day, To: last},
		StatisticTimeframe{ID: "y1", From: last - 365*day, To: last},
		StatisticTimeframe{ID: "y2", From: last - 2*365*day, To: last},
		StatisticTimeframe{ID: "y5", From: last - 5*365*day, To: last},
		StatisticTimeframe{ID: "fork", From: config.CHAINSPLIT_TIMESTAMP, To: last},
		StatisticTimeframe{ID: "genesis", From: 1231469665, To: last},
	}

	numSteps := uint64(72)
	presets := (*GetPresets())

	for _, coin := range coins {
		for _, tf := range timeframes {
			start := time.Now()
			packed := struct {
				Index  *[]uint64 `json:"index"`
				Values map[string]*[]Value `json:"values"`
			}{}

			packed.Values = make(map[string]*[]Value, 0)

			stepSize := GetStepSize(tf.From, tf.To, numSteps)
			index := GetCompacterIndex(tf.From, tf.To, stepSize)
			packed.Index = index

			for n, v := range presets {
				jobChan <- Job{
					Coin:       coin,
					StatName:   n,
					StatPreset: v,
					Timeframe:  tf,
					Step:       stepSize,
					Done:       resultChan,
				}
			}

			for n := 0; n < len(presets); n++ {
				result := <-resultChan
				packed.Values[result.StatName] = result.Values
			}

			j, err := json.Marshal(&packed)
			if err != nil {
				// ?
			}

			ioutil.WriteFile(fmt.Sprintf("/tmp/forklol_%s_%s.json", strings.ToLower(coin.Symbol), tf.ID), j, 0644)
			end := time.Now()

			log.Printf("\u2714 Packed %s %s statistics in %s\n", coin.Symbol, tf.ID, end.Sub(start))
		}
	}

	allEnd := time.Now()

	log.Printf("\u2714\u2714 Finished packing results after %s\n", allEnd.Sub(allStart))
}

type result struct {
	StatName string
	Values   *[]Value
}

type Job struct {
	Coin       bitcoin.Coin
	StatName   string
	StatPreset StatPreset
	Timeframe  StatisticTimeframe
	Step       uint64
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
				jobChannel <-job
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
					COMPACT_TIME,
					job.Timeframe.From,
					job.Timeframe.To,
					job.Step,
				)

				w.Done <- result{StatName: job.StatName, Values: s}
			}
		}
	}()
}
