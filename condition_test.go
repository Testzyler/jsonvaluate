package jsonvaluate

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestEvalSingleCondition_AllOperators(t *testing.T) {
	tm := time.Date(2024, 7, 1, 12, 0, 0, 0, time.UTC)
	data := map[string]interface{}{
		"age":       25,
		"country":   "TH",
		"score":     88.5,
		"tags":      []string{"a", "b", "c"},
		"desc":      "hello world",
		"empty":     "",
		"nil":       nil,
		"boolTrue":  true,
		"boolFalse": false,
		"date":      tm,
		"dateStr":   "2024-07-01T12:00:00Z",
	}

	tests := []struct {
		name   string
		key    string
		op     Operator
		value  interface{}
		expect bool
	}{
		{"eq true", "country", OperatorEq, "TH", true},
		{"eq false", "country", OperatorEq, "SG", false},
		{"neq true", "country", OperatorNeq, "SG", true},
		{"neq false", "country", OperatorNeq, "TH", false},
		{"gt true", "age", OperatorGt, 18, true},
		{"gt false", "age", OperatorGt, 30, false},
		{"gte true", "age", OperatorGte, 25, true},
		{"gte false", "age", OperatorGte, 30, false},
		{"lt true", "score", OperatorLt, 90, true},
		{"lt false", "score", OperatorLt, 80, false},
		{"lte true", "score", OperatorLte, 88.5, true},
		{"lte false", "score", OperatorLte, 80, false},
		{"in true", "country", OperatorIn, []interface{}{"TH", "SG"}, true},
		{"in false", "country", OperatorIn, []interface{}{"SG", "MY"}, false},
		{"nin true", "country", OperatorNin, []interface{}{"SG", "MY"}, true},
		{"nin false", "country", OperatorNin, []interface{}{"TH", "SG"}, false},
		{"contains true", "desc", OperatorContains, "hello", true},
		{"contains false", "desc", OperatorContains, "bye", false},
		{"ncontains true", "desc", OperatorNcontains, "bye", true},
		{"ncontains false", "desc", OperatorNcontains, "hello", false},
		{"isnull true", "nil", OperatorIsnull, nil, true},
		{"isnull false", "country", OperatorIsnull, nil, false},
		{"isnotnull true", "country", OperatorIsnotnull, nil, true},
		{"isnotnull false", "nil", OperatorIsnotnull, nil, false},
		{"isempty true", "empty", OperatorIsEmpty, nil, true},
		{"isempty false", "desc", OperatorIsEmpty, nil, false},
		{"isnotempty true", "desc", OperatorIsNotEmpty, nil, true},
		{"isnotempty false", "empty", OperatorIsNotEmpty, nil, false},
		{"istrue true", "boolTrue", OperatorIsTrue, nil, true},
		{"istrue false", "boolFalse", OperatorIsTrue, nil, false},
		{"isfalse true", "boolFalse", OperatorIsFalse, nil, true},
		{"isfalse false", "boolTrue", OperatorIsFalse, nil, false},
		{"like true", "desc", OperatorLike, "%hello%", true},
		{"like false", "desc", OperatorLike, "%bye%", false},
		{"ilike true", "desc", OperatorIlike, "%HELLO%", true},
		{"ilike false", "desc", OperatorIlike, "%BYE%", false},
		{"nlike true", "desc", OperatorNlike, "%bye%", true},
		{"nlike false", "desc", OperatorNlike, "%hello%", false},
		{"startswith true", "desc", OperatorStartsWith, "hello", true},
		{"startswith false", "desc", OperatorStartsWith, "world", false},
		{"endswith true", "desc", OperatorEndsWith, "world", true},
		{"endswith false", "desc", OperatorEndsWith, "hello", false},
		{"between true", "age", OperatorBetween, []interface{}{20, 30}, true},
		{"between false", "age", OperatorBetween, []interface{}{30, 40}, false},
		{"notbetween true", "age", OperatorNotBetween, []interface{}{30, 40}, true},
		{"notbetween false", "age", OperatorNotBetween, []interface{}{20, 30}, false},
		{"between time true", "date", OperatorBetween, []interface{}{tm.Add(-time.Hour), tm.Add(time.Hour)}, true},
		{"between time false", "date", OperatorBetween, []interface{}{tm.Add(time.Hour), tm.Add(2 * time.Hour)}, false},
		{"between time string true", "dateStr", OperatorBetween, []interface{}{"2024-06-01T00:00:00Z", "2024-08-01T00:00:00Z"}, true},
		{"between time string false", "dateStr", OperatorBetween, []interface{}{"2024-08-01T00:00:00Z", "2024-09-01T00:00:00Z"}, false},
		{"missing key", "notfound", OperatorEq, 1, false},
		{"type flexible comparison", "age", OperatorEq, "25", true},
		{"nil vs nil eq", "nil", OperatorEq, nil, true},
		{"nil vs value eq", "nil", OperatorEq, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := evalSingleCondition(tt.key, tt.op, tt.value, data)
			if result != tt.expect {
				t.Errorf("evalSingleCondition(%s, %s, %v) = %v, want %v", tt.key, tt.op, tt.value, result, tt.expect)
			}
		})
	}
}

func TestEvaluateCondition_GroupsAndNest(t *testing.T) {
	data := map[string]interface{}{
		"age":     25,
		"country": "TH",
		"status":  "active",
	}

	// Single condition
	single := Conditions{
		Key:      "age",
		Operator: OperatorGt,
		Value:    18,
	}
	if !EvaluateCondition(single, data) {
		t.Error("Single condition should be true")
	}

	// AND group
	andGroup := Conditions{
		Logic: LogicAnd,
		Children: []Conditions{
			{Key: "age", Operator: OperatorGt, Value: 18},
			{Key: "country", Operator: OperatorEq, Value: "TH"},
		},
	}
	if !EvaluateCondition(andGroup, data) {
		t.Error("AND group should be true")
	}

	// OR group
	orGroup := Conditions{
		Logic: LogicOr,
		Children: []Conditions{
			{Key: "country", Operator: OperatorEq, Value: "SG"},
			{Key: "status", Operator: OperatorEq, Value: "active"},
		},
	}
	if !EvaluateCondition(orGroup, data) {
		t.Error("OR group should be true")
	}

	// Nested group
	nested := Conditions{
		Logic: LogicAnd,
		Children: []Conditions{
			{Key: "age", Operator: OperatorGt, Value: 18},
			{
				Logic: LogicOr,
				Children: []Conditions{
					{Key: "country", Operator: OperatorEq, Value: "SG"},
					{Key: "status", Operator: OperatorEq, Value: "active"},
				},
			},
		},
	}
	if !EvaluateCondition(nested, data) {
		t.Error("Nested group should be true")
	}

	// AND group with one false
	andFalse := Conditions{
		Logic: LogicAnd,
		Children: []Conditions{
			{Key: "age", Operator: OperatorGt, Value: 18},
			{Key: "country", Operator: OperatorEq, Value: "SG"},
		},
	}
	if EvaluateCondition(andFalse, data) {
		t.Error("AND group with one false should be false")
	}

	// OR group with all false
	orFalse := Conditions{
		Logic: LogicOr,
		Children: []Conditions{
			{Key: "country", Operator: OperatorEq, Value: "SG"},
			{Key: "status", Operator: OperatorEq, Value: "inactive"},
		},
	}
	if EvaluateCondition(orFalse, data) {
		t.Error("OR group with all false should be false")
	}
}

func BenchmarkEvalSingleCondition(b *testing.B) {
	tm := time.Date(2024, 7, 1, 12, 0, 0, 0, time.UTC)
	data := map[string]interface{}{
		"age":      25,
		"country":  "TH",
		"score":    88.5,
		"desc":     "hello world",
		"boolTrue": true,
		"date":     tm,
	}
	conds := []struct {
		key   string
		op    Operator
		value interface{}
	}{
		{"age", OperatorGt, 18},
		{"country", OperatorEq, "TH"},
		{"score", OperatorLte, 100},
		{"desc", OperatorContains, "hello"},
		{"boolTrue", OperatorIsTrue, nil},
		{"date", OperatorBetween, []interface{}{tm.Add(-time.Hour), tm.Add(time.Hour)}},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, c := range conds {
			_ = evalSingleCondition(c.key, c.op, c.value, data)
		}
	}
}

func BenchmarkEvaluateCondition(b *testing.B) {
	data := map[string]interface{}{
		"age":     25,
		"country": "TH",
		"status":  "active",
		"score":   88.5,
	}
	cond := Conditions{
		Logic: LogicAnd,
		Children: []Conditions{
			{Key: "age", Operator: OperatorGt, Value: 18},
			{Key: "country", Operator: OperatorEq, Value: "TH"},
			{
				Logic: LogicOr,
				Children: []Conditions{
					{Key: "status", Operator: OperatorEq, Value: "active"},
					{Key: "score", Operator: OperatorGt, Value: 80},
				},
			},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EvaluateCondition(cond, data)
	}
}

func TestCustomOperators(t *testing.T) {
	// Clean up any existing custom operators
	for _, op := range GetRegisteredCustomOperators() {
		UnregisterCustomOperator(op)
	}

	data := map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   25,
		"score": 85.5,
	}

	// Test 1: Register a case-insensitive equality operator
	RegisterCustomOperator("case_insensitive_eq", func(fieldValue, expectedValue interface{}) bool {
		str1 := strings.ToLower(fmt.Sprintf("%v", fieldValue))
		str2 := strings.ToLower(fmt.Sprintf("%v", expectedValue))
		return str1 == str2
	})

	// Test case-insensitive equality
	cond1 := Conditions{
		Key:      "name",
		Operator: "case_insensitive_eq",
		Value:    "JOHN DOE",
	}
	if !EvaluateCondition(cond1, data) {
		t.Error("Case insensitive equality should be true")
	}

	cond2 := Conditions{
		Key:      "name",
		Operator: "case_insensitive_eq",
		Value:    "Jane Doe",
	}
	if EvaluateCondition(cond2, data) {
		t.Error("Case insensitive equality should be false")
	}

	// Test 2: Register an email domain validator
	RegisterCustomOperator("email_domain", func(fieldValue, expectedValue interface{}) bool {
		email := fmt.Sprintf("%v", fieldValue)
		domain := fmt.Sprintf("%v", expectedValue)

		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			return false
		}
		return parts[1] == domain
	})

	// Test email domain validation
	cond3 := Conditions{
		Key:      "email",
		Operator: "email_domain",
		Value:    "example.com",
	}
	if !EvaluateCondition(cond3, data) {
		t.Error("Email domain validation should be true")
	}

	cond4 := Conditions{
		Key:      "email",
		Operator: "email_domain",
		Value:    "gmail.com",
	}
	if EvaluateCondition(cond4, data) {
		t.Error("Email domain validation should be false")
	}

	// Test 3: Register a numeric range operator
	RegisterCustomOperator("in_range", func(fieldValue, expectedValue interface{}) bool {
		value, ok := toNumber(fieldValue)
		if !ok {
			return false
		}

		rv := reflect.ValueOf(expectedValue)
		if rv.Kind() != reflect.Slice || rv.Len() != 2 {
			return false
		}

		min, okMin := toNumber(rv.Index(0).Interface())
		max, okMax := toNumber(rv.Index(1).Interface())
		if !okMin || !okMax {
			return false
		}

		return value >= min && value <= max
	})

	// Test numeric range validation
	cond5 := Conditions{
		Key:      "age",
		Operator: "in_range",
		Value:    []interface{}{20, 30},
	}
	if !EvaluateCondition(cond5, data) {
		t.Error("Numeric range validation should be true")
	}

	cond6 := Conditions{
		Key:      "age",
		Operator: "in_range",
		Value:    []interface{}{30, 40},
	}
	if EvaluateCondition(cond6, data) {
		t.Error("Numeric range validation should be false")
	}

	// Test 4: Check registered operators
	registeredOps := GetRegisteredCustomOperators()
	expectedOps := []Operator{"case_insensitive_eq", "email_domain", "in_range"}
	if len(registeredOps) != len(expectedOps) {
		t.Errorf("Expected %d registered operators, got %d", len(expectedOps), len(registeredOps))
	}

	// Verify all expected operators are registered
	opMap := make(map[Operator]bool)
	for _, op := range registeredOps {
		opMap[op] = true
	}
	for _, expectedOp := range expectedOps {
		if !opMap[expectedOp] {
			t.Errorf("Expected operator '%s' to be registered", expectedOp)
		}
	}

	// Test 5: Unregister an operator
	UnregisterCustomOperator("email_domain")

	// This should now fail because the operator is unregistered
	if EvaluateCondition(cond3, data) {
		t.Error("Unregistered operator should return false")
	}

	// Verify operator was removed from registry
	registeredOpsAfter := GetRegisteredCustomOperators()
	if len(registeredOpsAfter) != 2 {
		t.Errorf("Expected 2 registered operators after unregistering, got %d", len(registeredOpsAfter))
	}

	// Test 6: Custom operator with missing key handling
	RegisterCustomOperator("handle_missing", func(fieldValue, expectedValue interface{}) bool {
		// This operator returns true if the field is missing and expected value is "missing"
		if fieldValue == nil && expectedValue == "missing" {
			return true
		}
		return false
	})

	cond7 := Conditions{
		Key:      "nonexistent",
		Operator: "handle_missing",
		Value:    "missing",
	}
	if !EvaluateCondition(cond7, data) {
		t.Error("Custom operator should handle missing keys")
	}

	// Test 7: Custom operator in complex condition
	complexCond := Conditions{
		Logic: LogicAnd,
		Children: []Conditions{
			{Key: "age", Operator: OperatorGt, Value: 18},
			{Key: "name", Operator: "case_insensitive_eq", Value: "john doe"},
		},
	}
	if !EvaluateCondition(complexCond, data) {
		t.Error("Complex condition with custom operator should be true")
	}

	// Clean up
	for _, op := range GetRegisteredCustomOperators() {
		UnregisterCustomOperator(op)
	}
}

func TestCustomOperatorEdgeCases(t *testing.T) {
	// Clean up any existing custom operators
	for _, op := range GetRegisteredCustomOperators() {
		UnregisterCustomOperator(op)
	}

	data := map[string]interface{}{
		"value": "test",
	}

	// Test 1: Custom operator that panics (should not crash the evaluation)
	RegisterCustomOperator("panic_operator", func(fieldValue, expectedValue interface{}) bool {
		panic("This operator panics!")
	})

	cond := Conditions{
		Key:      "value",
		Operator: "panic_operator",
		Value:    "anything",
	}

	// This should not panic the entire evaluation
	result := EvaluateCondition(cond, data)
	if result {
		t.Error("Panicking operator should return false")
	}

	// Test 2: Custom operator with nil validator (should panic when registering)
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Registering nil validator should panic")
			}
		}()
		RegisterCustomOperator("nil_operator", nil)
	}()

	// Test 3: Thread safety test
	RegisterCustomOperator("thread_safe", func(fieldValue, expectedValue interface{}) bool {
		return true
	})

	// Run concurrent access to test thread safety
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			GetRegisteredCustomOperators()
			EvaluateCondition(Conditions{
				Key:      "value",
				Operator: "thread_safe",
				Value:    "test",
			}, data)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 100; i++ {
		<-done
	}

	// Clean up
	for _, op := range GetRegisteredCustomOperators() {
		UnregisterCustomOperator(op)
	}
}

func TestQuickCustomOperatorDemo(t *testing.T) {
	// Clean up any existing custom operators
	for _, op := range GetRegisteredCustomOperators() {
		UnregisterCustomOperator(op)
	}

	// Register a custom operator
	RegisterCustomOperator("case_insensitive_eq", func(fieldValue, expectedValue interface{}) bool {
		str1 := strings.ToLower(fmt.Sprintf("%v", fieldValue))
		str2 := strings.ToLower(fmt.Sprintf("%v", expectedValue))
		return str1 == str2
	})

	data := map[string]interface{}{
		"name": "John Doe",
	}

	condition := Conditions{
		Key:      "name",
		Operator: "case_insensitive_eq",
		Value:    "JOHN DOE",
	}

	result := EvaluateCondition(condition, data)
	if !result {
		t.Error("Custom operator should return true for case insensitive match")
	}

	// Check registered operators
	ops := GetRegisteredCustomOperators()
	if len(ops) != 1 || ops[0] != "case_insensitive_eq" {
		t.Errorf("Expected 1 registered operator 'case_insensitive_eq', got %v", ops)
	}

	// Clean up
	UnregisterCustomOperator("case_insensitive_eq")

	fmt.Printf("✅ Custom operator demo test passed!\n")
	fmt.Printf("   - Registered custom operator: %v\n", ops[0])
	fmt.Printf("   - Evaluated condition successfully: %v\n", result)
}
