package main

import (
	"fmt"
	"net/netip"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	netaddr "github.com/dspinhirne/netaddr-go"
)

const (
	// v4splitBits is the size of the subnet to split down to for IPv4
	v4splitBits uint = 24
	// v6splitBits is the size of the subnet to split down to for IPv6
	v6splitBits uint = 64
	// v6Bytes is used to describe how many bytes an ipv6 address takes to store.
	v6Bytes int = 16
)

var cli struct {
	Generate Generate `cmd:"" help:"generate in-addr for a given subnet"`
}

type Generate struct {
	Subnet string `arg:"" help:"the subnet in cidr notation like 192.168.0.0/16 or 2001:db8:abcd:1234::1/64"`
}

// fatal logs to output an error.
func fatal(m ...string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "[FATAL] ")
	for _, s := range m {
		fmt.Fprintf(&sb, "%s ", s)
	}

	return sb.String()
}

// ipv6 process IPv6 addresses into subnets, then reverse them.
func ipv6(prefix netip.Prefix) ([]string, error) {
	ipv6, err := netaddr.ParseIPv6Net(prefix.String())
	if err != nil {
		return []string{}, err
	}

	ips := make([]netip.Addr, 0)
	switch ipv6.SubnetCount(v6splitBits) {
	case 0:
		ips = append(ips, prefix.Addr())
	default:
		for i := uint64(0); i < ipv6.SubnetCount(v6splitBits); i++ {
			ip, _ := netip.ParseAddr(ipv6.NthSubnet(v6splitBits, i).Network().String())
			ips = append(ips, ip)
		}
	}

	split := (v6Bytes / (128 / 64))
	results := make([]string, len(ips))
	for idx := 0; idx < len(ips); idx++ {
		var reversed strings.Builder
		bits := ips[idx].As16()

		for segment := len(bits[:split]) - 1; segment >= 0; segment-- {
			fmt.Fprintf(&reversed, "%x.%x.", bits[segment]&0b1111, bits[segment]>>4)
		}
		fmt.Fprintf(&reversed, "in-addr.arpa.")
		results[idx] = reversed.String()
	}
	return results, nil
}

// ipv4 process IPv4 addresses into subnets, then reverse them.
func ipv4(prefix netip.Prefix) ([]string, error) {
	ipv4, err := netaddr.ParseIPv4Net(prefix.String())
	if err != nil {
		return []string{}, err
	}

	ips := make([]netip.Addr, 0)
	switch ipv4.SubnetCount(v4splitBits) {
	case 0:
		ips = append(ips, prefix.Addr())
	default:
		for i := uint32(0); i < ipv4.SubnetCount(v4splitBits); i++ {
			ip, _ := netip.ParseAddr(ipv4.NthSubnet(v4splitBits, i).Network().String())
			ips = append(ips, ip)
		}
	}

	results := make([]string, len(ips))
	// using the slice of subnets we need to reverse the network addresses and print them to stdout.
	for idx := 0; idx < len(ips); idx++ {
		var reversed strings.Builder
		for segment := len(ips[idx].As4()) - 2; segment >= 0; segment-- {
			fmt.Fprintf(&reversed, "%d.", ips[idx].As4()[segment])
		}
		fmt.Fprintf(&reversed, "in-addr.arpa.")
		results[idx] = reversed.String()
	}
	return results, nil
}

// checkPrefixes checks the type of prefix and sends the task to the correct function.
func checkPrefixes(prefix netip.Prefix) ([]string, error) {
	networkAddress := prefix.Masked()
	switch prefix.Addr().BitLen() {
	case 128:
		// IPv6 magic
		return ipv6(networkAddress)
	case 32:
		// IPv4 magic
		return ipv4(networkAddress)
	default:
		return []string{}, fmt.Errorf("exception in checkPrefixes, bitsize not handled")
	}
}

// Run generates the in-addr syntax.
func (g *Generate) Run() error {
	// check if the input is a subnet, if it is an IP we append the subnet mask.
	if !strings.Contains(g.Subnet, "/") {
		p, err := netip.ParseAddr(g.Subnet)
		if err != nil {
			return fmt.Errorf(fatal(g.Subnet, "invalid input, does not match IP Address or Prefix"))
		}
		var withPrefix string
		switch p.BitLen() {
		case 32:
			withPrefix = fmt.Sprintf("%s/32", g.Subnet)

		case 128:
			withPrefix = fmt.Sprintf("%s/128", g.Subnet)
		}
		return fmt.Errorf(fatal(g.Subnet, "does not include a subnet mask, try", withPrefix))
	}
	prefix, err := netip.ParsePrefix(g.Subnet)
	if err != nil {
		return fmt.Errorf(fatal(err.Error()))
	}
	results, err := checkPrefixes(prefix)
	if err != nil {
		return err
	}
	for idx := 0; idx < len(results); idx++ {
		fmt.Println(results[idx])
	}
	return nil
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name(os.Args[0]),
		kong.Description("print ip address subnet in-addr lines to stdout"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)
	err := ctx.Run(ctx)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}
