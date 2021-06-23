package consulagt

import (
	"github.com/hashicorp/consul/api"
)

// ServicesByName returns for a given service name: the aggregated health status for all services
// having the specified name.
// - If no service is not found, will return status (critical, [], nil)
// - If the service is found, will return (critical|passing|warning), []api.AgentServiceChecksInfo, nil)
// - In all other cases, will return an error
func ServicesByName(serviceName string) (string, []api.AgentServiceChecksInfo, error) {
	return defaultClient().Agent().AgentHealthServiceByName(serviceName)
}

// Services returns the locally registered services
func Services() (map[string]*api.AgentService, error) {
	return defaultClient().Agent().Services()
}

// ServicesWithFilter returns a subset of the locally registered services that match
// the given filter expression
func ServicesWithFilter(filter string) (map[string]*api.AgentService, error) {
	return defaultClient().Agent().ServicesWithFilter(filter)
}

// CatalogServices is used to query for all known services
func CatalogServices() (map[string][]string, error) {
	lists, _, err := defaultClient().Catalog().Services(&api.QueryOptions{})
	return lists, err
}

// CatalogServicesByName is used to query catalog entries for a given service
func CatalogServicesByName(serviceName, tag string) ([]*api.CatalogService, error) {
	services, _, err := defaultClient().Catalog().Service(serviceName, tag, &api.QueryOptions{})
	return services, err
}
