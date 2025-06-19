package chainflow

type Handler[E EnrichableEntity] interface {
	Handle(entity E) (E, error)
	SetNext(nextHandler Handler[E])
}
