/*
 *  *******************************************************************************
 *  * Copyright (c) 2019 Edgeworx, Inc.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

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
