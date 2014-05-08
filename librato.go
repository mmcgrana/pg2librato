package main

import (
	"github.com/samuel/go-librato/librato"
)

func LibratoStart(libratoAuth []string, metricBatches <-chan []interface{}, stop chan bool) {
	Log("librato.start")
	lb := &librato.Client{libratoAuth[0], libratoAuth[1]}

	stopping := false
	for {
		select {
		case metricBatch := <-metricBatches:
			Log("librato.post.start")
			err := lb.PostMetrics(&librato.Metrics{
				Gauges: metricBatch,
			})
			if err != nil {
				Error(err)
			}
			Log("librato.post.finish")
		case <-stop:
			Log("librato.stop")
			stopping = true
		}
		if stopping && len(metricBatches) == 0 {
			Log("librato.exit")
			stop <- true
			return
		}
	}
}
