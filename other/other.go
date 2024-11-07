package other

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
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

// 辨别是否为cdn
func IsIPCDN(ip string) bool {
	url := fmt.Sprintf("https://www.ipip.net/ip/%s.html", ip)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败:", err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return false
	}

	html := string(body)

	// 在HTML中查找特定的关键词，判断是否属于CDN
	isCDN := strings.Contains(html, "CDN")
	return isCDN
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
