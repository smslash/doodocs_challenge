package tests

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"

	"github.com/smslash/doodocs_challenge/pkg/handler"
)

func TestCreateArchiveHandler(t *testing.T) {
	t.Run("Method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/archive", nil)
		rr := httptest.NewRecorder()

		handler.CreateArchiveHandler(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
		}
	})

	t.Run("Unsupported file type", func(t *testing.T) {
		var b bytes.Buffer
		writer := multipart.NewWriter(&b)
		writer.WriteField("files[]", "test file content")
		h := make(textproto.MIMEHeader)

		h.Set("Content-Disposition", `form-data; name="files[]"; filename="test.txt"`)
		h.Set("Content-Type", "text/plain")

		part, _ := writer.CreatePart(h)
		part.Write([]byte("test file content"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "http://example.com/archive", &b)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()

		handler.CreateArchiveHandler(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		if !strings.Contains(rr.Body.String(), "Unsupported file type") {
			t.Error("Expected 'Unsupported file type' in response")
		}
	})

	t.Run("Successfully create archive", func(t *testing.T) {
		var b bytes.Buffer
		writer := multipart.NewWriter(&b)
		h := make(textproto.MIMEHeader)

		h.Set("Content-Disposition", `form-data; name="files[]"; filename="test.jpeg"`)
		h.Set("Content-Type", "image/jpeg")

		part, _ := writer.CreatePart(h)
		part.Write([]byte("dummy jpeg content"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "http://example.com/archive", &b)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()

		handler.CreateArchiveHandler(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		if rr.Header().Get("Content-Type") != "application/zip" {
			t.Errorf("handler returned wrong content type: got %v want %v", rr.Header().Get("Content-Type"), "application/zip")
		}
	})
}
