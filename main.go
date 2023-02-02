package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	flag "github.com/spf13/pflag"
)

const (
	ERR_PARAM = 1
	ERR_RUN   = 2
)

var (
	count    uint32
	size     uint32
	interval uint16
	timeout  uint16
	target   string
	verbose  bool
)

func init() {
	flag.Uint32VarP(&count, "count", "c", 10, "num. of icmp packets to send")
	flag.Uint32VarP(&size, "size", "s", 24, "size of icmp packet payload to send(in bytes)")
	flag.Uint16VarP(&interval, "interval", "i", 1000, "time interval between each icmp packet(in ms)")
	flag.Uint16VarP(&timeout, "timeout", "t", 5000, "timeout for waiting for icmp response(in ms)")
	flag.StringVarP(&target, "host", "h", "", "target host to ping")
	flag.BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func main() {
	flag.Parse()
	if target == "" {
		fmt.Fprintln(os.Stderr, "  --host/-h must be provided")
		flag.Usage()
		os.Exit(ERR_PARAM)
	}

	pinger, err := probing.NewPinger(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "--host/-h paramter is not valid: %s\n", err.Error())
		os.Exit(ERR_PARAM)
	}
	pinger.SetPrivileged(true)
	pinger.Interval = time.Duration(interval) * time.Millisecond
	pinger.Timeout = time.Duration(timeout) * time.Millisecond
	pinger.Count = int(count)
	pinger.Size = int(size)

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		pinger.Stop()
	}()

	var (
		stats_previous *probing.Statistics
		pkt_previous   *probing.Packet
	)
	pinger.OnRecv = func(pkt *probing.Packet) {
		if verbose {
			stats := pinger.Statistics()
			timestamp := time.Now().Format(time.RFC3339Nano)
			if stats_previous == nil || pkt_previous == nil {
				stats_previous = stats
				pkt_previous = pkt
			}
			fmt.Printf(
				"%s: %d packets transmitted(delta %d), %d packets received(delta %d), %.2f%% packet loss, %v ms rtt(delta %v)\n",
				timestamp,
				stats.PacketsSent,
				stats.PacketsSent-stats_previous.PacketsSent,
				stats.PacketsRecv,
				stats.PacketsRecv-stats_previous.PacketsRecv,
				stats.PacketLoss,
				pkt.Rtt.Milliseconds(),
				math.Abs(float64((pkt.Rtt - pkt_previous.Rtt).Milliseconds())),
			)
			stats_previous = stats
			pkt_previous = pkt
		} else {
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n", pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
		}

	}
	pinger.OnFinish = func(stats *probing.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %.2f%% packet loss\n", stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n", stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fail to start icmp ping: %s\n", err.Error())
	}
}
