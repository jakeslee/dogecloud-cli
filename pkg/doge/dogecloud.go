package doge

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	AccessKey string
	SecretKey string
)

func Fetch(apiPath string, data map[string]interface{}, jsonMode bool, result interface{}) {
	// 这里替换为你的多吉云永久 AccessKey 和 SecretKey，可在用户中心 - 密钥管理中查看
	// 请勿在客户端暴露 AccessKey 和 SecretKey，那样恶意用户将获得账号完全控制权

	body := ""
	mime := ""
	if jsonMode {
		_body, err := json.Marshal(data)
		if err != nil {
			log.Fatalln(err)
		}
		body = string(_body)
		mime = "application/json"
	} else {
		values := url.Values{}
		for k, v := range data {
			values.Set(k, v.(string))
		}
		body = values.Encode()
		mime = "application/x-www-form-urlencoded"
	}

	signStr := apiPath + "\n" + body
	hmacObj := hmac.New(sha1.New, []byte(SecretKey))
	hmacObj.Write([]byte(signStr))
	sign := hex.EncodeToString(hmacObj.Sum(nil))
	Authorization := "TOKEN " + AccessKey + ":" + sign

	req, err := http.NewRequest("POST", "https://api.dogecloud.com"+apiPath, strings.NewReader(body))
	req.Header.Add("Content-Type", mime)
	req.Header.Add("Authorization", Authorization)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	} // 网络错误
	defer resp.Body.Close()
	r, err := io.ReadAll(resp.Body)

	json.Unmarshal(r, result)

	// Debug，正式使用时可以注释掉
	// fmt.Printf("[DogeCloudAPI] code: %d, msg: %s, data: %s\n", int(ret["code"].(float64)), ret["msg"], ret["data"])
	return
}

// FetchResult 调用多吉云的 API
// apiPath：是调用的 API 接口地址，包含 URL 请求参数 QueryString，例如：/console/vfetch/add.json?url=xxx&a=1&b=2
// data：POST 的数据，对象，例如 {a: 1, b: 2}，传递此参数表示不是 GET 请求而是 POST 请求
// jsonMode：数据 data 是否以 JSON 格式请求，默认为 false 则使用表单形式（a=1&b=2）
// 返回值 ret 是一个 map[string]，其中 ret["code"] 为 200 表示 api 请求成功
func FetchResult(apiPath string, data map[string]interface{}, jsonMode bool) (ret map[string]interface{}) {
	Fetch(apiPath, data, jsonMode, &ret)
	return
}

func ListDomains() ([]string, error) {
	r := FetchResult("/cdn/domain/list.json", make(map[string]interface{}), false)
	data := r["data"].(map[string]interface{})
	domains := data["domains"].([]interface{})

	var ret []string

	for _, d := range domains {
		dm := d.(map[string]interface{})
		if dm["status"] == "online" {
			ret = append(ret, dm["name"].(string))
		}

		log.Printf("[%d]%-23sstatus: %s\t ctime: %s\n", int(dm["id"].(float64)), dm["name"], dm["status"], dm["ctime"])
	}

	return ret, nil
}

func UploadCert(name, cert, key string) (int, error) {
	params := make(map[string]interface{})
	params["note"] = name
	params["cert"] = cert
	params["private"] = key

	r := FetchResult("/cdn/cert/upload.json", params, false)

	if r["code"].(float64) != 200 {
		return 0, fmt.Errorf("error: %s", r["msg"])
	}

	return int(r["data"].(map[string]interface{})["id"].(float64)), nil
}

type Cert struct {
	Id     int       `json:"id"`
	Name   string    `json:"name"`
	Count  int       `json:"domainCount"`
	Expire int       `json:"expire"`
	Info   *CertInfo `json:"info"`
}

type CertInfo struct {
	SAN    []string `json:"SAN"`
	Type   string   `json:"type"`
	Period string   `json:"period"`
}

func ListCerts() []*Cert {
	type Body struct {
		Data struct {
			Certs []*Cert `json:"certs"`
		} `json:"data"`
	}
	var result Body

	Fetch("/cdn/cert/list.json", make(map[string]interface{}), false, &result)

	for _, cert := range result.Data.Certs {
		log.Printf("[%d]%s SAN: %s", cert.Id, cert.Name, strings.Join(cert.Info.SAN, ","))
	}

	return result.Data.Certs
}

func DeleteCert(id int) {
	params := make(map[string]interface{})
	params["id"] = strconv.Itoa(id)

	r := FetchResult("/cdn/cert/delete.json", params, false)
	log.Printf("delete cert %d, result: %v", id, r)
}

func DomainCertDeploy(certId int, domains []string) {
	params := make(map[string]interface{})
	params["cert_id"] = strconv.Itoa(certId)

	for _, domain := range domains {
		path := fmt.Sprintf("/cdn/domain/config.json?domain=%s", domain)

		r := FetchResult(path, params, true)
		log.Printf("deploy cert name: %-23scert: %v\tresult: %v", domain, certId, r)
	}
}
