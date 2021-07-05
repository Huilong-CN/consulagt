package consulagt

import (
	"io"
	"net/http"
	"sync"
	"time"
)

//DefaultHealthy set healthy check by default
// 关闭健康检查
const DefaultHealthy = "default-healthy"

func addRegistered(serviceID string, svcMeta *ServiceMeta) {
	registeredSvc.Store(serviceID, svcMeta)
}

func markRegistered(serviceID string) {
	// registeredSvc[serviceID] = true
	if svc, ok := registeredSvc.Get(serviceID); ok {
		svc.RegistStatus = true
	}
}

func init() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "pong")
	})
	go registering()
}
func registering() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			registeredSvc.RangeDo(func(key interface{}, meta *ServiceMeta) {
				go func(meta *ServiceMeta) {
					regist2ConsulAgent(meta)
				}(meta)
			})
			loadServices()
		}
	}
}

func loadServices() error {
	newServiceCache := &ServiceStore{&sync.Map{}}
	list, err := Services()
	if err != nil {
		return err
	}
	for _, s := range list {
		newServiceCache.Store(s.ID, s)
	}
	serviceCache = newServiceCache
	return nil
}
