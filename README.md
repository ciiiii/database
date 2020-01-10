# database
CoreDNS plugin for db backend

Usage:

1. clone CoreDNS repo

```bash
$git clone https://github.com/coredns/coredns.git
$git checkout v1.6.6
```

1. edit `plugin.cfg`(add db dialects according [gorm doc](https://gorm.io/docs/connecting_to_the_database.html) )
```bash
database:github.com/ciiiii/database
database_postgres:github.com/jinzhu/gorm/dialects/postgres
```

2. build CoreDNS

```bash
$make
// generate binary coredns
```

3. edit Corefile

```bash
service.dns {
  database postgres {
    username user
    password password 
    host 127.0.0.1
    port 5432
    db coredns
    ssl disable
    debug
  }
}
```

4. run

```bash
$./coredns
```

5. insert dns record to your db

```sql
insert into services(name, host, ttl) values ('example.service.dns.', '127.0.0.1', 100);
```

6. test

```bash
$dig @127.0.0.1 example.service.dns

; <<>> DiG 9.10.6 <<>> @127.0.0.1 example.service.dns
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 52319
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
;; WARNING: recursion requested but not available

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;example.service.dns.		IN	A

;; ANSWER SECTION:
example.service.dns.	100	IN	A	127.0.0.1

;; Query time: 2 msec
;; SERVER: 127.0.0.1#53(127.0.0.1)
;; WHEN: Fri Jan 10 14:53:45 CST 2020
;; MSG SIZE  rcvd: 83
```

