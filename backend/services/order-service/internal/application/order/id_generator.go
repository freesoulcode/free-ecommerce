package order

type IDGenerator interface {
	NextID() (int64, error)
}
