package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/miekg/dns"
)

var (
	disableCache   = flag.Bool("no-cache", false, "disable the dns cache")
	interfaces     = flag.String("i", "", "comma separated list of interfaces to send mdns probes on, defaults to all")
	requestTimeout = flag.Duration("timeout", time.Second, "timeout for each request")
	listenAddr     = flag.String("listen", "0.0.0.0", "address to listen on for incoming DNS queries")
	zoneFlag       = flag.String("zone", "local", "zone to relay to mdns")
)

var zone string

func main() {
	flag.Parse()
	zone = strings.ToLower(strings.TrimSuffix(dns.Fqdn(*zoneFlag), "."))
	var err error
	mdnsInterfaces, err = getInterfaces()
	if err != nil {
		log.Fatalf("error setting interfaces: %s", err)
	}
	log.Printf("starting dns -> mdns bridge for %s", zone)
	err = startServer(context.Background())
	if err != nil {
		log.Fatalf("dns server: %s", err)
	}
}
