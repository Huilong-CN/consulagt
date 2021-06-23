package consulagt

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

const defaultSize = 5

//SvrAddrList 获取对应服务名的所有服务器
func SvrAddrList(svrName string) []string {
	svrlist := make([]string, 0, defaultSize)
	serviceCache.Range(func(_, val interface{}) bool {
		if agentSvr, ok := val.(*api.AgentService); ok && agentSvr.Service == svrName {
			svrlist = append(svrlist, connectStr(agentSvr))
		}
		return true
	})
	return svrlist
}

//SvrAddrWithZone 获取对应服务名以及大区的的所有服务器
func SvrAddrWithZone(svrName, zone string) []string {
	svrlist := make([]string, 0, defaultSize)
	serviceCache.Range(func(_, val interface{}) bool {
		if agentSvr, ok := val.(*api.AgentService); ok && agentSvr.Service == svrName && agentSvr.Meta["zone"] == zone {
			svrlist = append(svrlist, connectStr(agentSvr))
		}
		return true
	})
	return svrlist
}

//SvrAddrWithTags  获取对应服务名以及命中所有tag的的所有服务器
func SvrAddrWithTags(svrName string, tags ...string) []string {
	svrlist := make([]string, 0, defaultSize)
	serviceCache.Range(func(_, val interface{}) bool {
		if agentSvr, ok := val.(*api.AgentService); ok && agentSvr.Service == svrName && checkTags(agentSvr.Tags, tags) {
			svrlist = append(svrlist, connectStr(agentSvr))
		}
		return true
	})
	return svrlist
}

func checkTags(tags, targetTags []string) bool {
	for _, tt := range targetTags {
		if !hasTags(tags, tt) {
			return false
		}
	}
	return true
}

func hasTags(tags []string, target string) bool {
	for _, tag := range tags {
		if tag == target {
			return true
		}
	}
	return false
}

func connectStr(svr *api.AgentService) string {
	return fmt.Sprintf("%s:%d", svr.Address, svr.Port)
}
