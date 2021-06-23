package consulagt

//ServiceMeta desc a services
type ServiceMeta struct {
	ServiceName     string
	HTTPHealthyAddr string
	Metas           map[string]string
	ServicesAddr    string
	RegistStatus    bool
}
