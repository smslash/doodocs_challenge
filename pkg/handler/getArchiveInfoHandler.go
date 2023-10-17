package handler

import (
	"archive/zip"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/smslash/doodocs_challenge/pkg/entity"
)

func GetArchiveInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Invalid request method"}`, http.StatusMethodNotAllowed)
		return
	}

	archiveFile, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error": "Failed to get the file from the request"}`, http.StatusBadRequest)
		return
	}
	defer archiveFile.Close()

	zipReader, err := zip.NewReader(archiveFile, r.ContentLength)
	if err != nil {
		http.Error(w, `{"error": "Not ZIP file"}`, http.StatusInternalServerError)
		return
	}

	var totalSize float64
	var files []entity.FileInfo

	for _, zipFile := range zipReader.File {
		totalSize += float64(zipFile.UncompressedSize64)

		mimetypeFull := http.DetectContentType([]byte(zipFile.Name))
		mimetypeParts := strings.Split(mimetypeFull, ";")
		mimetype := mimetypeParts[0]

		files = append(files, entity.FileInfo{
			FilePath: zipFile.Name,
			Size:     float64(zipFile.UncompressedSize64),
			Mimetype: mimetype,
		})
	}

	archiveInfo := entity.ArchiveInfo{
		Filename:    header.Filename,
		ArchiveSize: float64(r.ContentLength),
		TotalSize:   totalSize,
		TotalFiles:  len(files),
		Files:       files,
	}

	w.Header().Set("Content-Type", "application/json")
	formattedJSON, err := json.MarshalIndent(archiveInfo, "", "    ")
	if err != nil {
		http.Error(w, `{"error": "Failed to format JSON"}`, http.StatusInternalServerError)
		return
	}

	w.Write(formattedJSON)
}
