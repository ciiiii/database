package database

import (
	"context"
	"errors"
	"time"

	"fmt"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/coredns/coredns/plugin/pkg/upstream"
	"github.com/coredns/coredns/request"
	"github.com/jinzhu/gorm"
	"github.com/miekg/dns"
)

const Name = "database"

var errKeyNotFound = errors.New("key not found")

type DBBackend struct {
	*gorm.DB
	Zones    []string
	Upstream *upstream.Upstream
	Next     plugin.Handler
}

func (backend *DBBackend) Name() string { return Name }
func (backend *DBBackend) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()
	opt := plugin.Options{}
	state := request.Request{W: w, Req: r}
	zone := plugin.Zones(backend.Zones).Matches(state.Name())
	if zone == "" {
		return plugin.NextOrFailure(backend.Name(), backend.Next, ctx, w, r)
	}

	var (
		records, extra []dns.RR
		err            error
	)
	switch state.QType() {
	case dns.TypeA:
		records, err = plugin.A(ctx, backend, zone, state, nil, opt)
	case dns.TypeAAAA:
		records, err = plugin.AAAA(ctx, backend, zone, state, nil, opt)
	case dns.TypeTXT:
		records, err = plugin.TXT(ctx, backend, zone, state, opt)
	case dns.TypeCNAME:
		records, err = plugin.CNAME(ctx, backend, zone, state, opt)
	case dns.TypePTR:
		records, err = plugin.PTR(ctx, backend, zone, state, opt)
	case dns.TypeMX:
		records, extra, err = plugin.MX(ctx, backend, zone, state, opt)
	case dns.TypeSRV:
		records, extra, err = plugin.SRV(ctx, backend, zone, state, opt)
	case dns.TypeSOA:
		records, err = plugin.SOA(ctx, backend, zone, state, opt)
	case dns.TypeNS:
		if state.Name() == zone {
			records, extra, err = plugin.NS(ctx, backend, zone, state, opt)
			break
		}
		fallthrough
	default:
		// Do a fake A lookup, so we can distinguish between NODATA and NXDOMAIN
		_, err = plugin.A(ctx, backend, zone, state, nil, opt)
	}
	if err != nil && backend.IsNameError(err) {
		//if backend.Fall.Through(state.Name()) {
		//	return plugin.NextOrFailure(backend.Name(), backend.Next, ctx, w, r)
		//}
		// Make err nil when returning here, so we don't log spam for NXDOMAIN.
		return plugin.BackendError(ctx, backend, zone, dns.RcodeNameError, state, nil /* err */, opt)
	}
	if err != nil {
		return plugin.BackendError(ctx, backend, zone, dns.RcodeServerFailure, state, err, opt)
	}

	if len(records) == 0 {
		return plugin.BackendError(ctx, backend, zone, dns.RcodeSuccess, state, err, opt)
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = append(m.Answer, records...)
	m.Extra = append(m.Extra, extra...)

	w.WriteMsg(m)
	return dns.RcodeSuccess, nil

	//return plugin.NextOrFailure(backend.Name(), backend.Next, ctx, w, r)

}

func (backend *DBBackend) Services(ctx context.Context, state request.Request, exact bool, opt plugin.Options) ([]msg.Service, error) {
	services, err := backend.Records(ctx, state, exact)
	if err != nil {
		return nil, err
	}
	return msg.Group(services), err
}

func (backend *DBBackend) Reverse(ctx context.Context, state request.Request, exact bool, opt plugin.Options) ([]msg.Service, error) {
	return backend.Services(ctx, state, exact, opt)
}

func (backend *DBBackend) Lookup(ctx context.Context, state request.Request, name string, typ uint16) (*dns.Msg, error) {
	return backend.Upstream.Lookup(ctx, state, name, typ)
}

func (backend *DBBackend) IsNameError(err error) bool {
	return err == errKeyNotFound
}

func (backend *DBBackend) Serial(state request.Request) uint32 {
	return uint32(time.Now().Unix())
}

func (backend *DBBackend) MinTTL(state request.Request) uint32 {
	return 30
}

func (backend *DBBackend) Transfer(ctx context.Context, state request.Request) (int, error) {
	return dns.RcodeServerFailure, nil
}

func (backend *DBBackend) Records(ctx context.Context, state request.Request, exact bool) ([]msg.Service, error) {
	name := state.Name()
	serviceList, err := backend.get(name)
	if err != nil {
		return nil, err
	}
	var services []msg.Service
	for _, service := range serviceList {
		services = append(services, msg.Service{
			Host:        service.Host,
			Port:        service.Port,
			Priority:    service.Priority,
			Weight:      service.Weight,
			Text:        service.Text,
			Mail:        service.Mail,
			TTL:         service.TTL,
			TargetStrip: 0,
			Group:       "",
			Key:         "",
		})
	}
	return services, nil
}

func (backend *DBBackend) get(name string) ([]Service, error) {
	var serviceList []Service
	fmt.Print(name)
	if err := backend.DB.Model(&Service{}).Where("name = ?", name).Scan(&serviceList).Debug().Error; err != nil {
		return nil, err
	}
	fmt.Print(len(serviceList))
	return serviceList, nil
}
