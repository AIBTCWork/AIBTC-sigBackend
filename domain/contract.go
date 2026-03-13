package domain

// 数据库模型 一对一映射
type contract struct {
	ID     uint
	addr   string
	amount string
	block  uint
}
