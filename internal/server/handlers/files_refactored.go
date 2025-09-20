package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/JIIL07/jcloud/internal/server/interfaces"
	"github.com/JIIL07/jcloud/internal/server/storage"
	"github.com/JIIL07/jcloud/internal/server/utils"
	"github.com/gorilla/mux"
)

// FileHandler handles file-related HTTP requests
type FileHandler struct {
	fileService interfaces.FileService
	response    interfaces.ResponseService
}

// NewFileHandler creates a new file handler
func NewFileHandler(fileService interfaces.FileService, response interfaces.ResponseService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
		response:    response,
	}
}

// AddFileHandler handles file upload requests
func (fh *FileHandler) AddFileHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var files []storage.File
	if err := json.NewDecoder(r.Body).Decode(&files); err != nil {
		fh.response.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := fh.fileService.UploadFiles(r.Context(), user.UserID, files); err != nil {
		fh.response.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to upload files: %v", err))
		return
	}

	fh.response.WriteSuccess(w, "Files uploaded successfully")
}

// DeleteFileHandler handles file deletion requests
func (fh *FileHandler) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		fh.response.WriteError(w, http.StatusBadRequest, "Filename is required")
		return
	}

	if err := fh.fileService.DeleteFile(r.Context(), user.UserID, filename); err != nil {
		fh.response.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete file: %v", err))
		return
	}

	fh.response.WriteSuccess(w, "File deleted successfully")
}

// DownloadFileHandler handles file download requests
func (fh *FileHandler) DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		fh.response.WriteError(w, http.StatusBadRequest, "Filename is required")
		return
	}

	// Remove extension from filename for lookup
	nameWithoutExt := strings.Split(filename, ".")[0]

	file, err := fh.fileService.DownloadFile(r.Context(), user.UserID, nameWithoutExt)
	if err != nil {
		fh.response.WriteError(w, http.StatusNotFound, fmt.Sprintf("File not found: %v", err))
		return
	}

	fullFilename := file.Metadata.Name + "." + file.Metadata.Extension
	if err := fh.response.WriteFile(w, fullFilename, file.Data); err != nil {
		fh.response.WriteError(w, http.StatusInternalServerError, "Failed to write file")
		return
	}
}

// ListFilesHandler handles file listing requests
func (fh *FileHandler) ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	files, err := fh.fileService.ListFiles(r.Context(), user.UserID)
	if err != nil {
		fh.response.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list files: %v", err))
		return
	}

	if err := fh.response.WriteJSON(w, http.StatusOK, files); err != nil {
		fh.response.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

// ImageGalleryHandler handles image gallery requests
func (fh *FileHandler) ImageGalleryHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	files, err := fh.fileService.GetImageFiles(r.Context(), user.UserID)
	if err != nil {
		fh.response.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve images: %v", err))
		return
	}

	html := "<html><body><h1>Image Gallery</h1><div style='display: flex; flex-wrap: wrap;'>"
	for _, file := range files {
		imageDataURL := fmt.Sprintf("data:image/%s;base64,%s", file.Metadata.Extension, base64.StdEncoding.EncodeToString(file.Data))
		html += fmt.Sprintf(
			"<div style='margin: 10px;'><img src='%s' alt='%s' style='width: 200px; height: auto;'></div>",
			imageDataURL,
			file.Metadata.Name,
		)
	}
	html += "</div></body></html>"

	if err := fh.response.WriteHTML(w, http.StatusOK, html); err != nil {
		fh.response.WriteError(w, http.StatusInternalServerError, "Failed to write HTML response")
		return
	}
}

// HashSumHandler handles file hash calculation requests
func (fh *FileHandler) HashSumHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		fh.response.WriteError(w, http.StatusBadRequest, "Filename is required")
		return
	}

	hash, err := fh.fileService.CalculateFileHash(r.Context(), user.UserID, filename)
	if err != nil {
		fh.response.WriteError(w, http.StatusNotFound, fmt.Sprintf("File not found: %v", err))
		return
	}

	response := map[string]string{
		"filename": filename,
		"checksum": hash,
	}

	if err := fh.response.WriteJSON(w, http.StatusOK, response); err != nil {
		fh.response.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

// FileInfoHandler handles file information requests
func (fh *FileHandler) FileInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		fh.response.WriteError(w, http.StatusBadRequest, "Filename is required")
		return
	}

	file, err := fh.fileService.GetFileInfo(r.Context(), user.UserID, filename)
	if err != nil {
		fh.response.WriteError(w, http.StatusNotFound, fmt.Sprintf("File not found: %v", err))
		return
	}

	if err := fh.response.WriteJSON(w, http.StatusOK, file); err != nil {
		fh.response.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

// UpdateMetadataHandler handles file metadata update requests
func (fh *FileHandler) UpdateMetadataHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if r.Method != http.MethodPatch {
		fh.response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Filename    string `json:"filename"`
		Extension   string `json:"extension"`
		Description string `json:"description"`
		OldName     string `json:"oldname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fh.response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	metadata := storage.FileMetadata{
		Name:        req.Filename,
		Extension:   req.Extension,
		Description: req.Description,
	}

	if err := fh.fileService.UpdateFileMetadata(r.Context(), user.UserID, req.OldName, metadata); err != nil {
		fh.response.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update metadata: %v", err))
		return
	}

	fh.response.WriteSuccess(w, "Metadata updated successfully")
}

// FileVersionHandler handles file version operations
type FileVersionHandler struct {
	storage  interfaces.StorageService
	response interfaces.ResponseService
}

// NewFileVersionHandler creates a new file version handler
func NewFileVersionHandler(storage interfaces.StorageService, response interfaces.ResponseService) *FileVersionHandler {
	return &FileVersionHandler{
		storage:  storage,
		response: response,
	}
}

// AddFileVersionHandler handles file version creation
func (fvh *FileVersionHandler) AddFileVersionHandler(w http.ResponseWriter, r *http.Request) {
	var version storage.FileVersion
	if err := json.NewDecoder(r.Body).Decode(&version); err != nil {
		fvh.response.WriteError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if err := fvh.storage.AddFileVersion(version); err != nil {
		fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to add file version")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetFileVersionsHandler handles file versions listing
func (fvh *FileVersionHandler) GetFileVersionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileIDStr := vars["filename"]
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		fvh.response.WriteError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	versions, err := fvh.storage.GetFileVersions(fileID)
	if err != nil {
		fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to fetch file versions")
		return
	}

	if err := fvh.response.WriteJSON(w, http.StatusOK, versions); err != nil {
		fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

// GetFileVersionHandler handles specific file version retrieval
func (fvh *FileVersionHandler) GetFileVersionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileIDStr := vars["filename"]
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		fvh.response.WriteError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	versionStr := vars["version"]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		fvh.response.WriteError(w, http.StatusBadRequest, "Invalid version")
		return
	}

	versionData, err := fvh.storage.GetFileVersion(fileID, version)
	if err != nil {
		fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to fetch file version")
		return
	}

	if err := fvh.response.WriteJSON(w, http.StatusOK, versionData); err != nil {
		fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

// GetLastFileVersionHandler handles last file version retrieval
func (fvh *FileVersionHandler) GetLastFileVersionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID, err := strconv.Atoi(vars["filename"])
	if err != nil {
		fvh.response.WriteError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	versionData, err := fvh.storage.GetLastFileVersion(fileID)
	if err != nil {
		fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to fetch last file version")
		return
	}

	if err := fvh.response.WriteJSON(w, http.StatusOK, versionData); err != nil {
		fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

// DeleteFileVersionHandler handles file version deletion
func (fvh *FileVersionHandler) DeleteFileVersionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID, err := strconv.Atoi(vars["filename"])
	if err != nil {
		fvh.response.WriteError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	version, err := strconv.Atoi(vars["version"])
	if err != nil {
		fvh.response.WriteError(w, http.StatusBadRequest, "Invalid version")
		return
	}

	if err := fvh.storage.DeleteFileVersion(fileID, version); err != nil {
		fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to delete file version")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteFileVersionsHandler handles all file versions deletion
func (fvh *FileVersionHandler) DeleteFileVersionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID, err := strconv.Atoi(vars["filename"])
	if err != nil {
		fvh.response.WriteError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	if err := fvh.storage.DeleteFileVersions(fileID); err != nil {
		fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to delete file versions")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RestoreFileToVersionHandler handles file restoration to specific version
func (fvh *FileVersionHandler) RestoreFileToVersionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileIDStr := vars["filename"]
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		fvh.response.WriteError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	targetVersion, err := strconv.Atoi(r.URL.Query().Get("version"))
	if err != nil {
		fvh.response.WriteError(w, http.StatusBadRequest, "Invalid version")
		return
	}

	fileContent, err := fvh.storage.RestoreFileToVersion(fileID, targetVersion)
	if err != nil {
		fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to restore file")
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(fileContent); err != nil {
		fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to write file content")
		return
	}
}

// GetFileHistoryHandler handles file history retrieval
func (fvh *FileVersionHandler) GetFileHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileIDStr := vars["filename"]
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		fvh.response.WriteError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	history, err := fvh.storage.GetFileHistory(fileID)
	if err != nil {
		fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to fetch file history")
		return
	}

	if err := fvh.response.WriteJSON(w, http.StatusOK, history); err != nil {
		fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}
