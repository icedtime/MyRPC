package impl

import (
	"errors"
	"reflect"
)

type Service struct {
	Name       string
	RefVal     reflect.Value
	RefType    reflect.Type
	MethodType map[string]*MethodType
}

func NewService(server interface{}, useName bool, serverName string) (*Service, error) {
	ser := &Service{
		RefVal:  reflect.ValueOf(server),
		RefType: reflect.TypeOf(server),
	}

	//通过反射获取结构体名字
	sName := reflect.Indirect(ser.RefVal).Type().Name()
	//TODO 判断结构体名字是否大写，用来判断是否暴露给外界
	//TODO 学习一下反射的内容
	if useName {
		if serverName == "" {
			return nil, errors.New("Server name is null")
		}
		sName = serverName
	}
	ser.Name = sName

	//获取该结构体的所有方法
	methods, err := ConstructMethods(ser.RefType)
	if err != nil {
		return nil, err
	}
	ser.MethodType = methods
	return ser, nil
}
