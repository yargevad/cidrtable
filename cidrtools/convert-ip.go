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
	"flag"
	"fmt"
	"github.com/yargevad/net/ip/convert"
	"log"
	"os"
	"strings"
)

var toBinary = flag.Bool("binary", false, "convert to binary representation")
var toHex = flag.Bool("hex", false, "convert to hexadecimal representation")

func main() {
	flag.Parse()

	if (*toBinary == false) && (*toHex == false) {
		log.Fatal("specify one of: binary hex")
		os.Exit(1)
	}

	var convertTo int
	if *toBinary {
		convertTo = convert.Binary
	} else if *toHex {
		convertTo = convert.Hex
	}

	/* If there is more input on the command line, process that and exit */
	if len(flag.Args()) > 0 {
		for _, v := range flag.Args() {
			b, err := convert.ConvertIP(v, convertTo)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			} else {
				fmt.Println(b)
			}
		}
		return
	}

	/* Process STDIN and exit when we get EOF or a blank line */
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		/* Blank line means we're done. */
		if line == "" {
			break
		}

		b, err := convert.ConvertIP(line, convertTo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		} else {
			fmt.Println(b)
		}

	}
}
