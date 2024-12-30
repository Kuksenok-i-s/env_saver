package storage

type Event struct {
	ID          int
	Title       string
	Description string
	Date        string
	Time        string
	CommitMsg   string
	CommitID    string
}

type GitConfig struct {
	RepoDir string
}
