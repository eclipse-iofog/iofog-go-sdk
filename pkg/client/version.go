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

	if len(clt.status.versionNums) != 3 {
		return 0, 0, 0, errors.New(errMsg)
	}

	major, majErr := strconv.Atoi(clt.status.versionNums[0])
	minor, minErr := strconv.Atoi(clt.status.versionNums[1])
	patch, patErr := strconv.Atoi(clt.status.versionNums[2])
	if majErr != nil || minErr != nil || patErr != nil {
		return 0, 0, 0, errors.New(errMsg)
	}

	return major, minor, patch, nil
}
