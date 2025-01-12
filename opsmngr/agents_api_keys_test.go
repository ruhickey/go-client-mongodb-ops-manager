// Copyright 2020 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package opsmngr

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
)

const projectID = "5e66185d917b220fbd8bb4d1"

func TestAgentsServiceOp_ListAgentAPIKeys(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	if _, _, err := client.Agents.ListAgentAPIKeys(ctx, ""); err == nil {
		t.Error("expected an error but got nil")
	}

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/agentapikeys", projectID), func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `
						[{
						  "_id" : "1",
						  "createdBy" : "PUBLIC_API",
						  "createdIpAddr" : "1",
						  "createdTime" : 1520458807291,
						  "createdUserId" : "21",
						  "desc" : "Agent API Key for this project",
						  "key" : "****************************8b87"
						}, {
						  "_id" : "2",
						  "createdBy" : "PROVISIONING",
						  "createdTime" : 1508871142864,
						  "desc" : "Generated by Provisioning",
						  "key" : "****************************39fe"
						}, {
						  "_id" : "3",
						  "createdBy" : "USER",
						  "createdIpAddr" : "1",
						  "createdTime" : 1507067499083,
						  "createdUserId" : "21",
						  "desc" : "Initial API Key",
						  "key" : "****************************70d7"
						}]
		`)
	})

	agentAPIKeys, _, err := client.Agents.ListAgentAPIKeys(ctx, projectID)
	if err != nil {
		t.Fatalf("Agents.ListAgentAPIKeys returned error: %v", err)
	}

	CreatedUserID := "21"
	CreatedIPAddr := "1"

	expected := []*AgentAPIKey{
		{
			ID:            "1",
			Key:           "****************************8b87",
			Desc:          "Agent API Key for this project",
			CreatedTime:   1520458807291,
			CreatedUserID: &CreatedUserID,
			CreatedIPAddr: &CreatedIPAddr,
			CreatedBy:     "PUBLIC_API",
		},
		{
			ID:          "2",
			Key:         "****************************39fe",
			Desc:        "Generated by Provisioning",
			CreatedTime: 1508871142864,
			CreatedBy:   "PROVISIONING",
		},
		{
			ID:            "3",
			Key:           "****************************70d7",
			Desc:          "Initial API Key",
			CreatedTime:   1507067499083,
			CreatedUserID: &CreatedUserID,
			CreatedIPAddr: &CreatedIPAddr,
			CreatedBy:     "USER",
		},
	}

	if diff := deep.Equal(agentAPIKeys, expected); diff != nil {
		t.Error(diff)
	}
}

func TestAgentsServiceOp_CreateAgentAPIKey(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	if _, _, err := client.Agents.CreateAgentAPIKey(ctx, "", &AgentAPIKeysRequest{}); err == nil {
		t.Error("expected an error but got nil")
	}

	if _, _, err := client.Agents.CreateAgentAPIKey(ctx, projectID, nil); err == nil {
		t.Error("expected an error but got nil")
	}

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/agentapikeys", projectID), func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{
						  "_id" : "1",
						  "createdBy" : "PUBLIC_API",
						  "createdIpAddr" : "1",
						  "createdTime" : 1520458807291,
						  "createdUserId" : "21",
						  "desc" : "TEST",
						  "key" : "****************************8b87"
						}`)
	})

	agentRequest := &AgentAPIKeysRequest{Desc: "TEST"}
	agentAPIKey, _, err := client.Agents.CreateAgentAPIKey(ctx, projectID, agentRequest)

	if err != nil {
		t.Fatalf("Agents.CreateAgentAPIKey returned error: %v", err)
	}

	CreatedUserID := "21"
	CreatedIPAddr := "1"

	expected := &AgentAPIKey{
		ID:            "1",
		Key:           "****************************8b87",
		Desc:          "TEST",
		CreatedTime:   1520458807291,
		CreatedUserID: &CreatedUserID,
		CreatedIPAddr: &CreatedIPAddr,
		CreatedBy:     "PUBLIC_API",
	}

	if diff := deep.Equal(agentAPIKey, expected); diff != nil {
		t.Error(diff)
	}
}

func TestAgentsServiceOp_DeleteAgentAPIKey(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	agentAPIKey := "1"

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/agentapikeys/%s", projectID, agentAPIKey), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	type args struct {
		projectID   string
		agentAPIKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful",
			args: args{
				projectID:   projectID,
				agentAPIKey: agentAPIKey,
			},
			wantErr: false,
		},
		{
			name: "missing projectID",
			args: args{
				projectID:   "",
				agentAPIKey: agentAPIKey,
			},
			wantErr: true,
		},
		{
			name: "missing agentAPIKey",
			args: args{
				projectID:   projectID,
				agentAPIKey: "",
			},
			wantErr: true,
		},
		{
			name: "missing projectID and agentAPIKey",
			args: args{
				projectID:   "",
				agentAPIKey: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		wantErr := tt.wantErr
		args := tt.args
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Agents.DeleteAgentAPIKey(ctx, args.projectID, args.agentAPIKey)
			if (err != nil) != wantErr {
				t.Errorf("DeleteAgentAPIKey() error = %v, wantErr %v", err, wantErr)
				return
			}
		})
	}
}
