package librato

type Metric struct {
	Name        string  `json:"name"`
	Value       float64 `json:"value"`
	MeasureTime int64   `json:"measure_time,omitempty"`
	Source      string  `json:"source,omitempty"`
}

type Gauge struct {
	Name        string  `json:"name"`
	MeasureTime int64   `json:"measure_time,omitempty"`
	Source      string  `json:"source,omitempty"`
	Count       uint64  `json:"count"`
	Sum         float64 `json:"sum"`
	Max         float64 `json:"max,omitempty"`
	Min         float64 `json:"min,omitempty"`
	SumSquares  float64 `json:"sum_squares,omitempty"`
}

type Metrics struct {
	MeasureTime int64         `json:"measure_time,omitempty"`
	Source      string        `json:"source,omitempty"`
	Counters    []Metric      `json:"counters,omitempty"`
	Gauges      []interface{} `json:"gauges,omitempty"` // Values can be either Metric or Gauge
}

// PostMetrics submits measurements for new or existing metrics.
// http://dev.librato.com/v1/post/metrics
func (cli *Client) PostMetrics(metrics *Metrics) error {
	if len(metrics.Counters) == 0 && len(metrics.Gauges) == 0 {
		return nil
	}
	return cli.request("POST", metricsURL, metrics, nil)
}
