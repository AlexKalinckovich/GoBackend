package validation

import (
	"testing"
)

func TestValidationErrorCreation(t *testing.T) {
	runValidationErrorSuite(t)
}

func runValidationErrorSuite(t *testing.T) {
	err := NewValidationError("email", "invalid format")
	assertValidationErrorCode(t, err)
	assertValidationErrorString(t, err)
	assertValidationErrorField(t, err)
	assertValidationErrorContext(t, err)
}

func assertValidationErrorCode(t *testing.T, err *Error) {
	expected := ErrorCode
	actual := err.Code()
	if actual != expected {
		t.Errorf("code mismatch: expected %s, got %s", expected, actual)
	}
}

func assertValidationErrorString(t *testing.T, err *Error) {
	expected := "invalid format"
	actual := err.Error()
	if actual != expected {
		t.Errorf("error string mismatch: expected %s, got %s", expected, actual)
	}
}

func assertValidationErrorField(t *testing.T, err *Error) {
	expected := "email"
	actual := err.Field()
	if actual != expected {
		t.Errorf("field mismatch: expected %s, got %s", expected, actual)
	}
}

func assertValidationErrorContext(t *testing.T, err *Error) {
	ctx := err.ContextData()
	verifyContextField(t, ctx)
	verifyContextMessage(t, ctx)
}

func verifyContextField(t *testing.T, ctx map[string]any) {
	value, ok := ctx["field"].(string)
	if !ok || value != "email" {
		t.Errorf("context field mismatch: got %v", ctx["field"])
	}
}

func verifyContextMessage(t *testing.T, ctx map[string]any) {
	value, ok := ctx["message"].(string)
	if !ok || value != "invalid format" {
		t.Errorf("context message mismatch: got %v", ctx["message"])
	}
}

func TestValidationAggregateErrorLifecycle(t *testing.T) {
	runAggregateLifecycleSuite(t)
}

func runAggregateLifecycleSuite(t *testing.T) {
	agg := NewValidationAggregateError()
	assertAggregateEmpty(t, agg)
	aggregateAddSingleField(t, agg)
	assertAggregateSingleField(t, agg)
	aggregateAddMultipleFields(t, agg)
	assertAggregateMultipleFields(t, agg)
	assertAggregateCode(t, agg)
	assertAggregateErrorString(t, agg)
}

func assertAggregateEmpty(t *testing.T, agg *AggregateError) {
	hasErrors := agg.HasErrors()
	if hasErrors {
		t.Errorf("expected empty aggregate to report no errors")
	}
}

func aggregateAddSingleField(t *testing.T, agg *AggregateError) {
	agg.AddField("email", "invalid format")
}

func assertAggregateSingleField(t *testing.T, agg *AggregateError) {
	hasErrors := agg.HasErrors()
	if !hasErrors {
		t.Errorf("expected aggregate to report errors after adding one")
	}
}

func aggregateAddMultipleFields(t *testing.T, agg *AggregateError) {
	agg.AddField("age", "must be positive")
}

func assertAggregateMultipleFields(t *testing.T, agg *AggregateError) {
	ctx := agg.ContextData()
	details := extractDetailsFromContext(ctx)
	verifyDetailsCount(t, details)
}

func extractDetailsFromContext(ctx map[string]any) map[string]string {
	details, ok := ctx["details"].(map[string]string)
	if !ok {
		return make(map[string]string)
	}
	return details
}

func verifyDetailsCount(t *testing.T, details map[string]string) {
	if len(details) != 2 {
		t.Errorf("expected 2 details, got %d", len(details))
	}
}

func assertAggregateCode(t *testing.T, agg *AggregateError) {
	expected := AggregateErrorCode
	actual := agg.Code()
	if actual != expected {
		t.Errorf("aggregate code mismatch: expected %s, got %s", expected, actual)
	}
}

func assertAggregateErrorString(t *testing.T, agg *AggregateError) {
	expected := "validation failed"
	actual := agg.Error()
	if actual != expected {
		t.Errorf("aggregate error string mismatch: expected %s, got %s", expected, actual)
	}
}
