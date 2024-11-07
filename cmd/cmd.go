package cmd

import (
	"bufio"
	"domain/fofa"
	"domain/other"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
)

const (
	banner = ` 
		      .___                    .__     ___________    .__         
		    __| _/____   _____ _____  |__| ___\__    ___/___ |__|_____   
		   / __ |/  _ \ /     \\__  \ |  |/    \|    | /  _ \|  \____ \  
		  / /_/ (  <_> )  Y Y  \/ __ \|  |   |  \    |(  <_> )  |  |_> >  
		  \____ |\____/|__|_|  (____  /__|___|  /____| \____/|__|   __/  
		       \/            \/     \/        \/                |__|     `
	usage = `
	-d string
		target  www.baidu.com
	-f string
		target  baidu.com
	-l string
		target domain.txt
	-h help  
`
)

type FlagConfig struct {
	domain     string
	fofaDomain string
	file       string
	help       string
}

func (c *FlagConfig) ParseFlags() {
	fmt.Println(banner)
	fofa.IsConfig()
	flag.StringVar(&c.domain, "d", "", "target www.baidu.com")
	flag.StringVar(&c.fofaDomain, "f", "", "target baidu.com")
	flag.StringVar(&c.file, "l", "", "target domain.txt")
	flag.StringVar(&c.help, "h", "", "help")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: \n  %s [options]\n", filepath.Base(os.Args[0]))
		fmt.Println("Options:")
		fmt.Print(usage)
	}
	flag.Parse()
	//var set []string
	if c.help != "" {
		flag.Usage()
		os.Exit(0)
	} else if c.domain == "" && c.file == "" && c.fofaDomain == "" {
		fmt.Println(fmt.Sprintf("Error: required flag(s) not set"))
		flag.Usage()
		os.Exit(1)
	}
	execdomain(*c)
	if c.domain == "" {
		println("以下是c段聚合：")
		result_cidr(ipsR1)
	}

}

var ipsR1 []net.IP

func execdomain(c FlagConfig) {

	if c.domain == "" && c.file == "" && c.fofaDomain == "" {
		fmt.Println("\n请提供单个域名或域名文件")
		os.Exit(1)
	}
	if c.domain != "" && c.fofaDomain == "" && c.file == "" {
		ips, err := net.LookupIP(c.domain)
		// fmt.Print(ips)
		if err != nil {
			fmt.Printf("域名解析失败 %s - %v\n", c.domain, err)
			os.Exit(1)
		}
		ipsR := other.DeduplicateIPs(ips)
		for _, ip := range ipsR {
			result(ip.String(), &c.domain)
		}
	}
	if c.file != "" && c.domain == "" && c.fofaDomain == "" {
		//打开文件
		file, err := os.Open(c.file)
		if err != nil {
			fmt.Println("打开文件失败:", err)
			os.Exit(1)
		}
		defer file.Close()
		// 创建扫描器
		scanner := bufio.NewScanner(file)
		// 逐行读取文件
		var domain1 *string
		//var ipsR1 []net.IP
		for scanner.Scan() {
			domain := scanner.Text()
			iscdns := other.IsCDN(domain)
			if iscdns == "true" {
				continue
			}

			ips, err := net.LookupIP(domain)
			if err != nil {
				fmt.Printf("域名解析失败: %s - %v\n", domain, err)
				continue
			}
			ipsR := other.DeduplicateIPs(ips)
			//fmt.Println(ipsR1)
			//cidrCounts := other.CountIPsByCIDR(ipsR)

			domain1 = &domain
			for _, ip := range ipsR {
				result(ip.String(), domain1)
				ipsR1 = append(ipsR1, ip)
			}

		}

		if scanner.Err() != nil {
			fmt.Println("读取文件失败:", scanner.Err())
			os.Exit(1)
		}

	}
	if c.fofaDomain != "" && c.domain == "" && c.file == "" {
		domain := fmt.Sprintf("domain=\"%s\"", c.fofaDomain)
		email, key := fofa.Config()
		domainNames := fofa.FofaApi(domain, email, key)
		for _, doma := range domainNames {
			iscdns := other.IsCDN(doma)
			if iscdns == "true" {
				continue
			}

			ips, err := net.LookupIP(doma)
			if err != nil {
				fmt.Printf("域名解析失败: %s - %v\n", doma, err)
				continue
			}
			ipsR := other.DeduplicateIPs(ips)
			//fmt.Println(ipsR1)
			//cidrCounts := other.CountIPsByCIDR(ipsR)

			//domain1 = &doma
			for _, ip := range ipsR {
				result(ip.String(), &doma)
				ipsR1 = append(ipsR1, ip)
			}

		}
	}
}

// 判断输出
func result(ip string, domainptr *string) {
	iscdn := other.IsCDN(*domainptr)
	//fmt.Println(iscdn)
	//var iscdns string
	//if iscdn {
	//	iscdns = "true"
	//} else {
	//	iscdns = "false"
	//}
	location, operator, ipType, err := other.GetIPLocation(ip)
	if err != nil {
		fmt.Print("查询失败：", err)
	}
	fmt.Printf("domain: %s ip: %s iscdn: %s country: %s %s iptype: %s\n", *domainptr, ip, iscdn, location, operator, ipType)
	//return *domainptr, ip, iscdns, location, operator, ipType

}

// 打印c段
func result_cidr(ipsR []net.IP) {
	cidrCounts := other.CountIPsByCIDR(ipsR)

	// 打印归类结果和IP数量
	for cidr, count := range cidrCounts {
		fmt.Printf("C段地址: %s.0/24  ", cidr)
		fmt.Printf("IP数量: %d\n", count)
	}
}
