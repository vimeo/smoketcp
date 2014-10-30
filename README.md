# Smoketcp

Smokeping-like tcp connectivity tester, reports to statsd so you get aggregate statistics into graphite.
Written in Golang.


# How to
Create a "targets" file that looks like:
```
<host>:<port>
<host>:<port>
```

then just `go build smoketcp.go` and `./smoketcp --statsd_host <statsd_host> --statsd_port <statsd_port> --bucket <bucket_prefix>`

The available flags are available with --help:
```
./smoketcp --help
Usage of ./smoketcp:
  -bucket="smoketcp": Graphite bucket prefix
  -debug=false: if true, turn on debugging output
  -interval=10: How often to run the tests
  -statsd_host="localhost": Statsd Hostname
  -statsd_port="8125": Statsd port
  -target_file="targets": File containing the list of targets, ex: server1:80
```

#Ex: 
`./smoketcp -- statsd_host statsd.example.com --statsd_port 8125 --bucket Location.for.smokeping.values --interval 1 --target_file /etc/smoketcp_targets.txt`

Every second it tests every entry (in parallel), and reports errors and time-to-connection to statsd.
Statsd then aggregates across the flushInterval (in our case 60s) and stores in graphite per target the errors rate,
and the mean, lower, upper, upper_90 etc values.
