package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	FolderCount   prometheus.Gauge
	FileCount     prometheus.GaugeVec
	RequestCount  prometheus.CounterVec
	ActiveRequest prometheus.Gauge
	RequestTime   prometheus.HistogramVec
	ResponseTime  prometheus.Summary
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		FolderCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "File_Server",
			Name:      "Folder_Count",
			Help:      "Number of folders in file server",
		}),
		FileCount: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "File_Server",
			Name:      "File_Count",
			Help:      "No of files in each folder",
		}, []string{"folder"}),
		RequestCount: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "File_Server",
			Name:      "Request_Count",
			Help:      "No of request on each path",
		}, []string{"path", "method"}),
		ActiveRequest: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "File_Server",
			Name:      "Active_Request",
			Help:      "No of active requests",
		}),
		RequestTime: *prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "File_Server",
			Name:      "Request_Time",
			Help:      "Time taken for each request to complete",
			Buckets:   []float64{5.0, 8.0, 10.0, 12.0, 15.0},
		}, []string{"path"}),
		ResponseTime: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace:  "File_Server",
			Name:       "Response_Size",
			Help:       "Size of the response",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),
	}
	reg.MustRegister(m.FolderCount, m.FileCount, m.RequestCount, m.ActiveRequest, m.RequestTime, m.ResponseTime)
	return m
}
