package storage

import (
	"database/sql"
	"log"

	"github.com/Kuksenok-i-s/env_saver/pkg/config"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// Forgive me for this
func InitEventsStorage(config *config.Config) (*Storage, error) {
	db, err := sql.Open("sqlite3", "./local.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS 'events' (
			'id' INT NOT NULL AUTO_INCREMENT,
			'title' VARCHAR(255) NOT NULL,
			'description' TEXT,
			'time' DATETIME NOT NULL,
			'commit_id' VARCHAR(255) NOT NULL,
			'commit_message' TEXT,
		)
    `)
	if err != nil {
		return nil, err
	}

	_, errcfg := db.Exec(`
		CREATE TABLE IF NOT EXISTS 'configs' (
			'id' INT NOT NULL AUTO_INCREMENT,
			'watch_dir' VARCHAR(255) NOT NULL,
			'watched_file_types' TEXT,
			'repository_url' VARCHAR(255) NOT NULL,
			'make_remote_backup' BOOLEAN NOT NULL,
			'make_tags' BOOLEAN NOT NULL,
		)
	`)
	if errcfg != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveEvent(event Event) error {
	_, err := s.db.Exec("INSERT INTO events (title, description, time, commit_id, commit_message) VALUES (?, ?, ?, ?, ?)",
		event.Title,
		event.Description,
		event.Date,
		event.CommitID,
		event.CommitMsg,
	)
	if err != nil {
		return err
	}
	return nil
}

// Let's just return all events for now, will use for logs later
func (s *Storage) GetEvents() ([]Event, error) {
	rows, err := s.db.Query("SELECT * FROM events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Events := []Event{}
	for raw := rows.Next(); raw; raw = rows.Next() {
		event, err := rowToEvent(rows)
		if err != nil {
			log.Fatal(err) // TODO: handle error properly
		}
		Events = append(Events, event)
	}
	return Events, nil
}

func (s *Storage) WriteConfig(config *config.ConfigDb) error {
	_, err := s.db.Exec("INSERT INTO configs (watch_dir, watched_file_types, repository_url, make_remote_backup, make_tags) VALUES (?, ?, ?, ?, ?)",
		config.WatchDir,
		config.WatchedFileTypes,
		config.RemoteRepo,
		config.MakeRemoteBackup,
		config.MakeTags,
	)
	if err != nil {
		return err
	}
	return nil
}

func rowToEvent(row *sql.Rows) (Event, error) {
	var event Event
	err := row.Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Date,
		&event.CommitID,
		&event.CommitMsg,
	)
	if err != nil {
		return Event{}, err
	}
	return event, nil
}

// Show all configs if needed
func (s *Storage) GetConfigs() ([]config.ConfigDb, error) {
	rows, err := s.db.Query("SELECT * FROM configs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Configs := []config.ConfigDb{}
	for raw := rows.Next(); raw; raw = rows.Next() {
		config, err := rowToConfig(rows)
		if err != nil {
			log.Fatal(err) // TODO: handle error properly
		}
		Configs = append(Configs, config)
	}
	return Configs, nil
}

func rowToConfig(row *sql.Rows) (config.ConfigDb, error) {
	var config config.ConfigDb
	err := row.Scan(
		&config.ID,
		&config.WatchDir,
		&config.WatchedFileTypes,
		&config.RemoteRepo,
		&config.MakeRemoteBackup,
		&config.MakeTags,
	)
	if err != nil {
		return config, err
	}
	return config, nil
}

func (s *Storage) GetConfig(name string) (config.ConfigDb, error) {
	row := s.db.QueryRow("SELECT * FROM configs WHERE watch_dir = ?", name)
	var config config.ConfigDb
	err := row.Scan(
		&config.ID,
		&config.WatchDir,
		&config.WatchedFileTypes,
		&config.RemoteRepo,
		&config.MakeRemoteBackup,
		&config.MakeTags,
	)
	if err != nil {
		return config, err
	}
	return config, nil
}

func (s *Storage) UpdateConfig(config config.ConfigDb) error {
	_, err := s.db.Exec("UPDATE configs SET watch_dir = ?, watched_file_types = ?, repository_url = ?, make_remote_backup = ?, make_tags = ? WHERE id = ?",
		config.WatchDir,
		config.WatchedFileTypes,
		config.RemoteRepo,
		config.MakeRemoteBackup,
		config.MakeTags,
		config.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
