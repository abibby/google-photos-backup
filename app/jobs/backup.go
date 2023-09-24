package jobs

import (
	"context"
	"errors"
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
	"github.com/abibby/google-photos-backup/services/iterable"
	"github.com/abibby/salusa/database/model"
	"modernc.org/sqlite"
)

var ErrAlreadyDownloaded = errors.New("already downloaded")

func BackupJob(e *events.BackupEvent) error {
	log.Print("starting backup")
	ctx := context.Background()

	users, err := models.UserQuery(ctx).Get(database.DB)
	if err != nil {
		return err
	}
	for _, u := range users {
		err = backupUser(u)
		if err != nil {
			log.Print(err)
		}
	}
	log.Print("finished backup")
	return nil
}
func backupUser(u *models.User) error {
	pf := iterable.NewPhotoFetcher(u)
	defer pf.Close()
	for pf.Next() {
		err := copyPhoto(u, pf.Value())
		if errors.Is(err, ErrAlreadyDownloaded) && u.FinishedInitialFetch {
			return nil
		} else if err != nil {
			log.Print(err)
			continue
		}
	}
	u.FinishedInitialFetch = true
	err := model.Save(database.DB, u)
	if err != nil {
		return err
	}
	return nil
}
func copyPhoto(u *models.User, item *gphotos.MediaItem) error {
	p, err := models.PhotoQuery(context.Background()).Where("photo_id", "=", item.ID).First(database.DB)
	if err != nil {
		return err
	}
	if p != nil {
		return ErrAlreadyDownloaded
	}
	url := fmt.Sprintf("%s=w%s-h%s", item.BaseURL, item.MediaMetadata.Height, item.MediaMetadata.Width)

	dir := path.Join(
		config.PhotoDir,
		u.Email,
		fmt.Sprintf("%02d/%02d", item.MediaMetadata.CreationTime.Year(), item.MediaMetadata.CreationTime.Month()),
	)
	err = os.MkdirAll(dir, 0755)
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
	if err != nil {
		return err
	}

	p = &models.Photo{
		UserID:  u.ID,
		PhotoID: item.ID,
	}

	err = model.Save(database.DB, p)
	var sqlErr *sqlite.Error
	if errors.As(err, &sqlErr) {
		if sqlErr.Code() == 2067 {
			return ErrAlreadyDownloaded
		}
	}
	if err != nil {
		return err
	}
	return nil
}
