package gmail

import "testing"

func TestBuildRawMessage(t *testing.T) {
	raw, err := BuildRawMessage("a@example.com", "b@example.com", "Hi", "hello", "")
	if err != nil {
		t.Fatal(err)
	}
	if raw == "" {
		t.Fatal("empty raw")
	}
	if _, err := BuildRawMessage("", "b@example.com", "Hi", "", ""); err == nil {
		t.Fatal("expected error when body missing")
	}
}
