package iterable

import (
	"log"

	"github.com/abibby/google-photos-backup/app/models"
	"github.com/abibby/google-photos-backup/services/gphotos"
)

type PhotoFetcher struct {
	responses chan *gphotos.ListMediaItemsResponse
	// mtx      *sync.RWMutex
	photos []*gphotos.MediaItem
	index  int
	user   *models.User
	value  *gphotos.MediaItem
	open   bool
	done   chan struct{}
}

func NewPhotoFetcher(u *models.User) *PhotoFetcher {
	p := &PhotoFetcher{
		responses: make(chan *gphotos.ListMediaItemsResponse, 1),
		user:      u,
		open:      true,
		done:      make(chan struct{}),
	}
	go p.fetch()
	return p
}

func (p *PhotoFetcher) Close() error {
	p.open = false
	p.done <- struct{}{}
	close(p.done)
	return nil
}
func (p *PhotoFetcher) Next() bool {

	p.index++

	if p.index >= len(p.photos) || p.photos == nil {
		resp, ok := <-p.responses
		if !ok {
			return false
		}

		p.photos = resp.MediaItems
		p.index = 0
	}

	p.value = p.photos[p.index]
	return true
}

func (p *PhotoFetcher) Value() *gphotos.MediaItem {
	return p.value
}

func (p *PhotoFetcher) fetch() error {
	c := gphotos.NewClient(p.user)
	req := &gphotos.ListMediaItemsRequest{
		PageSize: 100,
	}
	total := 0
	defer func() {
		close(p.responses)
	}()
	for p.open {
		items, err := c.ListMediaItems(req)
		if err != nil {
			return err
		}

		total += len(items.MediaItems)
		log.Printf("%d items total with %d new", total, len(items.MediaItems))
		select {
		case p.responses <- items:
		case <-p.done:
			return nil
		}
		// p.push(items.MediaItems)
		p.responses <- items

		if items.NextPageToken == "" {
			return nil
		}
		req.PageToken = items.NextPageToken
	}

	return nil
}
