package job

import "testing"

func TestCalculateBackoff(t *testing.T) {
	tests := []struct {
		retryCnt int
		expected int
	}{
		{1, 2},
		{2, 4},
		{3, 8},
	}

	for _, tt := range tests {
		got := calculateBackoff(tt.retryCnt)
		if got != tt.expected {
			t.Errorf("calculateBackoff(%d) = %d, want %d", tt.retryCnt, got, tt.expected)
		}
	}
}

func TestRetryDecision(t *testing.T) {
	tests := []struct {
		retryCount  int
		shouldRetry bool
	}{
		{0, true},
		{1, true},
		{MaxRetry - 1, false},
		{MaxRetry, false},
	}

	for _, tt := range tests {
		nextRetryCnt := tt.retryCount + 1
		shouldRetry := nextRetryCnt < MaxRetry
		if shouldRetry != tt.shouldRetry {
			t.Errorf("retryCount=%d: shouldRetry=%v, want %v", tt.retryCount, shouldRetry, tt.shouldRetry)
		}
	}
}
