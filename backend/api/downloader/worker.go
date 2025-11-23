package downloader

import (
	"beatbump-server/backend/db"
	"database/sql"
	"log"
	"time"
)

func StartWorker() {
	go func() {
		for {
			task, err := db.GetPendingTask()
			if err != nil {
				if err != sql.ErrNoRows {
					log.Printf("Error fetching pending task: %v", err)
				}
				// No pending tasks or DB error
				time.Sleep(5 * time.Second)
				continue
			}

			if task != nil {
				log.Printf("Processing task %d: %s", task.ID, task.Type)
				if task.Type == "playlist_download" || task.Type == "ongoing_download" {
					DownloadPlaylist(task.ReferenceID, task.ID)
				} else {
					log.Printf("Unknown task type: %s", task.Type)
					db.UpdateTaskStatus(task.ID, "failed")
				}
			}
		}
	}()
}
