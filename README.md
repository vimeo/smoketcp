# Smoketcp

Smokeping-like tcp connectivity tester, reports to statsd so you get aggregate statistics into graphite.
Written in Golang.


# How to
Create a "targets" file that looks like:
```
<host>:<port>
<host>:<port>
```

then just:
```
go build smoketcp.go
```

and:
```
./smoketcp --statsdHost <statsdHost> --statsdPort <statsdPort> --bucket <bucket_prefix>
```

The flags are available with --help:
```
./smoketcp --help
Usage of ./smoketcp:
  -bucket="smoketcp": Graphite bucket prefix
  -debug=false: if true, turn on debugging output
  -interval=10: How often to run the tests
  -statsdHost="localhost": Statsd Hostname
  -statsdPort="8125": Statsd port
  -targetFile="targets": File containing the list of targets, ex: server1:80
```

Every ten seconds (configurable with --interval) it tests every entry (in parallel), and reports errors and time-to-connection to statsd.
Statsd then aggregates across the flushInterval (in our case 60s) and stores in graphite per target the errors rate,
and the mean, lower, upper, upper_90 etc values.
