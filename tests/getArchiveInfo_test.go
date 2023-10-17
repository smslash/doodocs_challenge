package tests

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smslash/doodocs_challenge/pkg/entity"
	"github.com/smslash/doodocs_challenge/pkg/handler"
	"github.com/stretchr/testify/assert"
)

func TestGetArchiveInfoHandler(t *testing.T) {
	t.Run("Invalid request method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/archiveinfo", nil)
		rr := httptest.NewRecorder()

		handler.GetArchiveInfoHandler(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid request method")
	})

	t.Run("Successfully get archive info", func(t *testing.T) {
		var b bytes.Buffer
		writer := multipart.NewWriter(&b)
		part, _ := writer.CreateFormFile("file", "test.zip")

		zipBuffer := new(bytes.Buffer)
		zipWriter := zip.NewWriter(zipBuffer)
		fileWriter, _ := zipWriter.Create("test.txt")

		fileWriter.Write([]byte("test content"))
		zipWriter.Close()

		part.Write(zipBuffer.Bytes())
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "http://example.com/archiveinfo", &b)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rr := httptest.NewRecorder()

		handler.GetArchiveInfoHandler(rr, req)
		t.Logf("Response body: %s", rr.Body.String())
		assert.Equal(t, http.StatusOK, rr.Code)

		var response entity.ArchiveInfo
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, "test.zip", response.Filename)
		assert.Equal(t, 1, response.TotalFiles)
		assert.Equal(t, "test.txt", response.Files[0].FilePath)
	})
}
