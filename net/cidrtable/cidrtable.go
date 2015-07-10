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

/* Key is a stringified IP w/o mask. */
type IpTable map[string]IpRangeNode

type CidrTable struct {
	list    *IpRangeNode /* Beginning of the list of IP Ranges, sorted ascending by IP. */
	prevIps IpTable      /* List of IPs just before each Range, with a ref to the range. */
	nextIps IpTable      /* List of IPs just after each Range, with a ref to the range. */
}

func InitCidr() (*CidrTable, error) {
	var c CidrTable
	c.prevIps = make(IpTable)
	c.nextIps = make(IpTable)
	return &c, nil
}

func (c *CidrTable) AddCidr(cstr string) error {
	//ip, ipnet, err := net.ParseCIDR(cstr)
	_, _, err := net.ParseCIDR(cstr)
	if err != nil {
		return err
	}
	return nil
}
