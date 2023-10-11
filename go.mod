module dns2mdns

go 1.16

require (
	github.com/bluele/gcache v0.0.2
	github.com/lanrat/zeroconf v0.0.0-20220623173108-ae93e87713d3
	github.com/libp2p/go-reuseport v0.2.0 // indirect
	github.com/miekg/dns v1.1.50
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sync v0.1.0
)

//replace github.com/lanrat/zeroconf => ./zeroconf
