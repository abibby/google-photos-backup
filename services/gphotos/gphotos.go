package gphotos

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type PhotoMetadata struct {
	CameraMake      string  `json:"cameraMake"`
	CameraModel     string  `json:"cameraModel"`
	FocalLength     float32 `json:"focalLength"`
	ApertureFNumber float32 `json:"apertureFNumber"`
	IsoEquivalent   int     `json:"isoEquivalent"`
	ExposureTime    string  `json:"exposureTime"`
}

type MediaMetadata struct {
	CreationTime time.Time      `json:"creationTime"`
	Width        string         `json:"width"`
	Height       string         `json:"height"`
	Photo        *PhotoMetadata `json:"photo"`
}

type MediaItem struct {
	ID            string         `json:"id"`
	ProductURL    string         `json:"productUrl"`
	BaseURL       string         `json:"baseUrl"`
	MimeType      string         `json:"mimeType"`
	MediaMetadata *MediaMetadata `json:"mediaMetadata"`
	Filename      string         `json:"filename"`
}

type ListMediaItemsRequest struct {
	PageToken string
	PageSize  int
}

type ListMediaItemsResponse struct {
	MediaItems    []*MediaItem `json:"mediaItems"`
	NextPageToken string       `json:"nextPageToken"`
}

func (c *Client) ListMediaItems(req *ListMediaItemsRequest) (*ListMediaItemsResponse, error) {
	if req == nil {
		req = &ListMediaItemsRequest{}
	}
	query := url.Values{}
	if req.PageToken != "" {
		query.Add("pageToken", req.PageToken)
	}
	if req.PageSize != 0 {
		query.Add("pageSize", fmt.Sprint(req.PageSize))
	}
	resp, err := c.Get("https://photoslibrary.googleapis.com/v1/mediaItems?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		gerr := &GError{}
		err = json.NewDecoder(resp.Body).Decode(gerr)
		if err != nil {
			return nil, err
		}
		return nil, gerr
	}

	r := &ListMediaItemsResponse{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
