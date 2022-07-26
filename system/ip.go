package system

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

func GetPrintedIP() string {
	printedBuilder := strings.Builder{}
	publicIP := getPublicIP()
	innerIP := getlocalIPv4()
	printedBuilder.WriteString("public ip: ")
	printedBuilder.WriteString(publicIP)
	printedBuilder.WriteString("internal ip: ")
	printedBuilder.WriteString(innerIP + "\n")
	return printedBuilder.String()
}

// unused function
func getInterIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "获取内网ip失败"
	}
	defer conn.Close()

	res := conn.LocalAddr().String()
	res = strings.Split(res, ":")[0]
	return res
}

func getlocalIPv4() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err.Error()
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "获取内网ip失败"
}

func getPublicIP() string {
	c := http.Client{}
	r, _ := http.NewRequest(http.MethodGet, "http://ip.sb", nil)
	r.Header = http.Header{
		"User-Agent": {"curl/7.64.1"},
	}
	var resp *http.Response
	for i := 0; i < 10; i++ {
		resp, _ = c.Do(r)
		if resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(time.Microsecond * 100)
	}
	if resp.StatusCode != http.StatusOK {
		return "获取公网IP地址失败"
	}
	ip, _ := ioutil.ReadAll(resp.Body)
	return string(ip)
}
