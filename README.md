# evaluate-json-condition

A flexible Go library for evaluating complex JSON conditions against data structures. This library allows you to define conditional logic using operators and logical groupings that can be evaluated against JSON-like data.

## Features

- **Rich Operator Support**: Supports equality, comparison, inclusion, string operations, null checks, and more
- **Logical Grouping**: Combine conditions with AND/OR logic
- **Nested Conditions**: Build complex condition trees with unlimited nesting
- **Type Flexible**: Works with various Go types including strings, numbers, booleans, time, slices, and maps
- **JSON Serializable**: All condition structures can be marshaled/unmarshaled from JSON
- **High Performance**: Optimized for fast evaluation with minimal allocations

## Installation

```bash
go get github.com/Testzyler/evaluate-json-condition
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/Testzyler/evaluate-json-condition"
)

func main() {
    // Sample data to evaluate against
    data := map[string]interface{}{
        "age":     25,
        "country": "US",
        "status":  "active",
        "score":   88.5,
    }

    // Simple condition: age > 18
    condition := evaluate.Conditions{
        Key:      "age",
        Operator: evaluate.OperatorGt,
        Value:    18,
    }

    result := evaluate.EvaluateCondition(condition, data)
    fmt.Printf("Age > 18: %v\n", result) // Output: Age > 18: true
}
```

## Supported Operators

### Comparison Operators
- `==` (OperatorEq) - Equal to
- `!=` (OperatorNeq) - Not equal to
- `>` (OperatorGt) - Greater than
- `>=` (OperatorGte) - Greater than or equal to
- `<` (OperatorLt) - Less than
- `<=` (OperatorLte) - Less than or equal to

### Collection Operators
- `in` (OperatorIn) - Value is in collection
- `nin` (OperatorNin) - Value is not in collection

### String Operators
- `contains` (OperatorContains) - String contains substring
- `ncontains` (OperatorNcontains) - String does not contain substring
- `like` (OperatorLike) - SQL-like pattern matching (case sensitive)
- `ilike` (OperatorIlike) - SQL-like pattern matching (case insensitive)
- `nlike` (OperatorNlike) - NOT SQL-like pattern matching
- `startswith` (OperatorStartsWith) - String starts with prefix
- `endswith` (OperatorEndsWith) - String ends with suffix

### State Operators
- `isnull` (OperatorIsnull) - Value is null or doesn't exist
- `isnotnull` (OperatorIsnotnull) - Value is not null and exists
- `isempty` (OperatorIsEmpty) - Value is empty (empty string, array, etc.)
- `isnotempty` (OperatorIsNotEmpty) - Value is not empty
- `istrue` (OperatorIsTrue) - Value is true (boolean or truthy)
- `isfalse` (OperatorIsFalse) - Value is false (boolean or falsy)

### Range Operators
- `between` (OperatorBetween) - Value is between two bounds (inclusive)
- `notbetween` (OperatorNotBetween) - Value is not between two bounds

## Usage Examples

### Single Conditions

```go
// String equality
condition := evaluate.Conditions{
    Key:      "country",
    Operator: evaluate.OperatorEq,
    Value:    "US",
}

// Numeric comparison
condition := evaluate.Conditions{
    Key:      "age",
    Operator: evaluate.OperatorGte,
    Value:    21,
}

// Collection membership
condition := evaluate.Conditions{
    Key:      "country",
    Operator: evaluate.OperatorIn,
    Value:    []interface{}{"US", "CA", "UK"},
}

// String contains
condition := evaluate.Conditions{
    Key:      "description",
    Operator: evaluate.OperatorContains,
    Value:    "important",
}

// Range check
condition := evaluate.Conditions{
    Key:      "score",
    Operator: evaluate.OperatorBetween,
    Value:    []interface{}{80, 100},
}
```

### Logical Groups

```go
// AND group - all conditions must be true
andCondition := evaluate.Conditions{
    Logic: evaluate.LogicAnd,
    Children: []evaluate.Conditions{
        {Key: "age", Operator: evaluate.OperatorGte, Value: 18},
        {Key: "country", Operator: evaluate.OperatorEq, Value: "US"},
        {Key: "status", Operator: evaluate.OperatorEq, Value: "active"},
    },
}

// OR group - at least one condition must be true
orCondition := evaluate.Conditions{
    Logic: evaluate.LogicOr,
    Children: []evaluate.Conditions{
        {Key: "role", Operator: evaluate.OperatorEq, Value: "admin"},
        {Key: "permissions", Operator: evaluate.OperatorContains, Value: "write"},
    },
}
```

### Nested Conditions

```go
// Complex nested condition
nestedCondition := evaluate.Conditions{
    Logic: evaluate.LogicAnd,
    Children: []evaluate.Conditions{
        {Key: "age", Operator: evaluate.OperatorGte, Value: 18},
        {
            Logic: evaluate.LogicOr,
            Children: []evaluate.Conditions{
                {Key: "country", Operator: evaluate.OperatorEq, Value: "US"},
                {Key: "country", Operator: evaluate.OperatorEq, Value: "CA"},
            },
        },
        {
            Logic: evaluate.LogicAnd,
            Children: []evaluate.Conditions{
                {Key: "status", Operator: evaluate.OperatorEq, Value: "active"},
                {Key: "verified", Operator: evaluate.OperatorIsTrue, Value: nil},
            },
        },
    },
}
```

### Working with JSON

```go
import (
    "encoding/json"
    "github.com/Testzyler/evaluate-json-condition"
)

// Conditions can be serialized to/from JSON
conditionJSON := `{
    "logic": "AND",
    "children": [
        {"key": "age", "operator": ">=", "value": 18},
        {"key": "country", "operator": "==", "value": "US"}
    ]
}`

var condition evaluate.Conditions
err := json.Unmarshal([]byte(conditionJSON), &condition)
if err != nil {
    panic(err)
}

data := map[string]interface{}{
    "age":     25,
    "country": "US",
}

result := evaluate.EvaluateCondition(condition, data)
fmt.Printf("Result: %v\n", result) // Output: Result: true
```

### Time-based Conditions

```go
import "time"

data := map[string]interface{}{
    "created_at": time.Now(),
    "date_str":   "2024-01-15T10:30:00Z",
}

// Time comparison
timeCondition := evaluate.Conditions{
    Key:      "created_at",
    Operator: evaluate.OperatorGt,
    Value:    time.Now().Add(-24 * time.Hour), // 24 hours ago
}

// Time range
timeRangeCondition := evaluate.Conditions{
    Key:      "date_str",
    Operator: evaluate.OperatorBetween,
    Value:    []interface{}{"2024-01-01T00:00:00Z", "2024-12-31T23:59:59Z"},
}
```

## Type Handling

The library intelligently handles type conversions:

- **Numbers**: Supports all Go numeric types (int, float, etc.) with automatic conversion
- **Strings**: Automatic string conversion for comparisons
- **Booleans**: Smart boolean evaluation (true/false, "true"/"false", 1/0, etc.)
- **Time**: Supports time.Time and string time formats (RFC3339, etc.)
- **Collections**: Works with slices, arrays, and maps
- **Nil/Empty**: Proper handling of nil values and empty collections

## Performance

The library is optimized for performance:

```go
// Benchmark single condition evaluation
func BenchmarkSingleCondition(b *testing.B) {
    data := map[string]interface{}{"age": 25}
    condition := evaluate.Conditions{
        Key:      "age",
        Operator: evaluate.OperatorGt,
        Value:    18,
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        evaluate.EvaluateCondition(condition, data)
    }
}
```

## API Reference

### Types

#### `Conditions`
The main structure representing a condition tree.

```go
type Conditions struct {
    Logic    Logic        `json:"logic,omitempty"`    // "AND" or "OR" for groups
    Children []Conditions `json:"children,omitempty"` // Child conditions
    Key      string       `json:"key,omitempty"`      // Field key for single condition
    Operator Operator     `json:"operator,omitempty"` // Comparison operator
    Value    interface{}  `json:"value,omitempty"`    // Expected value
}
```

#### `Operator`
String type representing comparison operators.

#### `Logic`
String type representing logical operators ("AND", "OR").

### Functions

#### `EvaluateCondition(cond Conditions, data map[string]interface{}) bool`
Evaluates a condition tree against the provided data and returns the result.

## Publishing and Usage Instructions

### For Users wanting to use this library:

#### Installation
```bash
go get github.com/Testzyler/evaluate-json-condition
```

#### Quick Import and Usage
```go
import "github.com/Testzyler/evaluate-json-condition"

// Simple usage
data := map[string]interface{}{"age": 25}
condition := evaluate.NewSimpleCondition("age", evaluate.OperatorGt, 18)
result := evaluate.EvaluateCondition(condition, data) // true
```

### For Library Authors (Publishing):

#### Steps to publish this library:

1. **Ensure the repository is public on GitHub**
2. **Tag a release version:**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. **Verify the module is available:**
   ```bash
   GOPROXY=proxy.golang.org go list -m github.com/Testzyler/evaluate-json-condition@latest
   ```

#### Semantic Versioning
- **v1.0.0** - Initial stable release
- **v1.0.1** - Bug fixes
- **v1.1.0** - New features (backward compatible)  
- **v2.0.0** - Breaking changes

#### Update Documentation
After tagging, update the installation instructions and examples in README.md with the specific version if needed.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.