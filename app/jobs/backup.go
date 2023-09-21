package jobs

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/abibby/google-photos-backup/app/events"
	"github.com/abibby/google-photos-backup/app/models"
	"github.com/abibby/google-photos-backup/config"
	"github.com/abibby/google-photos-backup/database"
	"github.com/abibby/google-photos-backup/services/gphotos"
)

func BackupJob(e *events.BackupEvent) error {
	ctx := context.Background()

	users, err := models.UserQuery(ctx).Get(database.DB)
	if err != nil {
		return err
	}
	for _, u := range users {
		c := gphotos.NewClient(u)
		items, err := c.ListMediaItems()
		if err != nil {
			return err
		}

		for _, item := range items.MediaItems {
			err = copyPhoto(u, item)
			if err != nil {
				log.Print(err)
			}
		}
	}
	return nil
}

func copyPhoto(u *models.User, item *gphotos.MediaItem) error {
	url := fmt.Sprintf("%s=w%s-h%s", item.BaseURL, item.MediaMetadata.Height, item.MediaMetadata.Width)

	dir := path.Join(
		config.PhotoDir,
		u.Email,
		fmt.Sprintf("%02d/%02d", item.MediaMetadata.CreationTime.Year(), item.MediaMetadata.CreationTime.Month()),
	)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path.Join(dir, item.Filename), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
