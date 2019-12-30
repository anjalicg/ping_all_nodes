package main

import (
	"flag"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"time"
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
	Loss        float64
	IsReachable bool
}

func IPToInt(ip_obj net.IP) int {
	//(first octet * 256³) + (second octet * 256²) + (third octet * 256) + (fourth octet)
	ip_bytes := ip_obj.To4()
	ip_integer := 0
	multiplier := 1
	for i := len(ip_bytes) - 1; i >= 0; i-- {
		ip_integer += (multiplier * int(ip_bytes[i]))
		multiplier *= 256
	}
	fmt.Printf("ip_integer=%v\n", ip_integer)
	return ip_integer

}

func IPIntToIPv4(ip_int int) net.IP {
	var ip_byte = []uint8{0, 0, 0, 0}
	for i := len(ip_byte) - 1; ip_int > 0; i-- {
		ip_mod := ip_int % 256
		ip_int = ip_int / 256
		ip_byte[i] = uint8(ip_mod)
	}
	return net.IP(ip_byte)
}

func pingIP(ip_str string) string {
	cmd_str := fmt.Sprintf("ping")
	output, _ := exec.Command(cmd_str, "-c 1", ip_str).Output()
	// if err != nil {
	// 	fmt.Println("err while ping ", err)
	// }

	return string(output)
}

func parsePingOut(output string) PingOutput {
	var ping_out PingOutput
	ping_out.Ouput = output
	var rgx = regexp.MustCompile(`(\d+)\spackets transmitted,\s(\d+) packets received, (\d+\.\d+)% packet loss`)
	matchFound := rgx.FindStringSubmatch(output)
	ping_out.Transmitted, _ = strconv.Atoi(matchFound[1])
	ping_out.Received, _ = strconv.Atoi(matchFound[2])
	ping_out.Loss, _ = strconv.ParseFloat(matchFound[3], 64)
	if ping_out.Received >= 1 {
		ping_out.IsReachable = true
	} else {
		ping_out.IsReachable = false
	}
	// fmt.Println(ping_out)
	return ping_out
}
func checkConnectivity(ip string) {
	output := pingIP(ip)
	ping_out := parsePingOut(output)
	if ping_out.IsReachable {
		fmt.Printf("ip=%s is reachable\n", ip)
	}
}
func checkConnectivityDial(ip string) {
	conn, err := net.Dial("ip:icmp", ip)
	if err != nil {
		if conn != nil {
			fmt.Printf("ip=%s is reachable\n", ip)
		}
	} else {
		fmt.Println("Error happened!!", err)
	}
}

func main() {
	fmt.Println("Ping all nodes")
	// subnet_str_ptr := flag.String("subnet", "", "Subnet ip e.g 192.168.0.0/24")
	// subnet_len_ptr := flag.String("len", "", "Subnet e.g like 24")
	flag.Parse()

	ip_obj, ipnet_obj, _ := net.ParseCIDR("192.168.86.1/24")
	fmt.Println("ip_obj", ip_obj)
	fmt.Println("ipnet_obj", ipnet_obj)
	fmt.Printf("%T andd %T\n", ip_obj, ipnet_obj)
	x := IPToInt(ip_obj)
	fmt.Println("int of ip_obj", x)
	str_ip := IPIntToIPv4(x)
	fmt.Println("ip_obj back to string", str_ip)
	// ip_int := IPIntToIPv4(ip_obj)
	for ip_int := IPToInt(ip_obj); ipnet_obj.Contains(ip_obj); {
		str_ip := IPIntToIPv4(ip_int)
		// fmt.Println("IP as string:", str_ip.String())
		go checkConnectivity(str_ip.String())
		// go checkConnectivity(str_ip.String())

		ip_int += 1
		ip_obj = IPIntToIPv4(ip_int)

	}
	// go checkConnectivity("8.8.8.8")
	// go checkConnectivity("192.168.2.1")
	time.Sleep(60 * time.Second)
}
