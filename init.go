package consulagt

import (
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var registeredSvc map[string]*ServiceMeta

// var defaultHTTPHealthyURL string

// var agentServices *cache.Cache = cache.New(10*time.Second, 3*time.Second) //进程缓存值 有效期10s,检查间隔5s

//DefaultHealthy set healthy check by default
// 关闭健康检查
const DefaultHealthy = "default-healthy"

func addRegistered(serviceID string, svcMeta *ServiceMeta) {
	registeredSvc[serviceID] = svcMeta
}

func markRegistered(serviceID string) {
	// registeredSvc[serviceID] = true
	if svc, ok := registeredSvc[serviceID]; ok {
		svc.RegistStatus = true
	}
}

func init() {
	registeredSvc = make(map[string]*ServiceMeta)
	http.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "pong")
	})
	go registering()
}
func registering() {
	exitSignalChan := make(chan os.Signal, 1)
	signal.Notify(exitSignalChan)
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-exitSignalChan:
			for serviceID, svcMeta := range registeredSvc {
				if svcMeta.RegistStatus {
					deregisterByServiceID(serviceID)
				}
			}
			return
		case <-ticker.C:
			for _, svcMeta := range registeredSvc {
				go func(meta *ServiceMeta) {
					regist2ConsulAgent(meta)
				}(svcMeta)
			}
			// caching services
			loadServices()
		}
	}
}

func loadServices() error {
	list, err := Services()
	if err != nil {
		return err
	}
	for _, s := range list {
		serviceCache.Store(s.ID, s)
	}
	return nil
}
