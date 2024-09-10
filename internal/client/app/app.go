package app

import (
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/JIIL07/jcloud/internal/client/storage"
	"github.com/JIIL07/jcloud/internal/config"
	"github.com/JIIL07/jcloud/internal/logger"
	"github.com/JIIL07/jcloud/pkg/home"
	"log/slog"
)

type ClientContext struct {
	Cfg    *config.ClientConfig
	common service

	// Services
	File            *FileService
	Storage         *StorageService
	Paths           *PathsService
	Logger          *LoggerService
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

func NewAppContext(cfg *config.ClientConfig) (*ClientContext, error) {
	p := home.SetPaths()

	s := storage.MustInit()

	l := logger.NewClientLogger(p.JlogFile)

	context := &ClientContext{
		Cfg:    cfg,
		common: service{},
	}

	context.common.Context = context

	context.File = &FileService{
		service: &context.common,
		F:       &models.File{},
	}
	context.Storage = &StorageService{
		service: &context.common,
		S:       s,
	}
	context.Paths = &PathsService{
		service: &context.common,
		P:       p,
	}
	context.Logger = &LoggerService{
		service: &context.common,
		L:       l,
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
