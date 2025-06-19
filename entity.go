package chainflow

// EnrichableEntity is a marker interface for entities that can be processed
type EnrichableEntity interface {
	GetType() string
}

// BreakProcessChainError is used to gracefully terminate the processing chain
type BreakProcessChainError struct {
	Message string
}

func (e *BreakProcessChainError) Error() string {
	return e.Message
}
