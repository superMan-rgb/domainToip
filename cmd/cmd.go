package cmd

import (
	"bufio"
	o "domaintoip/other"
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
	usage = `-d string 
        target  www.baidu.com
	-f string
	target domain.txt`
)

type FlagConfig struct {
	domain string
	file   string
}

func (c *FlagConfig) ParseFlags() {
	fmt.Println(banner)
	flag.StringVar(&c.domain, "d", "", "target www.baidu.com")
	flag.StringVar(&c.file, "f", "", "target domain.txt")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: \n  %s [options]\n", filepath.Base(os.Args[0]))
		fmt.Println("Options:")
		fmt.Print(usage)
	}
	flag.Parse()
	execdomain(*c)
}
func execdomain(c FlagConfig) {

	if c.domain == "" && c.file == "" {
		fmt.Println("请提供单个域名或域名文件")
		os.Exit(1)
	}
	if c.domain != "" {
		ips, err := net.LookupIP(c.domain)
		// fmt.Print(ips)
		if err != nil {
			fmt.Printf("域名解析失败 %s - %v\n", c.domain, err)
			os.Exit(1)
		}
		ipsR := o.DeduplicateIPs(ips)
		for _, ip := range ipsR {
			o.result(ip.String(), c.domain)
		}
	}
	if c.file != "" {
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
		for scanner.Scan() {
			domain := scanner.Text()

			ips, err := net.LookupIP(domain)
			if err != nil {
				fmt.Printf("域名解析失败: %s - %v\n", domain, err)
				continue
			}
			ipsR := o.DeduplicateIPs(ips)
			domain1 = &domain
			for _, ip := range ipsR {
				o.result(ip.String(), domain1)
			}
		}

		if scanner.Err() != nil {
			fmt.Println("读取文件失败:", scanner.Err())
			os.Exit(1)
		}

	}
}
