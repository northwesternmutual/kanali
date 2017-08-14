package metrics

// Metric represent a single request metric
type Metric struct {
	Name  string
	Value string
	Index bool
}

// Metrics represents a list of Metrics associated with a request
type Metrics []Metric

// Add appends a metric to the list of metrics
func (m *Metrics) Add(metrics ...Metric) {
	for _, metric := range metrics {
		*m = append(*m, metric)
	}
}
