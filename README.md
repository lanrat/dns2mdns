# dns2mdns

A bridge to allow mdns unaware clients to query mdns (`.local`) host names over traditional dns.

This can help hosts on your LAN resolve mdns host names even if they do not support mdns, such as Android and Windows clients.

This is typically a bad idea, so only run this if you know what you are doing.

## Usage

```text
Usage of dns2mdns:
  -i string
        comma separated list of interfaces to send mdns probes on, defaults to all
  -listen string
        address to listen on for incoming DNS queries (default "0.0.0.0")
  -no-cache
        disable the dns cache
  -timeout duration
        timeout for each request (default 1s)
  -zone string
        zone to relay to mdns (default "local")
```

## Example

Start dns2mdns in one terminal:

```console
$ ./dns2mdns
2021/04/25 22:33:01 starting dns -> mdns bridge for local
2021/04/25 22:33:01 starting dns server on 0.0.0.0:53
```

In another terminal, use dig to lookup a local address:

```console
$ dig @127.0.0.1 -t A -q pfsense.local. 

;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 64158
;; flags: qr rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;pfsense.local.                        IN      A

;; ANSWER SECTION:
pfsense.local.         12      IN      A       192.168.1.1

;; Query time: 0 msec
;; SERVER: 127.0.0.1#53(127.0.0.1)
;; MSG SIZE  rcvd: 62
```

## Configuration

dns2mdns is not a recursive resolver.
You will need to forward the mdns `.local` zone from your default dns resolver to the host running dns2mdns.

### unbound configuration

Edit `unbound.conf` with the following options:

```conf
private-domain: "local."
domain-insecure: "local."

forward-zone:
        name: "local."
        forward-addr: IP_OF_DNS2MDNS_SERVER
```

### pfSense configuration

Services -> DNS Resolver -> General Settings

Add a Domain Override like below, using the IP of the host running dns2mdns

![pfSense configuration](https://user-images.githubusercontent.com/164192/116126803-e6e6b680-a67b-11eb-8d21-0dda30e4a83c.png)

## Docker

A [Dockerfile](Dockerfile) is provided for running dns2mdns as a Docker container. However due to the nature of DNS and mDNS queries, you will need to run the container with `--net=host` for it to work correctly.
Prebuilt images are on the [Docker Hub](https://hub.docker.com/repository/docker/lanrat/dns2mdns).

### Docker Compose example

```yaml
version: '3.7'

services:
  dns2mdns:
    container_name: dns2mdns
    image: lanrat/dns2mdns
    network_mode: host
    restart: unless-stopped
    command: -listen 192.168.1.155 -i eno0,eno0.2,eno0.5
```
