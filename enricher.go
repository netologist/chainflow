package chainflow

// Enricher is a specialized processor for enriching entities
type Enricher[E EnrichableEntity] struct {
	Processor[E]
}

// Enrich is the abstract method for enrichment logic
func (e *Enricher[E]) Enrich(entity E) (E, error) {
	panic("Enrich method must be implemented by concrete enrichers")
}

// Process implements the processor interface by calling Enrich
func (e *Enricher[E]) Process(entity E) (E, error) {
	return e.Enrich(entity)
}
