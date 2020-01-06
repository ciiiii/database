package database

import (
	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"
)

func init() {
	plugin.Register(Name, setup)
}

func setup(c *caddy.Controller) error {
	parseDBConfig(c)
	c.OnStartup(func() error {
		once.Do(func() {
			metrics.MustRegister(c, requestCount)
		})
		return nil
	})
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return DBBackend{Next: next, Debug: true}
	})

	return nil
}

func parseDBConfig(c *caddy.Controller) (*DBBackend, error) {
	//var (
	//	username string
	//	password string
	//	host string
	//	port string
	//	db string
	//	suffix string
	//)
	for c.Next() {
		fmt.Print(c.RemainingArgs())
		for c.NextBlock() {
			switch c.Val() {
			case "username":
				fmt.Print("username", c.RemainingArgs())
			case "password":
				fmt.Print("password", c.RemainingArgs())
			default:
				if c.Val() != "}" {

				}
			}
		}
	}
	//backend := DBBackend{}
	return &DBBackend{}, nil
}

func newDBClient(host, username, password, dbName string, port int) (*gorm.DB, error) {
	connArgs := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", host, port, username, dbName, password)
	db, err := gorm.Open("postgres", connArgs)
	if err != nil {
		return db, err
	}
	db.AutoMigrate(&Domain{}, &Record{})
	return db, nil
}
