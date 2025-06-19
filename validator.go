package chainflow

// Validator is a specialized processor for validating entities
type Validator[E EnrichableEntity] struct {
	Processor[E]
}

// Validate is the abstract method for validation logic
func (v *Validator[E]) Validate(entity E) error {
	panic("Validate method must be implemented by concrete validators")
}

// Process implements the processor interface by calling Validate
func (v *Validator[E]) Process(entity E) (E, error) {
	err := v.Validate(entity)
	if err != nil {
		return entity, err
	}
	return entity, nil
}
