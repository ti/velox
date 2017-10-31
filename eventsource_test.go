package velox

import (
	"bytes"
	"testing"
)

func TestEncoderEncode(t *testing.T) {
	table := []struct {
		event
		expected string
	}{
		{event{Type: "type"}, "event: type\ndata\n\n"},
		{event{ID: "123"}, "id: 123\ndata\n\n"},
		{event{Retry: "10000"}, "retry: 10000\ndata\n\n"},
		{event{Data: []byte("data")}, "data: data\n\n"},
	}
	for i, tt := range table {
		b := bytes.Buffer{}
		if err := writeEvent(&b, tt.event); err != nil {
			t.Fatalf("%d. write error: %q", i, err)
		}
		if b.String() != tt.expected {
			t.Errorf("%d. expected %q, got %q", i, tt.expected, b.String())
		}
	}
}
