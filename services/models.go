package services

type ApiMetrics struct {
	ServiceName     string
	HttpsLatency    float64
	TotalChecks     int
	SuccessChecks   int
	IndexingLatency float64
	Errors          []error
}

type DAppMetrics struct {
	ServiceName         string
	MainPageLoadLatency float64
	ApiLatency          float64
	IndexationLatency   *float64
	TotalChecks         int
	SuccessChecks       int
	Errors              []error
}

type BridgeMetrics struct {
	ServiceName     string
	TotalChecks     int
	SuccessChecks   int
	Errors          []error
	TransferLatency float64
	Reconnects      int
}

type TxPropagationMetrics struct {
	Latencies map[string]float64 // service -> seconds until tx appeared
	Errors    []error
}
