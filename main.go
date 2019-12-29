package main

import (
	"flag"
	"fmt"
	"net"
	"os/exec"
	"regexp"
)

/*
type IPNet struct {
    IP   IP     // network number
    Mask IPMask // network mask
}
func ParseCIDR(s string) (IP, *IPNet, error)
ParseCIDR parses s as a CIDR notation IP address and prefix length, like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291.
func (*IPNet) Contains
func (n *IPNet) Contains(ip IP) bool
Contains reports whether the network includes ip.

func (*IPNet) Network
func (n *IPNet) Network() string
Network returns the address's network name, "ip+net".

func (*IPNet) String
func (n *IPNet) String() string
*/
type PingOutput struct {
	Ouput       string
	Transmitted int
	Received    int
	Loss        int
	Status      bool
}

func IPToInt(ip_obj net.IP) int {
	//(first octet * 256³) + (second octet * 256²) + (third octet * 256) + (fourth octet)
	ip_bytes := ip_obj.To4()
	ip_integer := 0
	multiplier := 1
	for i := len(ip_bytes) - 1; i >= 0; i-- {
		// fmt.Printf("%v and %T\n", ip_bytes[i], ip_bytes[i])
		ip_integer += (multiplier * int(ip_bytes[i]))
		multiplier *= 256

	}
	fmt.Printf("ip_integer=%v\n", ip_integer)
	return ip_integer

}

func IPIntToIPv4(ip_int int) net.IP {
	var ip_byte = []uint8{0, 0, 0, 0}
	// fmt.Printf("ip_int=%d\n", ip_int)
	for i := len(ip_byte) - 1; ip_int > 0; i-- {
		ip_mod := ip_int % 256
		ip_int = ip_int / 256
		ip_byte[i] = uint8(ip_mod)
		// fmt.Println(ip_mod, ip_int)

	}

	return net.IP(ip_byte)
}

func pingIP(ip_str string) {
	cmd_str := fmt.Sprintf("ping")
	output, err := exec.Command(cmd_str, "-c 1", ip_str).Output()

	fmt.Println(err)
	fmt.Println(string(output))
}

func parsePingOut(output string) PingOutput {
	/*
			PING 192.168.30.1 (192.168.30.1): 56 data bytes

		--- 192.168.30.1 ping statistics ---
		1 packets transmitted, 0 packets received, 100.0% packet loss
	*/
	output = `PING 192.168.30.1 (192.168.30.1): 56 data bytes

	--- 192.168.30.1 ping statistics ---
	1 packets transmitted, 0 packets received, 100.0% packet loss`
	var ping_out PingOutput
	ping_out.Ouput = output

	var rgx = regexp.MustCompile(`(\d+)\spackets transmitted,\s(\d+) packets received, (\d+\.\d+)% packet loss`)
	matchFound := rgx.Find([]byte(output))
	fmt.Println("matchfound...")
	fmt.Println(string(matchFound))

	return ping_out

}

func main() {
	fmt.Println("Ping all nodes")
	// subnet_str_ptr := flag.String("subnet", "", "Subnet ip e.g 192.168.0.0/24")
	// subnet_len_ptr := flag.String("len", "", "Subnet e.g like 24")
	flag.Parse()

	ip_obj, ipnet_obj, _ := net.ParseCIDR("192.0.2.1/24")
	fmt.Println(ip_obj)
	fmt.Println(ipnet_obj)
	fmt.Printf("%T andd %T\n", ip_obj, ipnet_obj)
	x := IPToInt(ip_obj)
	fmt.Println(x)
	str_ip := IPIntToIPv4(x)
	fmt.Println(str_ip)
	// pingIP("192.168.30.1")
	parsePingOut("something")

	// fmt.Printf("Subnet=%v and subnet len=%v\n", *subnet_str_ptr, *subnet_len_ptr)
	// // fmt.Println(net.ParseIP("192.0.2.1")[3]++)
	// ip_bytes := net.ParseIP("192.0.2.1").To4()
	// fmt.Println(ip_bytes[0])
	// fmt.Println(ip_bytes[1])
	// fmt.Println(ip_bytes[2])
	// fmt.Println(ip_bytes[3])
	//(first octet * 256³) + (second octet * 256²) + (third octet * 256) + (fourth octet)
}
