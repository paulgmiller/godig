// file: rawdig.go
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/miekg/dns"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <hostname> [server]\n", os.Args[0])
		os.Exit(1)
	}
	name := dns.Fqdn(os.Args[1])

	server := "1.1.1.1:53"
	if len(os.Args) >= 3 {
		server = os.Args[2]
		if _, _, err := net.SplitHostPort(server); err != nil {
			server = net.JoinHostPort(server, "53")
		}
	}

	c := &dns.Client{
		Net:     "udp",
		Timeout: 3 * time.Second,
	}

	// A query
	m := new(dns.Msg)
	m.SetQuestion(name, dns.TypeA)

	start := time.Now()
	log.Print("Sent")
	in, rtt, err := c.Exchange(m, server)
	elapsed := time.Since(start)

	if err != nil {
		log.Printf("dns error: %v (elapsed: %s)\n", err, elapsed)
		os.Exit(1)
	}

	log.Printf("Server: %s (rtt: %s)\n", server, rtt)
	if in.Rcode != dns.RcodeSuccess {
		fmt.Printf("Rcode: %s\n", dns.RcodeToString[in.Rcode])
		os.Exit(1)
	}

	if len(in.Answer) == 0 {
		fmt.Println("No A records")
	}

	for _, ans := range in.Answer {
		if a, ok := ans.(*dns.A); ok {
			fmt.Println(a.A.String())
		}
	}
}
