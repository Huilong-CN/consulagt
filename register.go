package consulagt

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/hashicorp/consul/api"
)

//Register services to consul by default args
func Register(servicesName, httpHealthyAddr string, metas map[string]string, servicesAddrs ...string) error {
	for _, servicesAddr := range servicesAddrs {
		regist2ConsulAgent(&ServiceMeta{
			ServiceName:     servicesName,
			ServicesAddr:    servicesAddr,
			Metas:           metas,
			HTTPHealthyAddr: httpHealthyAddr,
		})
	}
	return nil
}

func regist2ConsulAgent(serviceMeta *ServiceMeta) error {
	consulCli, err := NewClient()
	if err != nil {
		log.Printf("register server error : %v", err)
		return err
	}
	addRegistered(serivcesID(serviceMeta.ServiceName, serviceMeta.ServicesAddr), serviceMeta)
	err = consulCli.Agent().ServiceRegister(genAgentServiceRegistration(serviceMeta.ServiceName, serviceMeta.ServicesAddr, serviceMeta.HTTPHealthyAddr, serviceMeta.Metas))
	if err != nil {
		log.Printf("regist services:%s servicesAddr:%s err:%v", serviceMeta.ServiceName, serviceMeta.ServicesAddr, err)
		return err
	}
	markRegistered(serivcesID(serviceMeta.ServiceName, serviceMeta.ServicesAddr))
	return nil
}

//Deregister services to consul by default args
func Deregister(servicesName, servicesAddr string) error {
	return deregisterByServiceID(serivcesID(servicesName, servicesAddr))
}

func deregisterByServiceID(serivcesID string) error {
	consulCli, err := NewClient()
	if err != nil {
		log.Printf("register server error : %+v ", err)
		return err
	}
	err = consulCli.Agent().ServiceDeregister(serivcesID)
	if err != nil {
		log.Printf("deregist servicesid:%s err:%v", serivcesID, err)
		return err
	}
	return nil
}

func convert(servicesAddr string) *url.URL {
	urlObj, err := url.Parse(servicesAddr)
	if err != nil {
		panic(err)
	}
	return urlObj
}
func port(servicesAddr string) int {
	obj := convert(servicesAddr)
	port, _ := strconv.Atoi(obj.Port())
	return port
}
func host(servicesAddr string) string {
	obj := convert(servicesAddr)
	return obj.Hostname()
}

func serivcesID(servicesName, servicesAddr string) string {
	urlobj := convert(servicesAddr)
	return fmt.Sprintf("%s#%s:%s", servicesName, urlobj.Hostname(), urlobj.Port())
}

func checkID(servicesName, servicesAddr string) string {
	return serivcesID(servicesName, servicesAddr) + "#healthy"
}

func genAgentServiceRegistration(servicesName, servicesAddr, httpPingAddr string, metas map[string]string) *api.AgentServiceRegistration {
	urlobj := convert(servicesAddr)
	var checker *api.AgentServiceCheck
	if httpPingAddr != DefaultHealthy {
		checker = &api.AgentServiceCheck{
			CheckID:                        checkID(servicesName, servicesAddr),
			Name:                           servicesName + "-healthy",
			Status:                         "passing",
			Interval:                       "5s",
			Timeout:                        "3s",
			HTTP:                           httpPingAddr,
			DeregisterCriticalServiceAfter: "30s",
		}
	}
	_metas := make(map[string]string)
	_metas["scheme"] = urlobj.Scheme
	_metas["zone"] = Zone()
	_metas["evn"] = ENV()
	for key, val := range metas {
		_metas[key] = val
	}
	_tags := make([]string, 0, len(_metas))
	for _, val := range _metas {
		_tags = append(_tags, val)
	}
	return &api.AgentServiceRegistration{
		Kind:    api.ServiceKindTypical,
		ID:      serivcesID(servicesName, servicesAddr),
		Name:    servicesName,
		Tags:    _tags,
		Meta:    _metas,
		Port:    port(servicesAddr),
		Address: urlobj.Hostname(),
		Check:   checker,
	}
}
