package iface

type Iserver interface {
	Start()
	Stop()
	Server()
	//AddService(methodName string)
	RegisterService(service interface{}, useName bool, name string) error
	RegisterUseName(service interface{}, name string) error
	RegisterWithoutName(Service interface{}) error
}
