package operator

import (
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

type Reconciliation struct {
	End     bool
	Requeue bool
	Delay   time.Duration
	Err     error
}

func (recon *Reconciliation) IsFinal() bool {
	return recon.Err != nil || recon.Requeue || recon.End
}

func (recon *Reconciliation) Result() (ctrl.Result, error) {
	if recon.Err != nil {
		return RequeueWithError(recon.Err)
	}
	if recon.End {
		return DoNotRequeue()
	}
	return RequeueWithDelay(recon.Delay)
}

func ReconcileWithError(err error) Reconciliation {
	return Reconciliation{
		Err: err,
	}
}

func ReconcileWithRequeue(delay time.Duration) Reconciliation {
	return Reconciliation{
		Requeue: true,
		Delay:   delay,
	}
}

func Continue() Reconciliation {
	return Reconciliation{}
}

func Reconcile() Reconciliation {
	return Reconciliation{
		End: true,
	}
}
