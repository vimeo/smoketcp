package main

import (
	"flag"
	"fmt"
	"github.com/cactus/go-statsd-client/statsd"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

func dieIfError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func doEvery(d time.Duration, f func(*statsd.Client, string, bool), s *statsd.Client, target_file string, debug bool) {
	f(s, target_file, debug)
	for _ = range time.Tick(d) {
		f(s, target_file, debug)
	}
}

func process_targets(s *statsd.Client, target_file string, debug bool) {
	content, err := ioutil.ReadFile(target_file)
	if err != nil {
		if debug {
			fmt.Println("couldn't open targets file:", target_file)
		}
		return
	}
	targets := strings.Split(string(content), "\n")
	for _, target := range targets {
		if len(target) < 1 {
			continue
		}
		go test(target, s, debug)
	}
}

func test(target string, s *statsd.Client, debug bool) {
	tuple := strings.Split(target, ":")
	host := tuple[0]
	port := tuple[1]
	subhost := strings.Replace(host, ".", "_", -1)

	pre := time.Now()
	conn, err := net.Dial("tcp", target)
	if err != nil {
		if debug {
			fmt.Println("connect error:", subhost, port)
		}
		s.Inc(fmt.Sprintf("%s.%s.dial_failed", subhost, port), 1, 1)
		return
	}
	duration := time.Since(pre)
	ms := int64(duration / time.Millisecond)
	if debug {
		fmt.Printf("%s.%s.duration %d\n", subhost, port, ms)
	}
	s.Timing(fmt.Sprintf("%s.%s", subhost, port), ms, 1)
	conn.Close()
}

func main() {
	var statsd_host = flag.String("statsd_host", "localhost", "Statsd Hostname")
	var statsd_port = flag.String("statsd_port", "8125", "Statsd port")
	var bucket = flag.String("bucket", "smoketcp", "Graphite bucket prefix")
	var target_file = flag.String("target_file", "targets", "File containing the list of targets, ex: server1:80")
	var debug = flag.Bool("debug", false, "if true, turn on debugging output")
	flag.Parse()

	s, err := statsd.Dial(fmt.Sprintf("%s:%s", *statsd_host, *statsd_port), fmt.Sprintf("%s", *bucket))
	dieIfError(err)
	defer s.Close()
	doEvery(time.Second, process_targets, s, *target_file, *debug)
}
