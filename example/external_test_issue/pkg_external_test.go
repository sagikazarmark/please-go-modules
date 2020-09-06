package external_test_issue_test

import (
	"testing"

	"github.com/sagikazarmark/please-go-modules/example/external_test_issue"
)

func TestHello(t *testing.T) {
	hello := external_test_issue.Hello()

	if hello != external_test_issue.CompareWith {
		t.Error("hello mismatch")
	}
}
