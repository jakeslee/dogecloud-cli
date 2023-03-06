package doge

type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
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

type OSSSource struct {
	Bucket int `json:"bucket"`
}

type DomainSource struct {
	Addr     string `json:"addr"`
	Base     string `json:"base"`
	Host     string `json:"host"`
	Protocol string `json:"protocol"`
}

type Source struct {
	Type string `json:"type"`
	OSSSource
	DomainSource
}
type Domains struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Cname       string `json:"cname"`
	Ctime       string `json:"ctime"`
	ServiceType string `json:"service_type"`
	Status      string `json:"status"`
	CertID      int    `json:"cert_id"`
	Source      Source `json:"source,omitempty"`
}
