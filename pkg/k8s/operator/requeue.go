package operator

import (
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

func DoNotRequeue() (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func RequeueWithError(err error) (ctrl.Result, error) {
	return ctrl.Result{}, err
}

func Requeue() (ctrl.Result, error) {
	return ctrl.Result{
		Requeue: true,
	}, nil
}

func RequeueWithDelay(delay time.Duration) (ctrl.Result, error) {
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: delay,
	}, nil
}
