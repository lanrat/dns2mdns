package main

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"dns2mdns/cache"

	"github.com/miekg/dns"
	"golang.org/x/sync/errgroup"
)

var dnsCache cache.Cache

// TODO set the correct TTL for each response, needs support from the mdns library
const defaultTTL = 60

func startServer(ctx context.Context) error {
	listen := net.JoinHostPort(*listenAddr, "53")
	log.Printf("starting dns server on %s", listen)
	var g errgroup.Group

	// start udp dns server
	g.Go(func() error {
		srv := &dns.Server{
			Addr:      listen,
			Net:       "udp",
			ReusePort: true,
		}
		srv.Handler = &handler{
			ctx: ctx,
		}
		return srv.ListenAndServe()
	})

	// start tcp dns server
	g.Go(func() error {
		srv := &dns.Server{
			Addr:      listen,
			Net:       "tcp",
			ReusePort: true,
		}
		srv.Handler = &handler{
			ctx: ctx,
		}
		return srv.ListenAndServe()
	})

	return g.Wait()
}

type handler struct {
	ctx context.Context
}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	clientIP, _, err := net.SplitHostPort(w.RemoteAddr().String())
	if err != nil {
		log.Printf("error: unable to get client IP from %q", w.RemoteAddr().String())
	}

	// we don't support multiple questions
	if len(r.Question) > 1 {
		log.Printf("warning: DNS query sent multiple questions [%s] %s", clientIP, r.String())
	}

	// check zone for authoritative responses
	domain := strings.ToLower(msg.Question[0].Name)
	labels := dns.SplitDomainName(domain)
	if len(labels) == 0 {
		msg.SetRcode(r, dns.RcodeRefused)
		log.Printf("non-authoritative question from %s rejecting empty zone", clientIP)
		err = w.WriteMsg(&msg)
		if err != nil {
			log.Printf("error on WriteMsg: %s", err)
		}
		return
	}
	qZone := labels[len(labels)-1]
	if qZone != zone {
		msg.SetRcode(r, dns.RcodeRefused)
		log.Printf("non-authoritative question from %s rejecting zone %q", clientIP, qZone)
		err = w.WriteMsg(&msg)
		if err != nil {
			log.Printf("error on WriteMsg: %s", err)
		}
		return
	}

	ctx, cancel := context.WithTimeout(h.ctx, *requestTimeout)
	defer cancel()
	msg.Authoritative = true
	if r.Question[0].Qtype == dns.TypeA || r.Question[0].Qtype == dns.TypeAAAA {
		if !*disableCache {
			resp, found := dnsCache.Get(&msg)
			if found {
				// TODO might want to make this log cleaner
				log.Printf("[%s] query: %q type %s cache-response: %v", clientIP, r.Question[0].Name, dns.Type(r.Question[0].Qtype).String(), dnsAnswerIPs(resp.Answer))
				err = w.WriteMsg(resp)
				if err != nil {
					log.Printf("error on WriteMsg: %s", err)
				}
				return
			}
		}
		ips, err := mDNSLookup(ctx, domain, r.Question[0].Qtype)
		if err != nil {
			log.Printf("error on mDNSLookup for [%s] %q %s: %s", clientIP, domain, dns.Type(r.Question[0].Qtype).String(), err)
			msg.SetRcode(r, dns.RcodeNameError)
			err = w.WriteMsg(&msg)
			if err != nil {
				log.Printf("error on WriteMsg: %s", err)
			}
			return
		}

		if len(ips) == 0 {
			log.Printf("no results found for [%s] %q %s", clientIP, domain, dns.Type(r.Question[0].Qtype).String())
			msg.SetRcode(r, dns.RcodeNameError)
			err = w.WriteMsg(&msg)
			if err != nil {
				log.Printf("error on WriteMsg: %s", err)
			}
			return
		}

		// turn each IP into the correct dns RR type
		for _, ip := range ips {
			if ip.To4() != nil {
				msg.Answer = append(msg.Answer, &dns.A{
					Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: defaultTTL},
					A:   ip,
				})
			} else {
				msg.Answer = append(msg.Answer, &dns.AAAA{
					Hdr:  dns.RR_Header{Name: domain, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: defaultTTL},
					AAAA: ip,
				})
			}
		}
		if !*disableCache {
			dnsCache.Set(&msg)
		}
		log.Printf("[%s] query: %q type %s mdns-response: %v", clientIP, r.Question[0].Name, dns.Type(r.Question[0].Qtype).String(), ips)
	} else if r.Question[0].Qtype == dns.TypeSOA {
		if domain == dns.Fqdn(zone) {
			msg.Answer = append(msg.Answer, &dns.SOA{
				Hdr:     dns.RR_Header{Name: domain, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: defaultTTL},
				Ns:      "localhost.localdomain.",
				Mbox:    "root.localhost.localdomain.",
				Serial:  uint32(time.Now().Unix()),
				Refresh: 60 * 60, // 1 hour
				Retry:   60 * 60, // 1 hour
				Expire:  60 * 60, // 1 hour
				Minttl:  defaultTTL,
			})
			log.Printf("[%s] query: %q type %s soa-response", clientIP, r.Question[0].Name, dns.Type(r.Question[0].Qtype).String())
		} else {
			log.Printf("no results found for [%s] %q %s", clientIP, domain, dns.Type(r.Question[0].Qtype).String())
			msg.SetRcode(r, dns.RcodeNameError)
		}
	} else {
		log.Printf("[%s] unsupported question: %q type %s", clientIP, r.Question[0].Name, dns.Type(r.Question[0].Qtype).String())
		msg.SetRcode(r, dns.RcodeNotImplemented)
	}
	err = w.WriteMsg(&msg)
	if err != nil {
		log.Printf("error on WriteMsg: %s", err)
	}
}

// dnsAnswerIPs gets a list of IPs from a dns response. Used in logging
func dnsAnswerIPs(answers []dns.RR) []net.IP {
	out := make([]net.IP, 0, len(answers))
	for _, answer := range answers {
		if a, ok := answer.(*dns.A); ok {
			out = append(out, a.A)
		}
		if aaaa, ok := answer.(*dns.AAAA); ok {
			out = append(out, aaaa.AAAA)
		}
	}
	return out
}
