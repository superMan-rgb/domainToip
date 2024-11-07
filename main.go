package main

import "domain/cmd"

func main() {
	config := cmd.FlagConfig{}
	config.ParseFlags()
	//config := cmd.FlagConfig{}
	//config.ParseFlags()
	//other.Fofa_search()
	//other.NewClient("root@dnslog.vip", "b7bdbebd6a48c1fcd0a92ee6596d99f5")
	//fofa3 := &other.Fofa{}
	//req := &other.GetDataReq{}
	//result, err := fofa3.Get(req)
	//if err != nil {
	//	fmt.Println("error", err)
	//	return
	//}
	//fmt.Println(result)
	//fofa.FofaApi()

	// 打印提取到的域名
	//visited := make(map[string]bool)
	//for _, domain := range domainNames {
	//	domain1 := strings.Index(domain, ":")
	//	if domain1 != -1 {
	//		domain = domain[:domain1]
	//	}
	//	if !visited[domain] {
	//		fmt.Println(domain)
	//		visited[domain] = true
	//	}
	//
	//}

}
