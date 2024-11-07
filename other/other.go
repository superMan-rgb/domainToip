package other

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 去重ips
func DeduplicateIPs(ips []net.IP) []net.IP {
	uniqueIPs := []net.IP{}
	visited := make(map[string]bool)

	for _, ip := range ips {
		ipString := ip.String()
		if !visited[ipString] {
			visited[ipString] = true
			uniqueIPs = append(uniqueIPs, ip)
		}
	}
	return uniqueIPs
}

// 判断是否为CDN
func IsCDN(domain string) string {
	cname, _ := net.LookupCNAME(domain)
	//fmt.Print(cname)
	cname = strings.TrimSuffix(cname, ".")
	domain = strings.TrimSuffix(domain, ".")
	if cname != domain {
		return "true"
	}
	return "false"
}

// ip138查归属
func GetIPLocation(ip string) (string, string, string, error) {
	// 发起HTTP GET请求
	url := fmt.Sprintf("https://www.ipshudi.com/%s.htm", ip)
	// fmt.Print(url)
	// 创建HTTP客户端
	client := &http.Client{}
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", "", err
	}
	// 设置User-Agent头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.203")
	// 发起请求
	response, err := client.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer response.Body.Close()

	// 使用goquery解析HTML
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 提取归属地信息
	location := doc.Find("td.th:contains('归属地') + td span a").Text()

	// 提取运营商信息
	provider := doc.Find("td.th:contains('运营商') + td span").Text()

	// 提取IP类型信息
	ipType := doc.Find("td.th:contains('iP类型') + td span").Text()

	// 去除提取到的信息中的空格和换行符
	location = strings.TrimSpace(location)
	provider = strings.TrimSpace(provider)
	ipType = strings.TrimSpace(ipType)
	return location, provider, ipType, err
}

// CountIPsByCIDR 将给定的IP地址列表按C段地址进行归类，并返回每个C段地址中的IP数量
func CountIPsByCIDR(ipAddresses []net.IP) map[string]int {
	cidrMap := make(map[string]int)

	for _, ip := range ipAddresses {
		if ip == nil {
			fmt.Printf("Invalid IP address\n")
			continue
		}

		// 获取IP的C段地址
		cidr := getCIDR(ip)
		if cidr == "" {
			fmt.Printf("Failed to get CIDR for IP address: %s\n", ip.String())
			continue
		}

		// 增加C段地址的IP数量计数
		cidrMap[cidr]++
	}

	return cidrMap
}

// 获取IP的C段地址
func getCIDR(ip net.IP) string {
	ip = ip.To4()
	if ip == nil {
		return ""
	}

	// 将IP地址的前三个八位数字转换为C段地址
	cidr := strconv.Itoa(int(ip[0])) + "." + strconv.Itoa(int(ip[1])) + "." + strconv.Itoa(int(ip[2]))

	return cidr
}
