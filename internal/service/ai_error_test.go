package service

import (
	"errors"
	"testing"

	chatgpt "github.com/basketikun/chatgpt2api/internal/upstream"
)

func TestNewImageErrorClassifiesPrepareRejections(t *testing.T) {
	for _, status := range []int{400, 422} {
		t.Run(string(rune('0'+status/100))+string(rune('0'+(status%100)/10))+string(rune('0'+status%10)), func(t *testing.T) {
			err := &chatgpt.UpstreamHTTPError{
				Context:    "/backend-api/f/conversation/prepare",
				StatusCode: status,
				Body:       "{\"detail\":\"Invalid conversation body\"}",
			}
			got := newImageError(err)
			if got.Status != 422 || got.Type != "invalid_request_error" || got.Code != "upstream_prepare_rejected" {
				t.Fatalf("unexpected prepare error classification: %#v", got)
			}
		})
	}
}

func TestNewImageErrorDoesNotMisclassifyOtherUpstreamErrors(t *testing.T) {
	err := &chatgpt.UpstreamHTTPError{
		Context:    "/backend-api/f/conversation",
		StatusCode: 422,
		Body:       "unprocessable",
	}
	got := newImageError(err)
	if got.Code == "upstream_prepare_rejected" {
		t.Fatalf("non-prepare error was misclassified: %#v", got)
	}
	if !errors.Is(got.Err, err) {
		t.Fatal("original upstream error was not preserved")
	}
}
