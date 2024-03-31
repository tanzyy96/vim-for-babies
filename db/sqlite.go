package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const DB_NAME = "vim_cheatsheet.db"

type CommandDb struct {
	db *sql.DB
}

type Command struct {
	Title       string
	Description string
}

type DB interface {
	Init() error
	Seed() error
	Add(item Command) error
	GetAll() ([]Command, error)
}

func New() CommandDb {
	db, err := sql.Open("sqlite3", DB_NAME)
	if err != nil {
		panic(err)
	}
	cdb := CommandDb{
		db: db,
	}
	if err := cdb.Init(); err != nil {
		panic(err)
	}
	return cdb
}

func (cdb CommandDb) Init() error {
	_, err := cdb.db.Exec(`CREATE TABLE IF NOT EXISTS commands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		description TEXT
	)`)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	var tableExists bool
	err = cdb.db.QueryRow("SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = 'commands')").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("error checking if table exists: %w", err)
	}

	var rowCount int
	if tableExists {
		if err := cdb.db.QueryRow("SELECT COUNT(*) FROM commands").Scan(&rowCount); err != nil {
			return fmt.Errorf("error checking if table is empty: %w", err)
		}
	}
	if rowCount > 0 {
		return nil
	}

	// Seed some data
	if err := cdb.Seed(); err != nil {
		return err
	}

	return nil
}

func (cdb CommandDb) Seed() error {
	_, err := cdb.db.Exec(`INSERT INTO commands (title, description) VALUES
		("Move cursor", "hjkl"),
		("Go to top of page", "gg"),
		("Go to bottom of page", "G"),
		("Move to next word", "w/W"),
		("Move to previous word", "b/B"),
		("Move to end of word", "e/E"),
		("Move to beginning of line", "0"),
		("Move to end of line", "$"),
		("Move to first non-whitespace character of line", "_"),
		("Move to top of screen", "H"),
		("Move to middle of screen", "M"),
		("Move to bottom of screen", "L"),
		("Move up half a page", "Ctrl+u"),
		("Move down half a page", "Ctrl+d"),
		("Move up a page", "Ctrl+b"),
		("Move down a page", "Ctrl+f"),
		("Move screen up one line", "Ctrl+e"),
		("Move screen down one line", "Ctrl+y"),
		("Replace char", "r"),
		("Replace line", "R"),
		("Insert before cursor", "i"),
		("Insert at beginning of line", "I"),
		("Append after cursor", "a"),
		("Append at end of line", "A"),
		("Insert new line below cursor", "o"),
		("Insert new line above cursor", "O"),
		("Delete char", "x"),
		("Delete line", "dd"),
		("Delete word", "dw"),
		("Delete to end of line", "D"),
		("Delete to end of word", "de"),
		("Delete to beginning of line", "d0"),
		("Delete to beginning of word", "db"),
		("Change in word", "ciw"),
		("Change in brackets", "ci{"),
		("Change to end of line", "C"),
		("Change line", "cc"),
		("Indent line", ">>/<<"),
		("Undo", "u"),
		("Redo", "Ctrl+r"),
		("Copy/yank", "y"),
		("Copy/yank line", "yy"),
		("Paste", "p"),
		("Paste before cursor", "P"),
		("Delete", "d"),
		("Delete line", "dd"),
		("Delete word", "dw"),
		("Delete to end of line", "D"),
		("Delete to end of word", "de"),
		("Delete to beginning of line", "d0"),
		("Delete to beginning of word", "db"),
		("Visual line mode", "V"),
		("Flip cursor in visual mode", "o"),
		("Select similar words and replace", "*"),
		("Search for word", "/word + n/N"),
		("Search for word backwards", "?word + n/N"),
		("Search for next word", "n"),
		("Search for previous word", "N")
	`)
	if err != nil {
		return fmt.Errorf("error seeding data: %w", err)
	}
	return nil
}

func (cdb CommandDb) Add(item Command) error {
	_, err := cdb.db.Exec("INSERT INTO commands (title, description) VALUES (?, ?)", item.Title, item.Description)
	if err != nil {
		return fmt.Errorf("error adding item: %w", err)
	}
	return nil
}

func (cdb CommandDb) GetAll() ([]Command, error) {
	rows, err := cdb.db.Query("SELECT title, description FROM commands")
	if err != nil {
		return nil, fmt.Errorf("error getting all items: %w", err)
	}
	defer rows.Close()

	var items []Command
	for rows.Next() {
		var item Command
		if err := rows.Scan(&item.Title, &item.Description); err != nil {
			return nil, fmt.Errorf("error scanning item: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}
