package services

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/JIIL07/jcloud/internal/server/storage"
)

// ValidationService handles input validation
type ValidationService struct{}

// NewValidationService creates a new validation service
func NewValidationService() *ValidationService {
	return &ValidationService{}
}

// ValidateFile validates file data
func (vs *ValidationService) ValidateFile(file *storage.File) error {
	if file == nil {
		return errors.New("file cannot be nil")
	}

	if err := vs.ValidateFileMetadata(file.Metadata); err != nil {
		return err
	}

	if len(file.Data) == 0 {
		return errors.New("file data cannot be empty")
	}

	if file.UserID <= 0 {
		return errors.New("user ID must be positive")
	}

	return nil
}

// ValidateUser validates user data
func (vs *ValidationService) ValidateUser(user *storage.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if strings.TrimSpace(user.Username) == "" {
		return errors.New("username cannot be empty")
	}

	if len(user.Username) < 3 || len(user.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}

	if !vs.isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	if len(user.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

// ValidateFileMetadata validates file metadata
func (vs *ValidationService) ValidateFileMetadata(metadata storage.FileMetadata) error {
	if err := vs.ValidateFilename(metadata.Name); err != nil {
		return err
	}

	if metadata.Size < 0 {
		return errors.New("file size cannot be negative")
	}

	if metadata.Size > 100*1024*1024 { // 100MB limit
		return errors.New("file size exceeds maximum limit of 100MB")
	}

	if strings.TrimSpace(metadata.HashSum) == "" {
		return errors.New("file hash sum cannot be empty")
	}

	return nil
}

// ValidateFilename validates filename
func (vs *ValidationService) ValidateFilename(filename string) error {
	if strings.TrimSpace(filename) == "" {
		return errors.New("filename cannot be empty")
	}

	if !utf8.ValidString(filename) {
		return errors.New("filename contains invalid characters")
	}

	// Check for dangerous characters
	dangerousChars := []string{"..", "/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range dangerousChars {
		if strings.Contains(filename, char) {
			return errors.New("filename contains invalid characters")
		}
	}

	if len(filename) > 255 {
		return errors.New("filename too long")
	}

	return nil
}

// isValidEmail validates email format
func (vs *ValidationService) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateFileVersion validates file version data
func (vs *ValidationService) ValidateFileVersion(version storage.FileVersion) error {
	if version.FileID <= 0 {
		return errors.New("file ID must be positive")
	}

	if version.UserID <= 0 {
		return errors.New("user ID must be positive")
	}

	if version.Version <= 0 {
		return errors.New("version must be positive")
	}

	if len(version.Delta) == 0 {
		return errors.New("version data cannot be empty")
	}

	validChangeTypes := []string{"upload", "edit", "delete", "restore"}
	valid := false
	for _, changeType := range validChangeTypes {
		if version.ChangeType == changeType {
			valid = true
			break
		}
	}

	if !valid {
		return errors.New("invalid change type")
	}

	return nil
}
