package storage

import (
	"github.com/JIIL07/jcloud/internal/server/types"
)

func ConvertFileToTypes(file *File) *types.File {
	if file == nil {
		return nil
	}
	return &types.File{
		ID:            file.ID,
		UserID:        file.UserID,
		LastVersionID: file.LastVersionID,
		Metadata: types.FileMetadata{
			Name:        file.Metadata.Name,
			Extension:   file.Metadata.Extension,
			Size:        file.Metadata.Size,
			HashSum:     file.Metadata.HashSum,
			Description: file.Metadata.Description,
		},
		Status:     file.Status,
		Data:       file.Data,
		CreatedAt:  file.CreatedAt,
		ModifiedAt: file.ModifiedAt,
	}
}

func ConvertFileFromTypes(file *types.File) *File {
	if file == nil {
		return nil
	}
	return &File{
		ID:            file.ID,
		UserID:        file.UserID,
		LastVersionID: file.LastVersionID,
		Metadata: FileMetadata{
			Name:        file.Metadata.Name,
			Extension:   file.Metadata.Extension,
			Size:        file.Metadata.Size,
			HashSum:     file.Metadata.HashSum,
			Description: file.Metadata.Description,
		},
		Status:     file.Status,
		Data:       file.Data,
		CreatedAt:  file.CreatedAt,
		ModifiedAt: file.ModifiedAt,
	}
}

func ConvertFilesToTypes(files []File) []types.File {
	result := make([]types.File, len(files))
	for i, file := range files {
		result[i] = *ConvertFileToTypes(&file)
	}
	return result
}

func ConvertFilesFromTypes(files []types.File) []File {
	result := make([]File, len(files))
	for i, file := range files {
		result[i] = *ConvertFileFromTypes(&file)
	}
	return result
}

func ConvertFileVersionToTypes(version *FileVersion) *types.FileVersion {
	if version == nil {
		return nil
	}
	return &types.FileVersion{

		FileID:      version.FileID,
		UserID:      version.UserID,
		Version:     version.Version,
		FullVersion: version.FullVersion,
		Delta:       version.Delta,
		ChangeType:  version.ChangeType,
		CreatedAt:   version.CreatedAt,
	}
}

func ConvertFileVersionFromTypes(version *types.FileVersion) *FileVersion {
	if version == nil {
		return nil
	}
	return &FileVersion{
		FileID:      version.FileID,
		UserID:      version.UserID,
		Version:     version.Version,
		FullVersion: version.FullVersion,
		Delta:       version.Delta,
		ChangeType:  version.ChangeType,
		CreatedAt:   version.CreatedAt,
	}
}

func ConvertFileVersionsToTypes(versions []FileVersion) []types.FileVersion {
	result := make([]types.FileVersion, len(versions))
	for i, version := range versions {
		result[i] = *ConvertFileVersionToTypes(&version)
	}
	return result
}

func ConvertUserToTypes(user *User) *types.User {
	if user == nil {
		return nil
	}
	return &types.User{
		UserID:       user.UserID,
		Username:     user.Username,
		Email:        user.Email,
		Password:     user.Password,
		HashProtocol: user.Protocol,   // storage.User has Protocol field, not HashProtocol
		Admin:        user.Admin != 0, // Convert int to bool
	}
}

func ConvertUserFromTypes(user *types.User) *User {
	if user == nil {
		return nil
	}
	adminInt := 0
	if user.Admin {
		adminInt = 1
	}
	return &User{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Protocol: user.HashProtocol, // storage.User has Protocol field, not HashProtocol
		Admin:    adminInt,          // Convert bool to int
	}
}
