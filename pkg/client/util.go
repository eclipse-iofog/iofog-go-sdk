package client

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/arch"
)

// ArchNameToID maps canonical architecture names to Controller/Edgelet integer codes.
var ArchNameToID = arch.NameToID

// ArchIDToName maps Controller/Edgelet integer codes to canonical architecture names.
var ArchIDToName = arch.IDToName

// RegistryTypeRegistryTypeIDDict Map from string registry type to numeric id
var RegistryTypeRegistryTypeIDDict = map[string]int{
	"remote": 1,
	"local":  2,
}

// RegistryTypeIDRegistryTypeDict Map from numeric id registry type to string
var RegistryTypeIDRegistryTypeDict = map[int]string{
	1: "remote",
	2: "local",
}

func getString(in io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(in); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func checkStatusCode(code int, method, url string, body io.Reader) error {
	if code < 200 || code >= 300 {
		bodyString, err := getString(body)
		if err != nil {
			return err
		}
		switch code {
		case 404:
			return NewNotFoundError(fmt.Sprintf("Received Not found from %s %s\n: %s\n", method, url, bodyString))
		default:
			return NewHTTPError(fmt.Sprintf("Received %d from %s %s\n%s", code, method, url, bodyString), code)
		}
	}
	return nil
}

func before(input, substr string) string {
	pos := strings.Index(input, substr)
	if pos == -1 {
		return input
	}
	return input[0:pos]
}
