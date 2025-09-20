package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/JIIL07/jcloud/internal/server/storage"
	"github.com/JIIL07/jcloud/internal/server/types"
	"github.com/JIIL07/jcloud/internal/server/utils"
	"github.com/gorilla/mux"
)

type ResponseService struct{}

func NewResponseService() *ResponseService {
	return &ResponseService{}
}

func (rs *ResponseService) WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func (rs *ResponseService) WriteError(w http.ResponseWriter, statusCode int, message string) error {
	errorResponse := map[string]string{
		"error":   http.StatusText(statusCode),
		"message": message,
	}
	return rs.WriteJSON(w, statusCode, errorResponse)
}

func (rs *ResponseService) WriteSuccess(w http.ResponseWriter, message string) error {
	successResponse := map[string]string{
		"message": message,
		"status":  "success",
	}
	return rs.WriteJSON(w, http.StatusOK, successResponse)
}

func (rs *ResponseService) WriteFile(w http.ResponseWriter, filename string, data []byte) error {
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", string(rune(len(data))))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	return err
}

func (rs *ResponseService) WriteHTML(w http.ResponseWriter, statusCode int, html string) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(html))
	return err
}

type ValidationService struct{}

func NewValidationService() *ValidationService {
	return &ValidationService{}
}

func (vs *ValidationService) ValidateFile(file *types.File) error {
	if file == nil {
		return fmt.Errorf("file cannot be nil")
	}
	return nil
}

func (vs *ValidationService) ValidateUser(user *types.User) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}
	return nil
}

func (vs *ValidationService) ValidateFileMetadata(metadata types.FileMetadata) error {
	if metadata.Name == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	return nil
}

func (vs *ValidationService) ValidateFilename(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	return nil
}

type StorageAdapter struct {
	*storage.Storage
}

func NewStorageAdapter(storage *storage.Storage) types.StorageService {
	return &StorageAdapter{Storage: storage}
}

var _ types.StorageService = (*StorageAdapter)(nil)

func (sa *StorageAdapter) BeginTx() (types.Tx, error) {
	tx, err := sa.Storage.BeginTx()
	if err != nil {
		return nil, err
	}
	return storage.NewTxAdapter(tx.Tx), nil
}

func (sa *StorageAdapter) AddFile(file *types.File) error {
	storageFile := storage.ConvertFileFromTypes(file)
	return sa.Storage.AddFile(storageFile)
}

func (sa *StorageAdapter) AddFileTx(tx types.Tx, file *types.File) error {
	storageFile := storage.ConvertFileFromTypes(file)
	storageTx := &storage.Tx{Tx: tx.(*storage.TxAdapter).Tx}
	return sa.Storage.AddFileTx(storageTx.Tx, storageFile)
}

func (sa *StorageAdapter) GetFile(userID int, filename string) (*types.File, error) {
	file, err := sa.Storage.GetFile(userID, filename)
	if err != nil {
		return nil, err
	}
	return storage.ConvertFileToTypes(file), nil
}

func (sa *StorageAdapter) GetAllFiles(userID int) ([]types.File, error) {
	files, err := sa.Storage.GetAllFiles(userID)
	if err != nil {
		return nil, err
	}
	return storage.ConvertFilesToTypes(files), nil
}

func (sa *StorageAdapter) GetImageFiles(userID int) ([]types.File, error) {
	files, err := sa.Storage.GetImageFiles(userID)
	if err != nil {
		return nil, err
	}
	return storage.ConvertFilesToTypes(files), nil
}

func (sa *StorageAdapter) UpdateFile(file *types.File, data []byte) error {
	storageFile := storage.ConvertFileFromTypes(file)
	return sa.Storage.UpdateFile(storageFile, data)
}

func (sa *StorageAdapter) DeleteFile(file *types.File) error {
	storageFile := storage.ConvertFileFromTypes(file)
	return sa.Storage.DeleteFile(storageFile)
}

func (sa *StorageAdapter) UpdateFileMetadata(userID int, req interface{}) error {
	updateReq, ok := req.(struct {
		Filename    string `json:"filename"`
		Extension   string `json:"extension"`
		Description string `json:"description"`
		OldName     string `json:"oldname"`
	})
	if !ok {
		return fmt.Errorf("invalid request type")
	}
	return sa.Storage.UpdateFileMetadata(userID, updateReq)
}

func (sa *StorageAdapter) AddFileVersion(version types.FileVersion) error {
	storageVersion := storage.ConvertFileVersionFromTypes(&version)
	return sa.Storage.AddFileVersion(*storageVersion)
}

func (sa *StorageAdapter) GetFileVersions(fileID int) ([]types.FileVersion, error) {
	versions, err := sa.Storage.GetFileVersions(fileID)
	if err != nil {
		return nil, err
	}
	return storage.ConvertFileVersionsToTypes(versions), nil
}

func (sa *StorageAdapter) GetFileVersion(fileID, version int) (*types.FileVersion, error) {
	fileVersion, err := sa.Storage.GetFileVersion(fileID, version)
	if err != nil {
		return nil, err
	}
	return storage.ConvertFileVersionToTypes(&fileVersion), nil
}

func (sa *StorageAdapter) GetLastFileVersion(fileID int) (*types.FileVersion, error) {
	fileVersion, err := sa.Storage.GetLastFileVersion(fileID)
	if err != nil {
		return nil, err
	}
	return storage.ConvertFileVersionToTypes(&fileVersion), nil
}

func (sa *StorageAdapter) DeleteFileVersion(fileID, version int) error {
	return sa.Storage.DeleteFileVersion(fileID, version)
}

func (sa *StorageAdapter) DeleteFileVersions(fileID int) error {
	return sa.Storage.DeleteFileVersions(fileID)
}

func (sa *StorageAdapter) RestoreFileToVersion(fileID, version int) ([]byte, error) {
	return sa.Storage.RestoreFileToVersion(fileID, version)
}

func (sa *StorageAdapter) GetFileHistory(fileID int) ([]types.FileVersion, error) {
	versions, err := sa.Storage.GetFileHistory(fileID)
	if err != nil {
		return nil, err
	}
	return storage.ConvertFileVersionsToTypes(versions), nil
}

func (sa *StorageAdapter) GetUser(username string) (*types.User, error) {
	user, err := sa.Storage.GetUser(username)
	if err != nil {
		return nil, err
	}
	return storage.ConvertUserToTypes(&user), nil
}

func (sa *StorageAdapter) CreateUser(user *types.User) error {
	storageUser := storage.ConvertUserFromTypes(user)
	return sa.Storage.SaveNewUser(storageUser)
}

func (sa *StorageAdapter) UpdateUser(user *types.User) error {

	updates := map[string]interface{}{
		"email":    user.Email,
		"password": user.Password,
	}
	return sa.Storage.UpdateUserInfo(user.Username, updates)
}

func (sa *StorageAdapter) DeleteUser(userID int) error {

	return fmt.Errorf("DeleteUser by ID not implemented yet")
}

type FileService struct {
	storage   types.StorageService
	validator types.ValidationService
	response  types.ResponseService
}

func NewFileService(
	storage types.StorageService,
	validator types.ValidationService,
	response types.ResponseService,
) *FileService {
	return &FileService{
		storage:   storage,
		validator: validator,
		response:  response,
	}
}

func (fs *FileService) UploadFiles(ctx context.Context, userID int, files []types.File) error {
	if len(files) == 0 {
		return fmt.Errorf("no files provided")
	}

	for i, file := range files {
		file.UserID = userID
		if err := fs.validator.ValidateFile(&file); err != nil {
			return fmt.Errorf("validation failed for file %d: %w", i, err)
		}
	}

	tx, err := fs.storage.BeginTx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for i, file := range files {
		if err := fs.storage.AddFileTx(tx, &file); err != nil {
			return fmt.Errorf("failed to add file %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (fs *FileService) DownloadFile(ctx context.Context, userID int, filename string) (*types.File, error) {
	if err := fs.validator.ValidateFilename(filename); err != nil {
		return nil, fmt.Errorf("invalid filename: %w", err)
	}

	file, err := fs.storage.GetFile(userID, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	if file == nil {
		return nil, fmt.Errorf("file not found")
	}

	return file, nil
}

func (fs *FileService) DeleteFile(ctx context.Context, userID int, filename string) error {
	if err := fs.validator.ValidateFilename(filename); err != nil {
		return fmt.Errorf("invalid filename: %w", err)
	}

	file := &types.File{
		UserID: userID,
		Metadata: types.FileMetadata{
			Name: filename,
		},
	}

	if err := fs.storage.DeleteFile(file); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (fs *FileService) ListFiles(ctx context.Context, userID int) ([]types.File, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	files, err := fs.storage.GetAllFiles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get files: %w", err)
	}

	return files, nil
}

func (fs *FileService) GetFileInfo(ctx context.Context, userID int, filename string) (*types.File, error) {
	if err := fs.validator.ValidateFilename(filename); err != nil {
		return nil, fmt.Errorf("invalid filename: %w", err)
	}

	file, err := fs.storage.GetFile(userID, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	if file == nil {
		return nil, fmt.Errorf("file not found")
	}

	return file, nil
}

func (fs *FileService) UpdateFileMetadata(ctx context.Context, userID int, filename string, metadata types.FileMetadata) error {
	if err := fs.validator.ValidateFilename(filename); err != nil {
		return fmt.Errorf("invalid filename: %w", err)
	}

	if err := fs.validator.ValidateFileMetadata(metadata); err != nil {
		return fmt.Errorf("invalid metadata: %w", err)
	}

	updateReq := struct {
		Filename    string `json:"filename"`
		Extension   string `json:"extension"`
		Description string `json:"description"`
		OldName     string `json:"oldname"`
	}{
		Filename:    metadata.Name,
		Extension:   metadata.Extension,
		Description: metadata.Description,
		OldName:     filename,
	}

	if err := fs.storage.UpdateFileMetadata(userID, updateReq); err != nil {
		return fmt.Errorf("failed to update metadata: %w", err)
	}

	return nil
}

func (fs *FileService) CalculateFileHash(ctx context.Context, userID int, filename string) (string, error) {
	file, err := fs.GetFileInfo(ctx, userID, filename)
	if err != nil {
		return "", err
	}

	return file.Metadata.HashSum, nil
}

func (fs *FileService) GetImageFiles(ctx context.Context, userID int) ([]types.File, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	files, err := fs.storage.GetImageFiles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get image files: %w", err)
	}

	return files, nil
}

type FileHandler struct {
	fileService types.FileService
	response    types.ResponseService
}

func NewFileHandler(fileService types.FileService, response types.ResponseService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
		response:    response,
	}
}

func (fh *FileHandler) AddFileHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		_ = fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var files []types.File
	if err := json.NewDecoder(r.Body).Decode(&files); err != nil {
		_ = fh.response.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := fh.fileService.UploadFiles(r.Context(), user.UserID, files); err != nil {
		_ = fh.response.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to upload files: %v", err))
		return
	}

	_ = fh.response.WriteSuccess(w, "Files uploaded successfully")
}

func (fh *FileHandler) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		_ = fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		_ = fh.response.WriteError(w, http.StatusBadRequest, "Filename is required")
		return
	}

	if err := fh.fileService.DeleteFile(r.Context(), user.UserID, filename); err != nil {
		_ = fh.response.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete file: %v", err))
		return
	}

	_ = fh.response.WriteSuccess(w, "File deleted successfully")
}

func (fh *FileHandler) DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		_ = fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		_ = fh.response.WriteError(w, http.StatusBadRequest, "Filename is required")
		return
	}

	nameWithoutExt := strings.Split(filename, ".")[0]

	file, err := fh.fileService.DownloadFile(r.Context(), user.UserID, nameWithoutExt)
	if err != nil {
		_ = fh.response.WriteError(w, http.StatusNotFound, fmt.Sprintf("File not found: %v", err))
		return
	}

	fullFilename := file.Metadata.Name + "." + file.Metadata.Extension
	if err := fh.response.WriteFile(w, fullFilename, file.Data); err != nil {
		_ = fh.response.WriteError(w, http.StatusInternalServerError, "Failed to write file")
		return
	}
}

func (fh *FileHandler) ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		_ = fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	files, err := fh.fileService.ListFiles(r.Context(), user.UserID)
	if err != nil {
		_ = fh.response.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list files: %v", err))
		return
	}

	if err := fh.response.WriteJSON(w, http.StatusOK, files); err != nil {
		_ = fh.response.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

func (fh *FileHandler) ImageGalleryHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		_ = fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	files, err := fh.fileService.GetImageFiles(r.Context(), user.UserID)
	if err != nil {
		_ = fh.response.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve images: %v", err))
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
		_ = fh.response.WriteError(w, http.StatusInternalServerError, "Failed to write HTML response")
		return
	}
}

func (fh *FileHandler) HashSumHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		_ = fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		_ = fh.response.WriteError(w, http.StatusBadRequest, "Filename is required")
		return
	}

	hash, err := fh.fileService.CalculateFileHash(r.Context(), user.UserID, filename)
	if err != nil {
		_ = fh.response.WriteError(w, http.StatusNotFound, fmt.Sprintf("File not found: %v", err))
		return
	}

	response := map[string]string{
		"filename": filename,
		"checksum": hash,
	}

	if err := fh.response.WriteJSON(w, http.StatusOK, response); err != nil {
		_ = fh.response.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

func (fh *FileHandler) FileInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		_ = fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		_ = fh.response.WriteError(w, http.StatusBadRequest, "Filename is required")
		return
	}

	file, err := fh.fileService.GetFileInfo(r.Context(), user.UserID, filename)
	if err != nil {
		_ = fh.response.WriteError(w, http.StatusNotFound, fmt.Sprintf("File not found: %v", err))
		return
	}

	if err := fh.response.WriteJSON(w, http.StatusOK, file); err != nil {
		_ = fh.response.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

func (fh *FileHandler) UpdateMetadataHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.ProvideUser(r, w)
	if user == nil {
		_ = fh.response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if r.Method != http.MethodPatch {
		_ = fh.response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Filename    string `json:"filename"`
		Extension   string `json:"extension"`
		Description string `json:"description"`
		OldName     string `json:"oldname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = fh.response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	metadata := types.FileMetadata{
		Name:        req.Filename,
		Extension:   req.Extension,
		Description: req.Description,
	}

	if err := fh.fileService.UpdateFileMetadata(r.Context(), user.UserID, req.OldName, metadata); err != nil {
		_ = fh.response.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update metadata: %v", err))
		return
	}

	_ = fh.response.WriteSuccess(w, "Metadata updated successfully")
}

type FileVersionHandler struct {
	storage  types.StorageService
	response types.ResponseService
}

func NewFileVersionHandler(storage types.StorageService, response types.ResponseService) *FileVersionHandler {
	return &FileVersionHandler{
		storage:  storage,
		response: response,
	}
}

func (fvh *FileVersionHandler) AddFileVersionHandler(w http.ResponseWriter, r *http.Request) {
	var version types.FileVersion
	if err := json.NewDecoder(r.Body).Decode(&version); err != nil {
		_ = fvh.response.WriteError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if err := fvh.storage.AddFileVersion(version); err != nil {
		_ = fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to add file version")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (fvh *FileVersionHandler) GetFileVersionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileIDStr := vars["filename"]
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		_ = fvh.response.WriteError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	versions, err := fvh.storage.GetFileVersions(fileID)
	if err != nil {
		_ = fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to fetch file versions")
		return
	}

	if err := fvh.response.WriteJSON(w, http.StatusOK, versions); err != nil {
		_ = fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

func (fvh *FileVersionHandler) GetFileVersionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileIDStr := vars["filename"]
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		_ = fvh.response.WriteError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	versionStr := vars["version"]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		_ = fvh.response.WriteError(w, http.StatusBadRequest, "Invalid version")
		return
	}

	versionData, err := fvh.storage.GetFileVersion(fileID, version)
	if err != nil {
		_ = fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to fetch file version")
		return
	}

	if err := fvh.response.WriteJSON(w, http.StatusOK, versionData); err != nil {
		_ = fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

func (fvh *FileVersionHandler) GetLastFileVersionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID, err := strconv.Atoi(vars["filename"])
	if err != nil {
		_ = fvh.response.WriteError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	versionData, err := fvh.storage.GetLastFileVersion(fileID)
	if err != nil {
		_ = fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to fetch last file version")
		return
	}

	if err := fvh.response.WriteJSON(w, http.StatusOK, versionData); err != nil {
		_ = fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

func (fvh *FileVersionHandler) DeleteFileVersionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID, err := strconv.Atoi(vars["filename"])
	if err != nil {
		_ = fvh.response.WriteError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	version, err := strconv.Atoi(vars["version"])
	if err != nil {
		_ = fvh.response.WriteError(w, http.StatusBadRequest, "Invalid version")
		return
	}

	if err := fvh.storage.DeleteFileVersion(fileID, version); err != nil {
		_ = fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to delete file version")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (fvh *FileVersionHandler) DeleteFileVersionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID, err := strconv.Atoi(vars["filename"])
	if err != nil {
		_ = fvh.response.WriteError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	if err := fvh.storage.DeleteFileVersions(fileID); err != nil {
		_ = fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to delete file versions")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (fvh *FileVersionHandler) RestoreFileToVersionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileIDStr := vars["filename"]
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		_ = fvh.response.WriteError(w, http.StatusBadRequest, "Invalid file ID")
		return
	}

	targetVersion, err := strconv.Atoi(r.URL.Query().Get("version"))
	if err != nil {
		_ = fvh.response.WriteError(w, http.StatusBadRequest, "Invalid version")
		return
	}

	fileContent, err := fvh.storage.RestoreFileToVersion(fileID, targetVersion)
	if err != nil {
		_ = fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to restore file")
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(fileContent); err != nil {
		_ = fvh.response.WriteError(w, http.StatusInternalServerError, "Failed to write file content")
		return
	}
}
