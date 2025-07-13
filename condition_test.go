package jsonvaluate

import (
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
