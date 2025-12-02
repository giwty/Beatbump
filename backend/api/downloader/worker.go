package downloader

import (
	"beatbump-server/backend/db"
	"log"
	"math/rand"
	"time"
)

func StartWorker() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// 1. Prioritize User Group Tasks (Playlists)
			// We only pick up tasks that are pending and source='user'
			groupTask, err := db.GetPendingGroupTask()
			if err == nil && groupTask != nil {
				// Mark as processing
				db.UpdateGroupTaskStatus(int(groupTask.ID), db.TaskStatusProcessing)

				if groupTask != nil {
					log.Printf("Processing group task %d: %s", groupTask.ID, groupTask.Type)
					if groupTask.Type == db.TaskTypePlaylistDownload {
						PopulateGroupTask(groupTask.ReferenceID, int(groupTask.ID))
					} else if groupTask.Type == db.TaskTypeSongMixDownload {
						PopulateSongMixTask(int(groupTask.ID))
					} else {
						log.Printf("Unknown task type: %s", groupTask.Type)
						db.UpdateGroupTaskStatus(int(groupTask.ID), db.TaskStatusFailed)
					}
				}
				continue
			}

			// 2. Process Individual Song Tasks (from any source, prioritized by user)
			songTasks, err := db.GetPendingSongTasks()
			if err == nil && len(songTasks) > 0 {
				// Process concurrent downloads
				// Limit concurrency to avoid rate limiting or system overload
				concurrencyLimit := 1
				sem := make(chan struct{}, concurrencyLimit)

				for _, songTask := range songTasks {
					sem <- struct{}{} // Acquire semaphore
					go func(task *db.SongTask) {
						defer func() { <-sem }() // Release semaphore
						HandleSongTask(task)
						time.Sleep(time.Duration(rand.Intn(8)) * time.Second)
					}(songTask)
				}

				// Wait for batch to finish
				for i := 0; i < concurrencyLimit; i++ {
					sem <- struct{}{}
				}
			}
		}
	}()
}
