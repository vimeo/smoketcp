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

func doEvery(d time.Duration, s *statsd.Client, bucket string) {
	process_targets(s, bucket)
	for _ = range time.Tick(d) {
		process_targets(s, bucket)
	}
}

func process_targets(s *statsd.Client, bucket string) {
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
		go test(target, s, bucket)
	}
}

func test(target string, s *statsd.Client, bucket string) {
	tuple := strings.Split(target, ":")
	host := tuple[0]
	port := tuple[1]
	hostport := strings.Join([]string{host, port}, ":")

	pre := time.Now()
	conn, err := net.Dial("tcp", hostport)
	if err != nil {
		fmt.Println("connect error", target)
		s.Inc(fmt.Sprintf("%s.%s.%s.dial_failed", bucket, host, port), 1, 1)
		return
	}
	duration := time.Since(pre)
	subhost := strings.Replace(host, ".", "_", -1)
	ms := int64(duration / time.Millisecond)
	fmt.Printf("%s.%s.%s.duration %d\n", bucket, subhost, port, ms)
	s.Timing(fmt.Sprintf("%s.%s.%s", bucket, host, port), ms, 1)
	conn.Close()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: smoketcp <statsd_host>:<statsd_port> <bucket_prefix>")
		fmt.Println("\nEx: smoketcp statsd.example.com:8125 Location.for.smokeping.values")
		os.Exit(1)
	}

	hostname, err := os.Hostname()
	dieIfError(err)
	s, err := statsd.Dial(os.Args[1], fmt.Sprintf("smoketcp.%s", hostname))
	dieIfError(err)
	bucket := os.Args[2]
	defer s.Close()
	doEvery(time.Second, s, bucket)
}
