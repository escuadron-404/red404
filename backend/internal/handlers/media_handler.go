package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/escuadron-404/red404/backend/internal/dto"
	"github.com/escuadron-404/red404/backend/internal/services"
	"github.com/go-playground/validator/v10"
)

const maxUploadBytes = 10 << 20 // 10 MB (10 * 1024 * 1024 bytes)

type MediaHandler struct {
	mediaService services.MediaService
	validator    *validator.Validate
}

func NewMediaUploadHandler(mediaService services.MediaService, mediaValidator *validator.Validate) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
		validator:    mediaValidator,
	}
}

func (m *MediaHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// onichan yamete
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)

	reader, err := r.MultipartReader()
	if err != nil {
		if err == http.ErrNotMultipart {
			slog.Warn("Received non-multipart request for file upload")
			http.Error(w, "Request must be multipart/form-data", http.StatusBadRequest)
		} else if err.Error() == "http: request body too large" {
			slog.Warn("Received request body too large", "max_bytes", maxUploadBytes)
			http.Error(w, fmt.Sprintf("Request body too large. Max allowed: %d bytes", maxUploadBytes), http.StatusRequestEntityTooLarge)
		} else {
			slog.Error("Failed to create multipart reader", "err", err)
			http.Error(w, "Internal server error processing request", http.StatusInternalServerError)
		}
		return
	}

	var mediaUploadResponse *dto.MediaUploadResponse
	foundFileUpload := false
	var actualContentType string

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("Failed to read next part of multipart form", "err", err)
			http.Error(w, "Internal server error processing upload", http.StatusInternalServerError)
			return
		}
		defer part.Close() // Ensure part is closed after processing

		if part.FormName() == "file_upload" && part.FileName() != "" {
			foundFileUpload = true
			originalFilename := part.FileName()
			initialContentType := part.Header.Get("Content-Type")

			slog.Info("Starting file upload stream", "filename", originalFilename, "initial_mime", initialContentType)

			dstWriter, mediaResp, err := m.mediaService.PrepareFileUpload(r.Context(), originalFilename, initialContentType)
			if err != nil {
				slog.Error("Error preparing file upload destination", "err", err, "filename", originalFilename)
				http.Error(w, fmt.Sprintf("Error preparing file: %s", err.Error()), http.StatusInternalServerError)
				return
			}
			defer dstWriter.Close()

			// sniff em real good boy
			buffer := make([]byte, 512)
			n, readErr := part.Read(buffer)
			if readErr != nil && readErr != io.EOF {
				slog.Error("Error reading initial bytes for content type sniffing", "err", readErr, "filename", originalFilename)
				http.Error(w, "Error reading file content", http.StatusInternalServerError)
				return
			}
			actualContentType = http.DetectContentType(buffer[:n])

			if !strings.HasPrefix(actualContentType, "image/") && !strings.HasPrefix(actualContentType, "video/") {
				slog.Warn("Disallowed file type detected after sniffing", "filename", originalFilename, "sniffed_type", actualContentType)
				http.Error(w, "Disallowed file type detected", http.StatusBadRequest)
				go func(id string) {
					m.mediaService.FinalizeFileUpload(context.Background(), id, 0, "failed_sniff")
				}(mediaResp.ID)
				return
			}

			// Create a MultiReader to prepend the buffer back to the part, so io.Copy reads from the start.
			prefixedReader := io.MultiReader(bytes.NewReader(buffer[:n]), part)

			// Stream-copy the file content from the multipart part to the destination writer.
			bytesWritten, copyErr := io.Copy(dstWriter, prefixedReader)
			if copyErr != nil {
				slog.Error("Error writing file to disk", "err", copyErr, "filename", originalFilename, "media_id", mediaResp.ID)
				// Attempt to mark as failed
				go func(id string) {
					m.mediaService.FinalizeFileUpload(context.Background(), id, bytesWritten, "failed_copy")
				}(mediaResp.ID)
				http.Error(w, "Error saving the file to disk", http.StatusInternalServerError)
				return
			}
			bytesWritten += int64(n) // Add the bytes read by sniff to total

			slog.Info("File successfully streamed to disk", "media_id", mediaResp.ID, "filename", originalFilename, "bytes_written", bytesWritten)

			// Finalize the upload (e.g., update DB record status and size)
			err = m.mediaService.FinalizeFileUpload(r.Context(), mediaResp.ID, bytesWritten, actualContentType)
			if err != nil {
				slog.Error("Error finalizing file upload metadata", "err", err, "media_id", mediaResp.ID)
				http.Error(w, "Error finalizing file upload", http.StatusInternalServerError)
				return
			}

			mediaUploadResponse = mediaResp
			break // TODO: maybe more?
		} else {
			_, _ = io.Copy(io.Discard, part)
			slog.Debug("Skipped non-'file_upload' part", "form_name", part.FormName(), "filename", part.FileName())
		}
	}

	if !foundFileUpload {
		slog.Warn("No 'file_upload' field found in multipart request or no file uploaded.")
		http.Error(w, "No 'file_upload' field found or no file uploaded", http.StatusBadRequest)
		return
	}

	if mediaUploadResponse == nil {
		slog.Error("Media response was nil after successful upload process, indicating an internal logic error.")
		http.Error(w, "Internal server error: Failed to get media response after upload", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(mediaUploadResponse); err != nil {
		slog.Error("Failed to encode media upload response to JSON", "err", err, "media_id", mediaUploadResponse.ID)
		http.Error(w, "File uploaded, but failed to return response", http.StatusInternalServerError)
	}
}
