package types

import (
	"encoding/json"
	"testing"
)

func TestToolGroup(t *testing.T) {
	t.Parallel()

	// Test struct creation
	group := ToolGroup{
		Name:          "test-group",
		IncludedTools: []string{"tool1", "tool2"},
	}

	if group.Name != "test-group" {
		t.Errorf("Expected Name to be 'test-group', got %s", group.Name)
	}
	if len(group.IncludedTools) != 2 {
		t.Errorf("Expected IncludedTools to have 2 items, got %d", len(group.IncludedTools))
	}
}

func TestToolGroupJSONMarshaling(t *testing.T) {
	t.Parallel()

	group := ToolGroup{
		Name:          "json-group",
		IncludedTools: []string{"json-tool1"},
		Description:   "Group for JSON testing",
	}

	data, err := json.Marshal(group)
	if err != nil {
		t.Fatalf("Failed to marshal ToolGroup: %v", err)
	}

	expected := `{"name":"json-group","included_tools":["json-tool1"],"description":"Group for JSON testing"}`
	if string(data) != expected {
		t.Errorf("Expected JSON %s, got %s", expected, string(data))
	}
}

func TestToolGroupJSONMarshalingWithNewFields(t *testing.T) {
	t.Parallel()

	group := ToolGroup{
		Name:            "advanced-group",
		IncludedTools:   []string{"manual-tool1"},
		IncludedServers: []string{"time", "deepwiki"},
		ExcludedTools:   []string{"time__convert_time"},
		Description:     "Group with server inclusion and exclusion",
	}

	data, err := json.Marshal(group)
	if err != nil {
		t.Fatalf("Failed to marshal ToolGroup: %v", err)
	}

	// Unmarshal to verify all fields are present
	var unmarshaled ToolGroup
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ToolGroup: %v", err)
	}

	if unmarshaled.Name != group.Name {
		t.Errorf("Expected Name %s, got %s", group.Name, unmarshaled.Name)
	}
	if len(unmarshaled.IncludedTools) != 1 || unmarshaled.IncludedTools[0] != "manual-tool1" {
		t.Errorf("Expected IncludedTools [manual-tool1], got %v", unmarshaled.IncludedTools)
	}
	if len(unmarshaled.IncludedServers) != 2 {
		t.Errorf("Expected 2 IncludedServers, got %v", unmarshaled.IncludedServers)
	}
	if len(unmarshaled.ExcludedTools) != 1 || unmarshaled.ExcludedTools[0] != "time__convert_time" {
		t.Errorf("Expected ExcludedTools [time__convert_time], got %v", unmarshaled.ExcludedTools)
	}
}

func TestCreateToolGroupResponse(t *testing.T) {
	t.Parallel()

	// Test struct creation
	response := CreateToolGroupResponse{
		ToolGroupEndpoints: &ToolGroupEndpoints{
			StreamableHTTPEndpoint: "/api/tool-groups/test-group/stream",
			SSEEndpoint:            "/api/tool-groups/test-group/sse",
			SSEMessageEndpoint:     "/api/tool-groups/test-group/sse/message",
		},
	}

	if response.StreamableHTTPEndpoint != "/api/tool-groups/test-group/stream" {
		t.Errorf("Expected StreamableHTTPEndpoint to be '/api/tool-groups/test-group/stream', got %s", response.StreamableHTTPEndpoint)
	}
}

func TestCreateToolGroupResponseJSONMarshaling(t *testing.T) {
	t.Parallel()

	response := CreateToolGroupResponse{
		ToolGroupEndpoints: &ToolGroupEndpoints{
			StreamableHTTPEndpoint: "/api/tool-groups/json-group/stream",
			SSEEndpoint:            "/api/tool-groups/json-group/sse",
			SSEMessageEndpoint:     "/api/tool-groups/json-group/sse/message",
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal CreateToolGroupResponse: %v", err)
	}

	expected := `{"streamable_http_endpoint":"/api/tool-groups/json-group/stream","sse_endpoint":"/api/tool-groups/json-group/sse","sse_message_endpoint":"/api/tool-groups/json-group/sse/message"}`
	if string(data) != expected {
		t.Errorf("Expected JSON %s, got %s", expected, string(data))
	}
}

func TestGetToolGroupResponse(t *testing.T) {
	t.Parallel()

	// Test struct creation
	toolGroup := &ToolGroup{
		Name:          "get-group",
		IncludedTools: []string{"get-tool1"},
		Description:   "Group for get testing",
	}

	response := GetToolGroupResponse{
		ToolGroup: toolGroup,
		ToolGroupEndpoints: &ToolGroupEndpoints{
			StreamableHTTPEndpoint: "/api/tool-groups/get-group/stream",
			SSEEndpoint:            "/api/tool-groups/get-group/sse",
			SSEMessageEndpoint:     "/api/tool-groups/get-group/sse/message",
		},
	}

	if response.ToolGroup != toolGroup {
		t.Error("Expected ToolGroup pointer to match")
	}
	if response.StreamableHTTPEndpoint != "/api/tool-groups/get-group/stream" {
		t.Errorf("Expected StreamableHTTPEndpoint to be '/api/tool-groups/get-group/stream', got %s", response.StreamableHTTPEndpoint)
	}
}
