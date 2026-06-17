package operator

import (
	"fmt"
	"testing"
	"time"
)

func TestReconciliation(t *testing.T) {
	if ReconcileWithError(fmt.Errorf("err")).Err == nil {
		t.Error("ReconcileWithError.Err is nil")
	}
	if !ReconcileWithRequeue(time.Second).Requeue {
		t.Error("ReconcileWithRequeue.Requeue is false")
	}
	if ReconcileWithRequeue(time.Second).Delay != time.Second {
		t.Error("ReconcileWithRequeue.Delay is wrong")
	}
	if Continue().Requeue || Continue().End {
		t.Error("Continue.Requeue or Continue.End is not false")
	}
	if !Reconcile().End {
		t.Error("Reconcile.End is false")
	}
}
