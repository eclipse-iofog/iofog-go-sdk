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

package resttest

import (
	"fmt"
	"testing"

	"github.com/eclipse-iofog/iofog-go-sdk/v2/pkg/client"
)

type testState struct {
	email           string
	password        string
	name            string
	surname         string
	host            string
	port            int
	agent           string
	uuid            string
	fogType         int64
	appTemplateName string
	appName         string
}

var state = testState{
	email:           "serge@edgeworx.io",
	password:        "wfhoi982bv1sfdjoi",
	name:            "Serge",
	surname:         "Radinovich",
	host:            "localhost",
	port:            51121,
	agent:           "agent-1",
	fogType:         1, // x86
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
		host:     "localhost",
		port:     51121,
	}
	opt := client.Options{
		Endpoint: fmt.Sprintf("%s:%d", existingState.host, existingState.port),
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
		Endpoint: fmt.Sprintf("%s:%d", state.host, state.port),
	}
	clt = client.New(opt)

	if err := clt.CreateUser(client.User{
		Email:    state.email,
		Password: state.password,
		Name:     state.name,
		Surname:  state.surname,
	}); err != nil {
		t.Fatalf(fmt.Sprintf("Failed to create user : %s", err.Error()))
	}

	_, err := clt.GetStatus()
	if err != nil {
		t.Fatalf(fmt.Sprintf("Failed to get status: %s", err.Error()))
	}

	if err = clt.Login(client.LoginRequest{
		Email:    state.email,
		Password: state.password,
	}); err != nil {
		t.Fatalf(fmt.Sprintf("Failed to login: %s", err.Error()))
	}
}

func TestCreateAgent(t *testing.T) {
	request := client.CreateAgentRequest{}
	request.FogType = &state.fogType
	request.Name = state.agent
	request.Host = &state.host

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

func TestCreateUpdatePatchAppTemplate(t *testing.T) {
	// Create
	request := client.ApplicationTemplateCreateRequest{
		Description: "test desc",
		Name:        state.appTemplateName,
		Variables: []client.TemplateVariable{
			{
				Key:          "testkey",
				Description:  "vartestdesc",
				DefaultValue: "testdefaultval",
			},
		},
		Application: &client.ApplicationTemplateInfo{
			Microservices: []client.MicroserviceCreateRequest{},
			Routes:        []client.ApplicationRouteCreateRequest{},
		},
	}
	createName := "test1"
	request.Name = createName
	request.Description = "test2"
	response, err := clt.CreateApplicationTemplate(&request)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Failed to create App Template: %s", err.Error()))
	}
	if response.Id == 0 {
		t.Fatalf("Create App Template returned 0 Id")
	}
	if response.Name != request.Name {
		t.Fatalf(fmt.Sprintf("Create App Template returned wrong name: %s", response.Name))
	}

	// Update new and old
	updateName := "test123"
	request.Name = updateName
	for _, desc := range []string{"first", "second"} {
		request.Description = desc
		updateResponse, err := clt.UpdateApplicationTemplate(&request)
		if err != nil {
			t.Fatalf(fmt.Sprintf("Failed to update App Template: %s", err.Error()))
		}
		if updateResponse.Name != updateName {
			t.Fatalf(fmt.Sprintf("Update App Template returned wrong name: %s", updateResponse.Name))
		}
		getUpdateResponse, err := clt.GetApplicationTemplate(updateName)
		if err != nil {
			t.Fatalf(fmt.Sprintf("Failed to get updated App Template: %s", err.Error()))
		}
		if getUpdateResponse.Description != request.Description {
			t.Fatalf(fmt.Sprintf("Get updated App Template returned wrong description: %s", getUpdateResponse.Description))
		}
	}

	// Patch created
	if err := clt.UpdateApplicationTemplateMetadata(createName, &client.ApplicationTemplateMetadataUpdateRequest{
		Name: &state.appTemplateName,
	}); err != nil {
		t.Fatalf(fmt.Sprintf("Patch App Template failed: %s", err.Error()))
	}
}

func TestListAppTemplates(t *testing.T) {
	response, err := clt.ListApplicationTemplates()
	if err != nil {
		t.Fatalf(fmt.Sprintf("List App Templates failed: %s", err.Error()))
	}
	if len(response.ApplicationTemplates) != 2 {
		t.Fatalf(fmt.Sprintf("List App Templates returned incorrect count: %d", len(response.ApplicationTemplates)))
	}
	if response.ApplicationTemplates[0].Name != state.appTemplateName {
		t.Fatalf(fmt.Sprintf("List App Templates returned incorrect name: %s", response.ApplicationTemplates[0].Name))
	}
}

func TestGetAppTemplate(t *testing.T) {
	response, err := clt.GetApplicationTemplate(state.appTemplateName)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Get App Template failed: %s", err.Error()))
	}
	if response.Name != state.appTemplateName {
		t.Fatalf(fmt.Sprintf("Get App Template returned incorrect name: %s", response.Name))
	}
}

func TestCreateTemplatedApp(t *testing.T) {
	request := client.ApplicationCreateRequest{
		Name: state.appName,
		Template: &client.ApplicationTemplate{
			Name: state.appTemplateName,
			Variables: []client.TemplateVariable{
				{
					Key:   "agent-name",
					Value: state.agent,
				},
			},
		},
	}
	_, err := clt.CreateApplication(&request)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Create Templated Application failed: %s", err.Error()))
	}
}

func TestDeleteAppTemplate(t *testing.T) {
	if err := clt.DeleteApplicationTemplate(state.appTemplateName); err != nil {
		t.Fatalf("Failed to delete App Template: %s", err.Error())
	}
}

func TestDeleteAgent(t *testing.T) {
	if err := clt.DeleteAgent(state.uuid); err != nil {
		t.Fatalf("Failed to delete Agent: %s", err.Error())
	}
}
