package chainflow

import "errors"

// Processor is the abstract base class for all processors
type Processor[E EnrichableEntity] struct {
	nextHandler Handler[E]
}

// Process is the abstract method that must be implemented by concrete processors
func (p *Processor[E]) Process(entity E) (E, error) {
	panic("Process method must be implemented by concrete processors")
}

// Handle implements the template method pattern
func (p *Processor[E]) Handle(entity E) (E, error) {
	// Process the entity
	processedEntity, err := p.Process(entity)
	if err != nil {
		// Check if it's a chain break exception
		var breakErr *BreakProcessChainError
		if errors.As(err, &breakErr) {
			return entity, nil // Return current entity and stop chain
		}
		return entity, err // Propagate other errors
	}

	// Call next handler in chain
	return p.callNext(processedEntity)
}

// SetNext sets the next handler in the chain
func (p *Processor[E]) SetNext(nextHandler Handler[E]) {
	p.nextHandler = nextHandler
}

// callNext calls the next handler if it exists
func (p *Processor[E]) callNext(entity E) (E, error) {
	if p.nextHandler != nil {
		return p.nextHandler.Handle(entity)
	}
	return entity, nil
}
