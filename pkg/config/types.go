package config

type Config struct {
	ID               int
	WatchDir         string
	WatchedFileTypes string
	RemoteRepo       string
	MakeRemoteBackup bool
	MakeTags         bool
}
