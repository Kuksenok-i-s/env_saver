CREATE TABLE IF NOT EXISTS `events` (
    'id' INT NOT NULL AUTO_INCREMENT,
    'title' VARCHAR(255) NOT NULL,
    'description' TEXT,
    'time' DATETIME NOT NULL,
    'commit_id' VARCHAR(255) NOT NULL,
    'commit_message' TEXT,
)
