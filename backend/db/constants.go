package db

// Task Statuses
const (
	TaskStatusPending    = "pending"
	TaskStatusProcessing = "processing"
	TaskStatusCompleted  = "completed"
	TaskStatusFailed     = "failed"
	TaskStatusNotStarted = "not_started"
	TaskStatusPaused     = "paused"
)

// Task Types
const (
	TaskTypePlaylistDownload = "playlist_download"
	TaskTypeSongMixDownload  = "song_mix_download"
	TaskTypeOngoingDownload  = "ongoing_download"
)

// Task Sources
const (
	TaskSourceUser   = "user"
	TaskSourceSystem = "system"
)

// settings
const (
	OngoingListeningEnabledSetting = "ongoing_listening_enabled"
	DownloadPathSetting            = "download_path"
)
