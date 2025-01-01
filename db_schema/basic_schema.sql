CREATE TABLE IF NOT EXISTS `events` (
    'id' INT NOT NULL AUTO_INCREMENT,
    'title' VARCHAR(255) NOT NULL,
    'description' TEXT,
    'time' DATETIME NOT NULL,
    'commit_id' VARCHAR(255) NOT NULL,
    'commit_message' TEXT,
)

CREATE TABLE IF NOT EXISTS 'configs' (
    'id' INT NOT NULL AUTO_INCREMENT,
    'name' VARCHAR(255) NOT NULL,
    'watch_dir' VARCHAR(255) NOT NULL,
    'watched_file_types' TEXT,
    'repository_url' VARCHAR(255) NOT NULL,
    'repository_dir' VARCHAR(255) NOT NULL,
    'make_remote_backup' BOOLEAN NOT NULL,
    'make_tags' BOOLEAN NOT NULL,
)