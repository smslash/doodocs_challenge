package handler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func CreateArchiveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	allowedMimeTypes := map[string]bool{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/xml": true,
		"image/jpeg":      true,
		"image/png":       true,
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	files, err := r.MultipartReader()
	if err != nil {
		http.Error(w, `{"error": "Failed to read multipart data"}`, http.StatusInternalServerError)
		return
	}

	for {
		part, err := files.NextPart()
		if err == io.EOF {
			break
		}

		if part.FormName() == "files[]" {
			fmt.Println("Received MIME-type:", part.Header.Get("Content-Type"))

			if !allowedMimeTypes[part.Header.Get("Content-Type")] {
				http.Error(w, `{"error": "Unsupported file type"}`, http.StatusBadRequest)
				return
			}

			fileWriter, err := zipWriter.Create(part.FileName())
			if err != nil {
				http.Error(w, `{"error": "Failed to add file to archive"}`, http.StatusInternalServerError)
				return
			}
			io.Copy(fileWriter, part)
		}
	}

	zipWriter.Close()

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="archive.zip"`)
	w.Write(buf.Bytes())
}
