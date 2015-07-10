package main

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
		line := scanner.Text() // removes newlines

		// TODO: IPv6 support

		// fixup ipv4 addresses
		switch strings.Count(line, ".") {
		case 0:
			line += ".0.0.0/8"
		case 1:
			line += ".0.0/16"
		case 2:
			line += ".0/24"
		case 3:
			// We've got something with 3 dots that might be an IP!
		default:
			log.Fatalf("ERROR: improperly formatted IP [%s]\n", line)
		}

		// has 3 dots, append /32
		if strings.IndexRune(line, '/') < 0 {
			line += "/32"
		}

		// use this instead of ParseIP because this way we get an error back
		ip, ipnet, err := net.ParseCIDR(line)
		if err != nil {
			log.Fatal(err)
		}

		if ipnet != nil {
			fmt.Printf("NOTICE: loaded [%s] as [%s] [%s]\n", line, ip.String(), ipnet.String())
		} else {
			fmt.Printf("NOTICE: loaded [%s] as [%s]\n", line, ip.String())
		}
	}

}
