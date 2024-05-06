package file

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

// FileContext holds the context for file operations with database and metadata.
type FileContext struct {
	DB   *sql.DB
	Info *Info
}

// Add inserts a new file record if it does not exist.
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
		query := `INSERT INTO files (filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?)`
		_, err = ctx.DB.Exec(query, ctx.Info.Filename, ctx.Info.Extension, ctx.Info.Filesize, ctx.Info.Status, ctx.Info.Data)
		return err
	}
	return nil
}

// Delete marks a file as deleted and moves it to a "deleted" table.
func (ctx *FileContext) Delete() error {
	if err := ctx.Info.GetNameExt(); err != nil {
		return err
	}

	rows, err := find(ctx.DB, ctx.Info.Filename, ctx.Info.Extension)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows != nil {
		tx, err := ctx.DB.Begin()
		if err != nil {
			return fmt.Errorf("begin transaction error: %v", err)
		}

		stmt, err := tx.Prepare("INSERT INTO deleted (filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("prepare insert statement error: %v", err)
		}
		defer stmt.Close()

		if _, err := stmt.Exec(ctx.Info.Filename, ctx.Info.Extension, ctx.Info.Filesize, Statuses[3], ctx.Info.Data); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert into deleted error: %v", err)
		}

		if _, err := tx.Exec("DELETE FROM files WHERE filename = ? AND extension = ?", ctx.Info.Filename, ctx.Info.Extension); err != nil {
			tx.Rollback()
			return fmt.Errorf("delete from files error: %v", err)
		}

		return tx.Commit()
	}
	return err
}

// List prints all entries from a specified table.
func (ctx *FileContext) List(tablename string) error {
	if !isValidTableName(tablename) {
		return fmt.Errorf("invalid table name: %s", tablename)
	}
	query := fmt.Sprintf("SELECT * FROM %s", tablename)
	rows, err := ctx.DB.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	id := 1
	for rows.Next() {
		if err := rows.Scan(new(interface{}), &ctx.Info.Filename, &ctx.Info.Extension,
			new(interface{}), new(interface{}), new(interface{})); err != nil {
			return err
		}
		fmt.Printf("id:%v %v.%v\n", id, ctx.Info.Filename, ctx.Info.Extension)
		id++
	}
	return nil
}

// DataIn retrieves and displays the data for a specified file from the database.
func (ctx *FileContext) DataIn() error {
	if err := ctx.Info.GetNameExt(); err != nil {
		return err
	}
	query := `SELECT data FROM files WHERE filename=? AND extension=?`
	rows, err := ctx.DB.Query(query, ctx.Info.Filename, ctx.Info.Extension)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&ctx.Info.Data); err != nil {
			return err
		}
		fmt.Printf("Data:\n%v\n", strings.TrimSpace(string(ctx.Info.Data)))
	}
	return nil
}

// Search finds a file in the database and displays its information.
func (ctx *FileContext) Search() error {
	if err := ctx.Info.GetNameExt(); err != nil {
		return err
	}

	rows, err := find(ctx.DB, ctx.Info.Filename, ctx.Info.Extension)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err = rows.Scan(new(interface{}), &ctx.Info.Filename, &ctx.Info.Extension,
		new(interface{}), new(interface{}), new(interface{})); err != nil {
		return err
	}
	fmt.Printf("Found: %v.%v\n", ctx.Info.Filename, ctx.Info.Extension)
	return nil
}

// WriteData updates the data for a specified file in the database.
func (ctx *FileContext) WriteData() error {
	if err := ctx.Info.PrepareInfo(); err != nil {
		return err
	}
	query := `UPDATE files SET data = ? WHERE filename = ? AND extension = ?`
	_, err := ctx.DB.Exec(query, ctx.Info.Data, ctx.Info.Filename, ctx.Info.Extension)
	return err
}

// AddFile guides the user through adding a file to the system including opening the file explorer,
func (ctx *FileContext) AddFile() error {
	tempDir, err := createTempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if err = openExplorer(tempDir); err != nil {
		return err
	}

	if err = waitFile(tempDir); err != nil {
		return err
	}

	if ctx.Info, err = processFile(tempDir); err != nil {
		return err
	}

	ctx.Info.Filesize = len(ctx.Info.Data)
	ctx.Info.Status = Statuses[0]
	query := "INSERT INTO files (filename, extension, filesize, status, data) VALUES (?, ?, ?, ?, ?)"
	_, err = ctx.DB.Exec(query, ctx.Info.Filename, ctx.Info.Extension, ctx.Info.Filesize, ctx.Info.Status, ctx.Info.Data)
	return err
}
