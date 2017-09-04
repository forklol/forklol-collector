package stats

import (
	"fmt"
	"forklol-collector/db"
	"log"
)

type CompactType string

const (
	COMPACT_TIME   CompactType = "time"
	COMPACT_HEIGHT CompactType = "height"
)

type ValueMethod string

const (
	METHOD_MIN  ValueMethod = "MIN"
	METHOD_MAX  ValueMethod = "MAX"
	METHOD_AVG  ValueMethod = "AVG"
	METHOD_LAST ValueMethod = "LAST"
)

type ValueType string

const (
	TYPE_FLOAT64 ValueType = "float64"
	TYPE_FLOAT32 ValueType = "float32"
	TYPE_UINT64  ValueType = "uint64"
	TYPE_INT64   ValueType = "int64"
	TYPE_INT32   ValueType = "int32"
)

type Statistic struct {
	Coin, Query, Property, Base string
}

var baseBlockQuery = "SELECT %s as value FROM blocks b WHERE b.coin = '%s' "
var baseDetailQuery = "SELECT %s as value FROM blocks b JOIN details d ON b.coin = d.coin AND b.height = d.height WHERE b.coin = '%s' "
var baseHashrateQuery = "SELECT %s as value FROM blocks b JOIN hashrates h ON b.coin = h.coin AND b.height = h.height WHERE b.coin = '%s' "

type BlockStatistic struct {
	Statistic
}

type DetailStatistic struct {
	Statistic
}

type HashrateStatistic struct {
	Statistic
}

type PriceStatistic struct {
	Statistic
}

func newStatistic(coin, prop, query, base string) Statistic {
	return Statistic{
		Coin:     coin,
		Query:    query,
		Property: prop,
		Base:     base,
	}
}

func NewBlockStatistic(coin, prop string) *BlockStatistic {
	stat := BlockStatistic{
		Statistic: newStatistic(coin, prop, baseBlockQuery, "b"),
	}

	return &stat
}

func NewDetailStatistic(coin, prop string) *BlockStatistic {
	stat := BlockStatistic{
		Statistic: newStatistic(coin, prop, baseDetailQuery, "d"),
	}

	return &stat
}

func NewHashrateStatistic(coin, prop string) *BlockStatistic {
	stat := BlockStatistic{
		Statistic: newStatistic(coin, prop, baseHashrateQuery, "h"),
	}

	return &stat
}

type Compacter struct {
	Type               CompactType
	From, To, StepSize uint64
	Base               string
}

func newCompacter(typ CompactType, from, to, stepsize uint64, base string) Compacter {
	compacter := Compacter{
		Type:     typ,
		From:     from,
		To:       to,
		StepSize: stepsize,
		Base:     base,
	}

	left := (to - from) % stepsize
	if left != 0 {
		compacter.From -= stepsize - left // stretch From back a bit so we have "full steps"
	}

	return compacter
}

func (c Compacter) Num() uint64 {
	return (c.To - c.From) / c.StepSize
}

func (c Compacter) GetRestriction(seq uint64) string {
	f := c.From + (seq * c.StepSize)
	t := f + c.StepSize

	return fmt.Sprintf(" AND b.%s > %d AND b.%s <= %d", c.Type, f, c.Type, t)
}

type Value interface{}

func (s Statistic) Compacter(typ CompactType, from, to, stepsize uint64) Compacter {
	return newCompacter(typ, from, to, stepsize, s.Base)
}

// GetValues returns an array with Value's. Value data type depends on the given ValueType and may contain <nil>'s
func (s Statistic) GetValues(compacter Compacter, method ValueMethod, typ ValueType) *[]Value {
	selec := fmt.Sprintf("%s(%s.%s)", method, s.Base, s.Property)
	qry := fmt.Sprintf(s.Query, selec, s.Coin)

	steps := compacter.Num()
	values := make([]Value, steps)

	for n := uint64(0); n < steps; n++ {
		q := qry + compacter.GetRestriction(n) + " GROUP BY b.coin"
		switch typ {
		case TYPE_FLOAT64:
			values[n] = getFloat64(q)
			break
		case TYPE_FLOAT32:
			values[n] = getFloat32(q)
			break
		case TYPE_UINT64:
			values[n] = getUint64(q)
			break
		case TYPE_INT64:
			values[n] = getInt64(q)
			break
		case TYPE_INT32:
			values[n] = getInt32(q)
			break
		}
	}

	return &values
}

func getFloat64(qry string) interface{} {
	var v float64
	if err := db.GetDB().Get(&v, qry); err != nil {
		return nil
	}

	return v
}

func getFloat32(qry string) interface{} {
	var v float32
	if err := db.GetDB().Get(&v, qry); err != nil {
		return nil
	}

	return v
}

func getUint64(qry string) interface{} {
	var v uint64
	if err := db.GetDB().Get(&v, qry); err != nil {
		return nil
	}

	return v
}

func getInt64(qry string) interface{} {
	var v int64
	if err := db.GetDB().Get(&v, qry); err != nil {
		log.Println(err.Error())
		return nil
	}

	return v
}

func getInt32(qry string) interface{} {
	var v int32
	if err := db.GetDB().Get(&v, qry); err != nil {
		return nil
	}

	return v
}
