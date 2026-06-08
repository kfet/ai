package ai

import "testing"

func TestIsRateLimitText(t *testing.T) {
	for _, s := range []string{"429 Too Many Requests", "rate limit exceeded", "resource exhausted", "overloaded"} {
		if !IsRateLimitText(s) {
			t.Errorf("IsRateLimitText(%q) = false, want true", s)
		}
	}
	if IsRateLimitText("everything is fine") {
		t.Error("IsRateLimitText benign = true")
	}
}

func TestIsTransientServerError(t *testing.T) {
	for _, s := range []string{"Internal server error", "502 Bad Gateway", "503", "504"} {
		if !IsTransientServerError(s) {
			t.Errorf("IsTransientServerError(%q) = false, want true", s)
		}
	}
	if IsTransientServerError("ok") {
		t.Error("IsTransientServerError benign = true")
	}
}

func TestIsTransientNetworkError(t *testing.T) {
	for _, s := range []string{"connection reset by peer", "broken pipe", "tls handshake timeout", "unexpected EOF", "http2: server sent GOAWAY", "EOF"} {
		if !IsTransientNetworkError(s) {
			t.Errorf("IsTransientNetworkError(%q) = false, want true", s)
		}
	}
	if IsTransientNetworkError("clean shutdown") {
		t.Error("IsTransientNetworkError benign = true")
	}
}

func TestIsRetryableError(t *testing.T) {
	// One representative of each contributing predicate, plus a negative.
	for _, s := range []string{"429", "internal server error", "connection refused"} {
		if !IsRetryableError(s) {
			t.Errorf("IsRetryableError(%q) = false, want true", s)
		}
	}
	if IsRetryableError("benign message") {
		t.Error("IsRetryableError benign = true")
	}
}
