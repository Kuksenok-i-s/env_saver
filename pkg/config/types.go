package config

type Config struct {
	ID               int
	WatchDir         string
	WatchedFileTypes string
	RemoteRepo       string
	RepositoryDir    string
	MakeRemoteBackup bool
	MakeTags         bool
}
