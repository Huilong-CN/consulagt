package consulagt

import "sync"

type ServiceStore struct {
	*sync.Map
}

var serviceCache = &ServiceStore{&sync.Map{}}
