## TLS exposing by Man-In-The-Middle attack in Go

Gomitm is a transparent proxy.
<br>
It leverages NAT for packets redirection.
Originally destined for ports 80, 443 are redirected to 8888.
Linux netfilter(iptables) and Macos pf are supported.
<br>
Gomitm terminates TLS connections, pretending to be legitimate server.
After detecting original IP addr and SNI from TLS Handshake it establishes connection to original server, pretending to be a legitimate client. 
Then it forges server certificate for legitimate server and presents it to legitimate client.
<br>
This requires RootCA to be trusted at clients side for this trick to succeed.

Start Demo (terminal 1):
```bash
cd examples/demo
sh start_macos.sh # or linux
```

Simulate client (terminal 2):
```bash
cd examples/demo
sh client.sh
```

Finish Demo (terminal 1):
```bash
^C
sh finish_macos.sh # or linux
```