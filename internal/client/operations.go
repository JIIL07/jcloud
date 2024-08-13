package cloudfiles

import (
	"database/sql"
	"fmt"

	"github.com/JIIL07/jcloud/internal/client/config"
	"github.com/JIIL07/jcloud/internal/client/models"
	"github.com/JIIL07/jcloud/internal/client/util"
)

type FileContext struct {
	DB   *sql.DB
	Info *models.Info
}

func (ctx *FileContext) FileExists() (bool, error) {
	return util.Exists(ctx.DB, ctx.Info.Metadata.Filename, ctx.Info.Metadata.Extension)
}

func (ctx *FileContext) AddFile() error {
	if err := ctx.Info.PrepareInfo(); err != nil {
		return fmt.Errorf("failed to prepare info: %w", err)
	}

	fileExists, err := ctx.FileExists()
	if err != nil {
		return fmt.Errorf("failed to check if file exists: %w", err)
	}

	if !fileExists {
		ctx.Info.Filesize = len(ctx.Info.Data)
		ctx.Info.Status = config.Statuses[0]
		_, err = ctx.DB.Exec(`INSERT INTO files (filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?)`,
			ctx.Info.Filename, ctx.Info.Extension, ctx.Info.Filesize, ctx.Info.Status, ctx.Info.Data)
		return err
	}
	return nil
}

func (ctx *FileContext) DeleteFile() error {
	if err := ctx.Info.GetNameExt(); err != nil {
		return err
	}

	_, err := ctx.DB.Exec(`DELETE FROM files WHERE filename = ? AND extension = ?`, ctx.Info.Filename, ctx.Info.Extension)
	return err
}

func (ctx *FileContext) ListFiles(tablename string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tablename)
	rows, err := ctx.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return util.ScanRows(rows)
}

func (ctx *FileContext) DataInFile() error {
	if err := ctx.Info.GetNameExt(); err != nil {
		return err
	}

	query := `SELECT data FROM files WHERE filename=? AND extension=?`
	rows, err := ctx.DB.Query(query, ctx.Info.Filename, ctx.Info.Extension)
	if err != nil {
		return err
	}
	defer rows.Close()

	return util.WriteData(rows, ctx.Info)
}

func (ctx *FileContext) SearchFile() error {
	if err := ctx.Info.GetNameExt(); err != nil {
		return err
	}

	rows, err := util.Find(ctx.DB, ctx.Info.Filename, ctx.Info.Extension)
	if err != nil {
		return err
	}
	defer rows.Close()

	item, err := util.ScanRow(rows)
	if err != nil {
		return err
	}

	fmt.Printf("Found: %v.%v\n", item["filename"], item["extension"])
	return nil
}
