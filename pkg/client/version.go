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

package client

import (
	"errors"
	"fmt"
	"strconv"
)

func (clt *Client) GetVersion() string {
	return clt.status.version
}

func (clt *Client) GetVersionNumbers() (major, minor, patch int, err error) {
	errMsg := fmt.Sprintf("Controller did not return a valid API version: %s", clt.status.version)

	// Split version string
	if len(clt.status.versionNums) != 3 {
		err = errors.New(errMsg)
		return
	}

	// Convert to int
	major, majErr := strconv.Atoi(clt.status.versionNums[0])
	minor, minErr := strconv.Atoi(clt.status.versionNums[1])
	patch, patErr := strconv.Atoi(clt.status.versionNums[2])
	if majErr != nil || minErr != nil || patErr != nil {
		err = errors.New(errMsg)
		return
	}

	return
}

func (clt *Client) IsDevVersion() bool {
	return clt.status.version == "0.0.0-dev"
}

func (clt *Client) IsEdgeResourceCapable() error {
	if clt.IsDevVersion() {
		return nil
	}
	// Decode version
	major, minor, _, err := clt.GetVersionNumbers()
	if err != nil {
		return err
	}
	// Supported
	if major >= 2 && minor >= 1 {
		return nil
	}
	// Unsupported
	return fmt.Errorf("controller version %s does not support Edge Resources", clt.status.version)
}
