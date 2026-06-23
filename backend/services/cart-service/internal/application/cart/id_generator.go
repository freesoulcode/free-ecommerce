package cart

type IDGenerator interface {
	NextID() (int64, error)
}
