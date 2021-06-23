package consulagt

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
)

const (
	CONSUL_ADDR  = "CONSUL_ADDR"
	CONSUL_TOKEN = "CONSUL_TOKEN"
	CONSUL_DC    = "CONSUL_DC"
	CONSUL_ZONE  = "CONSUL_ZONE"
	CONSUL_ENV   = "CONSUL_ENV"
)

//NewClient new conusl api client by env
func NewClient() (*api.Client, error) {
	config := api.DefaultConfig()
	consulAddr := "127.0.0.1:8500" // dfault
	consulAddr, ok := os.LookupEnv(CONSUL_ADDR)
	if ok { //set by env
		config.Address = consulAddr
	}
	consulToken, ok := os.LookupEnv(CONSUL_TOKEN)
	if ok {
		config.Token = consulToken
	}
	consulDc, ok := os.LookupEnv(CONSUL_DC)
	if ok {
		config.Datacenter = consulDc
	}
	config.HttpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   15 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxConnsPerHost:       64,
			MaxIdleConnsPerHost:   64,
			MaxIdleConns:          128,
		},
		Timeout: 30 * time.Second,
	}
	return api.NewClient(config)
}

var defaultCli *api.Client

func defaultClient() *api.Client {
	var err error
	if defaultCli == nil {
		defaultCli, err = NewClient()
		if err != nil {
			panic(err)
		}
	}
	return defaultCli
}

//Zone 获取环境变量的区信息 区分相同服务的不同集群
func Zone() string {
	Zone, ok := os.LookupEnv(CONSUL_ZONE)
	if ok {
		return Zone
	}
	return "UNKOWN"
}

//ENV 环境 TEST|PRODUCT|ALPHA|BETA
func ENV() string {
	environment, ok := os.LookupEnv(CONSUL_ENV)
	if ok {
		return environment
	}
	return "TEST"
}
