package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"

	"github.com/0xrinful/funlock/internal/models"
	"github.com/0xrinful/funlock/migrations"
)

type config struct {
	dataDir  string
	XPFactor float64
}

type application struct {
	config config
	models *models.Models
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("\033[31mERROR: \033[0m")
}

func main() {
	cfg := config{
		dataDir:  filepath.Join(os.Getenv("XDG_DATA_HOME"), "funlock"),
		XPFactor: 1,
	}

	db, err := openDB(cfg.dataDir)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := application{
		config: cfg,
		models: models.NewModels(db),
	}

	cli := &cli.App{
		Name:     "funlock",
		Usage:    "Earn XP by working, spend XP to unlock fun apps!",
		Commands: app.commands(),
	}

	err = cli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func initDB(db *sql.DB) error {
	_, err := db.Exec(migrations.InitSQL)
	if err != nil {
		return fmt.Errorf("failed to execute init.sql: %w", err)
	}
	return nil
}

func openDB(dataDir string) (*sql.DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dataDir, "funlock.db")

	init := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		init = true
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if init {
		if err := initDB(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}
