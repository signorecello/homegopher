package service


type ServiceCall interface {
	GetServiceOpts() ServiceOpts
	SetServiceOpts(ServiceOpts) ServiceOpts
	GetService() string
	SetService(string)
}

type ServiceOpts interface {
	SetEntityID(string)
}


