package db

// Task Statuses
const (
	TaskStatusPending    = "pending"
	TaskStatusProcessing = "processing"
	TaskStatusCompleted  = "completed"
	TaskStatusFailed     = "failed"
	TaskStatusNotStarted = "not_started"
)

// Task Types
const (
	TaskTypePlaylistDownload = "playlist_download"
	TaskTypeOngoingDownload  = "ongoing_download"
)

// Task Sources
const (
	TaskSourceUser   = "user"
	TaskSourceSystem = "system"
)
