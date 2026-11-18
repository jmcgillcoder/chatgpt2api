package upstream

import (
	"reflect"
	"testing"
)

func TestPreparePayloadImageUploadIDBatchLimit(t *testing.T) {
	for _, count := range []int{0, 1, 2, 3, 4, 10} {
		t.Run(string(rune('0'+count)), func(t *testing.T) {
			ids := make([]string, count)
			for i := range ids {
				ids[i] = "upload-" + string(rune('a'+i))
			}

			payload := preparePayload(
				"gpt-image-2.6",
				[]string{"picture_v2"},
				[]string{"image/png"},
				ids,
			)

			got, exists := payload["prepare_image_upload_ids"]
			if count > 0 && count <= 2 {
				if !exists {
					t.Fatal("expected prepare_image_upload_ids for batches of one or two")
				}
				if !reflect.DeepEqual(got, ids) {
					t.Fatalf("upload IDs changed: got %#v, want %#v", got, ids)
				}
			} else if exists {
				t.Fatalf("prepare_image_upload_ids must be omitted for batch size %d, got %#v", count, got)
			}

			if !reflect.DeepEqual(payload["attachment_mime_types"], []string{"image/png"}) {
				t.Fatalf("attachment MIME metadata was not preserved: %#v", payload["attachment_mime_types"])
			}
		})
	}
}

func TestConversationAttachmentsKeepsAllImageAttachments(t *testing.T) {
	attachments := make([]any, 10)
	wantIDs := make([]string, 10)
	for i := range attachments {
		id := "upload-" + string(rune('a'+i))
		wantIDs[i] = id
		attachments[i] = map[string]any{
			"id":        id,
			"mime_type": "image/png",
		}
	}

	mimeTypes, uploadIDs := conversationAttachments([]map[string]any{
		{
			"metadata": map[string]any{
				"attachments": attachments,
			},
		},
	})

	if !reflect.DeepEqual(uploadIDs, wantIDs) {
		t.Fatalf("conversation attachments were truncated or reordered: got %#v, want %#v", uploadIDs, wantIDs)
	}
	if !reflect.DeepEqual(mimeTypes, []string{"image/png"}) {
		t.Fatalf("unexpected MIME types: %#v", mimeTypes)
	}
}
