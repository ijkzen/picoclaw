package providers

import (
	"errors"
	"testing"
)

func TestIsMaxTokensOutOfRangeError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "max_tokens must be between",
			err:  errors.New("400: max_tokens must be between 1 and 8192"),
			want: true,
		},
		{
			name: "max_completion_tokens param json",
			err: errors.New(
				`{"error":{"message":"invalid value","param":"max_completion_tokens"}}`,
			),
			want: true,
		},
		{
			name: "max_output_tokens out of range",
			err:  errors.New("max_output_tokens out of range"),
			want: true,
		},
		{
			name: "context overflow should not match",
			err:  errors.New("context length exceeded maximum tokens"),
			want: false,
		},
		{
			name: "param present without range hint should not match",
			err:  errors.New("unknown field max_tokens in request"),
			want: false,
		},
		{
			name: "range hint without param should not match",
			err:  errors.New("value must be between 1 and 8192"),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsMaxTokensOutOfRangeError(tc.err)
			if got != tc.want {
				t.Fatalf("IsMaxTokensOutOfRangeError() = %v, want %v", got, tc.want)
			}
		})
	}
}
