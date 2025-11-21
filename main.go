package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <hostname> [server]\n", os.Args[0])
		os.Exit(1)
	}
	name := os.Args[1]
	server := "1.1.1.1:53"
	if len(os.Args) >= 3 {
		server = os.Args[2]
		if _, _, err := net.SplitHostPort(server); err != nil {
			// allow "1.1.1.1" shorthand
			server = net.JoinHostPort(server, "53")
		}
	}

	r := &net.Resolver{
		PreferGo: true, // force Go's resolver, skip libc
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 2 * time.Second,
			}
			// always talk to our chosen server over UDP
			return d.DialContext(ctx, "udp", server)
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	log.Printf("Sent")
	ips, err := r.LookupHost(ctx, name)
	elapsed := time.Since(start)

	if err != nil {
		log.Printf("lookup error: %v (elapsed: %s)\n", err, elapsed)
		os.Exit(1)
	}

	log.Printf("Server: %s\n", server)
	log.Printf("Query time: %s\n", elapsed)
	for _, ip := range ips {
		fmt.Println(ip)
	}
}
