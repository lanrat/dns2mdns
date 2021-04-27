module dns2mdns

go 1.16

require (
	github.com/bluele/gcache v0.0.2
	github.com/lanrat/zeroconf v0.0.0-20210427155953-5493dd57f5a4
	github.com/miekg/dns v1.1.41
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210426230700-d19ff857e887 // indirect
)

//replace github.com/lanrat/zeroconf => ./zeroconf
