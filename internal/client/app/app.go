package app

import (
	"github.com/JIIL07/jcloud/internal/client/config"
	"github.com/JIIL07/jcloud/internal/client/lib/home"
	jlog "github.com/JIIL07/jcloud/internal/client/lib/logger"
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/JIIL07/jcloud/internal/client/storage"
	"log/slog"
)

type ClientContext struct {
	cfg    *config.Config
	common service

	// Services
	FileService     *FileService
	StorageService  *StorageService
	PathsService    *PathsService
	LoggerService   *LoggerService
	AnchorService   *AnchorService
	DeltaService    *DeltaService
	SnapshotService *SnapshotService
}

type service struct {
	Context *ClientContext
}

type FileService struct {
	*service
	F *models.File
}

func NewFileService() *FileService {
	fs := &FileService{}
	return fs
}

type StorageService struct {
	*service
	S *storage.SQLite
}

type PathsService struct {
	*service
	P *home.Paths
}

type LoggerService struct {
	*service
	L *slog.Logger
}

type AnchorService struct {
	*service
}

type DeltaService struct {
	*service
}

type SnapshotService struct {
	*service
}

func NewAppContext(cfg *config.Config) (*ClientContext, error) {
	s := storage.MustInit()

	p := home.SetPaths()

	logger := jlog.NewLogger(p.JlogFile)

	context := &ClientContext{
		cfg:    cfg,
		common: service{},
	}

	context.common.Context = context

	context.FileService = &FileService{
		service: &context.common,
		F:       &models.File{},
	}
	context.StorageService = &StorageService{
		service: &context.common,
		S:       s,
	}
	context.PathsService = &PathsService{
		service: &context.common,
		P:       p,
	}
	context.LoggerService = &LoggerService{
		service: &context.common,
		L:       logger,
	}
	context.AnchorService = &AnchorService{
		service: &context.common,
	}
	context.DeltaService = &DeltaService{
		service: &context.common,
	}
	context.SnapshotService = &SnapshotService{
		service: &context.common,
	}

	return context, nil
}
