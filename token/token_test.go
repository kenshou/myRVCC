package token

import "testing"

func TestTokenKind_String(t *testing.T) {
	kind := keywords["return"]
	t.Log(kind.String())
}
