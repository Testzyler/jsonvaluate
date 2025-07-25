// Package evaluate provides a flexible JSON condition evaluation library.
// It allows you to define complex conditional logic using operators and logical groupings
// that can be evaluated against JSON-like data structures.
package jsonvaluate

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Operator represents the type of comparison operation to perform.
type Operator string

// Available operators for condition evaluation
const (
	OperatorEq         Operator = "=="         // Equal to
	OperatorNeq        Operator = "!="         // Not equal to
	OperatorGt         Operator = ">"          // Greater than
	OperatorGte        Operator = ">="         // Greater than or equal to
	OperatorLt         Operator = "<"          // Less than
	OperatorLte        Operator = "<="         // Less than or equal to
	OperatorIn         Operator = "in"         // Value is in collection
	OperatorNin        Operator = "nin"        // Value is not in collection
	OperatorContains   Operator = "contains"   // String contains substring
	OperatorNcontains  Operator = "ncontains"  // String does not contain substring
	OperatorIsnull     Operator = "isnull"     // Value is null or doesn't exist
	OperatorIsnotnull  Operator = "isnotnull"  // Value is not null and exists
	OperatorIsEmpty    Operator = "isempty"    // Value is empty (empty string, array, etc.)
	OperatorIsNotEmpty Operator = "isnotempty" // Value is not empty
	OperatorIsTrue     Operator = "istrue"     // Value is true (boolean or truthy)
	OperatorIsFalse    Operator = "isfalse"    // Value is false (boolean or falsy)
	OperatorLike       Operator = "like"       // SQL-like pattern matching (case sensitive)
	OperatorIlike      Operator = "ilike"      // SQL-like pattern matching (case insensitive)
	OperatorNlike      Operator = "nlike"      // NOT SQL-like pattern matching
	OperatorStartsWith Operator = "startswith" // String starts with prefix
	OperatorEndsWith   Operator = "endswith"   // String ends with suffix
	OperatorBetween    Operator = "between"    // Value is between two bounds (inclusive)
	OperatorNotBetween Operator = "notbetween" // Value is not between two bounds
)

// Logic represents the logical operation for combining multiple conditions.
type Logic string

// Available logical operators
const (
	LogicAnd Logic = "AND" // All conditions must be true
	LogicOr  Logic = "OR"  // At least one condition must be true
)

// Conditions represents a condition tree that can be either a single condition
// or a group of conditions combined with logical operators (AND/OR).
//
// For single conditions, use Key, Operator, and Value fields.
// For group conditions, use Logic and Children fields.
//
// Example single condition:
//
//	cond := Conditions{
//	    Key:      "age",
//	    Operator: OperatorGt,
//	    Value:    18,
//	}
//
// Example group condition:
//
//	cond := Conditions{
//	    Logic: LogicAnd,
//	    Children: []Conditions{
//	        {Key: "age", Operator: OperatorGt, Value: 18},
//	        {Key: "country", Operator: OperatorEq, Value: "US"},
//	    },
//	}
type Conditions struct {
	Logic    Logic        `json:"logic,omitempty"`    // "AND" or "OR" for group, empty for single
	Children []Conditions `json:"children,omitempty"` // Child conditions for group

	Key      string      `json:"key,omitempty"`      // Field key for single condition
	Operator Operator    `json:"operator,omitempty"` // Comparison operator for single condition
	Value    interface{} `json:"value,omitempty"`    // Expected value for single condition
}

// CustomOperatorValidator defines the function signature for custom operator validation.
// It takes the field value from the data and the expected value from the condition,
// and returns true if the condition is satisfied.
type CustomOperatorValidator func(fieldValue, expectedValue interface{}) bool

// Thread-safe registry for custom operators
var (
	customOperators = make(map[Operator]CustomOperatorValidator)
	customOpsMutex  sync.RWMutex
)

// RegisterCustomOperator registers a new custom operator with its validation function.
// The operator name should be unique and not conflict with built-in operators.
// The validator function will be called with the field value and expected value.
//
// Example:
//
//	RegisterCustomOperator("case_insensitive_eq", func(fieldValue, expectedValue interface{}) bool {
//	    str1 := strings.ToLower(fmt.Sprintf("%v", fieldValue))
//	    str2 := strings.ToLower(fmt.Sprintf("%v", expectedValue))
//	    return str1 == str2
//	})
func RegisterCustomOperator(operator Operator, validator CustomOperatorValidator) {
	if validator == nil {
		panic("custom operator validator cannot be nil")
	}

	customOpsMutex.Lock()
	defer customOpsMutex.Unlock()
	customOperators[operator] = validator
}

// UnregisterCustomOperator removes a custom operator from the registry.
// Built-in operators cannot be unregistered.
func UnregisterCustomOperator(operator Operator) {
	customOpsMutex.Lock()
	defer customOpsMutex.Unlock()
	delete(customOperators, operator)
}

// GetRegisteredCustomOperators returns a list of all registered custom operators.
func GetRegisteredCustomOperators() []Operator {
	customOpsMutex.RLock()
	defer customOpsMutex.RUnlock()

	operators := make([]Operator, 0, len(customOperators))
	for op := range customOperators {
		operators = append(operators, op)
	}
	return operators
}

// EvaluateCondition evaluates a condition tree against the provided data.
// It returns true if the condition is satisfied, false otherwise.
//
// The data parameter should be a map where keys correspond to the field names
// used in the conditions, and values are the actual data to evaluate against.
//
// For group conditions (with Logic field set), it evaluates all children:
//   - AND logic: returns true only if ALL children evaluate to true
//   - OR logic: returns true if ANY child evaluates to true
//
// For single conditions, it compares the data field value against the expected
// value using the specified operator.
//
// Example usage:
//
//	data := map[string]interface{}{
//	    "age":     25,
//	    "country": "US",
//	}
//
//	condition := Conditions{
//	    Key:      "age",
//	    Operator: OperatorGt,
//	    Value:    18,
//	}
//
//	result := EvaluateCondition(condition, data) // returns true
func EvaluateCondition(cond Conditions, data map[string]interface{}) bool {
	// Handle group conditions (AND/OR logic)
	if cond.Logic != "" && len(cond.Children) > 0 {
		switch cond.Logic {
		case LogicAnd:
			for _, child := range cond.Children {
				if !EvaluateCondition(child, data) {
					return false
				}
			}
			return true
		case LogicOr:
			for _, child := range cond.Children {
				if EvaluateCondition(child, data) {
					return true
				}
			}
			return false
		}
	}

	// Handle single conditions
	if cond.Key != "" && cond.Operator != "" {
		return evalSingleCondition(cond.Key, cond.Operator, cond.Value, data)
	}

	// Default case for empty conditions
	return true
}

// evalSingleCondition evaluates a single condition against the data
func evalSingleCondition(key string, op Operator, value interface{}, data map[string]interface{}) bool {
	v, exists := data[key]

	switch op {
	case OperatorIsnull:
		return !exists || v == nil
	case OperatorIsnotnull:
		return exists && v != nil
	case OperatorIsEmpty:
		return isEmpty(v)
	case OperatorIsNotEmpty:
		return !isEmpty(v)
	case OperatorIsTrue:
		return toBool(v)
	case OperatorIsFalse:
		return !toBool(v)
	}

	// For other built-in operators, the key must exist
	if !exists {
		// Check if this is a custom operator first
		customOpsMutex.RLock()
		validator, isCustom := customOperators[op]
		customOpsMutex.RUnlock()

		if isCustom {
			// Handle panics in custom operators gracefully
			defer func() {
				if r := recover(); r != nil {
					// Custom operator panicked, return false
				}
			}()
			return validator(v, value) // v will be nil for missing keys
		}

		return false
	}

	switch op {
	case OperatorEq:
		return isEqual(v, value)
	case OperatorNeq:
		return !isEqual(v, value)
	case OperatorGt:
		return compareValues(v, value) > 0
	case OperatorGte:
		return compareValues(v, value) >= 0
	case OperatorLt:
		return compareValues(v, value) < 0
	case OperatorLte:
		return compareValues(v, value) <= 0
	case OperatorIn:
		return isIn(v, value)
	case OperatorNin:
		return !isIn(v, value)
	case OperatorContains:
		return contains(v, value)
	case OperatorNcontains:
		return !contains(v, value)
	case OperatorLike:
		return like(v, value, false)
	case OperatorIlike:
		return like(v, value, true)
	case OperatorNlike:
		return !like(v, value, false)
	case OperatorStartsWith:
		return startsWith(v, value)
	case OperatorEndsWith:
		return endsWith(v, value)
	case OperatorBetween:
		return between(v, value)
	case OperatorNotBetween:
		return !between(v, value)
	default:
		// Check for custom operators
		customOpsMutex.RLock()
		validator, exists := customOperators[op]
		customOpsMutex.RUnlock()

		if exists {
			// Handle panics in custom operators gracefully
			defer func() {
				if r := recover(); r != nil {
					// Custom operator panicked, return false
				}
			}()
			return validator(v, value)
		}

		return false
	}
}

// Helper functions

// isEmpty checks if a value is considered empty
func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		return rv.String() == ""
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return rv.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return rv.IsNil()
	default:
		return false
	}
}

// toBool converts various types to boolean
func toBool(v interface{}) bool {
	if v == nil {
		return false
	}

	switch val := v.(type) {
	case bool:
		return val
	case string:
		return strings.ToLower(val) == "true"
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(val).Int() != 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(val).Uint() != 0
	case float32, float64:
		return reflect.ValueOf(val).Float() != 0
	default:
		return !isEmpty(v)
	}
}

// isEqual checks equality between two values
func isEqual(v1, v2 interface{}) bool {
	if v1 == nil && v2 == nil {
		return true
	}
	if v1 == nil || v2 == nil {
		return false
	}

	// Try direct comparison first
	if reflect.DeepEqual(v1, v2) {
		return true
	}

	// Try numeric comparison
	if n1, ok1 := toNumber(v1); ok1 {
		if n2, ok2 := toNumber(v2); ok2 {
			return n1 == n2
		}
	}

	// Try string comparison
	return toString(v1) == toString(v2)
}

// compareValues compares two values and returns -1, 0, or 1
func compareValues(v1, v2 interface{}) int {

	// Try numeric comparison first
	if n1, ok1 := toNumber(v1); ok1 {
		if n2, ok2 := toNumber(v2); ok2 {
			if n1 < n2 {
				return -1
			} else if n1 > n2 {
				return 1
			}
			return 0
		}
	}

	// Try time comparison
	if t1, ok1 := toTime(v1); ok1 {
		if t2, ok2 := toTime(v2); ok2 {
			// Debug output

			if t1.Before(t2) {
				return -1
			} else if t1.After(t2) {
				return 1
			}
			return 0
		} else {
		}
	} else {
	}

	// Fall back to string comparison
	s1, s2 := toString(v1), toString(v2)
	if s1 < s2 {
		return -1
	} else if s1 > s2 {
		return 1
	}
	return 0
}

// toNumber converts various types to float64
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
		if f, err := parseFloat(val); err == nil {
			return f, true
		}
	}
	return 0, false
}

// parseFloat parses a string to float64 with strict validation
func parseFloat(s string) (float64, error) {
	// Use strconv.ParseFloat for proper validation
	// This will only succeed if the entire string is a valid number
	return strconv.ParseFloat(s, 64)
}

// toString converts any value to string
func toString(v interface{}) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case fmt.Stringer:
		return val.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}

// toTime converts various types to time.Time
func toTime(v interface{}) (time.Time, bool) {
	switch val := v.(type) {
	case time.Time:
		return val, true
	case string:
		// Try common time formats
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02 15:04:05",
			"2006-01-02",
			"15:04:05",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, val); err == nil {
				return t, true
			}
		}
	case int64:
		return time.Unix(val, 0), true
	}
	return time.Time{}, false
}

// isIn checks if value is in the collection
func isIn(v, collection interface{}) bool {
	if collection == nil {
		return false
	}

	cv := reflect.ValueOf(collection)
	switch cv.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < cv.Len(); i++ {
			if isEqual(v, cv.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, key := range cv.MapKeys() {
			if isEqual(v, key.Interface()) {
				return true
			}
		}
	case reflect.String:
		return strings.Contains(cv.String(), toString(v))
	}
	return false
}

// contains checks if haystack contains needle
func contains(haystack, needle interface{}) bool {
	if haystack == nil || needle == nil {
		return false
	}

	haystackStr := toString(haystack)
	needleStr := toString(needle)
	return strings.Contains(haystackStr, needleStr)
}

// like performs SQL-like pattern matching
func like(v, pattern interface{}, caseInsensitive bool) bool {
	if v == nil || pattern == nil {
		return false
	}

	str := toString(v)
	pat := toString(pattern)

	if caseInsensitive {
		str = strings.ToLower(str)
		pat = strings.ToLower(pat)
	}

	// Convert SQL LIKE pattern to regex
	// % matches any sequence of characters
	// _ matches any single character
	regexPattern := strings.ReplaceAll(pat, "%", ".*")
	regexPattern = strings.ReplaceAll(regexPattern, "_", ".")
	regexPattern = "^" + regexPattern + "$"

	matched, err := regexp.MatchString(regexPattern, str)
	return err == nil && matched
}

// startsWith checks if string starts with prefix
func startsWith(v, prefix interface{}) bool {
	if v == nil || prefix == nil {
		return false
	}

	str := toString(v)
	pre := toString(prefix)
	return strings.HasPrefix(str, pre)
}

// endsWith checks if string ends with suffix
func endsWith(v, suffix interface{}) bool {
	if v == nil || suffix == nil {
		return false
	}

	str := toString(v)
	suf := toString(suffix)
	return strings.HasSuffix(str, suf)
}

// between checks if value is between two bounds (inclusive)
func between(v, bounds interface{}) bool {
	if v == nil || bounds == nil {
		return false
	}

	// bounds should be a slice with 2 elements [min, max]
	boundsSlice := reflect.ValueOf(bounds)
	if boundsSlice.Kind() != reflect.Slice || boundsSlice.Len() != 2 {
		return false
	}

	min := boundsSlice.Index(0).Interface()
	max := boundsSlice.Index(1).Interface()

	return compareValues(v, min) >= 0 && compareValues(v, max) <= 0
}

// ConditionGroup represents a more flexible condition structure that allows
// different logical operations between different pairs of conditions.
type ConditionGroup struct {
	Conditions []ConditionWithLogic `json:"conditions"`
}

// ConditionWithLogic represents a single condition with an optional logical operator
// that connects it to the next condition.
type ConditionWithLogic struct {
	// Single condition fields
	Key      string      `json:"key,omitempty"`      // Field key for condition
	Operator Operator    `json:"operator,omitempty"` // Comparison operator
	Value    interface{} `json:"value,omitempty"`    // Expected value

	// Group condition (alternative to single condition)
	Group *ConditionGroup `json:"group,omitempty"` // Nested group of conditions

	// Logic operator to connect to the next condition
	NextLogic Logic `json:"next_logic,omitempty"` // "AND" or "OR" to connect to next condition
}

// EvaluateConditionGroup evaluates a ConditionGroup against the provided data.
// This allows for more flexible logical expressions between conditions.
//
// Example usage:
//
//	group := ConditionGroup{
//	    Conditions: []ConditionWithLogic{
//	        {Key: "sum_insured", Operator: OperatorGte, Value: 200000, NextLogic: LogicAnd},
//	        {
//	            Group: &ConditionGroup{
//	                Conditions: []ConditionWithLogic{
//	                    {Key: "amount", Operator: OperatorGte, Value: 100000, NextLogic: LogicOr},
//	                    {Key: "amount", Operator: OperatorLte, Value: 1000000},
//	                },
//	            },
//	            NextLogic: LogicAnd,
//	        },
//	        {Key: "percent_of_sum_insured", Operator: "%of", Value: 20},
//	    },
//	}
func EvaluateConditionGroup(group ConditionGroup, data map[string]interface{}) bool {
	if len(group.Conditions) == 0 {
		return true
	}

	// Evaluate first condition
	result := evaluateConditionWithLogic(group.Conditions[0], data)

	// Process remaining conditions with their logic operators
	for i := 1; i < len(group.Conditions); i++ {
		prevCondition := group.Conditions[i-1]
		currentResult := evaluateConditionWithLogic(group.Conditions[i], data)

		// Apply the logic operator from the previous condition
		switch prevCondition.NextLogic {
		case LogicAnd:
			result = result && currentResult
		case LogicOr:
			result = result || currentResult
		default:
			// If no logic specified, default to AND
			result = result && currentResult
		}
	}

	return result
}

// evaluateConditionWithLogic evaluates a single ConditionWithLogic
func evaluateConditionWithLogic(condition ConditionWithLogic, data map[string]interface{}) bool {
	// If it's a group condition, evaluate the group
	if condition.Group != nil {
		return EvaluateConditionGroup(*condition.Group, data)
	}

	// Otherwise, evaluate as a single condition
	return evalSingleCondition(condition.Key, condition.Operator, condition.Value, data)
}

// Helper functions for creating common condition patterns

// NewSimpleCondition creates a simple condition with key, operator, and value.
// This is a convenience function for creating single conditions.
func NewSimpleCondition(key string, operator Operator, value interface{}) Conditions {
	return Conditions{
		Key:      key,
		Operator: operator,
		Value:    value,
	}
}

// NewAndGroup creates an AND group condition from a list of child conditions.
// All child conditions must evaluate to true for the group to be true.
func NewAndGroup(children ...Conditions) Conditions {
	return Conditions{
		Logic:    LogicAnd,
		Children: children,
	}
}

// NewOrGroup creates an OR group condition from a list of child conditions.
// At least one child condition must evaluate to true for the group to be true.
func NewOrGroup(children ...Conditions) Conditions {
	return Conditions{
		Logic:    LogicOr,
		Children: children,
	}
}

// Helper functions for creating flexible condition patterns

// NewConditionGroup creates a new ConditionGroup with the specified conditions.
func NewConditionGroup(conditions ...ConditionWithLogic) ConditionGroup {
	return ConditionGroup{
		Conditions: conditions,
	}
}

// NewConditionWithLogic creates a single condition with logic operator for the next condition.
func NewConditionWithLogic(key string, operator Operator, value interface{}, nextLogic Logic) ConditionWithLogic {
	return ConditionWithLogic{
		Key:       key,
		Operator:  operator,
		Value:     value,
		NextLogic: nextLogic,
	}
}

// NewGroupConditionWithLogic creates a group condition with logic operator for the next condition.
func NewGroupConditionWithLogic(group ConditionGroup, nextLogic Logic) ConditionWithLogic {
	return ConditionWithLogic{
		Group:     &group,
		NextLogic: nextLogic,
	}
}

// ConvertToConditionGroup converts the traditional nested Conditions structure
// to the new flexible ConditionGroup structure.
func ConvertToConditionGroup(conditions Conditions) ConditionGroup {
	// If it's a single condition
	if conditions.Key != "" {
		return ConditionGroup{
			Conditions: []ConditionWithLogic{
				{
					Key:      conditions.Key,
					Operator: conditions.Operator,
					Value:    conditions.Value,
				},
			},
		}
	}

	// If it's a group condition
	if len(conditions.Children) == 0 {
		return ConditionGroup{}
	}

	var conditionsWithLogic []ConditionWithLogic
	for i, child := range conditions.Children {
		var nextLogic Logic
		// For all conditions except the last one, use the group's logic
		if i < len(conditions.Children)-1 {
			nextLogic = conditions.Logic
		}

		// If child is a single condition
		if child.Key != "" {
			conditionsWithLogic = append(conditionsWithLogic, ConditionWithLogic{
				Key:       child.Key,
				Operator:  child.Operator,
				Value:     child.Value,
				NextLogic: nextLogic,
			})
		} else {
			// If child is a nested group, convert it recursively
			childGroup := ConvertToConditionGroup(child)
			conditionsWithLogic = append(conditionsWithLogic, ConditionWithLogic{
				Group:     &childGroup,
				NextLogic: nextLogic,
			})
		}
	}

	return ConditionGroup{
		Conditions: conditionsWithLogic,
	}
}

// EvaluateFlexibleCondition evaluates either the traditional Conditions structure
// or the new ConditionGroup structure against the provided data.
func EvaluateFlexibleCondition(conditions interface{}, data map[string]interface{}) bool {
	switch cond := conditions.(type) {
	case Conditions:
		return EvaluateCondition(cond, data)
	case ConditionGroup:
		return EvaluateConditionGroup(cond, data)
	case *Conditions:
		return EvaluateCondition(*cond, data)
	case *ConditionGroup:
		return EvaluateConditionGroup(*cond, data)
	default:
		return false
	}
}
