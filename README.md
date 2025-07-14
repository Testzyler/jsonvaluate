# evaluate-json-condition

A flexible Go library for evaluating complex JSON conditions against data structures. This library allows you to define conditional logic using operators and logical groupings that can be evaluated against JSON-like data.

## Features

- **Rich Operator Support**: Supports equality, comparison, inclusion, string operations, null checks, and more
- **Custom Operators**: Extend functionality with your own validation logic using custom operators
- **Flexible Logic**: Choose different logical operators (AND/OR) between different condition pairs
- **Logical Grouping**: Combine conditions with AND/OR logic in traditional nested structure
- **Nested Conditions**: Build complex condition trees with unlimited nesting
- **Type Flexible**: Works with various Go types including strings, numbers, booleans, time, slices, and maps
- **JSON Serializable**: All condition structures can be marshaled/unmarshaled from JSON
- **Thread Safe**: Safe for concurrent use across goroutines
- **High Performance**: Optimized for fast evaluation with minimal allocations

## Installation

```bash
go get github.com/Testzyler/jsonvaluate
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/Testzyler/jsonvaluate"
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
    condition := jsonvaluate.Conditions{
        Key:      "age",
        Operator: jsonvaluate.OperatorGt,
        Value:    18,
    }

    result := jsonvaluate.EvaluateCondition(condition, data)
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

## Custom Operators

The library supports custom operators that allow you to extend the built-in functionality with your own validation logic.

### Quick Custom Operator Example

```go
// Register a case-insensitive equality operator
jsonvaluate.RegisterCustomOperator("iequal", func(fieldValue, expectedValue interface{}) bool {
    str1 := strings.ToLower(fmt.Sprintf("%v", fieldValue))
    str2 := strings.ToLower(fmt.Sprintf("%v", expectedValue))
    return str1 == str2
})

// Use the custom operator
condition := jsonvaluate.Conditions{
    Key:      "name",
    Operator: "iequal",
    Value:    "JOHN DOE",
}

result := jsonvaluate.EvaluateCondition(condition, data) // true for "John Doe"
```

### Custom Operator Functions

- **RegisterCustomOperator(operator, validator)** - Register a new custom operator
- **UnregisterCustomOperator(operator)** - Remove a custom operator
- **GetRegisteredCustomOperators()** - List all registered custom operators

For detailed examples and best practices, see [CUSTOM_OPERATORS.md](CUSTOM_OPERATORS.md).

## Flexible Logic Conditions

The library now supports flexible logical expressions where you can specify different logical operators (AND/OR) between different condition pairs, similar to SQL expressions.

### Traditional vs Flexible Structure

**Traditional nested structure:**
```json
{
    "logic": "AND",
    "children": [
        {"key": "age", "operator": ">=", "value": 18},
        {
            "logic": "OR",
            "children": [
                {"key": "status", "operator": "==", "value": "active"},
                {"key": "role", "operator": "==", "value": "admin"}
            ]
        }
    ]
}
```

**New flexible structure:**
```json
{
    "conditions": [
        {
            "key": "age",
            "operator": ">=", 
            "value": 18,
            "next_logic": "AND"
        },
        {
            "key": "status",
            "operator": "==",
            "value": "active", 
            "next_logic": "OR"
        },
        {
            "key": "role",
            "operator": "==",
            "value": "admin"
        }
    ]
}
```

### Flexible Logic Example

```go
// Expression: sum_insured >= 200000 AND (amount >= 100000 OR amount <= 1000000) AND percent >= 20
flexibleCondition := jsonvaluate.ConditionGroup{
    Conditions: []jsonvaluate.ConditionWithLogic{
        {
            Key:       "sum_insured",
            Operator:  jsonvaluate.OperatorGte,
            Value:     200000,
            NextLogic: jsonvaluate.LogicAnd,
        },
        {
            Group: &jsonvaluate.ConditionGroup{
                Conditions: []jsonvaluate.ConditionWithLogic{
                    {
                        Key:       "amount",
                        Operator:  jsonvaluate.OperatorGte,
                        Value:     100000,
                        NextLogic: jsonvaluate.LogicOr,
                    },
                    {
                        Key:      "amount",
                        Operator: jsonvaluate.OperatorLte,
                        Value:    1000000,
                    },
                },
            },
            NextLogic: jsonvaluate.LogicAnd,
        },
        {
            Key:      "percent",
            Operator: jsonvaluate.OperatorGte,
            Value:    20,
        },
    },
}

result := jsonvaluate.EvaluateConditionGroup(flexibleCondition, data)
```

### Helper Functions for Flexible Logic

```go
// Create conditions using helper functions
group := jsonvaluate.NewConditionGroup(
    jsonvaluate.NewConditionWithLogic("age", ">=", 18, "AND"),
    jsonvaluate.NewConditionWithLogic("status", "==", "active", "OR"),
    jsonvaluate.NewConditionWithLogic("role", "==", "admin", ""),
)

// Convert traditional structure to flexible
traditionalCondition := jsonvaluate.Conditions{...}
flexibleGroup := jsonvaluate.ConvertToConditionGroup(traditionalCondition)

// Universal evaluation (works with both structures)
result := jsonvaluate.EvaluateFlexibleCondition(anyConditionStructure, data)
```

## Usage Examples

### Single Conditions

```go
// String equality
condition := jsonvaluate.Conditions{
    Key:      "country",
    Operator: jsonvaluate.OperatorEq,
    Value:    "US",
}

// Numeric comparison
condition := jsonvaluate.Conditions{
    Key:      "age",
    Operator: jsonvaluate.OperatorGte,
    Value:    21,
}

// Collection membership
condition := jsonvaluate.Conditions{
    Key:      "country",
    Operator: jsonvaluate.OperatorIn,
    Value:    []interface{}{"US", "CA", "UK"},
}

// String contains
condition := jsonvaluate.Conditions{
    Key:      "description",
    Operator: jsonvaluate.OperatorContains,
    Value:    "important",
}

// Range check
condition := jsonvaluate.Conditions{
    Key:      "score",
    Operator: jsonvaluate.OperatorBetween,
    Value:    []interface{}{80, 100},
}
```

### Logical Groups

```go
// AND group - all conditions must be true
andCondition := jsonvaluate.Conditions{
    Logic: jsonvaluate.LogicAnd,
    Children: []jsonvaluate.Conditions{
        {Key: "age", Operator: jsonvaluate.OperatorGte, Value: 18},
        {Key: "country", Operator: jsonvaluate.OperatorEq, Value: "US"},
        {Key: "status", Operator: jsonvaluate.OperatorEq, Value: "active"},
    },
}

// OR group - at least one condition must be true
orCondition := jsonvaluate.Conditions{
    Logic: jsonvaluate.LogicOr,
    Children: []jsonvaluate.Conditions{
        {Key: "role", Operator: jsonvaluate.OperatorEq, Value: "admin"},
        {Key: "permissions", Operator: jsonvaluate.OperatorContains, Value: "write"},
    },
}
```

### Nested Conditions

```go
// Complex nested condition
nestedCondition := jsonvaluate.Conditions{
    Logic: jsonvaluate.LogicAnd,
    Children: []jsonvaluate.Conditions{
        {Key: "age", Operator: jsonvaluate.OperatorGte, Value: 18},
        {
            Logic: jsonvaluate.LogicOr,
            Children: []jsonvaluate.Conditions{
                {Key: "country", Operator: jsonvaluate.OperatorEq, Value: "US"},
                {Key: "country", Operator: jsonvaluate.OperatorEq, Value: "CA"},
            },
        },
        {
            Logic: jsonvaluate.LogicAnd,
            Children: []jsonvaluate.Conditions{
                {Key: "status", Operator: jsonvaluate.OperatorEq, Value: "active"},
                {Key: "verified", Operator: jsonvaluate.OperatorIsTrue, Value: nil},
            },
        },
    },
}
```

### Working with JSON

```go
import (
    "encoding/json"
    "github.com/Testzyler/jsonvaluate"
)

// Conditions can be serialized to/from JSON
conditionJSON := `{
    "logic": "AND",
    "children": [
        {"key": "age", "operator": ">=", "value": 18},
        {"key": "country", "operator": "==", "value": "US"}
    ]
}`

var condition jsonvaluate.Conditions
err := json.Unmarshal([]byte(conditionJSON), &condition)
if err != nil {
    panic(err)
}

data := map[string]interface{}{
    "age":     25,
    "country": "US",
}

result := jsonvaluate.EvaluateCondition(condition, data)
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
timeCondition := jsonvaluate.Conditions{
    Key:      "created_at",
    Operator: jsonvaluate.OperatorGt,
    Value:    time.Now().Add(-24 * time.Hour), // 24 hours ago
}

// Time range
timeRangeCondition := jsonvaluate.Conditions{
    Key:      "date_str",
    Operator: jsonvaluate.OperatorBetween,
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
    condition := jsonvaluate.Conditions{
        Key:      "age",
        Operator: jsonvaluate.OperatorGt,
        Value:    18,
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        jsonvaluate.EvaluateCondition(condition, data)
    }
}
```

## API Reference

### Core Types

#### `Conditions`
The main structure representing a condition tree (traditional nested structure).

```go
type Conditions struct {
    Logic    Logic        `json:"logic,omitempty"`    // "AND" or "OR" for groups
    Children []Conditions `json:"children,omitempty"` // Child conditions
    Key      string       `json:"key,omitempty"`      // Field key for single condition
    Operator Operator     `json:"operator,omitempty"` // Comparison operator
    Value    interface{}  `json:"value,omitempty"`    // Expected value
}
```

#### `ConditionGroup`
New flexible structure for expressing mixed logical operations.

```go
type ConditionGroup struct {
    Conditions []ConditionWithLogic `json:"conditions"`
}

type ConditionWithLogic struct {
    Key       string           `json:"key,omitempty"`       // Field key for condition
    Operator  Operator         `json:"operator,omitempty"`  // Comparison operator
    Value     interface{}      `json:"value,omitempty"`     // Expected value
    Group     *ConditionGroup  `json:"group,omitempty"`     // Nested group (alternative)
    NextLogic Logic            `json:"next_logic,omitempty"` // Logic to connect to next condition
}
```

#### `Operator`
String type representing comparison operators.

#### `Logic`
String type representing logical operators ("AND", "OR").

#### `CustomOperatorValidator`
Function type for custom operator validation logic.

```go
type CustomOperatorValidator func(fieldValue, expectedValue interface{}) bool
```

### Core Functions

#### `EvaluateCondition(cond Conditions, data map[string]interface{}) bool`
Evaluates a traditional condition tree against the provided data.

#### `EvaluateConditionGroup(group ConditionGroup, data map[string]interface{}) bool`
Evaluates a flexible condition group against the provided data.

#### `EvaluateFlexibleCondition(conditions interface{}, data map[string]interface{}) bool`
Universal evaluation function that works with both Conditions and ConditionGroup structures.

### Helper Functions

#### `NewSimpleCondition(key, operator, value) Conditions`
Creates a simple condition with key, operator, and value.

#### `NewAndGroup(children ...Conditions) Conditions`
Creates an AND group condition from child conditions.

#### `NewOrGroup(children ...Conditions) Conditions`
Creates a OR group condition from child conditions.

#### `NewConditionGroup(conditions ...ConditionWithLogic) ConditionGroup`
Creates a new flexible condition group.

#### `NewConditionWithLogic(key, operator, value, nextLogic) ConditionWithLogic`
Creates a single condition with specified logic for the next condition.

#### `ConvertToConditionGroup(conditions Conditions) ConditionGroup`
Converts traditional nested structure to flexible structure.

### Custom Operator Functions

#### `RegisterCustomOperator(operator Operator, validator CustomOperatorValidator)`
Registers a new custom operator with validation logic.

#### `UnregisterCustomOperator(operator Operator)`
Removes a custom operator from the registry.

#### `GetRegisteredCustomOperators() []Operator`
Returns a list of all registered custom operators.

## Publishing and Usage Instructions

### For Users wanting to use this library:

#### Installation
```bash
go get github.com/Testzyler/jsonvaluate
```

#### Quick Import and Usage
```go
import "github.com/Testzyler/jsonvaluate"

// Simple usage
data := map[string]interface{}{"age": 25}
condition := jsonvaluate.NewSimpleCondition("age", jsonvaluate.OperatorGt, 18)
result := jsonvaluate.EvaluateCondition(condition, data) // true

// Custom operator usage
jsonvaluate.RegisterCustomOperator("custom_op", func(field, expected interface{}) bool {
    return fmt.Sprintf("%v", field) == fmt.Sprintf("%v", expected)
})

// Flexible logic usage
group := jsonvaluate.NewConditionGroup(
    jsonvaluate.NewConditionWithLogic("age", ">=", 18, "AND"),
    jsonvaluate.NewConditionWithLogic("status", "==", "active", ""),
)
result = jsonvaluate.EvaluateConditionGroup(group, data)
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
   GOPROXY=proxy.golang.org go list -m github.com/Testzyler/jsonvaluate@latest
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

## Advanced Features

### When to Use Which Structure

#### Traditional Nested Structure (`Conditions`)
- **Best for**: Simple, uniform logic (all AND or all OR)
- **Use when**: Your conditions follow a clear hierarchical pattern
- **Example**: User permissions - must meet ALL criteria OR be in ANY admin role

#### Flexible Logic Structure (`ConditionGroup`)
- **Best for**: Complex expressions with mixed logic operators
- **Use when**: You need different logic between different condition pairs
- **Example**: Insurance eligibility - multiple criteria with different relationships

#### Custom Operators
- **Best for**: Domain-specific validation logic
- **Use when**: Built-in operators don't cover your specific business rules
- **Example**: Email domain validation, regex matching, complex calculations

### Real-World Use Cases

#### E-commerce Product Filtering
```go
// "Show products that are (in stock AND price < $100) OR (featured AND rating >= 4)"
filter := jsonvaluate.ConditionGroup{
    Conditions: []jsonvaluate.ConditionWithLogic{
        {
            Group: &jsonvaluate.ConditionGroup{
                Conditions: []jsonvaluate.ConditionWithLogic{
                    {Key: "in_stock", Operator: "==", Value: true, NextLogic: "AND"},
                    {Key: "price", Operator: "<", Value: 100},
                },
            },
            NextLogic: "OR",
        },
        {
            Group: &jsonvaluate.ConditionGroup{
                Conditions: []jsonvaluate.ConditionWithLogic{
                    {Key: "featured", Operator: "==", Value: true, NextLogic: "AND"},
                    {Key: "rating", Operator: ">=", Value: 4},
                },
            },
        },
    },
}
```

#### User Access Control
```go
// Register custom operator for role checking
jsonvaluate.RegisterCustomOperator("has_role", func(userRoles, requiredRole interface{}) bool {
    roles, ok := userRoles.([]string)
    if !ok { return false }
    required := fmt.Sprintf("%v", requiredRole)
    for _, role := range roles {
        if role == required { return true }
    }
    return false
})

// "Allow if (admin OR moderator) AND account_active AND (trial_expired == false OR subscription_active)"
accessControl := jsonvaluate.ConditionGroup{
    Conditions: []jsonvaluate.ConditionWithLogic{
        {
            Group: &jsonvaluate.ConditionGroup{
                Conditions: []jsonvaluate.ConditionWithLogic{
                    {Key: "roles", Operator: "has_role", Value: "admin", NextLogic: "OR"},
                    {Key: "roles", Operator: "has_role", Value: "moderator"},
                },
            },
            NextLogic: "AND",
        },
        {Key: "account_active", Operator: "==", Value: true, NextLogic: "AND"},
        {
            Group: &jsonvaluate.ConditionGroup{
                Conditions: []jsonvaluate.ConditionWithLogic{
                    {Key: "trial_expired", Operator: "==", Value: false, NextLogic: "OR"},
                    {Key: "subscription_active", Operator: "==", Value: true},
                },
            },
        },
    },
}
```

#### Financial Risk Assessment
```go
// Register custom operators for financial calculations
jsonvaluate.RegisterCustomOperator("debt_to_income_ratio", func(debt, income interface{}) bool {
    d, _ := jsonvaluate.ToNumber(debt)
    i, _ := jsonvaluate.ToNumber(income)
    if i == 0 { return false }
    return (d / i) < 0.4 // Less than 40% debt-to-income ratio
})

// Complex financial approval logic
approval := jsonvaluate.ConditionGroup{
    Conditions: []jsonvaluate.ConditionWithLogic{
        {Key: "credit_score", Operator: ">=", Value: 650, NextLogic: "AND"},
        {Key: "annual_income", Operator: ">=", Value: 50000, NextLogic: "AND"},
        {
            Group: &jsonvaluate.ConditionGroup{
                Conditions: []jsonvaluate.ConditionWithLogic{
                    {Key: "debt", Operator: "debt_to_income_ratio", Value: "annual_income", NextLogic: "OR"},
                    {Key: "collateral_value", Operator: ">=", Value: 100000},
                },
            },
            NextLogic: "AND",
        },
        {Key: "employment_verified", Operator: "==", Value: true},
    },
}
```

## Version Compatibility

### Current Version Features

#### v1.x.x Features:
- ✅ Traditional nested condition structures
- ✅ All built-in operators (comparison, string, state, range)
- ✅ Custom operator support with thread-safe registry
- ✅ Flexible logic conditions with mixed AND/OR operations
- ✅ Type-safe validation and panic recovery
- ✅ JSON serialization/deserialization
- ✅ Helper functions for easy condition creation
- ✅ Backward compatibility with existing code

#### Migration Guide

**From v0.x to v1.x:** No breaking changes - all existing code continues to work.

**New features are additive:**
```go
// Existing code works unchanged
oldCondition := jsonvaluate.Conditions{...}
result := jsonvaluate.EvaluateCondition(oldCondition, data)

// New features available
customOp := jsonvaluate.RegisterCustomOperator("custom", validator)
flexibleCondition := jsonvaluate.ConditionGroup{...}
result2 := jsonvaluate.EvaluateConditionGroup(flexibleCondition, data)
```