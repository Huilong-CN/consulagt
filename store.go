package consulagt

import "sync"

type ServiceStore struct {
	*sync.Map
}

var serviceCache = &ServiceStore{&sync.Map{}}

type SvcRegistedStore struct {
	*sync.Map
}

func (srs *SvcRegistedStore) Get(key string) (*ServiceMeta, bool) {
	val, ok := srs.Load(key)
	if !ok {
		return nil, false
	}
	meta, ok := val.(*ServiceMeta)
	if !ok {
		return nil, false
	}
	return meta, true
}

func (srs *SvcRegistedStore) RangeDo(f func(key interface{}, val *ServiceMeta)) {
	srs.Range(func(key interface{}, value interface{}) bool {
		meta, ok := value.(*ServiceMeta)
		if ok {
			f(key, meta)
		}
		return true
	})
}

var registeredSvc = &SvcRegistedStore{&sync.Map{}}
