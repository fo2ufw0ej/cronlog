package notify

import "fmt"

// ParseCondition converts a string to a Condition, returning an error for unknown values.
func ParseCondition(s string) (Condition, error) {
	switch Condition(s) {
	case ConditionAlways, ConditionFailure, ConditionNever:
		return Condition(s), nil
	case "":
		return ConditionAlways, nil
	}
	return "", fmt.Errorf("notify: unknown condition %q (want always|failure|never)", s)
}

// String returns the string representation of a Condition.
func (c Condition) String() string {
	return string(c)
}
