# Smoketcp

Smokeping-like tcp connectivity tester, reports to statsd so you get aggregate statistics into graphite.
Written in Golang.


# How to
Create a "targets" file that looks like:
```
<host>:<port>
<host>:<port>
```
then just `go build smoketcp.go` and `./smoketcp <statsd_host>:<statsd_port> <bucket_prefix>` and boom.

#Ex: 
`./smoketcp statsd.example.com:8125 Location.for.smokeping.values`

Every second it tests every entry (in parallel), and reports errors and time-to-connection to statsd.
Statsd then aggregates across the flushInterval (in our case 60s) and stores in graphite per target the errors rate,
and the mean, lower, upper, upper_90 etc values.
