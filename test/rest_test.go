package resttest

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/arch"
	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client"
)

type testState struct {
	email           string
	password        string
	name            string
	surname         string
	url             *url.URL
	agent           string
	uuid            string
	archID          int64
	appTemplateName string
	appName         string
}

var state = testState{
	email:           "serge@edgeworx.io",
	password:        "wfhoi982bv1sfdjoi",
	name:            "Serge",
	surname:         "Radinovich",
	agent:           "agent-1",
	archID:          arch.AMD64,
	appTemplateName: "apptemplate1",
	appName:         "app-1",
}

var clt *client.Client

func TestNewAndLogin(t *testing.T) {
	// client.SetVerbosity(true)

	var existingState = testState{
		email:    "user@domain.com",
		password: "g9hr823rhuoi",
		name:     "Foo",
		surname:  "Bar",
	}
	var err error
	state.url, err = url.Parse("http://localhost:51121/api/v3")
	if err != nil {
		t.Error(err)
	}
	opt := client.Options{
		BaseURL: state.url,
	}

	clt, err := client.NewAndLogin(opt, existingState.email, existingState.password)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Failed to create client and login: %s", err.Error()))
	}

	_, err = clt.GetStatus()
	if err != nil {
		t.Fatalf(fmt.Sprintf("Failed to get status: %s", err.Error()))
	}
}

func TestNewAndCreate(t *testing.T) {
	opt := client.Options{
		BaseURL: state.url,
	}
	adminClt, err := client.NewAndLogin(opt, "user@domain.com", "g9hr823rhuoi")
	if err != nil {
		t.Fatalf(fmt.Sprintf("Failed to login as admin: %s", err.Error()))
	}

	if _, err := adminClt.CreateAuthUser(client.AuthUserCreateRequest{
		Email:    state.email,
		Password: state.password,
	}); err != nil {
		t.Fatalf(fmt.Sprintf("Failed to create user: %s", err.Error()))
	}

	clt, err = client.NewAndLogin(opt, state.email, state.password)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Failed to login: %s", err.Error()))
	}

	_, err = clt.GetStatus()
	if err != nil {
		t.Fatalf(fmt.Sprintf("Failed to get status: %s", err.Error()))
	}
}

func TestCreateAgent(t *testing.T) {
	request := &client.CreateAgentRequest{}
	archID := state.archID
	request.ArchID = &archID
	request.Name = state.agent
	host := "localhost"
	request.Host = &host

	response, err := clt.CreateAgent(request)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Failed to create Agent: %s", err.Error()))
	}

	getResponse, err := clt.GetAgentByID(response.UUID)
	if err != nil {
		t.Fatalf((fmt.Sprintf("Failed to get Agent by UUID: %s", err.Error())))
	}

	if getResponse.Name != request.Name {
		t.Fatalf(fmt.Sprintf("Controller returned unexpected Agent name: %s", getResponse.Name))
	}

	nameInfo, err := clt.GetAgentByName(state.agent, false)
	if err != nil {
		t.Fatalf("Failed to get Agent by name: %s", err.Error())
	}
	idInfo, err := clt.GetAgentByID(nameInfo.UUID)
	if err != nil {
		t.Fatalf("Failed to get Agent by UUID: %s", err.Error())
	}
	state.uuid = idInfo.UUID
}

func TestDeleteAgent(t *testing.T) {
	if err := clt.DeleteAgent(state.uuid); err != nil {
		t.Fatalf("Failed to delete Agent: %s", err.Error())
	}
}
