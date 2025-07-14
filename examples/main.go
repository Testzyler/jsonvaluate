package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Testzyler/jsonvaluate"
)

func main() {
	// Example data to evaluate against
	data := map[string]interface{}{
		"name":     "John Doe",
		"email":    "john.doe@example.com",
		"age":      25,
		"username": "john_doe123",
		"phone":    "+1-555-123-4567",
		"score":    85.5,
		"tags":     []string{"developer", "golang", "backend"},
		"active":   true,
	}

	fmt.Println("=== Custom Operator Examples ===")
	fmt.Println()

	// Example 1: Case-insensitive string equality
	fmt.Println("1. Case-insensitive string equality:")
	jsonvaluate.RegisterCustomOperator("iequal", func(fieldValue, expectedValue interface{}) bool {
		str1 := strings.ToLower(fmt.Sprintf("%v", fieldValue))
		str2 := strings.ToLower(fmt.Sprintf("%v", expectedValue))
		return str1 == str2
	})

	condition1 := jsonvaluate.Conditions{
		Key:      "name",
		Operator: "iequal",
		Value:    "JOHN DOE",
	}
	result1 := jsonvaluate.EvaluateCondition(condition1, data)
	fmt.Printf("   name iequal 'JOHN DOE': %v\n", result1)

	// Example 2: Email domain validation
	fmt.Println("\n2. Email domain validation:")
	jsonvaluate.RegisterCustomOperator("email_domain", func(fieldValue, expectedValue interface{}) bool {
		email := fmt.Sprintf("%v", fieldValue)
		domain := fmt.Sprintf("%v", expectedValue)

		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			return false
		}
		return parts[1] == domain
	})

	condition2 := jsonvaluate.Conditions{
		Key:      "email",
		Operator: "email_domain",
		Value:    "example.com",
	}
	result2 := jsonvaluate.EvaluateCondition(condition2, data)
	fmt.Printf("   email email_domain 'example.com': %v\n", result2)

	// Example 3: Regex pattern matching
	fmt.Println("\n3. Regex pattern matching:")
	jsonvaluate.RegisterCustomOperator("regex", func(fieldValue, expectedValue interface{}) bool {
		str := fmt.Sprintf("%v", fieldValue)
		pattern := fmt.Sprintf("%v", expectedValue)

		matched, err := regexp.MatchString(pattern, str)
		return err == nil && matched
	})

	condition3 := jsonvaluate.Conditions{
		Key:      "phone",
		Operator: "regex",
		Value:    `^\+1-\d{3}-\d{3}-\d{4}$`,
	}
	result3 := jsonvaluate.EvaluateCondition(condition3, data)
	fmt.Printf("   phone regex '^\\+1-\\d{3}-\\d{3}-\\d{4}$': %v\n", result3)

	// Example 4: Username validation (alphanumeric + underscore, 3-20 chars)
	fmt.Println("\n4. Username validation:")
	jsonvaluate.RegisterCustomOperator("valid_username", func(fieldValue, expectedValue interface{}) bool {
		username := fmt.Sprintf("%v", fieldValue)

		// Length check
		if len(username) < 3 || len(username) > 20 {
			return false
		}

		// Character check (alphanumeric + underscore only)
		matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
		return matched
	})

	condition4 := jsonvaluate.Conditions{
		Key:      "username",
		Operator: "valid_username",
		Value:    nil, // Not used for this validation
	}
	result4 := jsonvaluate.EvaluateCondition(condition4, data)
	fmt.Printf("   username valid_username: %v\n", result4)

	// Example 5: Age group classification
	fmt.Println("\n5. Age group classification:")
	jsonvaluate.RegisterCustomOperator("age_group", func(fieldValue, expectedValue interface{}) bool {
		ageFloat, ok := toNumber(fieldValue)
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

	condition5 := jsonvaluate.Conditions{
		Key:      "age",
		Operator: "age_group",
		Value:    "adult",
	}
	result5 := jsonvaluate.EvaluateCondition(condition5, data)
	fmt.Printf("   age age_group 'adult': %v\n", result5)

	// Example 6: Array contains any of the specified values
	fmt.Println("\n6. Array contains any of:")
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

	condition6 := jsonvaluate.Conditions{
		Key:      "tags",
		Operator: "contains_any",
		Value:    []string{"python", "golang", "java"},
	}
	result6 := jsonvaluate.EvaluateCondition(condition6, data)
	fmt.Printf("   tags contains_any ['python', 'golang', 'java']: %v\n", result6)

	// Example 7: Score grade calculation
	fmt.Println("\n7. Score grade calculation:")
	jsonvaluate.RegisterCustomOperator("grade", func(fieldValue, expectedValue interface{}) bool {
		scoreFloat, ok := toNumber(fieldValue)
		if !ok {
			return false
		}

		expectedGrade := fmt.Sprintf("%v", expectedValue)

		var actualGrade string
		switch {
		case scoreFloat >= 90:
			actualGrade = "A"
		case scoreFloat >= 80:
			actualGrade = "B"
		case scoreFloat >= 70:
			actualGrade = "C"
		case scoreFloat >= 60:
			actualGrade = "D"
		default:
			actualGrade = "F"
		}

		return actualGrade == expectedGrade
	})

	condition7 := jsonvaluate.Conditions{
		Key:      "score",
		Operator: "grade",
		Value:    "B",
	}
	result7 := jsonvaluate.EvaluateCondition(condition7, data)
	fmt.Printf("   score grade 'B': %v\n", result7)

	// Example 8: Complex condition using multiple custom operators
	fmt.Println("\n8. Complex condition with custom operators:")
	complexCondition := jsonvaluate.Conditions{
		Logic: jsonvaluate.LogicAnd,
		Children: []jsonvaluate.Conditions{
			{Key: "username", Operator: "valid_username", Value: nil},
			{Key: "age", Operator: "age_group", Value: "adult"},
			{Key: "email", Operator: "email_domain", Value: "example.com"},
			{
				Logic: jsonvaluate.LogicOr,
				Children: []jsonvaluate.Conditions{
					{Key: "score", Operator: "grade", Value: "A"},
					{Key: "score", Operator: "grade", Value: "B"},
				},
			},
		},
	}

	complexResult := jsonvaluate.EvaluateCondition(complexCondition, data)
	fmt.Printf("   Complex condition result: %v\n", complexResult)

	// Display registered custom operators
	fmt.Println("\n=== Registered Custom Operators ===")
	customOps := jsonvaluate.GetRegisteredCustomOperators()
	for i, op := range customOps {
		fmt.Printf("%d. %s\n", i+1, op)
	}

	// Example of error handling with invalid data
	fmt.Println("\n=== Error Handling Examples ===")

	invalidData := map[string]interface{}{
		"age":   "not a number",
		"email": "invalid-email",
	}

	condition8 := jsonvaluate.Conditions{
		Key:      "age",
		Operator: "age_group",
		Value:    "adult",
	}
	result8 := jsonvaluate.EvaluateCondition(condition8, invalidData)
	fmt.Printf("Invalid age data: %v\n", result8)

	condition9 := jsonvaluate.Conditions{
		Key:      "email",
		Operator: "email_domain",
		Value:    "example.com",
	}
	result9 := jsonvaluate.EvaluateCondition(condition9, invalidData)
	fmt.Printf("Invalid email data: %v\n", result9)

	fmt.Println("\n=== Cleanup ===")
	// Clean up registered operators (optional)
	for _, op := range jsonvaluate.GetRegisteredCustomOperators() {
		jsonvaluate.UnregisterCustomOperator(op)
		fmt.Printf("Unregistered operator: %s\n", op)
	}

	// demonstrateCustomOperators shows how to use custom operators
	func() {
		fmt.Println("\n=== Custom Operators Demo ===")

		// Register custom operators
		fmt.Println("Registering custom operators...")

		// Case-insensitive string comparison
		jsonvaluate.RegisterCustomOperator("case_insensitive_eq", func(fieldValue, expectedValue interface{}) bool {
			str1 := strings.ToLower(fmt.Sprintf("%v", fieldValue))
			str2 := strings.ToLower(fmt.Sprintf("%v", expectedValue))
			return str1 == str2
		})

		// Email domain validation
		jsonvaluate.RegisterCustomOperator("email_domain", func(fieldValue, expectedValue interface{}) bool {
			email := fmt.Sprintf("%v", fieldValue)
			domain := fmt.Sprintf("%v", expectedValue)

			parts := strings.Split(email, "@")
			if len(parts) != 2 {
				return false
			}
			return parts[1] == domain
		})

		// String length validation
		jsonvaluate.RegisterCustomOperator("min_length", func(fieldValue, expectedValue interface{}) bool {
			str := fmt.Sprintf("%v", fieldValue)
			minLen, ok := expectedValue.(int)
			if !ok {
				return false
			}
			return len(str) >= minLen
		})

		// Test data
		data := map[string]interface{}{
			"name":     "John Doe",
			"email":    "john@example.com",
			"username": "johndoe123",
		}

		// Test case-insensitive comparison
		cond1 := jsonvaluate.Conditions{
			Key:      "name",
			Operator: "case_insensitive_eq",
			Value:    "JOHN DOE",
		}
		fmt.Printf("Case insensitive match (John Doe == JOHN DOE): %v\n", jsonvaluate.EvaluateCondition(cond1, data))

		// Test email domain validation
		cond2 := jsonvaluate.Conditions{
			Key:      "email",
			Operator: "email_domain",
			Value:    "example.com",
		}
		fmt.Printf("Email domain match (john@example.com domain is example.com): %v\n", jsonvaluate.EvaluateCondition(cond2, data))

		// Test string length validation
		cond3 := jsonvaluate.Conditions{
			Key:      "username",
			Operator: "min_length",
			Value:    5,
		}
		fmt.Printf("Username min length (johndoe123 >= 5 chars): %v\n", jsonvaluate.EvaluateCondition(cond3, data))

		// Complex condition with custom operators
		complexCond := jsonvaluate.Conditions{
			Logic: jsonvaluate.LogicAnd,
			Children: []jsonvaluate.Conditions{
				{Key: "name", Operator: "case_insensitive_eq", Value: "john doe"},
				{Key: "email", Operator: "email_domain", Value: "example.com"},
				{Key: "username", Operator: "min_length", Value: 5},
			},
		}
		fmt.Printf("Complex condition (all conditions must pass): %v\n", jsonvaluate.EvaluateCondition(complexCond, data))

		// List registered custom operators
		fmt.Printf("Registered custom operators: %v\n", jsonvaluate.GetRegisteredCustomOperators())

		// Clean up
		fmt.Println("Cleaning up custom operators...")
		jsonvaluate.UnregisterCustomOperator("case_insensitive_eq")
		jsonvaluate.UnregisterCustomOperator("email_domain")
		jsonvaluate.UnregisterCustomOperator("min_length")

		fmt.Printf("After cleanup: %v\n", jsonvaluate.GetRegisteredCustomOperators())
	}()
}

// Helper function to expose the toNumber function for examples
// Note: This assumes the toNumber function is exported or we need to implement it
func toNumber(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}
