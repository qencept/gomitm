## TLS exposing by Man-In-The-Middle attack in Go

Gomitm acts as a transparent proxy.
<br>
It leverages NAT for packets redirection.
Originally destined for ports 80, 443 are redirected to 8888 in demo.
Linux netfilter(iptables) and Macos pf are supported.
Gomitm terminates TLS connections, pretending to be legitimate server.
After detecting original IP addr and SNI from TLS Handshake it establishes connection to original server, pretending to be a legitimate client. 
Then it forges server certificate for legitimate server and presents it to legitimate client.
This requires RootCA to be trusted at client side for this trick to succeed.

Verified on GOOS/GOARCH:
- linux/amd64 (Intel Core i7)
- linux/arm (Raspberry Pi 3)
- darwin/amd64 (Intel Core i5)
- darwin/arm64 (Apple Silicon M1)

## Demo

Start (generate certs, configure NAT, run proxy)
```bash
[terminal1]
cd examples/demo
sh start_macos.sh # or linux
```
 
Simulate Client (initiate HTTP, HTTP(S), DoH)
```bash
[terminal2]
cd examples/demo
sh client.sh
```

Finish (reset NAT to default):
```bash
[terminal1]
^C
sh finish_macos.sh # or linux
```
Note:
- Dumps can be found in dirs "doh", "http", "session".
- Still runs at root, since Macos requires privileges to access "/dev/pf"