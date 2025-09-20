package services

import (
	"context"
	"fmt"

	"github.com/JIIL07/jcloud/internal/server/interfaces"
	"github.com/JIIL07/jcloud/internal/server/storage"
	jhash "github.com/JIIL07/jcloud/pkg/hash"
)

// FileService handles file business logic
type FileService struct {
	storage   interfaces.StorageService
	validator interfaces.ValidationService
	response  interfaces.ResponseService
}

// NewFileService creates a new file service
func NewFileService(
	storage interfaces.StorageService,
	validator interfaces.ValidationService,
	response interfaces.ResponseService,
) *FileService {
	return &FileService{
		storage:   storage,
		validator: validator,
		response:  response,
	}
}

// UploadFiles uploads multiple files
func (fs *FileService) UploadFiles(ctx context.Context, userID int, files []storage.File) error {
	if len(files) == 0 {
		return fmt.Errorf("no files provided")
	}

	// Validate all files first
	for i, file := range files {
		file.UserID = userID
		if err := fs.validator.ValidateFile(&file); err != nil {
			return fmt.Errorf("validation failed for file %d: %w", i, err)
		}
	}

	// Start transaction
	tx, err := fs.storage.BeginTx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Add all files in transaction
	for i, file := range files {
		if err := fs.storage.AddFileTx(tx, &file); err != nil {
			return fmt.Errorf("failed to add file %d: %w", i, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DownloadFile downloads a file
func (fs *FileService) DownloadFile(ctx context.Context, userID int, filename string) (*storage.File, error) {
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

// DeleteFile deletes a file
func (fs *FileService) DeleteFile(ctx context.Context, userID int, filename string) error {
	if err := fs.validator.ValidateFilename(filename); err != nil {
		return fmt.Errorf("invalid filename: %w", err)
	}

	file := &storage.File{
		UserID: userID,
		Metadata: storage.FileMetadata{
			Name: filename,
		},
	}

	if err := fs.storage.DeleteFile(file); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// ListFiles lists all files for a user
func (fs *FileService) ListFiles(ctx context.Context, userID int) ([]storage.File, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	files, err := fs.storage.GetAllFiles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get files: %w", err)
	}

	return files, nil
}

// GetFileInfo gets file information
func (fs *FileService) GetFileInfo(ctx context.Context, userID int, filename string) (*storage.File, error) {
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

// UpdateFileMetadata updates file metadata
func (fs *FileService) UpdateFileMetadata(ctx context.Context, userID int, filename string, metadata storage.FileMetadata) error {
	if err := fs.validator.ValidateFilename(filename); err != nil {
		return fmt.Errorf("invalid filename: %w", err)
	}

	if err := fs.validator.ValidateFileMetadata(metadata); err != nil {
		return fmt.Errorf("invalid metadata: %w", err)
	}

	// Create update request
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

// CalculateFileHash calculates file hash
func (fs *FileService) CalculateFileHash(ctx context.Context, userID int, filename string) (string, error) {
	file, err := fs.GetFileInfo(ctx, userID, filename)
	if err != nil {
		return "", err
	}

	hash := jhash.Hash(file.Data)
	return hash, nil
}

// GetImageFiles gets all image files for a user
func (fs *FileService) GetImageFiles(ctx context.Context, userID int) ([]storage.File, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	files, err := fs.storage.GetImageFiles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get image files: %w", err)
	}

	return files, nil
}
