module dns2mdns

go 1.16

require (
	github.com/bluele/gcache v0.0.2
	github.com/lanrat/zeroconf v1.0.1-0.20210426172419-2c6c839d1006
	github.com/miekg/dns v1.1.41
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

//replace github.com/lanrat/zeroconf => ./zeroconf
