package httpcontext

const (
	// InitMethod 表示初始化函数名称
	InitMethod = "Init"
	// InvokeMethodName 表示路由里面待调用的函数名称的键值
	// /api/:method 表示 如果匹配到 /api/MethodName -> 会调用对象的 MethodName 方法（如果有的话）
	InvokeMethodName = "method"
)
