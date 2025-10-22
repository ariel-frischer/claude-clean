package main

// Message types from Claude stream
type StreamMessage struct {
	Type              string          `json:"type"`
	Subtype           string          `json:"subtype,omitempty"`
	Message           *MessageContent `json:"message,omitempty"`
	SessionID         string          `json:"session_id,omitempty"`
	ParentToolUseID   string          `json:"parent_tool_use_id,omitempty"`
	CWD               string          `json:"cwd,omitempty"`
	Tools             []string        `json:"tools,omitempty"`
	Model             string          `json:"model,omitempty"`
	ClaudeCodeVersion string          `json:"claude_code_version,omitempty"`
	// Result message fields
	IsError           bool                   `json:"is_error,omitempty"`
	DurationMS        int                    `json:"duration_ms,omitempty"`
	DurationAPIMS     int                    `json:"duration_api_ms,omitempty"`
	NumTurns          int                    `json:"num_turns,omitempty"`
	Result            string                 `json:"result,omitempty"`
	TotalCostUSD      float64                `json:"total_cost_usd,omitempty"`
	Usage             *Usage                 `json:"usage,omitempty"`
	ModelUsage        map[string]interface{} `json:"modelUsage,omitempty"`
	PermissionDenials []interface{}          `json:"permission_denials,omitempty"`
}

type MessageContent struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         string         `json:"role"`
	Model        string         `json:"model"`
	Content      []ContentBlock `json:"content"`
	StopReason   *string        `json:"stop_reason"`
	StopSequence *string        `json:"stop_sequence"`
	Usage        *Usage         `json:"usage"`
}

type ContentBlock struct {
	Type      string                 `json:"type"`
	Text      string                 `json:"text,omitempty"`
	ID        string                 `json:"id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	Input     map[string]interface{} `json:"input,omitempty"`
	ToolUseID string                 `json:"tool_use_id,omitempty"`
	Content   interface{}            `json:"content,omitempty"`
	IsError   bool                   `json:"is_error,omitempty"`
}

type Usage struct {
	InputTokens              int                  `json:"input_tokens"`
	OutputTokens             int                  `json:"output_tokens"`
	CacheCreationInputTokens int                  `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int                  `json:"cache_read_input_tokens,omitempty"`
	CacheCreation            *CacheCreationDetail `json:"cache_creation,omitempty"`
	ServiceTier              string               `json:"service_tier,omitempty"`
}

type CacheCreationDetail struct {
	Ephemeral5mInputTokens int `json:"ephemeral_5m_input_tokens"`
	Ephemeral1hInputTokens int `json:"ephemeral_1h_input_tokens"`
}
