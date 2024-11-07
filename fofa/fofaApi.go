package fofa

import (
	"domain/other"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func FofaApi(domain, email, key string) []string {
	//email, key := fofa.Config()
	var client = other.NewClient(email, key)
	builder := other.NewGetDataReqBuilder().Query(domain).Fields("host").Build()
	result, err := client.Get(builder)
	if err != nil {
		return nil
	}

	jsonstr, err := json.Marshal(result.Data)
	//_, res, err := bufio.ScanLines(jsonstr, true)
	//err = json.Unmarshal(jsonstr, &pro)
	if err != nil {
		fmt.Println("错误")
		return nil
	}
	//var res = string(jsonstr)
	domainNames, err := extractDomainNameFromJSON(string(jsonstr))
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return domainNames

}

type Host struct {
	Host string `json:"host"`
}

func extractDomainNameFromJSON(jsonStr string) ([]string, error) {
	var hosts []Host
	err := json.Unmarshal([]byte(jsonStr), &hosts)
	if err != nil {
		return nil, err
	}

	var domainNames []string
	visited := make(map[string]bool)
	for _, host := range hosts {
		if checkHTTP(host.Host) {
			doma := extractDomain(host.Host)
			if !visited[doma] {
				domainNames = append(domainNames, doma)
				visited[doma] = true
			}
		} else {
			doma1 := "http://" + host.Host
			//fmt.Println(doma1)
			doma := extractDomain(doma1)
			// 检查域名是否已经存在于visited中，如果不存在则添加到domainNames中
			if !visited[doma] {
				domainNames = append(domainNames, doma)
				visited[doma] = true
			}
		}

	}
	return domainNames, err

}

// 使用net/url包来提取域名
func extractDomain(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	return u.Hostname()
}

// 检查头是否为http或https
func checkHTTP(str string) bool {
	return strings.Contains(str, "http") || strings.Contains(str, "https")
}
