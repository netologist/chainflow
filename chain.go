package chainflow

import "errors"

// ProcessorChainBuilder helps build processor chains
type ProcessorChainBuilder[E EnrichableEntity] struct {
	processors []Handler[E]
}

// NewProcessorChainBuilder creates a new builder
func NewProcessorChainBuilder[E EnrichableEntity]() *ProcessorChainBuilder[E] {
	return &ProcessorChainBuilder[E]{
		processors: make([]Handler[E], 0),
	}
}

// AddProcessor adds a processor to the chain
func (b *ProcessorChainBuilder[E]) AddProcessor(processor Handler[E]) *ProcessorChainBuilder[E] {
	b.processors = append(b.processors, processor)
	return b
}

// Build constructs the processor chain
func (b *ProcessorChainBuilder[E]) Build() (Handler[E], error) {
	if len(b.processors) == 0 {
		return nil, errors.New("at least one processor must be added")
	}

	// Link processors together
	for i := 0; i < len(b.processors)-1; i++ {
		b.processors[i].SetNext(b.processors[i+1])
	}

	return b.processors[0], nil
}

// ProcessorChain manages the execution of a processor chain
type ProcessorChain[E EnrichableEntity] struct {
	firstProcessor Handler[E]
}

// NewProcessorChain creates a new processor chain
func NewProcessorChain[E EnrichableEntity](firstProcessor Handler[E]) *ProcessorChain[E] {
	return &ProcessorChain[E]{
		firstProcessor: firstProcessor,
	}
}

// ExecuteChain executes the entire processor chain
func (pc *ProcessorChain[E]) ExecuteChain(entity E) (E, error) {
	return pc.firstProcessor.Handle(entity)
}
