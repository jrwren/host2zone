package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var ns, zone string
	var ttl, refresh, retry, expire, negative_cache int
	flag.StringVar(&ns, "ns", "ns.example.com", "the name of the ns server")
	flag.StringVar(&zone, "zone", "example.com", "the name of the dns zone")
	flag.IntVar(&ttl, "ttl", 120, "the ttl of the dns zone")
	flag.IntVar(&refresh, "refresh", 604800, "the secondary refresh of the dns zone")
	flag.IntVar(&retry, "retry", 86400, "the secondary retry of the dns zone")
	flag.IntVar(&expire, "expire", 2419200, "the secondary expire of the dns zone")
	flag.IntVar(&negative_cache, "negative_cache", 604800, "the negative_cache of the dns zone")
	flag.Parse()
	var err error
	var r io.ReadCloser
	if len(flag.Args()) == 0 {
		r = os.Stdin
	} else {
		r, err = os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
	}
	serial := time.Now().Format("0601021504") // YYMMDDHHmmSS
	fmt.Printf(soaTemplate, ttl, zone, zone, serial, refresh, retry, expire,
		negative_cache, ns)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		fields := strings.Split(line, " ")
		if len(fields) < 2 {
			continue
		}
		ip := fields[0]
		for _, host := range fields[1:] {
			if host == "" || strings.Contains(ip, "%") {
				continue
			}
			if strings.Contains(ip, ":") {
				fmt.Printf("%s IN AAAA %s\n", host, ip)
				continue
			}
			host = strings.Replace(host, "."+zone, "", -1)
			fmt.Printf("%s IN A %s\n", host, ip)
		}
	}
}

var soaTemplate string = `$TTL    %d
@       IN      SOA     %s. hostmaster@%s. (
                             %s         ; Serial
                             %d         ; Refresh
                             %d         ; Retry
                             %d         ; Expire
                             %d )       ; Negative Cache TTL
;
@                       IN      NS              %s.
`
