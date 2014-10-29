package main

import (
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

func doEvery(d time.Duration, f func(*statsd.Client), s *statsd.Client) {
	f(s)
	for _ = range time.Tick(d) {
		f(s)
	}
}

func process_targets(s *statsd.Client) {
	content, err := ioutil.ReadFile("targets")
	if err != nil {
		fmt.Println("couldn't open targets file")
		return
	}
	targets := strings.Split(string(content), "\n")
	for _, target := range targets {
		if len(target) < 1 {
			continue
		}
		go test(target, s)
	}
}

func test(target string, s *statsd.Client) {
	tuple := strings.Split(target, ":")
	host := tuple[0]
	port := tuple[1]
	subhost := strings.Replace(host, ".", "_", -1)

	pre := time.Now()
	conn, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Println("connect error", target)
		s.Inc(fmt.Sprintf("%s.%s.dial_failed", subhost, port), 1, 1)
		return
	}
	duration := time.Since(pre)
	ms := int64(duration / time.Millisecond)
	fmt.Printf("%s.%s.duration %d\n", subhost, port, ms)
	s.Timing(fmt.Sprintf("%s.%s", subhost, port), ms, 1)
	conn.Close()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: smoketcp <statsd_host>:<statsd_port> <bucket_prefix>")
		fmt.Println("\nEx: smoketcp statsd.example.com:8125 Location.for.smokeping.values")
		os.Exit(1)
	}

	s, err := statsd.Dial(os.Args[1], fmt.Sprintf("%s", os.Args[2]))
	dieIfError(err)
	defer s.Close()
	doEvery(time.Second, process_targets, s)
}
