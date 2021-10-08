module github.com/ciiiii/database

go 1.13

require (
	github.com/caddyserver/caddy v1.0.4
	github.com/coredns/coredns v1.8.6
	github.com/jinzhu/gorm v1.9.16
	github.com/miekg/dns v1.1.31
	github.com/prometheus/client_golang v1.3.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
)

replace github.com/coredns/coredns => github.com/coredns/coredns v1.6.6
