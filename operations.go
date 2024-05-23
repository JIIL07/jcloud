package cloudfiles

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

type FileContext struct {
	DB   *sql.DB
	Info *Info
}

func (ctx *FileContext) Add() error {
	if err := ctx.Info.PrepareInfo(); err != nil {
		return err
	}

	exists, err := exists(ctx.DB, ctx.Info.Filename, ctx.Info.Extension)
	if err != nil {
		return err
	}

	if !exists {
		ctx.Info.Filesize = len(ctx.Info.Data)
		ctx.Info.Status = Statuses[0]
		_, err = ctx.DB.Exec(`INSERT INTO files (filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?)`,
			ctx.Info.Filename, ctx.Info.Extension, ctx.Info.Filesize, ctx.Info.Status, ctx.Info.Data)
		return err
	}
	return nil
}

func (ctx *FileContext) Delete() error {
	if err := ctx.Info.GetNameExt(); err != nil {
		return err
	}

	err := ctx.DB.QueryRow("SELECT * FROM files WHERE filename = ? AND extension = ?",
		ctx.Info.Filename, ctx.Info.Extension).Scan(&ctx.Info.Id, &ctx.Info.Filename,
		&ctx.Info.Extension, &ctx.Info.Filesize, &ctx.Info.Status, &ctx.Info.Data)

	if err != nil {
		return fmt.Errorf("query row error: %v", err)
	}

	tx, err := ctx.DB.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction error: %v", err)
	}

	_, err = tx.Exec("INSERT INTO deleted (id, filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?, ?)",
		ctx.Info.Id, ctx.Info.Filename, ctx.Info.Extension, ctx.Info.Filesize, "Deleted", ctx.Info.Data)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("insert error: %v", err)
	}

	_, err = tx.Exec("DELETE FROM files WHERE filename = ? AND extension = ?",
		ctx.Info.Filename, ctx.Info.Extension)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("delete error: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction error: %v", err)
	}
	return nil
}

func (ctx *FileContext) List(tablename string) error {
	rows, err := ctx.DB.Query(fmt.Sprintf("SELECT * FROM %s", tablename))
	if err != nil {
		return err
	}
	defer rows.Close()
	id := 1
	for rows.Next() {
		err := rows.Scan(new(interface{}), &ctx.Info.Filename, &ctx.Info.Extension,
			new(interface{}), new(interface{}), new(interface{}))
		if err != nil {
			return err
		}
		fmt.Printf("id:%v %v.%v\n", id, ctx.Info.Filename, ctx.Info.Extension)
		id++
	}
	return nil
}

func (ctx *FileContext) DataIn() error {
	if err := ctx.Info.GetNameExt(); err != nil {
		return err
	}

	rows, err := ctx.DB.Query(`SELECT data FROM files WHERE filename=? AND extension=?`,
		ctx.Info.Filename, ctx.Info.Extension)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&ctx.Info.Data)
		if err != nil {
			return err
		}
		fmt.Printf("Data:\n%v\n", strings.TrimSpace(string(ctx.Info.Data)))
	}
	return nil
}

func (ctx *FileContext) Search() error {
	if err := ctx.Info.GetNameExt(); err != nil {
		return err
	}

	rows, err := find(ctx.DB, ctx.Info.Filename, ctx.Info.Extension)
	if err != nil {
		return err
	}
	defer rows.Close()

	err = rows.Scan(new(interface{}), &ctx.Info.Filename, &ctx.Info.Extension,
		new(interface{}), new(interface{}), new(interface{}))
	if err != nil {
		return err
	}

	fmt.Printf("Found: %v.%v\n", ctx.Info.Filename, ctx.Info.Extension)
	return nil
}

func (ctx *FileContext) WriteData() error {
	if err := ctx.Info.PrepareInfo(); err != nil {
		return err
	}

	_, err := ctx.DB.Exec(`UPDATE files SET data = ? WHERE filename = ? AND extension = ?`,
		ctx.Info.Data, ctx.Info.Filename, ctx.Info.Extension)
	return err
}

func (ctx *FileContext) AddFile() error {
	tempDir, err := createTempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	err = openExplorer(tempDir)
	if err != nil {
		return err
	}

	err = waitFile(tempDir)
	if err != nil {
		return err
	}

	ctx.Info, err = processFile(tempDir)
	if err != nil {
		return err
	}

	ctx.Info.Filesize = len(ctx.Info.Data)
	ctx.Info.Status = Statuses[0]
	_, err = ctx.DB.Exec("INSERT INTO files (filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?)",
		ctx.Info.Filename, ctx.Info.Extension, ctx.Info.Filesize, ctx.Info.Status, ctx.Info.Data)
	return err
}
