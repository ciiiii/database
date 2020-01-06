package database

import (
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/miekg/dns"
	"github.com/jinzhu/gorm"
	"context"
)

const Name = "database"

type DBBackend struct {
	*gorm.DB
	Debug bool
	Next  plugin.Handler
}

func (backend DBBackend) Name() string { return Name }
func (backend DBBackend) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	return plugin.NextOrFailure(backend.Name(), backend.Next, ctx, w, r)

}
