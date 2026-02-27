package tools

import "testing"

func TestShouldSuppressScheduledCommandErrorForFeishu(t *testing.T) {
	tests := []struct {
		name    string
		channel string
		result  *ToolResult
		want    bool
	}{
		{
			name:    "suppress feishu safety guard outside working dir error",
			channel: "feishu",
			result: &ToolResult{
				IsError: true,
				ForLLM:  "Command blocked by safety guard (path outside working dir)",
			},
			want: true,
		},
		{
			name:    "do not suppress feishu other error",
			channel: "feishu",
			result: &ToolResult{
				IsError: true,
				ForLLM:  "timeout",
			},
			want: false,
		},
		{
			name:    "do not suppress non-feishu",
			channel: "telegram",
			result: &ToolResult{
				IsError: true,
				ForLLM:  "Command blocked by safety guard (path outside working dir)",
			},
			want: false,
		},
		{
			name:    "do not suppress success",
			channel: "feishu",
			result: &ToolResult{
				IsError: false,
				ForLLM:  "ok",
			},
			want: false,
		},
		{
			name:    "do not suppress nil result",
			channel: "feishu",
			result:  nil,
			want:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldSuppressScheduledCommandErrorForFeishu(tc.channel, tc.result)
			if got != tc.want {
				t.Fatalf("shouldSuppressScheduledCommandErrorForFeishu() = %v, want %v", got, tc.want)
			}
		})
	}
}
