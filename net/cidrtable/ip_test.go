package cidrtable

import (
	"bytes"
	"net"
	"testing"
)

/* Figure out how to compare IPs: */
func TestIpCompare(t *testing.T) {
	cases := []struct {
		a, b net.IP
		cmp  int
	}{
		{net.IPv4(1, 2, 3, 4), net.IPv4(1, 2, 3, 4), 0},
		{net.IPv4(1, 2, 3, 4), net.IPv4(1, 2, 3, 3), 1},
		{net.IPv4(1, 2, 3, 4), net.IPv4(1, 2, 3, 5), -1},
	}
	for _, c := range cases {
		got := bytes.Compare([]byte(c.a), []byte(c.b))
		if got != c.cmp {
			t.Errorf("[%s] <=> [%s] = %d (%d)\n", c.a.String(), c.b.String(), got, c.cmp)
		}
	}
}
