package main

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
)

// TestParseSystemMessage tests parsing of a system message from the mock data
func TestParseSystemMessage(t *testing.T) {
	jsonLine := `{"type":"system","subtype":"init","cwd":"/home/ari/repos/claude-code-clean-output","session_id":"test-session","model":"claude-sonnet-4-5-20250929","claude_code_version":"2.0.25"}`

	var msg StreamMessage
	err := json.Unmarshal([]byte(jsonLine), &msg)
	if err != nil {
		t.Fatalf("Failed to unmarshal system message: %v", err)
	}

	if msg.Type != "system" {
		t.Errorf("Expected type 'system', got '%s'", msg.Type)
	}
	if msg.Subtype != "init" {
		t.Errorf("Expected subtype 'init', got '%s'", msg.Subtype)
	}
	if msg.Model != "claude-sonnet-4-5-20250929" {
		t.Errorf("Expected model 'claude-sonnet-4-5-20250929', got '%s'", msg.Model)
	}
	if msg.ClaudeCodeVersion != "2.0.25" {
		t.Errorf("Expected version '2.0.25', got '%s'", msg.ClaudeCodeVersion)
	}
}

// TestParseAssistantMessage tests parsing of an assistant message
func TestParseAssistantMessage(t *testing.T) {
	jsonLine := `{"type":"assistant","message":{"model":"claude-sonnet-4-5-20250929","id":"msg_test","type":"message","role":"assistant","content":[{"type":"text","text":"Hello, world!"}],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":100,"output_tokens":20}},"session_id":"test-session"}`

	var msg StreamMessage
	err := json.Unmarshal([]byte(jsonLine), &msg)
	if err != nil {
		t.Fatalf("Failed to unmarshal assistant message: %v", err)
	}

	if msg.Type != "assistant" {
		t.Errorf("Expected type 'assistant', got '%s'", msg.Type)
	}
	if msg.Message == nil {
		t.Fatal("Message content is nil")
	}
	if msg.Message.Role != "assistant" {
		t.Errorf("Expected role 'assistant', got '%s'", msg.Message.Role)
	}
	if len(msg.Message.Content) != 1 {
		t.Fatalf("Expected 1 content block, got %d", len(msg.Message.Content))
	}
	if msg.Message.Content[0].Type != "text" {
		t.Errorf("Expected content type 'text', got '%s'", msg.Message.Content[0].Type)
	}
	if msg.Message.Content[0].Text != "Hello, world!" {
		t.Errorf("Expected text 'Hello, world!', got '%s'", msg.Message.Content[0].Text)
	}
}

// TestParseToolUse tests parsing of a tool_use content block
func TestParseToolUse(t *testing.T) {
	jsonLine := `{"type":"assistant","message":{"model":"claude-sonnet-4-5-20250929","id":"msg_test","type":"message","role":"assistant","content":[{"type":"tool_use","id":"toolu_test","name":"Read","input":{"file_path":"/test/path.go"}}],"usage":{"input_tokens":50,"output_tokens":10}},"session_id":"test-session"}`

	var msg StreamMessage
	err := json.Unmarshal([]byte(jsonLine), &msg)
	if err != nil {
		t.Fatalf("Failed to unmarshal tool_use message: %v", err)
	}

	if len(msg.Message.Content) != 1 {
		t.Fatalf("Expected 1 content block, got %d", len(msg.Message.Content))
	}

	block := msg.Message.Content[0]
	if block.Type != "tool_use" {
		t.Errorf("Expected type 'tool_use', got '%s'", block.Type)
	}
	if block.Name != "Read" {
		t.Errorf("Expected name 'Read', got '%s'", block.Name)
	}
	if block.ID != "toolu_test" {
		t.Errorf("Expected ID 'toolu_test', got '%s'", block.ID)
	}
	if block.Input == nil {
		t.Fatal("Input is nil")
	}
	if filePath, ok := block.Input["file_path"].(string); !ok || filePath != "/test/path.go" {
		t.Errorf("Expected file_path '/test/path.go', got '%v'", block.Input["file_path"])
	}
}

// TestParseResultMessage tests parsing of a result message
func TestParseResultMessage(t *testing.T) {
	jsonLine := `{"type":"result","is_error":false,"duration_ms":5000,"num_turns":3,"result":"Task completed successfully","total_cost_usd":0.0012,"usage":{"input_tokens":1000,"output_tokens":200}}`

	var msg StreamMessage
	err := json.Unmarshal([]byte(jsonLine), &msg)
	if err != nil {
		t.Fatalf("Failed to unmarshal result message: %v", err)
	}

	if msg.Type != "result" {
		t.Errorf("Expected type 'result', got '%s'", msg.Type)
	}
	if msg.IsError {
		t.Error("Expected is_error to be false")
	}
	if msg.DurationMS != 5000 {
		t.Errorf("Expected duration 5000ms, got %d", msg.DurationMS)
	}
	if msg.NumTurns != 3 {
		t.Errorf("Expected 3 turns, got %d", msg.NumTurns)
	}
	if msg.Result != "Task completed successfully" {
		t.Errorf("Expected result text, got '%s'", msg.Result)
	}
	if msg.TotalCostUSD != 0.0012 {
		t.Errorf("Expected cost 0.0012, got %f", msg.TotalCostUSD)
	}
}

// TestParseUsageWithCache tests parsing of usage statistics with cache information
func TestParseUsageWithCache(t *testing.T) {
	jsonLine := `{"type":"assistant","message":{"model":"test","id":"msg_test","type":"message","role":"assistant","content":[{"type":"text","text":"test"}],"usage":{"input_tokens":100,"output_tokens":20,"cache_creation_input_tokens":50,"cache_read_input_tokens":30,"cache_creation":{"ephemeral_5m_input_tokens":50,"ephemeral_1h_input_tokens":0},"service_tier":"standard"}},"session_id":"test"}`

	var msg StreamMessage
	err := json.Unmarshal([]byte(jsonLine), &msg)
	if err != nil {
		t.Fatalf("Failed to unmarshal usage with cache: %v", err)
	}

	usage := msg.Message.Usage
	if usage == nil {
		t.Fatal("Usage is nil")
	}
	if usage.InputTokens != 100 {
		t.Errorf("Expected input tokens 100, got %d", usage.InputTokens)
	}
	if usage.OutputTokens != 20 {
		t.Errorf("Expected output tokens 20, got %d", usage.OutputTokens)
	}
	if usage.CacheCreationInputTokens != 50 {
		t.Errorf("Expected cache creation tokens 50, got %d", usage.CacheCreationInputTokens)
	}
	if usage.CacheReadInputTokens != 30 {
		t.Errorf("Expected cache read tokens 30, got %d", usage.CacheReadInputTokens)
	}
	if usage.ServiceTier != "standard" {
		t.Errorf("Expected service tier 'standard', got '%s'", usage.ServiceTier)
	}
	if usage.CacheCreation == nil {
		t.Fatal("CacheCreation is nil")
	}
	if usage.CacheCreation.Ephemeral5mInputTokens != 50 {
		t.Errorf("Expected ephemeral 5m tokens 50, got %d", usage.CacheCreation.Ephemeral5mInputTokens)
	}
}

// TestParseMockFile tests parsing the actual mock file
func TestParseMockFile(t *testing.T) {
	file, err := os.Open("../mocks/claude-stream-json-simple.jsonl")
	if err != nil {
		t.Fatalf("Failed to open mock file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	messageTypes := make(map[string]int)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var msg StreamMessage
		err := json.Unmarshal([]byte(line), &msg)
		if err != nil {
			t.Errorf("Failed to parse line %d: %v", lineCount+1, err)
			continue
		}

		lineCount++
		messageTypes[msg.Type]++
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading mock file: %v", err)
	}

	if lineCount == 0 {
		t.Error("No lines parsed from mock file")
	}

	t.Logf("Parsed %d lines from mock file", lineCount)
	t.Logf("Message types: %v", messageTypes)

	// Verify we have at least some expected message types
	if messageTypes["system"] == 0 {
		t.Error("Expected at least one system message")
	}
	if messageTypes["assistant"] == 0 {
		t.Error("Expected at least one assistant message")
	}
}

// TestParseToolResult tests parsing of tool_result in user messages
func TestParseToolResult(t *testing.T) {
	jsonLine := `{"type":"user","message":{"role":"user","content":[{"tool_use_id":"toolu_test","type":"tool_result","content":"File contents here","is_error":false}]}}`

	var msg StreamMessage
	err := json.Unmarshal([]byte(jsonLine), &msg)
	if err != nil {
		t.Fatalf("Failed to unmarshal tool_result: %v", err)
	}

	if msg.Type != "user" {
		t.Errorf("Expected type 'user', got '%s'", msg.Type)
	}
	if len(msg.Message.Content) != 1 {
		t.Fatalf("Expected 1 content block, got %d", len(msg.Message.Content))
	}

	block := msg.Message.Content[0]
	if block.Type != "tool_result" {
		t.Errorf("Expected type 'tool_result', got '%s'", block.Type)
	}
	if block.ToolUseID != "toolu_test" {
		t.Errorf("Expected tool_use_id 'toolu_test', got '%s'", block.ToolUseID)
	}
	if block.IsError {
		t.Error("Expected is_error to be false")
	}
}
