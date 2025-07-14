package jsonvaluate

import (
	"encoding/json"
	"fmt"
	"log"
)

// Assuming this would be your jsonvaluate package
// For demo purposes, I'll show the JSON structures

func main() {
	fmt.Println("=== Flexible Condition Logic Demo ===\n")

	// Sample data
	// data := map[string]interface{}{
	// 	"sum_insured":            250000,
	// 	"amount":                 150000,
	// 	"percent_of_sum_insured": 25,
	// 	"age":                    30,
	// 	"status":                 "active",
	// }

	// Example 1: Traditional nested structure (your "before")
	traditionalJSON := `{
		"logic": "AND",
		"children": [
			{
				"key": "sum_insured",
				"operator": ">=",
				"value": 200000
			},
			{
				"logic": "OR", 
				"children": [
					{
						"key": "amount",
						"operator": ">=",
						"value": 100000
					},
					{
						"key": "amount",
						"operator": "<=",
						"value": 1000000
					}
				]
			},
			{
				"key": "percent_of_sum_insured",
				"operator": "%of",
				"value": 20
			}
		]
	}`

	// Example 2: New flexible structure (your "after")
	flexibleJSON := `{
		"conditions": [
			{
				"key": "sum_insured",
				"operator": ">=",
				"value": 200000,
				"next_logic": "AND"
			},
			{
				"group": {
					"conditions": [
						{
							"key": "amount",
							"operator": ">=", 
							"value": 100000,
							"next_logic": "OR"
						},
						{
							"key": "amount",
							"operator": "<=",
							"value": 1000000
						}
					]
				},
				"next_logic": "AND"
			},
			{
				"key": "percent_of_sum_insured",
				"operator": "%of",
				"value": 20
			}
		]
	}`

	// Example 3: More complex flexible logic
	complexFlexibleJSON := `{
		"conditions": [
			{
				"key": "age",
				"operator": ">",
				"value": 25,
				"next_logic": "OR"
			},
			{
				"key": "status", 
				"operator": "==",
				"value": "active",
				"next_logic": "AND"
			},
			{
				"key": "sum_insured",
				"operator": ">=",
				"value": 200000,
				"next_logic": "AND"
			},
			{
				"group": {
					"conditions": [
						{
							"key": "amount",
							"operator": ">=",
							"value": 100000,
							"next_logic": "OR"
						},
						{
							"key": "amount", 
							"operator": "<=",
							"value": 1000000
						}
					]
				}
			}
		]
	}`

	fmt.Println("1. Traditional Nested Structure:")
	fmt.Println("Expression: sum_insured >= 200000 AND (amount >= 100000 OR amount <= 1000000) AND percent_of_sum_insured %of 20")
	fmt.Println("JSON:")
	printFormattedJSON(traditionalJSON)

	fmt.Println("\n2. New Flexible Structure (same logic):")
	fmt.Println("Expression: sum_insured >= 200000 AND (amount >= 100000 OR amount <= 1000000) AND percent_of_sum_insured %of 20")
	fmt.Println("JSON:")
	printFormattedJSON(flexibleJSON)

	fmt.Println("\n3. Complex Mixed Logic:")
	fmt.Println("Expression: age > 25 OR status == 'active' AND sum_insured >= 200000 AND (amount >= 100000 OR amount <= 1000000)")
	fmt.Println("JSON:")
	printFormattedJSON(complexFlexibleJSON)

	fmt.Println("\n=== Key Benefits ===")
	fmt.Println("✅ More natural expression of complex logic")
	fmt.Println("✅ Different logic operators between different condition pairs")
	fmt.Println("✅ Reduced nesting compared to traditional structure")
	fmt.Println("✅ Easier to read and maintain")
	fmt.Println("✅ Backward compatible with existing nested structure")

	fmt.Println("\n=== Usage Examples ===")
	fmt.Println("// Using helper functions:")
	fmt.Println(`group := NewConditionGroup(
    NewConditionWithLogic("sum_insured", ">=", 200000, "AND"),
    NewGroupConditionWithLogic(
        NewConditionGroup(
            NewConditionWithLogic("amount", ">=", 100000, "OR"),
            NewConditionWithLogic("amount", "<=", 1000000, ""),
        ), "AND"),
    NewConditionWithLogic("percent_of_sum_insured", "%of", 20, ""),
)

result := EvaluateConditionGroup(group, data)`)

	fmt.Println("\n// Converting from traditional structure:")
	fmt.Println(`traditionalCondition := Conditions{...}
flexibleGroup := ConvertToConditionGroup(traditionalCondition)
result := EvaluateConditionGroup(flexibleGroup, data)`)

	fmt.Println("\n// Universal evaluation:")
	fmt.Println(`result := EvaluateFlexibleCondition(anyConditionStructure, data)`)
}

func printFormattedJSON(jsonStr string) {
	var obj interface{}
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		log.Fatal(err)
	}

	formatted, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(formatted))
}
