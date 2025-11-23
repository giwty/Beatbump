package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

type Task struct {
	ID          int
	Type        string
	ReferenceID string
	Status      string
	Payload     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func InitDB() {
	var err error
	// Enable WAL mode and set busy timeout to 5000ms to handle concurrency
	DB, err = sql.Open("sqlite", "file:beatbump.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatal(err)
	}

	createTables()

	err = ResetStuckTasks()
	if err != nil {
		log.Printf("Failed to reset stuck tasks: %v", err)
	}
}

type TaskTrack struct {
	TaskID       int
	VideoID      string
	Status       string
	Title        string
	Artist       string
	Album        string
	ThumbnailURL string
	FilePath     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func createTables() {
	taskTable := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT,
		reference_id TEXT,
		status TEXT,
		payload TEXT,
		created_at DATETIME,
		updated_at DATETIME
	);`

	taskTracksTable := `CREATE TABLE IF NOT EXISTS task_tracks (
		task_id INTEGER,
		video_id TEXT,
		status TEXT,
		title TEXT,
		artist TEXT,
		album TEXT,
		thumbnail_url TEXT,
		file_path TEXT,
		created_at DATETIME,
		updated_at DATETIME,
		PRIMARY KEY (task_id, video_id),
		FOREIGN KEY(task_id) REFERENCES tasks(id)
	);`

	settingsTable := `CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT
	);`

	_, err := DB.Exec(taskTable)
	if err != nil {
		log.Fatal("Error creating tasks table:", err)
	}

	_, err = DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_tasks_reference_id ON tasks(reference_id)")
	if err != nil {
		log.Fatal("Error creating unique index on tasks:", err)
	}

	_, err = DB.Exec(taskTracksTable)
	if err != nil {
		log.Fatal("Error creating task_tracks table:", err)
	}

	// Attempt to add thumbnail_url column if it doesn't exist (for existing DBs)
	_, _ = DB.Exec("ALTER TABLE task_tracks ADD COLUMN thumbnail_url TEXT")
	// Attempt to add file_path column if it doesn't exist
	_, _ = DB.Exec("ALTER TABLE task_tracks ADD COLUMN file_path TEXT")

	_, err = DB.Exec(settingsTable)
	if err != nil {
		log.Fatal("Error creating settings table:", err)
	}
}

func AddTaskTrack(taskID int, videoID, title, artist, album, thumbnailURL string) error {
	stmt, err := DB.Prepare("INSERT OR IGNORE INTO task_tracks(task_id, video_id, status, title, artist, album, thumbnail_url, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(taskID, videoID, "not_started", title, artist, album, thumbnailURL, time.Now(), time.Now())
	return err
}

func GetTaskTracks(taskID int) ([]TaskTrack, error) {
	rows, err := DB.Query("SELECT task_id, video_id, status, title, artist, album, thumbnail_url, file_path, created_at, updated_at FROM task_tracks WHERE task_id = ? ORDER BY created_at ASC", taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tracks := []TaskTrack{}
	for rows.Next() {
		var track TaskTrack
		var thumb sql.NullString
		var filePath sql.NullString
		err := rows.Scan(&track.TaskID, &track.VideoID, &track.Status, &track.Title, &track.Artist, &track.Album, &thumb, &filePath, &track.CreatedAt, &track.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if thumb.Valid {
			track.ThumbnailURL = thumb.String
		}
		if filePath.Valid {
			track.FilePath = filePath.String
		}
		tracks = append(tracks, track)
	}
	return tracks, nil
}

func UpdateTaskTrackStatus(taskID int, videoID, status string) error {
	stmt, err := DB.Prepare("UPDATE task_tracks SET status = ?, updated_at = ? WHERE task_id = ? AND video_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, time.Now(), taskID, videoID)
	return err
}

func MarkTrackCompleted(taskID int, videoID, filePath string) error {
	stmt, err := DB.Prepare("UPDATE task_tracks SET status = 'completed', file_path = ?, updated_at = ? WHERE task_id = ? AND video_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(filePath, time.Now(), taskID, videoID)
	return err
}

func HasTaskTracks(taskID int) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM task_tracks WHERE task_id = ?", taskID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func ResetStuckTasks() error {
	// Reset tasks that are stuck in 'processing' state to 'pending' on server start
	_, err := DB.Exec("UPDATE tasks SET status = 'pending' WHERE status = 'processing'")
	return err
}

func AddTask(taskType, referenceID, payload string) error {
	stmt, err := DB.Prepare("INSERT INTO tasks(type, reference_id, status, payload, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(taskType, referenceID, "pending", payload, time.Now(), time.Now())
	return err
}

func GetPendingTask() (*Task, error) {
	row := DB.QueryRow("SELECT id, type, reference_id, status, payload, created_at, updated_at FROM tasks WHERE status = 'pending' ORDER BY created_at ASC LIMIT 1")

	var task Task
	err := row.Scan(&task.ID, &task.Type, &task.ReferenceID, &task.Status, &task.Payload, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func GetTask(id int) (*Task, error) {
	row := DB.QueryRow("SELECT id, type, reference_id, status, payload, created_at, updated_at FROM tasks WHERE id = ?", id)

	var task Task
	err := row.Scan(&task.ID, &task.Type, &task.ReferenceID, &task.Status, &task.Payload, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func GetTaskByReferenceID(referenceID string) (*Task, error) {
	row := DB.QueryRow("SELECT id, type, reference_id, status, payload, created_at, updated_at FROM tasks WHERE reference_id = ?", referenceID)

	var task Task
	err := row.Scan(&task.ID, &task.Type, &task.ReferenceID, &task.Status, &task.Payload, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateTaskStatus(id int, status string) error {
	stmt, err := DB.Prepare("UPDATE tasks SET status = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, time.Now(), id)
	return err
}

func RetryTask(taskID int) error {
	// Reset failed tracks to not_started
	_, err := DB.Exec("UPDATE task_tracks SET status = 'not_started', updated_at = ? WHERE task_id = ? AND status = 'failed'", time.Now(), taskID)
	if err != nil {
		return err
	}

	// Reset task status to pending so worker picks it up
	_, err = DB.Exec("UPDATE tasks SET status = 'pending', updated_at = ? WHERE id = ?", time.Now(), taskID)
	return err
}

func GetAllTasks() ([]Task, error) {
	rows, err := DB.Query("SELECT id, type, reference_id, status, payload, created_at, updated_at FROM tasks ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Type, &task.ReferenceID, &task.Status, &task.Payload, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetSetting(key string) (string, error) {
	var value string
	err := DB.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	return value, err
}

func SetSetting(key, value string) error {
	_, err := DB.Exec("INSERT OR REPLACE INTO settings(key, value) VALUES(?, ?)", key, value)
	return err
}

type TaskTrackWithPlaylist struct {
	TaskID       int
	VideoID      string
	Status       string
	Title        string
	Artist       string
	Album        string
	ThumbnailURL string
	FilePath     string
	PlaylistID   string
	PlaylistName string
}

func GetAllTaskTracksWithPlaylistInfo() ([]TaskTrackWithPlaylist, error) {
	query := `
		SELECT 
			tt.task_id, tt.video_id, tt.status, tt.title, tt.artist, tt.album, tt.thumbnail_url, tt.file_path,
			t.reference_id, t.payload
		FROM task_tracks tt
		JOIN tasks t ON tt.task_id = t.id
		ORDER BY tt.created_at DESC
	`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []TaskTrackWithPlaylist
	for rows.Next() {
		var t TaskTrackWithPlaylist
		var payload string
		var thumb sql.NullString
		var filePath sql.NullString
		err := rows.Scan(&t.TaskID, &t.VideoID, &t.Status, &t.Title, &t.Artist, &t.Album, &thumb, &filePath, &t.PlaylistID, &payload)
		if err != nil {
			return nil, err
		}
		if thumb.Valid {
			t.ThumbnailURL = thumb.String
		}
		if filePath.Valid {
			t.FilePath = filePath.String
		}

		var payloadMap map[string]interface{}
		if payload != "" {
			if err := json.Unmarshal([]byte(payload), &payloadMap); err == nil {
				if name, ok := payloadMap["playlistName"].(string); ok {
					t.PlaylistName = name
				}
			}
		}

		tracks = append(tracks, t)
	}
	return tracks, nil
}

type TaskWithStats struct {
	ID           int
	Type         string
	ReferenceID  string
	Status       string
	Payload      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	PlaylistName string
	TotalTracks  int
	Processed    int
	Failed       int
}

func GetAllTasksWithStats() ([]TaskWithStats, error) {
	query := `
		SELECT 
			t.id, t.type, t.reference_id, t.status, t.payload, t.created_at, t.updated_at,
			COUNT(tt.video_id) as total_tracks,
			SUM(CASE WHEN tt.status = 'completed' THEN 1 ELSE 0 END) as processed,
			SUM(CASE WHEN tt.status = 'failed' THEN 1 ELSE 0 END) as failed
		FROM tasks t
		LEFT JOIN task_tracks tt ON t.id = tt.task_id
		GROUP BY t.id
		ORDER BY t.created_at DESC
	`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []TaskWithStats
	for rows.Next() {
		var t TaskWithStats
		err := rows.Scan(&t.ID, &t.Type, &t.ReferenceID, &t.Status, &t.Payload, &t.CreatedAt, &t.UpdatedAt, &t.TotalTracks, &t.Processed, &t.Failed)
		if err != nil {
			return nil, err
		}

		if t.Payload != "" {
			var payloadMap map[string]interface{}
			if err := json.Unmarshal([]byte(t.Payload), &payloadMap); err == nil {
				if name, ok := payloadMap["playlistName"].(string); ok {
					t.PlaylistName = name
				}
			}
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
