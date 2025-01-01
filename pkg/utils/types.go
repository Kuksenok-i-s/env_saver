package utils

import "time"

type FileUpdateEvent struct {
	FileName     string
	EventType    string
	EventMessage string
	Time         time.Time
}
