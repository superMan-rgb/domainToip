package other

import (
	"domain/plugin"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"strconv"
	"strings"
)

type Req struct {
	HttpMethod  string
	ApiPath     string
	Body        interface{}
	QueryParams *plugin.QueryParams
	//PathParams  *plugin.PathParams
}

type Result struct {
	Query string              `json:"query"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
	Data  []map[string]string `json:"items"`
}

type allFields struct {
	Ip              string `json:"ip,omitempty"`               //ip	    			ip地址				权限：无
	Port            string `json:"port,omitempty"`             //port				端口					权限：无
	Protocol        string `json:"protocol,omitempty"`         //protocol			协议名				权限：无
	Country         string `json:"country,omitempty"`          //country			国家代码				权限：无
	CountryName     string `json:"country_name,omitempty"`     //country_name		国家名				权限：无
	Region          string `json:"region,omitempty"`           //region				区域					权限：无
	City            string `json:"city,omitempty"`             //city				城市					权限：无
	Longitude       string `json:"longitude,omitempty"`        //longitude			地理位置 经度			权限：无
	Latitude        string `json:"latitude,omitempty"`         //latitude			地理位置 纬度			权限：无
	AsNumber        string `json:"as_number,omitempty"`        //as_number			asn编号				权限：无
	AsOrganization  string `json:"as_organization,omitempty"`  //as_organization	asn组织				权限：无
	Host            string `json:"host,omitempty"`             //host				主机名				权限：无
	Domain          string `json:"domain,omitempty"`           //domain				域名					权限：无
	Os              string `json:"os,omitempty"`               //os				    操作系统				权限：无
	Server          string `json:"server,omitempty"`           //server				网站server			权限：无
	Icp             string `json:"icp,omitempty"`              //icp				icp备案号			权限：无
	Title           string `json:"title,omitempty"`            //title				网站标题				权限：无
	Jarm            string `json:"jarm,omitempty"`             //jarm				jarm 指纹			权限：无
	Header          string `json:"header,omitempty"`           //header				网站header			权限：无
	Banner          string `json:"banner,omitempty"`           //banner				协议 banner			权限：无
	Cert            string `json:"cert,omitempty"`             //cert				证书					权限：无
	BaseProtocol    string `json:"base_protocol,omitempty"`    //base_protocol		基础协议，比如tcp/udp	权限：无
	Link            string `json:"link,omitempty"`             //link				资产的URL链接			权限：无
	Product         string `json:"product,omitempty"`          //product			产品名				权限：专业版本及以上
	ProductCategory string `json:"product_category,omitempty"` //product_category   产品分类			    权限：专业版本及以上
	Version         string `json:"version,omitempty"`          //version			版本号				权限：专业版本及以上
	LastUpdateTime  string `json:"lastupdatetime,omitempty"`   //lastupdatetime		FOFA最后更新时间	    权限：专业版本及以上
	Cname           string `json:"cname,omitempty"`            //cname				域名cname			权限：专业版本及以上
	IconHash        string `json:"icon_hash,omitempty"`        //icon_hash		    返回的icon_hash值	权限：商业版本及以上
	CertsValid      string `json:"certs_valid,omitempty"`      //certs_valid		证书是否有效			权限：商业版本及以上
	CnameDomain     string `json:"cname_domain,omitempty"`     //cname_domain		cname的域名			权限：商业版本及以上
	Body            string `json:"body,omitempty"`             //body				网站正文内容			权限：商业版本及以上
	Icon            string `json:"icon,omitempty"`             //icon				icon 图标			权限：企业会员
	Fid             string `json:"fid,omitempty"`              //fid	     		fid					权限：企业会员
	Structinfo      string `json:"structinfo,omitempty"`       //structinfo			结构化信息 (部分协议支持、比如elastic、mongodb)	权限：企业会员
}

func (f *Fofa) Get(req *GetDataReq) (*Result, error) {
	//NewClient("root@dnslog.vip", "b7bdbebd6a48c1fcd0a92ee6596d99f5")
	req.req.QueryParams.Set("email", f.email)
	req.req.QueryParams.Set("key", f.key)
	url := fmt.Sprintf("%v?%s", plugin.FofaApiUrl, req.req.QueryParams.Encode())
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := f.http.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, &plugin.NonStatusOK{}
	}
	//body, err := ghttp.GetResponseBody(response.Body)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("failed to read response body", err)
		return nil, err
	}
	var tmpResponse struct {
		Error           bool   `json:"error"`
		ErrMsg          string `json:"errmsg"`
		ConsumedFpoint  int    `json:"consumed_fpoint"`
		RequiredFpoints int    `json:"required_fpoints"`
		Size            int64  `json:"size"`
		Page            int    `json:"page"`
		Mode            string `json:"mode"`
		Query           string `json:"query"`
	}
	err = json.Unmarshal(body, &tmpResponse)
	if err != nil {
		return nil, err
	}
	if tmpResponse.Error {
		return nil, errors.New(tmpResponse.ErrMsg)
	}

	items, _, _, err := jsonparser.Get(body, "results")
	if items == nil {
		return nil, nil
	}

	// 把不带键的 [[],[]...] 转成带键的 [{},{}...]
	fields := req.req.QueryParams.Get("fields")
	formatItems, err := f.formatItems(items, fields)
	if err != nil {
		return nil, err
	}

	var resp Result
	resp.Data = formatItems
	resp.Query = tmpResponse.Query
	resp.Page = tmpResponse.Page
	size, _ := strconv.Atoi(req.req.QueryParams.Get("size"))
	resp.Size = size
	resp.Total = tmpResponse.Size
	return &resp, nil
}

// User struct for fofa user
type User struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Category        string `json:"category"`
	Fcoin           int    `json:"fcoin"`
	FofaPoint       int    `json:"fofa_point"`
	RemainFreePoint int    `json:"remain_free_point"`
	RemainAPIQuery  int    `json:"remain_api_query"`
	RemainAPIData   int    `json:"remain_api_data"`
	IsVip           bool   `json:"isvip"`
	VipLevel        int    `json:"vip_level"`
	IsVerified      bool   `json:"is_verified"`
	Avatar          string `json:"avatar"`
	Message         string `json:"message"`
	FofacliVer      string `json:"fofacli_ver"`
	FofaServer      bool   `json:"fofa_server"`
}

func (f *Fofa) User() (*User, error) {
	params := netUrl.Values{}
	params.Add("Email", f.email)
	params.Add("key", f.key)
	url := fmt.Sprintf("%v?%s", plugin.FofaUserApiUrl, params.Encode())
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := f.http.Do(request)
	if err != nil {
		return nil, err
	}
	//body, err := ghttp.GetResponseBody(response.Body)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("failed to read response body", err)
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	var tmpResponse struct {
		User
		Error  bool   `json:"error"`
		ErrMsg string `json:"errmsg,omitempty'"`
	}
	if err = json.Unmarshal(body, &tmpResponse); err != nil {
		return nil, err
	}
	if tmpResponse.Error {
		return nil, errors.New(tmpResponse.ErrMsg)
	}
	return &tmpResponse.User, nil
}

// fofa查询结果是不带键的，会按照你查询时给定的fields的先后顺序输出，每个item是 [] 此处是把每个item附上键变成 {"field":"value"}，键名为fields的每个字段
func (f *Fofa) formatItems(items []byte, fields string) ([]map[string]string, error) {
	var keys []string
	var newItems []map[string]string
	//对应两种返回结果 field不唯一和field唯一的情况
	var tmpListList [][]string //field不唯一
	var tmpStringList []string //field唯一
	err1 := json.Unmarshal(items, &tmpListList)
	err2 := json.Unmarshal(items, &tmpStringList)
	if err1 != nil && err2 != nil {
		return nil, errors.New("format fofa items failed：" + err1.Error())
	}
	var listLen int
	if err1 == nil {
		listLen = len(tmpListList)
	} else {
		listLen = len(tmpStringList)
	}
	for _, v := range strings.Split(fields, ",") {
		if v != "" {
			keys = append(keys, v)
		}
	}

	//fields 为空 fofa默认按照 host,ip,port先后顺序返回
	if len(keys) == 0 {
		for i := 0; i < listLen; i++ {
			var tmpItem = map[string]string{}
			for index, key := range []string{"host", "ip", "port"} {
				tmpItem[key] = tmpListList[i][index]
			}
			newItems = append(newItems, tmpItem)
		}
	} else if len(keys) == 1 { //传入的field只有一个的话，只会返回 字符串列表 不会返回 字典列表
		for i := 0; i < listLen; i++ {
			var _item = map[string]string{}
			_item[keys[0]] = tmpStringList[i]
			newItems = append(newItems, _item)
		}
	} else {
		for i := 0; i < listLen; i++ {
			var _item = map[string]string{}
			for j := 0; j < len(keys); j++ {
				if j > len(tmpListList[i]) {
					_item[keys[j]] = ""
				} else {
					_item[keys[j]] = tmpListList[i][j]
				}
			}
			newItems = append(newItems, _item)
		}
	}
	return newItems, nil
}
