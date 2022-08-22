package impl

import "reflect"

type MethodType struct {
	Method       reflect.Method
	RequestType  reflect.Type
	ResponseType reflect.Type
}

func ConstructMethods(typ reflect.Type) (map[string]*MethodType, error) {
	//使用一个结构体的反射类型就可以基本获取所有的信息
	//可以把refTyp理解为结构体实体
	res := map[string]*MethodType{}

	//结构体方法数:typ.NumMethod()
	//遍历结构体全部方法
	for i := 0; i < typ.NumMethod(); i++ {
		//获取方法实体
		method := typ.Method(i)
		//获取方法名
		mName := method.Name
		//获取方法类型,在反射中类型才是最重要的，可以获取很多东西
		mType := method.Type

		//TODO 对于每一个方法，要做以下的事情
		//方法名要检查是否开头大写
		//检查参数四个，因为提供的服务方法，其方法签名是固定的
		//检查每一个参数是否对应

		//requestType responseType其实就是方法中的，req,res参数
		//获取方法中的这两个参数
		//这些方法其实本质上就两个参数，为什么要获取，我个人想法是，在通过方法名调用的时候，可以直接赋值调用

		//参数列表不是从0开始的，从1
		requestType := mType.In(1)
		ResponseType := mType.In(2)

		res[mName] = &MethodType{
			Method:       method,
			RequestType:  requestType,
			ResponseType: ResponseType,
		}
	}

	return res, nil
}
