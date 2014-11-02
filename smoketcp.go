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

func doEvery(d time.Duration, f func(*statsd.Client, string, bool), s *statsd.Client, targetFile string, debug bool) {
	f(s, targetFile, debug)
	for _ = range time.Tick(d) {
		f(s, targetFile, debug)
	}
}

func processTargets(s *statsd.Client, targetFile string, debug bool) {
	content, err := ioutil.ReadFile(targetFile)
	if err != nil {
		if debug {
			fmt.Println("couldn't open targets file:", targetFile)
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
	conn.Close()
	duration := time.Since(pre)
	ms := int64(duration / time.Millisecond)
	if debug {
		fmt.Printf("%s.%s.duration %d\n", subhost, port, ms)
	}
	s.Timing(fmt.Sprintf("%s.%s", subhost, port), ms, 1)
}

func main() {
	var statsdHost = flag.String("statsdHost", "localhost", "Statsd Hostname")
	var statsdPort = flag.String("statsdPort", "8125", "Statsd port")
	var bucket = flag.String("bucket", "smoketcp", "Graphite bucket prefix")
	var targetFile = flag.String("targetFile", "targets", "File containing the list of targets, ex: server1:80")
	var debug = flag.Bool("debug", false, "if true, turn on debugging output")
	var interval = flag.Int("interval", 10, "How often to run the tests")
	flag.Parse()

	s, err := statsd.Dial(fmt.Sprintf("%s:%s", *statsdHost, *statsdPort), fmt.Sprintf("%s", *bucket))
	dieIfError(err)
	defer s.Close()
	doEvery(time.Duration(*interval)*time.Second, processTargets, s, *targetFile, *debug)
}
