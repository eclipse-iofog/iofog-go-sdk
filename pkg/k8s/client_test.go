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

package k8s

import (
	"testing"
)

func TestCreation(t *testing.T) {
	// Here just to test compilation
	client := &Client{}
	if client == nil {
		t.Error("This is impossible")
	}
}
