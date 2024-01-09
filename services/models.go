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
