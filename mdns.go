package main

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/lanrat/zeroconf"
	"github.com/miekg/dns"
)

// getInterfaces returns a list of interfaces from the flags, if none are provided the list is nil
func getInterfaces() ([]net.Interface, error) {
	if len(*interfaces) == 0 {
		return nil, nil
	}
	names := strings.Split(*interfaces, ",")
	out := make([]net.Interface, 0, 1)
	for _, name := range names {
		iface, err := net.InterfaceByName(name)
		if err != nil {
			return nil, fmt.Errorf("interface %s: %s", name, err)
		}
		out = append(out, *iface)
	}

	return out, nil
}

var mdnsInterfaces []net.Interface

func mDNSLookup(ctx context.Context, name string, qType uint16) ([]net.IP, error) {
	// TODO flag for zeroconf.SelectIPTraffic(zeroconf.IPv4) and zeroconf.IPv6
	// TODO initalize resolver once?
	resolver, err := zeroconf.NewResolver(zeroconf.SelectIfaces(mdnsInterfaces))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize resolver: %w", err)
	}

	name = dns.Fqdn(name)
	v("performing mDNS lookup for %s %q", dns.Type(qType).String(), name)
	IPs, err := resolver.ResolveOnce(ctx, name, qType)
	if err != nil {
		return nil, err
	}
	v("recieved mDNS response for %s %q: %+v", dns.Type(qType).String(), name, IPs)

	return IPs, nil
}
