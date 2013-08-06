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
	pre := time.Now()
	conn, err := net.Dial("tcp", target)
	if err != nil {
		fmt.Println("connect error", target)
		s.Inc(fmt.Sprintf("error.%s.dial_failed", target), 1, 1)
		return
	}
	duration := time.Since(pre)
	tuple := strings.Split(target, ":")
	host := strings.Replace(tuple[0], ".", "_", -1)
	port := tuple[1]
	ms := int64(duration / time.Millisecond)
	fmt.Printf("%s.%s.duration %d\n", host, port, ms)
	s.Timing(fmt.Sprintf("dial.%s.%s", host, port), ms, 1)
	conn.Close()
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: smoketcp <statsd_host>:<statsd_port>")
		os.Exit(1)
	}

	hostname, err := os.Hostname()
	dieIfError(err)
	s, err := statsd.Dial(os.Args[1], fmt.Sprintf("smoketcp.%s", hostname))
	dieIfError(err)
	defer s.Close()
	doEvery(time.Second, process_targets, s)
}
