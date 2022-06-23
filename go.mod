module dns2mdns

go 1.16

require (
	github.com/bluele/gcache v0.0.2
	github.com/lanrat/zeroconf v0.0.0-20220623173108-ae93e87713d3
	github.com/libp2p/go-reuseport v0.2.0 // indirect
	github.com/miekg/dns v1.1.50
	golang.org/x/net v0.0.0-20220622184535-263ec571b305 // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f
	golang.org/x/sys v0.0.0-20220622161953-175b2fd9d664 // indirect
	golang.org/x/tools v0.1.11 // indirect
)

//replace github.com/lanrat/zeroconf => ./zeroconf
