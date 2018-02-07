package Action

type Action interface {
	Process() error
}
