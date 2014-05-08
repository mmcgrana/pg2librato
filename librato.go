package main

import (
	"github.com/samuel/go-librato/librato"
)

func LibratoStart(libratoAuth []string, metricBatches <-chan []interface{}) {
	Log("librato.start")
	lb := &librato.Client{libratoAuth[0], libratoAuth[1]}
	for {
		metricBatch := <-metricBatches
		Log("librato.post.start")
		err := lb.PostMetrics(&librato.Metrics{
			Gauges: metricBatch,
		})
		if err != nil {
			Error(err)
		}
		Log("librato.post.finish")
	}
}
