package database

import (
	"fmt"
	"io"
	"os"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

const Name = "database"

type DBBackend struct {
	//*gorm.DB
	Debug bool
	Next  plugin.Handler
}

func (backend DBBackend) Name() string { return Name }
func (backend DBBackend) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	pw := NewResponsePrinter(w)
	requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	return plugin.NextOrFailure(backend.Name(), backend.Next, ctx, pw, r)

}

type ResponsePrinter struct {
	dns.ResponseWriter
}

// NewResponsePrinter returns ResponseWriter.
func NewResponsePrinter(w dns.ResponseWriter) *ResponsePrinter {
	return &ResponsePrinter{ResponseWriter: w}
}

// WriteMsg calls the underlying ResponseWriter's WriteMsg method and prints "example" to standard output.
func (r *ResponsePrinter) WriteMsg(res *dns.Msg) error {
	fmt.Fprintln(out, "example")
	return r.ResponseWriter.WriteMsg(res)
}

// Make out a reference to os.Stdout so we can easily overwrite it for testing.
var out io.Writer = os.Stdout
