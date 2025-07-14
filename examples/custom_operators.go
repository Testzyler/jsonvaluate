package main

import (
	"fmt"
	"strings"

	"github.com/Testzyler/jsonvaluate"
)

func main() {
	// Example 1: Case-insensitive string comparison
	jsonvaluate.RegisterCustomOperator("case_insensitive_eq", func(fieldValue, expectedValue interface{}) bool {
		str1 := strings.ToLower(fmt.Sprintf("%v", fieldValue))
		str2 := strings.ToLower(fmt.Sprintf("%v", expectedValue))
		return str1 == str2
	})

	// Example 2: Email domain validation
	jsonvaluate.RegisterCustomOperator("email_domain", func(fieldValue, expectedValue interface{}) bool {
		email := fmt.Sprintf("%v", fieldValue)
		domain := fmt.Sprintf("%v", expectedValue)

		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			return false
		}
		return parts[1] == domain
	})

	// Example 3: String length validation
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
	fmt.Printf("Case insensitive match: %v\n", jsonvaluate.EvaluateCondition(cond1, data))

	// Test email domain validation
	cond2 := jsonvaluate.Conditions{
		Key:      "email",
		Operator: "email_domain",
		Value:    "example.com",
	}
	fmt.Printf("Email domain match: %v\n", jsonvaluate.EvaluateCondition(cond2, data))

	// Test string length validation
	cond3 := jsonvaluate.Conditions{
		Key:      "username",
		Operator: "min_length",
		Value:    5,
	}
	fmt.Printf("Username min length: %v\n", jsonvaluate.EvaluateCondition(cond3, data))

	// Complex condition with custom operator
	complexCond := jsonvaluate.Conditions{
		Logic: jsonvaluate.LogicAnd,
		Children: []jsonvaluate.Conditions{
			{Key: "name", Operator: "case_insensitive_eq", Value: "john doe"},
			{Key: "email", Operator: "email_domain", Value: "example.com"},
			{Key: "username", Operator: "min_length", Value: 5},
		},
	}
	fmt.Printf("Complex condition: %v\n", jsonvaluate.EvaluateCondition(complexCond, data))

	// List registered custom operators
	fmt.Printf("Registered custom operators: %v\n", jsonvaluate.GetRegisteredCustomOperators())

	// Clean up
	jsonvaluate.UnregisterCustomOperator("case_insensitive_eq")
	jsonvaluate.UnregisterCustomOperator("email_domain")
	jsonvaluate.UnregisterCustomOperator("min_length")

	fmt.Printf("After cleanup: %v\n", jsonvaluate.GetRegisteredCustomOperators())
}
