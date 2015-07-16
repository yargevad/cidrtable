package main

/*
 * The idea is to accept IPs as input and convert to various types of output:
 * - decimal (32, 64)
 * - decimal, dots between octets
 * - hexadecimal
 * - binary
 */

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		/* Blank line means we're done. */
		if line == "" {
			break
		}

		/* Get index of forward slash (indicating subnet), < 0 if not there. */
		fwdSlashIdx := strings.IndexRune(line, '/')
		var ip net.IP
		//var ipnet *net.IPNet
		var err interface{}

		if fwdSlashIdx < 0 {
			/* No forward slash, use net.ParseIP */
			ip = net.ParseIP(line)
			if ip == nil {
				err = fmt.Sprintf("couldn't parse ip [%s]", line)
			}
		} else {
			/* We have a forward slash, use net.ParseIP */
			ip, _, err = net.ParseCIDR(line)
		}

		if err != nil {
			log.Printf("WARNING: %s\n", err)
		} else {
			p := ip
			if p4 := p.To4(); len(p4) == net.IPv4len {
				for i := 0; i < 4; i++ {
					if i > 0 {
						fmt.Print(".")
					}
					fmt.Printf("%08b", uint(p4[i]))
				}
			}
			fmt.Println("")
		}
	}
}
