package middleware

// DefaultMiddle 默认的中间件
var DefaultMiddle = Chain(WrapResponseMiddleware, LoggerMiddleware, RecoveryMiddleware)
