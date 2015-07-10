package cidrtable

import (
	"net"
)

/*
 * TODO: optionally collapse across 8-bit boundaries (wishlist)
 */

type IpRange struct {
	start net.IP    /* Beginning of the IP Range. */
	end   net.IP    /* End of the IP Range. */
	cidr  net.IPNet /* The IP Range in CIDR notation (<ip>/32, etc). */
}

type IpRangeNode struct {
	data     IpRange      /* IP and CIDR info. */
	prevNode *IpRangeNode /* Previous node. */
	nextNode *IpRangeNode /* Next node. */
}

type IpTable map[[]byte]IpRangeNode

type CidrTable struct {
	list    IpRangeNode /* Beginning of the list of IP Ranges, sorted ascending by IP. */
	prevIps IpTable     /* List of IPs just before each Range, with a ref to the range. */
	nextIps IpTable     /* List of IPs just after each Range, with a ref to the range. */
}
