package database

import (
	"fmt"
	"strconv"

	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/coredns/coredns/plugin/pkg/upstream"
	"github.com/jinzhu/gorm"
)

func init() {
	plugin.Register(name, setup)
}

func setup(c *caddy.Controller) error {
	backend, err := parseDBConfig(c)
	if err != nil {
		return err
	}
	c.OnStartup(func() error {
		once.Do(func() {
			metrics.MustRegister(c, requestCount)
		})
		return nil
	})
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		backend.Next = next
		return backend
	})

	return nil
}

func parseDBConfig(c *caddy.Controller) (*DBBackend, error) {
	var (
		dialect  string
		err      error
		username string
		password string
		host     string
		port     int
		dbName   string
		ssl      string
		debug    bool
	)
	backend := DBBackend{}
	debug = false
	backend.Zones = make([]string, len(c.ServerBlockKeys))
	copy(backend.Zones, c.ServerBlockKeys)
	for i, str := range backend.Zones {
		backend.Zones[i] = plugin.Host(str).Normalize()
	}
	backend.Upstream = upstream.New()
	for c.Next() {
		args := c.RemainingArgs()
		if len(args) == 0 {
			dialect = "postgres"
		} else {
			dialect = args[0]
		}
		for c.NextBlock() {
			switch c.Val() {
			case "username":
				username = c.RemainingArgs()[0]
			case "password":
				password = c.RemainingArgs()[0]
			case "host":
				host = c.RemainingArgs()[0]
			case "port":
				port, err = strconv.Atoi(c.RemainingArgs()[0])
				if err != nil {
					return &backend, c.Errf("port should be int '%s'", c.Val())
				}
			case "db":
				dbName = c.RemainingArgs()[0]
			case "ssl":
				ssl = c.RemainingArgs()[0]
			case "debug":
				debug = true
			default:
				if c.Val() != "}" {
					return &backend, c.Errf("unknown property '%s'", c.Val())
				}
			}
		}
	}
	db, err := newDBClient(dialect, host, username, password, dbName, ssl, port)
	if err != nil {
		return &backend, c.Errf("db connect error '%s:%s@tcp(%s:%d)/%s'\n%s", username, password, host, port, dbName, err.Error())
	}
	if debug {
		backend.DB = db.Debug()
	} else {
		backend.DB = db
	}
	return &backend, nil
}

func newDBClient(dialect, host, username, password, dbName, ssl string, port int) (*gorm.DB, error) {
	connArgs := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", host, port, username, dbName, password, ssl)
	db, err := gorm.Open(dialect, connArgs)
	if err != nil {
		return db, err
	}
	db.AutoMigrate(&Service{})
	return db, nil
}
