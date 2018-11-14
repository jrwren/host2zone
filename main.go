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
	var zone string
	flag.StringVar(&zone, "zone", "example.com", "the name of the dns zone")
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
	fmt.Printf(soaTemplate, zone, zone, serial, zone)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		fields := strings.Split(line, " ")
		if len(fields) < 2 {
			continue
		}
		ip := fields[0]
		for _, host := range fields[1:] {
			if host == "" {
				continue
			}
			fmt.Printf("%s IN A %s\n", host, ip)
		}
	}
}

var soaTemplate string = `$TTL    604800
@       IN      SOA     %s. hostmaster@%s. (
                             %s          ; Serial
                         604800         ; Refresh
                          86400         ; Retry
                        2419200         ; Expire
                         604800 )       ; Negative Cache TTL
;
@                       IN      NS              %s.
`
