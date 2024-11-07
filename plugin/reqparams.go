package plugin

import "net/url"

type PathParams map[string]string

func (u PathParams) Get(key string) string {
	vs := u[key]
	if len(vs) == 0 {
		return ""
	}
	return vs
}
func (u PathParams) Set(key, value string) {
	u[key] = value
}

type QueryParams map[string][]string

func (u QueryParams) Get(key string) string {
	vs := u[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}
func (u QueryParams) Set(key, value string) {
	u[key] = []string{value}
}

func (u QueryParams) Encode() string {
	return url.Values(u).Encode()
}

func (u QueryParams) Add(key, value string) {
	u[key] = append(u[key], value)
}

var FofaApiUrl = "https://fofa.info/api/v1/search/all"
var FofaUserApiUrl = "https://fofa.info/api/v1/info/my"

type NonStatusOK struct {
	// 定义类型的字段（如果有需要）
}

func (e NonStatusOK) Error() string {
	return "NonStatusOK error occurred"
}
