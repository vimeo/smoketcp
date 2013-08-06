# Smoketcp

Smokeping like tcp connectivity tester, reports to statsd so you get aggregate statistics into graphite.
Written in Golang.


# How to
Create a "targets" file that looks like:
```
<host>:<port>
<host>:<port>
```
then just `go build smoketcp.go` and `./smoketcp <statsd_host>:<statsd_port>` and boom.
