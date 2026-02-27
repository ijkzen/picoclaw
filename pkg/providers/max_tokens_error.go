package providers

import "strings"

// IsMaxTokensOutOfRangeError returns true when an error clearly indicates
// the max_tokens parameter itself is out of range.
//
// It intentionally requires both:
//  1. a max-tokens parameter hint (max_tokens/max_completion_tokens/max_output_tokens)
//  2. a range-validation hint (out of range/must be/between/exceed/...)
//
// This avoids misclassifying generic context-window overflows as parameter-range errors.
func IsMaxTokensOutOfRangeError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	compact := strings.NewReplacer(" ", "", "\n", "", "\t", "", "\r", "").Replace(msg)

	hasParamHint := strings.Contains(msg, "max_tokens") ||
		strings.Contains(msg, "max_completion_tokens") ||
		strings.Contains(msg, "max_output_tokens") ||
		strings.Contains(msg, "max completion tokens") ||
		strings.Contains(msg, "max output tokens") ||
		strings.Contains(compact, `"param":"max_tokens"`) ||
		strings.Contains(compact, `"param":"max_completion_tokens"`) ||
		strings.Contains(compact, `"param":"max_output_tokens"`)
	if !hasParamHint {
		return false
	}

	hasRangeHint := strings.Contains(msg, "out of range") ||
		strings.Contains(msg, "must be") ||
		strings.Contains(msg, "between") ||
		strings.Contains(msg, "less than or equal") ||
		strings.Contains(msg, "greater than or equal") ||
		strings.Contains(msg, "too large") ||
		strings.Contains(msg, "cannot exceed") ||
		strings.Contains(msg, "exceed") ||
		strings.Contains(msg, "invalid")

	return hasRangeHint
}
