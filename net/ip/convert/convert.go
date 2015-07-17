package convert

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

const (
	Binary = iota
	Hex
)

func ConvertIP(inIp string, convertTo int) (string, error) {
	var buf bytes.Buffer
	var ip net.IP
	//var ipnet *net.IPNet
	var err error

	/* Get index of forward slash (indicating subnet), < 0 if not there. */
	fwdSlashIdx := strings.IndexRune(inIp, '/')
	if fwdSlashIdx < 0 {
		/* No forward slash, use net.ParseIP */
		ip = net.ParseIP(inIp)
		if ip == nil {
			err = fmt.Errorf("couldn't parse ip [%s]", inIp)
		}
	} else {
		/* We have a forward slash, use net.ParseIP */
		ip, _, err = net.ParseCIDR(inIp)
	}

	if err != nil {
		return "", err
	} else if len(ip) == 0 {
		return "", fmt.Errorf("zero-length IP parsed from [%s]", inIp)
	}

	/* TODO: optional zero-padding in format strings */
	p := ip

	if p4 := p.To4(); len(p4) == net.IPv4len {
		for i := 0; i < 4; i++ {
			if i > 0 {
				buf.WriteString(".")
			}
			if convertTo == Binary {
				buf.WriteString(fmt.Sprintf("%08b", uint(p4[i])))
			} else if convertTo == Hex {
				buf.WriteString(fmt.Sprintf("%02x", uint(p4[i])))
			}
		}
		return buf.String(), nil
	}

	if len(p) != net.IPv6len {
		return "", fmt.Errorf("WARNING: unexpected ip length %d for [%s]\n", len(p), inIp)
	}

	/* TODO: optionally collapse zeroes */
	for i := 0; i < net.IPv6len; i += 2 {
		if i > 0 {
			buf.WriteString(fmt.Sprintf(":"))
		}
		if convertTo == Binary {
			buf.WriteString(fmt.Sprintf("%016b", (uint32(p[i])<<8)|uint32(p[i+1])))
		} else if convertTo == Hex {
			buf.WriteString(fmt.Sprintf("%04x", (uint32(p[i])<<8)|uint32(p[i+1])))
		}
	}
	return buf.String(), nil

}
