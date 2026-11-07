# ChainFlow

[![Go Version](https://img.shields.io/badge/Go-1.20+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A robust, extensible framework for building sequential processing pipelines that can enrich and validate entities through a chain of responsibility pattern, handling complex business logic while maintaining clean separation of concerns.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Core Components](#core-components)
- [Design Patterns](#design-patterns)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage Examples](#usage-examples)
- [API Reference](#api-reference)
- [Best Practices](#best-practices)
- [Contributing](#contributing)

## Overview

The ChainFlow provides a flexible way to build processing pipelines where:

- Each processor has a single responsibility
- Processors can be chained in any order
- Chain execution can be terminated gracefully
- Type safety is maintained throughout the pipeline
- Error handling is consistent and predictable

### Key Benefits

- **🔗 Flexible Chaining**: Combine processors in any sequence
- **🛡️ Type Safety**: Generic constraints ensure compile-time safety
- **🚫 Graceful Termination**: Break chains without error propagation
- **🔄 Reusable Components**: Same processors in different chains
- **⚡ Fail Fast**: Early validation prevents unnecessary processing
- **🧪 Testable**: Each processor can be tested independently

## Architecture

```mermaid
graph TB
    subgraph "Core Framework"
        EE[EnrichableEntity Interface]
        H[Handler Interface]
        P[Processor Struct]
        E[Enricher Struct]
        V[Validator Struct]
        PCB[ProcessorChainBuilder]
        PC[ProcessorChain]
    end
    
    subgraph "Concrete Implementations"
        TE[TestEntity]
        TEN[TestEnricher]
        TV[TestValidator]
        CE[ConditionalEnricher]
    end
    
    subgraph "Error Handling"
        BPCE[BreakProcessChainError]
    end
    
    EE --> TE
    H --> P
    P --> E
    P --> V
    E --> TEN
    E --> CE
    V --> TV
    PCB --> PC
    PC --> H
    
    style EE fill:#e1f5fe
    style H fill:#e8f5e8
    style P fill:#fff3e0
    style BPCE fill:#ffebee
```

## Core Components

### Component Hierarchy

```mermaid
classDiagram
    class EnrichableEntity {
        <<interface>>
        +GetType() string
    }
    
    class Handler~E~ {
        <<interface>>
        +Handle(entity E) (E, error)
        +SetNext(nextHandler Handler~E~)
    }
    
    class Processor~E~ {
        -nextHandler Handler~E~
        +Process(entity E) (E, error)*
        +Handle(entity E) (E, error)
        +SetNext(nextHandler Handler~E~)
        #callNext(entity E) (E, error)
    }
    
    class Enricher~E~ {
        +Enrich(entity E) (E, error)*
        +Process(entity E) (E, error)
    }
    
    class Validator~E~ {
        +Validate(entity E) error*
        +Process(entity E) (E, error)
    }
    
    class ProcessorChainBuilder~E~ {
        -processors []Handler~E~
        +AddProcessor(processor Handler~E~) ProcessorChainBuilder~E~
        +Build() (Handler~E~, error)
    }
    
    class ProcessorChain~E~ {
        -firstProcessor Handler~E~
        +ExecuteChain(entity E) (E, error)
    }
    
    class BreakProcessChainError {
        +Message string
        +Error() string
    }
    
    Handler~E~ <|-- Processor~E~
    Processor~E~ <|-- Enricher~E~
    Processor~E~ <|-- Validator~E~
    ProcessorChainBuilder~E~ --> ProcessorChain~E~
    ProcessorChain~E~ --> Handler~E~
    
    note for Processor~E~ "Template Method Pattern\nImplements chain linking"
    note for Enricher~E~ "Modifies entity data"
    note for Validator~E~ "Validates without modification"
```

### Interface Definitions

```mermaid
graph LR
    subgraph EE["EnrichableEntity"]
        EE_METHOD["GetType() string"]
    end

    subgraph HD["Handler#lt;E#gt;"]
        H_HANDLE["Handle(entity E) E, error"]
        H_SETNEXT["SetNext(nextHandler Handler#lt;E#gt;)"]
    end

    subgraph CM["Core Methods"]
        PROCESS["Process abstract method"]
        ENRICH["Enrich abstract method"]
        VALIDATE["Validate abstract method"]
    end

    EE_METHOD --> H_HANDLE
    H_HANDLE --> PROCESS
    H_SETNEXT --> PROCESS
    PROCESS --> ENRICH
    PROCESS --> VALIDATE

    style EE_METHOD fill:#e1f5fe
    style H_HANDLE fill:#e8f5e8
    style H_SETNEXT fill:#e8f5e8
    style PROCESS fill:#fff3e0
    style ENRICH fill:#f3e5f5
    style VALIDATE fill:#fce4ec
```

## Design Patterns

### Template Method Pattern

```mermaid
sequenceDiagram
    participant Client
    participant Processor
    participant ConcreteProcessor
    
    Client->>+Processor: Handle(entity)
    
    rect rgb(255, 248, 220)
        Note over Processor: Template Method
        Processor->>+ConcreteProcessor: Process(entity)
        ConcreteProcessor-->>-Processor: processedEntity
        
        alt No Exception
            Processor->>Processor: callNext(processedEntity)
        else BreakProcessChainError
            Processor-->>Client: entity (unchanged)
        end
    end
    
    Processor-->>-Client: finalEntity
    
    Note over Processor,ConcreteProcessor: Template defines algorithm structure<br/>Concrete classes implement specific steps
```

### Chain of Responsibility

```mermaid
graph LR
    Client[Client] --> PC[ProcessorChain]
    PC --> P1[Processor 1]
    P1 --> P2[Processor 2]
    P2 --> P3[Processor 3]
    P3 --> Result[Final Result]
    
    P1 -.->|nextHandler| P2
    P2 -.->|nextHandler| P3
    P3 -.->|nextHandler| NULL
    
    style Client fill:#e3f2fd
    style PC fill:#f1f8e9
    style P1 fill:#fff3e0
    style P2 fill:#fff3e0
    style P3 fill:#fff3e0
    style Result fill:#e8f5e8
```

## Installation

```bash
# Clone the repository
git clone https://github.com/netologist/chainflow.git

# Navigate to the project directory
cd processor-chain-go

# Initialize Go module
go mod init processor-chain

# Run the example
go run main.go
```

### Requirements

- Go 1.18 or higher (for generics support)
- No external dependencies

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
)

func main() {
    // Create processors
    enricher1 := NewTestEnricher(" - Step1")
    validator := NewTestValidator("validation", true)
    enricher2 := NewTestEnricher(" - Step2")

    // Build the chain
    builder := NewProcessorChainBuilder[*TestEntity]()
    firstProcessor, err := builder.
        AddProcessor(enricher1).
        AddProcessor(validator).
        AddProcessor(enricher2).
        Build()

    if err != nil {
        panic(err)
    }

    // Create and execute chain
    chain := NewProcessorChain(firstProcessor)
    entity := &TestEntity{Data: "Initial Data"}
    
    result, err := chain.ExecuteChain(entity)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Result: %s\n", result)
    // Output: Result: TestEntity{Data: Initial Data - Step1 - Step2}
}
```

## Usage Examples

### 1. Normal Chain Execution

```mermaid
sequenceDiagram
    participant Client
    participant Chain as ProcessorChain
    participant V as Validator
    participant P as Processor
    participant E as Enricher
    
    Client->>+Chain: ExecuteChain(entity)
    Chain->>+V: Handle(entity)
    
    V->>V: Validate(entity)
    Note right of V: Validation passes
    
    V->>+P: Handle(entity)
    P->>P: Process(entity)
    
    P->>+E: Handle(modifiedEntity)
    E->>E: Enrich(entity)
    
    E-->>-P: enrichedEntity
    P-->>-V: enrichedEntity
    V-->>-Chain: enrichedEntity
    Chain-->>-Client: enrichedEntity
    
    Note over Client,E: All processors execute successfully
```

### 2. Chain Breaking Scenario

```mermaid
sequenceDiagram
    participant Client
    participant Chain as ProcessorChain
    participant V as Validator
    participant P as Processor
    participant E as Enricher
    
    Client->>+Chain: ExecuteChain(entity)
    Chain->>+V: Handle(entity)
    
    V->>V: Validate(entity)
    V->>V: throw BreakProcessChainError
    
    rect rgb(255, 235, 238)
        Note over V: Template Method catches exception
        V->>V: return entity (unchanged)
    end
    
    V-->>-Chain: entity (unchanged)
    Chain-->>-Client: entity (unchanged)
    
    Note over P,E: These processors are never called
    
    style V fill:#ffcdd2
```

### 3. Validation Failure

```go
// Example: Validation that fails
validator := NewTestValidator("strict validation", false)

chain := NewProcessorChainBuilder[*TestEntity]().
    AddProcessor(NewTestEnricher(" - Before")).
    AddProcessor(validator).
    AddProcessor(NewTestEnricher(" - After")).
    Build()

result, err := chain.ExecuteChain(entity)
if err != nil {
    fmt.Printf("Validation failed: %v\n", err)
    // Output: Validation failed: invalid entity validator: strict validation
}
```

### 4. Conditional Processing

```go
// Custom processor that conditionally breaks the chain
type ConditionalProcessor struct {
    Processor[*TestEntity]
}

func (cp *ConditionalProcessor) Process(entity *TestEntity) (*TestEntity, error) {
    if entity.Data == "STOP" {
        return entity, &BreakProcessChainError{Message: "Processing stopped"}
    }
    return &TestEntity{Data: entity.Data + " - Processed"}, nil
}
```

## API Reference

### Core Types

#### EnrichableEntity
```go
type EnrichableEntity interface {
    GetType() string
}
```

#### Handler[E EnrichableEntity]
```go
type Handler[E EnrichableEntity] interface {
    Handle(entity E) (E, error)
    SetNext(nextHandler Handler[E])
}
```

#### Processor[E EnrichableEntity]
```go
type Processor[E EnrichableEntity] struct {
    nextHandler Handler[E]
}

func (p *Processor[E]) Handle(entity E) (E, error)
func (p *Processor[E]) SetNext(nextHandler Handler[E])
func (p *Processor[E]) Process(entity E) (E, error) // Abstract method
```

### Builder Pattern

#### ProcessorChainBuilder[E EnrichableEntity]
```go
func NewProcessorChainBuilder[E EnrichableEntity]() *ProcessorChainBuilder[E]
func (b *ProcessorChainBuilder[E]) AddProcessor(processor Handler[E]) *ProcessorChainBuilder[E]
func (b *ProcessorChainBuilder[E]) Build() (Handler[E], error)
```

### Error Handling

#### BreakProcessChainError
```go
type BreakProcessChainError struct {
    Message string
}

func (e *BreakProcessChainError) Error() string
```

## Execution Flow

### Complete Chain Execution Flow

```mermaid
flowchart TD
    Start([Start]) --> CreateChain[Create ProcessorChain]
    CreateChain --> Execute[ExecuteChain entity]
    Execute --> FirstProcessor{First Processor}
    
    FirstProcessor --> HandleMethod[Call Handle method]
    HandleMethod --> TryCatch{Try-Catch Block}
    
    TryCatch --> ProcessMethod[Call Process method]
    ProcessMethod --> ProcessorType{Processor Type?}
    
    ProcessorType -->|Validator| ValidateMethod[Call Validate]
    ProcessorType -->|Enricher| EnrichMethod[Call Enrich]
    ProcessorType -->|Custom| CustomMethod[Call Custom Process]
    
    ValidateMethod --> ValidationResult{Validation Result?}
    ValidationResult -->|Pass| CallNext[Call callNext]
    ValidationResult -->|Fail| ThrowError[Throw Error]
    
    EnrichMethod --> CallNext
    CustomMethod --> CallNext
    
    CallNext --> HasNext{Has Next Processor?}
    HasNext -->|Yes| NextProcessor[Next Processor Handle]
    HasNext -->|No| ReturnResult[Return Final Entity]
    
    NextProcessor --> HandleMethod
    
    ThrowError --> ErrorType{Error Type?}
    ErrorType -->|BreakChainError| ReturnCurrent[Return Current Entity]
    ErrorType -->|Other Error| PropagateError[Propagate Error]
    
    ReturnCurrent --> End([End - Graceful])
    ReturnResult --> End
    PropagateError --> End([End - Error])
    
    style Start fill:#e8f5e8
    style End fill:#ffcdd2
    style TryCatch fill:#fff3e0
    style ProcessorType fill:#e1f5fe
    style ValidationResult fill:#f3e5f5
```

## Best Practices

### 1. Design Guidelines

```mermaid
mindmap
  root((Best Practices))
    Design
      Single Responsibility
        One processor, one purpose
        Clear, focused functionality
      Immutable Entities
        Create new instances
        Avoid modifying input
      Error Handling
        Use appropriate error types
        Fail fast when possible
    
    Performance
      Early Validation
        Place validators first
        Catch issues quickly
      Chain Reuse
        Build once, use many times
        Cache expensive chains
      Memory Management
        Use pointers for large entities
        Let GC handle cleanup
    
    Testing
      Unit Testing
        Test each processor independently
        Mock dependencies
      Integration Testing
        Test complete chains
        Verify error scenarios
      Performance Testing
        Benchmark chain execution
        Profile memory usage
```

### 2. Processor Ordering Strategy

```mermaid
graph TD
    A[Input Entity] --> B[Quick Validators]
    B --> C[Expensive Validators]
    C --> D[Light Enrichers]
    D --> E[Heavy Enrichers]
    E --> F[Final Validators]
    F --> G[Output Entity]
    
    B -.->|Fail Fast| ERROR1[Early Error]
    C -.->|Business Rules| ERROR2[Validation Error]
    F -.->|Final Check| ERROR3[Final Error]
    
    style A fill:#e8f5e8
    style G fill:#e8f5e8
    style ERROR1 fill:#ffcdd2
    style ERROR2 fill:#ffcdd2
    style ERROR3 fill:#ffcdd2
```

### 3. Error Handling Strategy

```mermaid
flowchart LR
    subgraph "Error Types"
        VE[Validation Errors]
        PE[Processing Errors]
        BE[Business Errors]
        SE[System Errors]
    end
    
    subgraph "Handling Strategy"
        VE --> |Immediate Return| Client1[Return Error to Client]
        PE --> |Log and Continue| Log1[Log Error + Continue]
        BE --> |Break Chain| Break1[Graceful Chain Break]
        SE --> |Panic Recovery| Recover1[Recover and Log]
    end
    
    style VE fill:#ffcdd2
    style PE fill:#fff3e0
    style BE fill:#e1f5fe
    style SE fill:#f3e5f5
```

## Performance Considerations

### Memory Usage Patterns

```mermaid
graph TB
    subgraph "Memory Optimization"
        A[Use Pointers for Large Entities]
        B[Avoid Deep Copying]
        C[Reuse Chain Instances]
        D[Pool Entity Objects]
    end
    
    subgraph "Performance Metrics"
        E[Chain Build Time: O1]
        F[Chain Execution: On]
        G[Memory Usage: Linear]
        H[GC Pressure: Low]
    end
    
    A --> E
    B --> F
    C --> G
    D --> H
    
    style A fill:#e8f5e8
    style B fill:#e8f5e8
    style C fill:#e8f5e8
    style D fill:#e8f5e8
```

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Fork and clone the repository
git clone https://github.com/netologist/chainflow.git

# Create a feature branch
git checkout -b feature/your-feature-name

# Make your changes and add tests
go test ./...

# Submit a pull request
```

### Code Style

- Follow Go conventions and use `gofmt`
- Write comprehensive tests for new features
- Document public APIs with Go doc comments
- Use meaningful commit messages

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for details about changes in each version.

---

**Questions?** Open an issue or start a discussion in our [GitHub Discussions](https://github.com/netologist/chainflow/discussions).
