package interfaces

import (
	"context"
	"net/http"

	"github.com/JIIL07/jcloud/internal/server/storage"
)

// StorageService defines the interface for storage operations
type StorageService interface {
	// User operations
	GetUser(username string) (*storage.User, error)
	CreateUser(user *storage.User) error
	UpdateUser(user *storage.User) error
	DeleteUser(userID int) error

	// File operations
	AddFile(file *storage.File) error
	AddFileTx(tx *storage.Tx, file *storage.File) error
	GetFile(userID int, filename string) (*storage.File, error)
	GetAllFiles(userID int) ([]storage.File, error)
	GetImageFiles(userID int) ([]storage.File, error)
	UpdateFile(file *storage.File, data []byte) error
	DeleteFile(file *storage.File) error
	UpdateFileMetadata(userID int, req interface{}) error

	// File version operations
	AddFileVersion(version storage.FileVersion) error
	GetFileVersions(fileID int) ([]storage.FileVersion, error)
	GetFileVersion(fileID, version int) (*storage.FileVersion, error)
	GetLastFileVersion(fileID int) (*storage.FileVersion, error)
	DeleteFileVersion(fileID, version int) error
	DeleteFileVersions(fileID int) error
	RestoreFileToVersion(fileID, version int) ([]byte, error)
	GetFileHistory(fileID int) ([]storage.FileVersion, error)

	// Transaction operations
	BeginTx() (*storage.Tx, error)
	CloseDatabase() error
}

// AuthService defines the interface for authentication operations
type AuthService interface {
	AuthenticateUser(username, password string) (*storage.User, error)
	CreateSession(user *storage.User) (string, error)
	ValidateSession(sessionID string) (*storage.User, error)
	DestroySession(sessionID string) error
}

// FileService defines the interface for file business logic
type FileService interface {
	UploadFiles(ctx context.Context, userID int, files []storage.File) error
	DownloadFile(ctx context.Context, userID int, filename string) (*storage.File, error)
	DeleteFile(ctx context.Context, userID int, filename string) error
	ListFiles(ctx context.Context, userID int) ([]storage.File, error)
	GetFileInfo(ctx context.Context, userID int, filename string) (*storage.File, error)
	UpdateFileMetadata(ctx context.Context, userID int, filename string, metadata storage.FileMetadata) error
	CalculateFileHash(ctx context.Context, userID int, filename string) (string, error)
	GetImageFiles(ctx context.Context, userID int) ([]storage.File, error)
}

// MiddlewareService defines the interface for middleware operations
type MiddlewareService interface {
	StorageMiddleware(storage StorageService) func(http.Handler) http.Handler
	UserMiddleware(next http.Handler) http.Handler
	AuthMiddleware(next http.Handler) http.Handler
	LoggingMiddleware(next http.Handler) http.Handler
}

// ResponseService defines the interface for HTTP response handling
type ResponseService interface {
	WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error
	WriteError(w http.ResponseWriter, statusCode int, message string) error
	WriteSuccess(w http.ResponseWriter, message string) error
	WriteFile(w http.ResponseWriter, filename string, data []byte) error
	WriteHTML(w http.ResponseWriter, statusCode int, html string) error
}

// ValidationService defines the interface for input validation
type ValidationService interface {
	ValidateFile(file *storage.File) error
	ValidateUser(user *storage.User) error
	ValidateFileMetadata(metadata storage.FileMetadata) error
	ValidateFilename(filename string) error
}
