package stats

type StatPreset struct {
	Kind     StatKind
	Method   ValueMethod
	Type     ValueType
	Property string
}

type StatKind string

const (
	KIND_BLOCK    StatKind = "block"
	KIND_DETAIL   StatKind = "detail"
	KIND_HASHRATE StatKind = "hashrate"
	KIND_PRICE    StatKind = "price"
)

var presets map[string]StatPreset

func GetPresets() *map[string]StatPreset {
	return &presets
}

func LoadPresets() {

	presets = make(map[string]StatPreset, 0)

	/* Block statistics */
	presets["BLK_DIFFICULTY_AVG"] = StatPreset{Kind: KIND_BLOCK, Property: "difficulty", Method: METHOD_AVG, Type: TYPE_FLOAT64}
	presets["BLK_WORK_AVG"] = StatPreset{Kind: KIND_BLOCK, Property: "work", Method: METHOD_AVG, Type: TYPE_FLOAT64}
	presets["BLK_TIME_AVG"] = StatPreset{Kind: KIND_BLOCK, Property: "time", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["BLK_HEIGHT_AVG"] = StatPreset{Kind: KIND_BLOCK, Property: "height", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["BLK_HEIGHT_MIN"] = StatPreset{Kind: KIND_BLOCK, Property: "height", Method: METHOD_MIN, Type: TYPE_INT32}
	presets["BLK_HEIGHT_MAX"] = StatPreset{Kind: KIND_BLOCK, Property: "height", Method: METHOD_MAX, Type: TYPE_INT32}
	presets["BLK_COUNT_COUNT"] = StatPreset{Kind: KIND_BLOCK, Property: "coin", Method: METHOD_COUNT, Type: TYPE_INT32}

	/* Hashrate statistics */
	presets["HASHRATE_H3_AVG"] = StatPreset{Kind: KIND_HASHRATE, Property: "h3", Method: METHOD_AVG, Type: TYPE_FLOAT64}
	presets["HASHRATE_H6_AVG"] = StatPreset{Kind: KIND_HASHRATE, Property: "h6", Method: METHOD_AVG, Type: TYPE_FLOAT64}
	presets["HASHRATE_H12_AVG"] = StatPreset{Kind: KIND_HASHRATE, Property: "h12", Method: METHOD_AVG, Type: TYPE_FLOAT64}
	presets["HASHRATE_D1_AVG"] = StatPreset{Kind: KIND_HASHRATE, Property: "d1", Method: METHOD_AVG, Type: TYPE_FLOAT64}
	presets["HASHRATE_D3_AVG"] = StatPreset{Kind: KIND_HASHRATE, Property: "d3", Method: METHOD_AVG, Type: TYPE_FLOAT64}
	presets["HASHRATE_D7_AVG"] = StatPreset{Kind: KIND_HASHRATE, Property: "d7", Method: METHOD_AVG, Type: TYPE_FLOAT64}
	presets["HASHRATE_D30_AVG"] = StatPreset{Kind: KIND_HASHRATE, Property: "d30", Method: METHOD_AVG, Type: TYPE_FLOAT64}

	presets["HASHRATE_H3_MAX"] = StatPreset{Kind: KIND_HASHRATE, Property: "h3", Method: METHOD_MAX, Type: TYPE_FLOAT64}
	presets["HASHRATE_H6_MAX"] = StatPreset{Kind: KIND_HASHRATE, Property: "h6", Method: METHOD_MAX, Type: TYPE_FLOAT64}
	presets["HASHRATE_H12_MAX"] = StatPreset{Kind: KIND_HASHRATE, Property: "h12", Method: METHOD_MAX, Type: TYPE_FLOAT64}
	presets["HASHRATE_D1_MAX"] = StatPreset{Kind: KIND_HASHRATE, Property: "d1", Method: METHOD_MAX, Type: TYPE_FLOAT64}
	presets["HASHRATE_D3_MAX"] = StatPreset{Kind: KIND_HASHRATE, Property: "d3", Method: METHOD_MAX, Type: TYPE_FLOAT64}
	presets["HASHRATE_D7_MAX"] = StatPreset{Kind: KIND_HASHRATE, Property: "d7", Method: METHOD_MAX, Type: TYPE_FLOAT64}
	presets["HASHRATE_D30_MAX"] = StatPreset{Kind: KIND_HASHRATE, Property: "d30", Method: METHOD_MAX, Type: TYPE_FLOAT64}

	presets["HASHRATE_H3_MIN"] = StatPreset{Kind: KIND_HASHRATE, Property: "h3", Method: METHOD_MIN, Type: TYPE_FLOAT64}
	presets["HASHRATE_H6_MIN"] = StatPreset{Kind: KIND_HASHRATE, Property: "h6", Method: METHOD_MIN, Type: TYPE_FLOAT64}
	presets["HASHRATE_H12_MIN"] = StatPreset{Kind: KIND_HASHRATE, Property: "h12", Method: METHOD_MIN, Type: TYPE_FLOAT64}
	presets["HASHRATE_D1_MIN"] = StatPreset{Kind: KIND_HASHRATE, Property: "d1", Method: METHOD_MIN, Type: TYPE_FLOAT64}
	presets["HASHRATE_D3_MIN"] = StatPreset{Kind: KIND_HASHRATE, Property: "d3", Method: METHOD_MIN, Type: TYPE_FLOAT64}
	presets["HASHRATE_D7_MIN"] = StatPreset{Kind: KIND_HASHRATE, Property: "d7", Method: METHOD_MIN, Type: TYPE_FLOAT64}
	presets["HASHRATE_D30_MIN"] = StatPreset{Kind: KIND_HASHRATE, Property: "d30", Method: METHOD_MIN, Type: TYPE_FLOAT64}

	/*	presets["HASHRATE_H3_LAST"] = StatPreset{Kind: KIND_HASHRATE, Property: "h3", Method: METHOD_LAST, Type: TYPE_FLOAT64}
		presets["HASHRATE_H6_LAST"] = StatPreset{Kind: KIND_HASHRATE, Property: "h6", Method: METHOD_LAST, Type: TYPE_FLOAT64}
		presets["HASHRATE_H12_LAST"] = StatPreset{Kind: KIND_HASHRATE, Property: "h12", Method: METHOD_LAST, Type: TYPE_FLOAT64}
		presets["HASHRATE_D1_LAST"] = StatPreset{Kind: KIND_HASHRATE, Property: "d1", Method: METHOD_LAST, Type: TYPE_FLOAT64}
		presets["HASHRATE_D3_LAST"] = StatPreset{Kind: KIND_HASHRATE, Property: "d3", Method: METHOD_LAST, Type: TYPE_FLOAT64}
		presets["HASHRATE_D7_LAST"] = StatPreset{Kind: KIND_HASHRATE, Property: "d7", Method: METHOD_LAST, Type: TYPE_FLOAT64}
		presets["HASHRATE_D30_LAST"] = StatPreset{Kind: KIND_HASHRATE, Property: "d30", Method: METHOD_LAST, Type: TYPE_FLOAT64}*/

	/* Detail statistics */
	presets["DETAIL_AVGFEE_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "avgfee", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_MAXFEE_MAX"] = StatPreset{Kind: KIND_DETAIL, Property: "maxfee", Method: METHOD_MAX, Type: TYPE_INT32}
	presets["DETAIL_MINFEE_MIN"] = StatPreset{Kind: KIND_DETAIL, Property: "minfee", Method: METHOD_MIN, Type: TYPE_INT32}
	presets["DETAIL_MEDIANFEE_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "medianfee", Method: METHOD_AVG, Type: TYPE_FLOAT32}

	presets["DETAIL_AVGFEERATE_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "avgfeerate", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_VAVGFEERATE_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "vavgfeerate", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_MINFEERATE_MIN"] = StatPreset{Kind: KIND_DETAIL, Property: "minfeerate", Method: METHOD_MIN, Type: TYPE_INT32}
	presets["DETAIL_VMINFEERATE_MIN"] = StatPreset{Kind: KIND_DETAIL, Property: "vminfeerate", Method: METHOD_MIN, Type: TYPE_INT32}
	presets["DETAIL_MAXFEERATE_MAX"] = StatPreset{Kind: KIND_DETAIL, Property: "maxfeerate", Method: METHOD_MAX, Type: TYPE_INT32}
	presets["DETAIL_VMAXFEERATE_MAX"] = StatPreset{Kind: KIND_DETAIL, Property: "vmaxfeerate", Method: METHOD_MAX, Type: TYPE_INT32}
	presets["DETAIL_MEDIANFEERATE_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "medianfeerate", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_VMEDIANFEERATE_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "vmedianfeerate", Method: METHOD_AVG, Type: TYPE_FLOAT32}

	presets["DETAIL_INPUTS_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "inputs", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_OUTPUTS_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "outputs", Method: METHOD_AVG, Type: TYPE_FLOAT32}

	presets["DETAIL_MEDIANTIME_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "mediantime", Method: METHOD_AVG, Type: TYPE_FLOAT32}

	presets["DETAIL_REWARD_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "reward", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_REWARD_MIN"] = StatPreset{Kind: KIND_DETAIL, Property: "reward", Method: METHOD_MIN, Type: TYPE_INT64}
	presets["DETAIL_REWARD_MAX"] = StatPreset{Kind: KIND_DETAIL, Property: "reward", Method: METHOD_MAX, Type: TYPE_INT64}

	presets["DETAIL_FEE_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "fee", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_FEE_MIN"] = StatPreset{Kind: KIND_DETAIL, Property: "fee", Method: METHOD_MIN, Type: TYPE_INT64}
	presets["DETAIL_FEE_MAX"] = StatPreset{Kind: KIND_DETAIL, Property: "fee", Method: METHOD_MAX, Type: TYPE_INT64}

	presets["DETAIL_TOTALSIZE_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "totalsize", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_TOTALSIZE_MIN"] = StatPreset{Kind: KIND_DETAIL, Property: "totalsize", Method: METHOD_MIN, Type: TYPE_INT32}
	presets["DETAIL_TOTALSIZE_MAX"] = StatPreset{Kind: KIND_DETAIL, Property: "totalsize", Method: METHOD_MAX, Type: TYPE_INT32}

	presets["DETAIL_SWTOTALSIZE_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "swtotalsize", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_SWTOTALSIZE_MIN"] = StatPreset{Kind: KIND_DETAIL, Property: "swtotalsize", Method: METHOD_MIN, Type: TYPE_INT32}
	presets["DETAIL_SWTOTALSIZE_MAX"] = StatPreset{Kind: KIND_DETAIL, Property: "swtotalsize", Method: METHOD_MAX, Type: TYPE_INT32}

	presets["DETAIL_SWTOTALWEIGHT_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "swtotalweight", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_SWTOTALWEIGHT_MIN"] = StatPreset{Kind: KIND_DETAIL, Property: "swtotalweight", Method: METHOD_MIN, Type: TYPE_INT32}
	presets["DETAIL_SWTOTALWEIGHT_MAX"] = StatPreset{Kind: KIND_DETAIL, Property: "swtotalweight", Method: METHOD_MAX, Type: TYPE_INT32}

	presets["DETAIL_TXS_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "txs", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_TXS_MIN"] = StatPreset{Kind: KIND_DETAIL, Property: "txs", Method: METHOD_MIN, Type: TYPE_INT32}
	presets["DETAIL_TXS_MAX"] = StatPreset{Kind: KIND_DETAIL, Property: "txs", Method: METHOD_MAX, Type: TYPE_INT32}

	presets["DETAIL_SWTXS_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "swtxs", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_SWTXS_MIN"] = StatPreset{Kind: KIND_DETAIL, Property: "swtxs", Method: METHOD_MIN, Type: TYPE_INT32}
	presets["DETAIL_SWTXS_MAX"] = StatPreset{Kind: KIND_DETAIL, Property: "swtxs", Method: METHOD_MAX, Type: TYPE_INT32}

	presets["DETAIL_TOTALOUT_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "totalout", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_TOTALOUT_MAX"] = StatPreset{Kind: KIND_DETAIL, Property: "totalout", Method: METHOD_MAX, Type: TYPE_UINT64}
	presets["DETAIL_TOTALOUT_MIN"] = StatPreset{Kind: KIND_DETAIL, Property: "totalout", Method: METHOD_MIN, Type: TYPE_UINT64}

	presets["DETAIL_TOTALWEIGHT_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "totalweight", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_TOTALWEIGHT_MAX"] = StatPreset{Kind: KIND_DETAIL, Property: "totalweight", Method: METHOD_MAX, Type: TYPE_INT32}
	presets["DETAIL_TOTALWEIGHT_MIN"] = StatPreset{Kind: KIND_DETAIL, Property: "totalweight", Method: METHOD_MIN, Type: TYPE_INT32}

	presets["DETAIL_UTXOINC_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "utxoinc", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_UTXOINC_MAX"] = StatPreset{Kind: KIND_DETAIL, Property: "utxoinc", Method: METHOD_MAX, Type: TYPE_INT32}
	presets["DETAIL_UTXOINC_MIN"] = StatPreset{Kind: KIND_DETAIL, Property: "utxoinc", Method: METHOD_MIN, Type: TYPE_INT32}

	presets["DETAIL_UTXOSIZEINC_AVG"] = StatPreset{Kind: KIND_DETAIL, Property: "utxosizeinc", Method: METHOD_AVG, Type: TYPE_FLOAT32}
	presets["DETAIL_UTXOSIZEINC_MAX"] = StatPreset{Kind: KIND_DETAIL, Property: "utxosizeinc", Method: METHOD_MAX, Type: TYPE_INT32}
	presets["DETAIL_UTXOSIZEINC_MIN"] = StatPreset{Kind: KIND_DETAIL, Property: "utxosizeinc", Method: METHOD_MIN, Type: TYPE_INT32}
}
