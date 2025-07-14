# Custom Operators

The jsonvaluate library supports custom operators that allow you to extend the built-in functionality with your own validation logic.

## Overview

Custom operators enable you to:
- Add domain-specific validation logic
- Create reusable validation patterns
- Extend the library without modifying core code
- Maintain type safety with custom validation functions

## Quick Start

```go
package main

import (
    "fmt"
    "strings"
    "github.com/Testzyler/jsonvaluate"
)

func main() {
    // Register a case-insensitive equality operator
    jsonvaluate.RegisterCustomOperator("iequal", func(fieldValue, expectedValue interface{}) bool {
        str1 := strings.ToLower(fmt.Sprintf("%v", fieldValue))
        str2 := strings.ToLower(fmt.Sprintf("%v", expectedValue))
        return str1 == str2
    })

    // Use the custom operator
    data := map[string]interface{}{
        "name": "John Doe",
    }

    condition := jsonvaluate.Conditions{
        Key:      "name",
        Operator: "iequal",
        Value:    "JOHN DOE",
    }

    result := jsonvaluate.EvaluateCondition(condition, data)
    fmt.Println(result) // Output: true
}
```

## API Reference

### RegisterCustomOperator

Registers a new custom operator with its validation function.

```go
func RegisterCustomOperator(operator Operator, validator CustomOperatorValidator)
```

**Parameters:**
- `operator`: Unique identifier for the custom operator
- `validator`: Function that implements the validation logic

**Panics:** If validator is nil

### UnregisterCustomOperator

Removes a custom operator from the registry.

```go
func UnregisterCustomOperator(operator Operator)
```

### GetRegisteredCustomOperators

Returns a list of all registered custom operators.

```go
func GetRegisteredCustomOperators() []Operator
```

### CustomOperatorValidator

Function type for custom operator validators.

```go
type CustomOperatorValidator func(fieldValue, expectedValue interface{}) bool
```

**Parameters:**
- `fieldValue`: The actual value from the data being evaluated
- `expectedValue`: The expected value from the condition

**Returns:** `true` if condition is satisfied, `false` otherwise

## Helper Functions

### ToNumber

Converts various types to float64 for numeric operations in custom operators.

```go
func ToNumber(v interface{}) (float64, bool)
```

### ToString

Converts any value to string for string operations in custom operators.

```go
func ToString(v interface{}) string
```

## Examples

### 1. Email Domain Validation

```go
jsonvaluate.RegisterCustomOperator("email_domain", func(fieldValue, expectedValue interface{}) bool {
    email := fmt.Sprintf("%v", fieldValue)
    domain := fmt.Sprintf("%v", expectedValue)
    
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return false
    }
    return parts[1] == domain
})

// Usage
condition := jsonvaluate.Conditions{
    Key:      "email",
    Operator: "email_domain",
    Value:    "example.com",
}
```

### 2. Regex Pattern Matching

```go
jsonvaluate.RegisterCustomOperator("regex", func(fieldValue, expectedValue interface{}) bool {
    str := fmt.Sprintf("%v", fieldValue)
    pattern := fmt.Sprintf("%v", expectedValue)
    
    matched, err := regexp.MatchString(pattern, str)
    return err == nil && matched
})

// Usage
condition := jsonvaluate.Conditions{
    Key:      "phone",
    Operator: "regex",
    Value:    `^\+1-\d{3}-\d{3}-\d{4}$`,
}
```

### 3. Age Group Classification

```go
jsonvaluate.RegisterCustomOperator("age_group", func(fieldValue, expectedValue interface{}) bool {
    ageFloat, ok := jsonvaluate.ToNumber(fieldValue)
    if !ok {
        return false
    }
    age := int(ageFloat)
    
    group := fmt.Sprintf("%v", expectedValue)
    
    switch group {
    case "child":
        return age < 13
    case "teen":
        return age >= 13 && age < 20
    case "adult":
        return age >= 20 && age < 65
    case "senior":
        return age >= 65
    default:
        return false
    }
})

// Usage
condition := jsonvaluate.Conditions{
    Key:      "age",
    Operator: "age_group",
    Value:    "adult",
}
```

### 4. Array Contains Any

```go
jsonvaluate.RegisterCustomOperator("contains_any", func(fieldValue, expectedValue interface{}) bool {
    // Convert fieldValue to slice of strings
    var fieldSlice []string
    switch v := fieldValue.(type) {
    case []string:
        fieldSlice = v
    case []interface{}:
        for _, item := range v {
            fieldSlice = append(fieldSlice, fmt.Sprintf("%v", item))
        }
    default:
        return false
    }
    
    // Convert expectedValue to slice of strings
    var expectedSlice []string
    switch v := expectedValue.(type) {
    case []string:
        expectedSlice = v
    case []interface{}:
        for _, item := range v {
            expectedSlice = append(expectedSlice, fmt.Sprintf("%v", item))
        }
    default:
        return false
    }
    
    // Check if any expected value is in the field slice
    for _, expected := range expectedSlice {
        for _, field := range fieldSlice {
            if field == expected {
                return true
            }
        }
    }
    return false
})

// Usage
condition := jsonvaluate.Conditions{
    Key:      "tags",
    Operator: "contains_any",
    Value:    []string{"python", "golang", "java"},
}
```

## Best Practices

### 1. Error Handling

Custom operators should handle edge cases gracefully:

```go
jsonvaluate.RegisterCustomOperator("safe_divide", func(fieldValue, expectedValue interface{}) bool {
    dividend, ok1 := jsonvaluate.ToNumber(fieldValue)
    divisor, ok2 := jsonvaluate.ToNumber(expectedValue)
    
    // Handle type conversion errors
    if !ok1 || !ok2 {
        return false
    }
    
    // Handle division by zero
    if divisor == 0 {
        return false
    }
    
    return dividend/divisor > 1.0
})
```

### 2. Type Safety

Check types explicitly when needed:

```go
jsonvaluate.RegisterCustomOperator("string_length", func(fieldValue, expectedValue interface{}) bool {
    str, ok := fieldValue.(string)
    if !ok {
        return false
    }
    
    length, ok := jsonvaluate.ToNumber(expectedValue)
    if !ok {
        return false
    }
    
    return float64(len(str)) == length
})
```

### 3. Missing Key Handling

Custom operators receive the raw field value, including `nil` for missing keys:

```go
jsonvaluate.RegisterCustomOperator("key_exists", func(fieldValue, expectedValue interface{}) bool {
    shouldExist, ok := expectedValue.(bool)
    if !ok {
        return false
    }
    
    exists := fieldValue != nil
    return exists == shouldExist
})
```

### 4. Complex Logic

Break down complex validation into smaller functions:

```go
func isValidEmail(email string) bool {
    // Email validation logic
    return strings.Contains(email, "@") && len(email) > 3
}

func isValidDomain(domain string) bool {
    // Domain validation logic
    return strings.Contains(domain, ".") && len(domain) > 3
}

jsonvaluate.RegisterCustomOperator("valid_email_domain", func(fieldValue, expectedValue interface{}) bool {
    email := fmt.Sprintf("%v", fieldValue)
    domain := fmt.Sprintf("%v", expectedValue)
    
    if !isValidEmail(email) {
        return false
    }
    
    if !isValidDomain(domain) {
        return false
    }
    
    parts := strings.Split(email, "@")
    return len(parts) == 2 && parts[1] == domain
})
```

## Thread Safety

The custom operator registry is thread-safe and supports concurrent:
- Registration and unregistration of operators
- Evaluation of conditions using custom operators
- Retrieval of registered operator lists

## Error Recovery

Custom operators that panic are handled gracefully - the evaluation returns `false` instead of crashing the application.

```go
jsonvaluate.RegisterCustomOperator("might_panic", func(fieldValue, expectedValue interface{}) bool {
    // If this panics, evaluation returns false
    panic("Something went wrong!")
})
```

## Performance Considerations

- Custom operators are checked first in the evaluation process
- Registry access is optimized with read-write locks
- Consider caching complex calculations within custom operators
- Avoid heavy I/O operations in custom operators

## Integration with Complex Conditions

Custom operators work seamlessly with logical groups (AND/OR):

```go
complexCondition := jsonvaluate.Conditions{
    Logic: jsonvaluate.LogicAnd,
    Children: []jsonvaluate.Conditions{
        {Key: "email", Operator: "email_domain", Value: "example.com"},
        {Key: "age", Operator: "age_group", Value: "adult"},
        {
            Logic: jsonvaluate.LogicOr,
            Children: []jsonvaluate.Conditions{
                {Key: "username", Operator: "valid_username", Value: nil},
                {Key: "phone", Operator: "regex", Value: `^\+1-\d{3}-\d{3}-\d{4}$`},
            },
        },
    },
}
```
