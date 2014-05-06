package main

import (
	"github.com/samuel/go-librato/librato"
)

func libratoConnect(libratoAuth []string) *librato.Client {
	Log("librato.connect.start")
	lb := &librato.Client{libratoAuth[0], libratoAuth[1]}
	Log("librato.connect.finish")
	return lb
}

func LibratoStart(libratoAuth []string, metricBatches <-chan []interface{}, stop <-chan bool) {
	Log("librato.start")
	lb := libratoConnect(libratoAuth)

	for {
		select {
		case <-stop:
			Log("librato.stop")
			return
		default:
		}

		select {
		case metricBatch := <-metricBatches:
			Log("librato.post.start")
			err := lb.PostMetrics(&librato.Metrics{
				Gauges: metricBatch,
			})
			if err != nil {
				panic(err)
			}
			Log("librato.post.finish")
		default:
		}
	}
}
