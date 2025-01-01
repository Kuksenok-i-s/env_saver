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

type LocalStorageInterface interface {
	SaveEvent(event Event) error
	GetEvents() ([]Event, error)
	WriteConfig(config *config.Config) error
	GetConfigs() ([]config.Config, error)
	GetConfig(name string) (config.Config, error)
	UpdateConfig(config config.Config) error
}

func NewLocalStorage() LocalStorageInterface {
	db, err := sql.Open("sqlite3", "./local_storage.db")
	if err != nil {
		log.Fatal(err)
	}
	initLocalStorage(db)
	return &Storage{
		db: db,
	}
}

func initLocalStorage(db *sql.DB) (*Storage, error) {

	_, err := db.Exec(`
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
			'repository_dir' VARCHAR(255) NOT NULL,
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

func (s *Storage) WriteConfig(config *config.Config) error {
	_, err := s.db.Exec("INSERT INTO configs (watch_dir, watched_file_types, repository_url, repository_dir, make_remote_backup, make_tags) VALUES (?, ?, ?, ?, ?, ?)",
		config.WatchDir,
		config.WatchedFileTypes,
		config.RemoteRepo,
		config.RepositoryDir,
		config.MakeRemoteBackup,
		config.MakeTags,
	)
	if err != nil {
		return err
	}
	return nil
}

// Show all configs if needed
func (s *Storage) GetConfigs() ([]config.Config, error) {
	rows, err := s.db.Query("SELECT * FROM configs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	Configs := []config.Config{}
	for raw := rows.Next(); raw; raw = rows.Next() {
		config, err := rowToConfig(rows)
		if err != nil {
			log.Fatal(err) // TODO: handle error properly
		}
		Configs = append(Configs, config)
	}
	return Configs, nil
}

func rowToConfig(row *sql.Rows) (config.Config, error) {
	var config config.Config
	err := row.Scan(
		&config.ID,
		&config.WatchDir,
		&config.WatchedFileTypes,
		&config.RemoteRepo,
		&config.RepositoryDir,
		&config.MakeRemoteBackup,
		&config.MakeTags,
	)
	if err != nil {
		return config, err
	}
	return config, nil
}

func (s *Storage) GetConfig(name string) (config.Config, error) {
	row := s.db.QueryRow("SELECT * FROM configs WHERE watch_dir = ?", name)
	var config config.Config
	err := row.Scan(
		&config.ID,
		&config.WatchDir,
		&config.WatchedFileTypes,
		&config.RemoteRepo,
		&config.RepositoryDir,
		&config.MakeRemoteBackup,
		&config.MakeTags,
	)
	if err != nil {
		return config, err
	}
	return config, nil
}

func (s *Storage) UpdateConfig(config config.Config) error {
	_, err := s.db.Exec("UPDATE configs SET watch_dir = ?, watched_file_types = ?, repository_dir = ?, repository_url = ?, make_remote_backup = ?, make_tags = ? WHERE id = ?",
		config.WatchDir,
		config.WatchedFileTypes,
		config.RemoteRepo,
		config.RepositoryDir,
		config.MakeRemoteBackup,
		config.MakeTags,
		config.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
