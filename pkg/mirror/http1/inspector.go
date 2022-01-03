package http1

type Inspector interface {
	Inspect(params *Parameters) error
}
