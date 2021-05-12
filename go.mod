module dns2mdns

go 1.16

require (
	github.com/bluele/gcache v0.0.2
	github.com/lanrat/zeroconf v0.0.0-20210427155953-5493dd57f5a4
	github.com/miekg/dns v1.1.42
	golang.org/x/net v0.0.0-20210510120150-4163338589ed // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210511113859-b0526f3d8744 // indirect
)

//replace github.com/lanrat/zeroconf => ./zeroconf
