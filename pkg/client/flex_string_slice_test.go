package client

import (
	"encoding/json"
	"testing"
)

func TestFlexStringSliceUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want FlexStringSlice
	}{
		{name: "null", raw: `null`, want: nil},
		{name: "array", raw: `["crun","runc"]`, want: FlexStringSlice{"crun", "runc"}},
		{name: "stringified array", raw: `"[\"crun\"]"`, want: FlexStringSlice{"crun"}},
		{name: "plain string fallback", raw: `"crun"`, want: FlexStringSlice{"crun"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var got FlexStringSlice
			if err := json.Unmarshal([]byte(tt.raw), &got); err != nil {
				t.Fatalf("UnmarshalJSON: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestAgentInfoAvailableRuntimesStringifiedArray(t *testing.T) {
	t.Parallel()

	raw := `{"name":"edge-1","availableRuntimes":"[\"crun\"]"}`

	var agent AgentInfo
	if err := json.Unmarshal([]byte(raw), &agent); err != nil {
		t.Fatalf("Unmarshal AgentInfo: %v", err)
	}
	if len(agent.AvailableRuntimes) != 1 || agent.AvailableRuntimes[0] != "crun" {
		t.Fatalf("AvailableRuntimes = %v", agent.AvailableRuntimes)
	}
}
