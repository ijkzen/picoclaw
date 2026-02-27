//go:build amd64 || arm64 || riscv64 || mips64 || ppc64

package channels

import (
	"context"
	"testing"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func TestShouldReactToFeishuMessage(t *testing.T) {
	tests := []struct {
		name        string
		messageID   string
		messageType string
		want        bool
	}{
		{name: "normal text", messageID: "om_1", messageType: "text", want: true},
		{name: "empty message id", messageID: "", messageType: "text", want: false},
		{name: "system message", messageID: "om_1", messageType: "system", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldReactToFeishuMessage(tt.messageID, tt.messageType)
			if got != tt.want {
				t.Fatalf("shouldReactToFeishuMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReactToMessageSkipPaths(t *testing.T) {
	ch := &FeishuChannel{}

	if err := ch.reactToMessage(context.Background(), "", "text"); err != nil {
		t.Fatalf("reactToMessage() empty id returned error: %v", err)
	}
	if err := ch.reactToMessage(context.Background(), "om_1", "system"); err != nil {
		t.Fatalf("reactToMessage() system type returned error: %v", err)
	}
}

func TestExtractFeishuSenderID(t *testing.T) {
	sender := &larkim.EventSender{
		SenderId: &larkim.UserId{
			OpenId: strPtr("ou_123"),
		},
	}

	got := extractFeishuSenderID(sender)
	if got != "ou_123" {
		t.Fatalf("extractFeishuSenderID() = %q, want %q", got, "ou_123")
	}
}

func strPtr(v string) *string {
	return &v
}
