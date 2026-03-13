package entity

// 领域模型 业务逻辑处理 用于 Handler → Service、Service → Handler 之间传递 业务语义明确 的参数/结果。
type contractEntity struct {
	addr   string
	amount string
}
