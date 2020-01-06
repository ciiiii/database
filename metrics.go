package database

import (
	"sync"

	"github.com/coredns/coredns/plugin"
	"github.com/prometheus/client_golang/prometheus"
)

var requestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: "database",
	Name:      "request_count_total",
	Help:      "Counter of requests made.",
}, []string{"server"})

var once sync.Once
