package db

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// ...

var DB *gorm.DB

type GroupTask struct {
	ID               uint `gorm:"primaryKey"`
	Type             string
	ReferenceID      string `gorm:"uniqueIndex"`
	Status           string
	PlaylistName     string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Source           string `gorm:"default:user"`

	// Ignored fields for GORM, used for UI
	TotalTracks int `gorm:"-"`
	Processed   int `gorm:"-"`
	Failed      int `gorm:"-"`

	// New field for song download limit
	MaxTracks int `gorm:"default:0"`
}

type SongTask struct {
	GroupTaskID  uint   `gorm:"primaryKey;autoIncrement:false"`
	VideoID      string `gorm:"primaryKey"`
	Status       string
	Title        string
	Artist       string
	Album        string
	ThumbnailURL string
	FilePath     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Setting struct {
	Key   string `gorm:"primaryKey"`
	Value string
}

func InitDB() {
	var err error
	dbPath := os.Getenv("BEATBUMP_DB_PATH")
	dsn := "file:"+filepath.Join(dbPath, "beatbump.db")+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"

	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate
	err = DB.AutoMigrate(&GroupTask{}, &SongTask{}, &Setting{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	err = ResetStuckTasks()
	if err != nil {
		log.Printf("Failed to reset stuck tasks: %v", err)
	}
}

// Group Task Functions

func AddGroupTask(taskType, referenceID, playlistName, source string, maxTracks int) error {
	task := GroupTask{
		Type:         taskType,
		ReferenceID:  referenceID,
		Status:       TaskStatusPending,
		PlaylistName: playlistName,
		Source:       source,
		MaxTracks:    maxTracks,
	}
	return DB.Create(&task).Error
}

func GetGroupTask(id int) (*GroupTask, error) {
	var task GroupTask
	err := DB.First(&task, id).Error
	return &task, err
}

func GetGroupTaskByReferenceID(referenceID string) (*GroupTask, error) {
	var task GroupTask
	err := DB.Where("reference_id = ?", referenceID).First(&task).Error
	return &task, err
}

func GetActiveSessionGroupTask(timeout time.Duration) (*GroupTask, error) {
	var task GroupTask
	// Fetch the latest ongoing session task
	err := DB.Where("type = ? AND reference_id LIKE ?", TaskTypeOngoingDownload, "ongoing:songs:%").
		Order("updated_at DESC").
		First(&task).Error

	if err != nil {
		log.Printf("GetActiveSessionGroupTask: No recent session found (Error: %v)", err)
		return nil, err
	}

	// Check if it's within the timeout window
	if time.Since(task.UpdatedAt) > timeout {
		log.Printf("GetActiveSessionGroupTask: Found task %d but it's too old (Updated: %v, Threshold: %v)", task.ID, task.UpdatedAt, timeout)
		return nil, gorm.ErrRecordNotFound
	}

	log.Printf("GetActiveSessionGroupTask: Found active task %d (Updated: %v)", task.ID, task.UpdatedAt)
	return &task, nil
}

func GetPendingGroupTask() (*GroupTask, error) {
	var task GroupTask
	err := DB.Where("status = ? AND source = ?", TaskStatusPending, TaskSourceUser).
		Order("created_at ASC").
		First(&task).Error
	return &task, err
}

func GetAllGroupTasks() ([]GroupTask, error) {
	var tasks []GroupTask
	if err := DB.Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, err
	}

	// Optimized: Use SQL aggregation to get counts per group task
	type TaskCounts struct {
		GroupTaskID uint
		Total       int
		Processed   int
		Failed      int
	}

	var counts []TaskCounts
	err := DB.Model(&SongTask{}).
		Select(`
			group_task_id,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as processed,
			SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed
		`).
		Group("group_task_id").
		Scan(&counts).Error

	if err != nil {
		return nil, err
	}

	// Create a map for O(1) lookup
	countsMap := make(map[uint]TaskCounts)
	for _, c := range counts {
		countsMap[c.GroupTaskID] = c
	}

	// Assign counts to tasks
	for i := range tasks {
		if c, ok := countsMap[tasks[i].ID]; ok {
			tasks[i].TotalTracks = c.Total
			tasks[i].Processed = c.Processed
			tasks[i].Failed = c.Failed
		}
	}

	return tasks, nil
}

func UpdateGroupTaskStatus(id int, status string) error {
	return DB.Model(&GroupTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error
}

func RetryGroupTask(id int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Reset failed songs to not_started
		if err := tx.Model(&SongTask{}).
			Where("group_task_id = ? AND status = ?", id, TaskStatusFailed).
			Updates(map[string]interface{}{
				"status":     TaskStatusNotStarted,
				"updated_at": time.Now(),
			}).Error; err != nil {
			return err
		}

		// Reset group task status to pending
		if err := tx.Model(&GroupTask{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"status":     TaskStatusPending,
				"updated_at": time.Now(),
			}).Error; err != nil {
			return err
		}
		return nil
	})
}

// Song Task Functions

func AddSongTask(groupTaskID int, videoID, title, artist, album, thumbnailURL string) error {
	task := SongTask{
		GroupTaskID:  uint(groupTaskID),
		VideoID:      videoID,
		Status:       TaskStatusNotStarted,
		Title:        title,
		Artist:       artist,
		Album:        album,
		ThumbnailURL: thumbnailURL,
	}

	// Use Clauses to handle "INSERT OR IGNORE" equivalent (OnConflict Do Nothing)
	err := DB.Clauses(clause.OnConflict{DoNothing: true}).Omit("GroupTask").Create(&task).Error
	if err != nil {
		return err
	}

	// Update parent task updated_at
	DB.Model(&GroupTask{}).Where("id = ?", groupTaskID).Update("updated_at", time.Now())
	return nil
}

func GetSongTasks(groupTaskID int) ([]SongTask, error) {
	var songs []SongTask
	err := DB.Where("group_task_id = ?", groupTaskID).Order("created_at ASC").Find(&songs).Error
	return songs, err
}

func GetSongTask(groupTaskID int, videoID string) (*SongTask, error) {
	var task SongTask
	err := DB.Where("group_task_id = ? AND video_id = ?", groupTaskID, videoID).First(&task).Error
	return &task, err
}

func GetPendingSongTasks() ([]*SongTask, error) {
	var tasks []*SongTask

	// Join with GroupTask to prioritize user tasks
	// GORM Joins
	err := DB.Table("song_tasks").
		Select("song_tasks.*, group_tasks.source").
		Joins("JOIN group_tasks ON song_tasks.group_task_id = group_tasks.id").
		Where("song_tasks.status = ?", TaskStatusNotStarted).
		Where("group_tasks.status != ?", TaskStatusPaused).
		Order("CASE WHEN group_tasks.source = '" + TaskSourceUser + "' THEN 0 ELSE 1 END").
		Order("song_tasks.created_at ASC").
		Scan(&tasks).Error

	return tasks, err
}

func UpdateSongTaskStatus(groupTaskID int, videoID, status string) error {
	return DB.Model(&SongTask{}).
		Where("group_task_id = ? AND video_id = ?", groupTaskID, videoID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func MarkSongTaskCompleted(groupTaskID int, videoID, filePath string) error {
	return DB.Model(&SongTask{}).
		Where("group_task_id = ? AND video_id = ?", groupTaskID, videoID).
		Updates(map[string]interface{}{
			"status":     TaskStatusCompleted,
			"file_path":  filePath,
			"updated_at": time.Now(),
		}).Error
}

func CheckGroupCompletion(groupTaskID int) (bool, error) {
	var total int64
	var completed int64

	err := DB.Model(&SongTask{}).Where("group_task_id = ?", groupTaskID).Count(&total).Error
	if err != nil {
		return false, err
	}

	err = DB.Model(&SongTask{}).Where("group_task_id = ? AND status = ?", groupTaskID, TaskStatusCompleted).Count(&completed).Error
	if err != nil {
		return false, err
	}

	return total > 0 && total == completed, nil
}

// Helpers

func ResetStuckTasks() error {
	return DB.Model(&GroupTask{}).
		Where("status = ?", TaskStatusProcessing).
		Update("status", TaskStatusPending).Error
}

func GetSetting(key string) (string, error) {
	var setting Setting
	err := DB.First(&setting, "key = ?", key).Error
	return setting.Value, err
}

func SetSetting(key, value string) error {
	setting := Setting{Key: key, Value: value}
	return DB.Save(&setting).Error // Save handles Insert or Update (Upsert)
}
