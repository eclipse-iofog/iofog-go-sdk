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
	"net/url"
	"testing"
)

func TestCreation(t *testing.T) {
	baseURL, err := url.Parse("http://localhost:51121/api/v3")
	if err != nil {
		t.Error(err)
	}
	client := New(Options{
		BaseURL: baseURL,
	})
	if client == nil {
		t.Error("Client pointer is nil")
	}
}

func TestGenerateListAgentsURL(t *testing.T) {
	request := ListAgentsRequest{
		System: true,
		Filters: []AgentListFilter{
			{
				Key:       "first",
				Value:     "second",
				Condition: "third",
			},
		},
	}
	url := generateListAgentURL(request)
	if url != "/iofog-list?system=true&filters[0][key]=first&filters[0][value]=second&filters[0][condition]=third" {
		t.Errorf("Failed to generate List Agents URL: %s", url)
	}
	url = generateListAgentURL(ListAgentsRequest{})
	if url != "/iofog-list?system=false" {
		t.Errorf("Failed to generate List Agents URL: %s", url)
	}
}
