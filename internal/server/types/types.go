package types

import (
	"context"
	"net/http"
	"time"
)

type File struct {
	ID            int `db:"id" json:"id"`
	UserID        int `db:"user_id" json:"user_id"`
	LastVersionID int `db:"last_version_id" json:"last_version_id"`
	Metadata      FileMetadata
	Status        string    `db:"status" json:"status"`
	Data          []byte    `db:"data" json:"data"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	ModifiedAt    time.Time `db:"last_modified_at" json:"modified_at"`
}

type FileMetadata struct {
	Name        string `db:"filename" json:"filename"`
	Extension   string `db:"extension" json:"extension"`
	Size        int    `db:"filesize" json:"filesize"`
	HashSum     string `db:"hash_sum" json:"hash_sum"`
	Description string `db:"description,omitempty" json:"description,omitempty"`
}

type FileVersion struct {
	ID          int       `db:"id" json:"id"`
	FileID      int       `db:"file_id" json:"file_id"`
	UserID      int       `db:"user_id" json:"user_id"`
	Version     int       `db:"version" json:"version"`
	FullVersion bool      `db:"full_version" json:"full_version"`
	Delta       []byte    `db:"delta" json:"delta"`
	ChangeType  string    `db:"change_type" json:"change_type"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type User struct {
	UserID       int    `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	Email        string `db:"email" json:"email"`
	Password     string `db:"password" json:"password"`
	HashProtocol string `db:"hashprotocol" json:"hashprotocol"`
	Admin        bool   `db:"admin" json:"admin"`
}

type Tx interface {
	Exec(query string, args ...interface{}) (interface{}, error)
	Commit() error
	Rollback() error
}

type StorageService interface {
	GetUser(username string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(userID int) error

	AddFile(file *File) error
	AddFileTx(tx Tx, file *File) error
	GetFile(userID int, filename string) (*File, error)
	GetAllFiles(userID int) ([]File, error)
	GetImageFiles(userID int) ([]File, error)
	UpdateFile(file *File, data []byte) error
	DeleteFile(file *File) error
	UpdateFileMetadata(userID int, req interface{}) error

	AddFileVersion(version FileVersion) error
	GetFileVersions(fileID int) ([]FileVersion, error)
	GetFileVersion(fileID, version int) (*FileVersion, error)
	GetLastFileVersion(fileID int) (*FileVersion, error)
	DeleteFileVersion(fileID, version int) error
	DeleteFileVersions(fileID int) error
	RestoreFileToVersion(fileID, version int) ([]byte, error)
	GetFileHistory(fileID int) ([]FileVersion, error)

	BeginTx() (Tx, error)
	CloseDatabase() error
}

type AuthService interface {
	AuthenticateUser(username, password string) (*User, error)
	CreateSession(user *User) (string, error)
	ValidateSession(sessionID string) (*User, error)
	DestroySession(sessionID string) error
}

type FileService interface {
	UploadFiles(ctx context.Context, userID int, files []File) error
	DownloadFile(ctx context.Context, userID int, filename string) (*File, error)
	DeleteFile(ctx context.Context, userID int, filename string) error
	ListFiles(ctx context.Context, userID int) ([]File, error)
	GetFileInfo(ctx context.Context, userID int, filename string) (*File, error)
	UpdateFileMetadata(ctx context.Context, userID int, filename string, metadata FileMetadata) error
	CalculateFileHash(ctx context.Context, userID int, filename string) (string, error)
	GetImageFiles(ctx context.Context, userID int) ([]File, error)
}

type MiddlewareService interface {
	StorageMiddleware(storage StorageService) func(http.Handler) http.Handler
	UserMiddleware(next http.Handler) http.Handler
	AuthMiddleware(next http.Handler) http.Handler
	LoggingMiddleware(next http.Handler) http.Handler
}

type ResponseService interface {
	WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error
	WriteError(w http.ResponseWriter, statusCode int, message string) error
	WriteSuccess(w http.ResponseWriter, message string) error
	WriteFile(w http.ResponseWriter, filename string, data []byte) error
	WriteHTML(w http.ResponseWriter, statusCode int, html string) error
}

type ValidationService interface {
	ValidateFile(file *File) error
	ValidateUser(user *User) error
	ValidateFileMetadata(metadata FileMetadata) error
	ValidateFilename(filename string) error
}
