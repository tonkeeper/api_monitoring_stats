package services

type ApiMetrics struct {
	ServiceName     string
	HttpsLatency    float64
	TotalChecks     int
	SuccessChecks   int
	IndexingLatency float64
	Errors          []error
}
