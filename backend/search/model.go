package search

import (
	"time"

	"github.com/inokone/photostorage/collection"
	"github.com/inokone/photostorage/photo"
)

// QuickSearchResp is a JSON type fro results of sitewide quick search
type QuickSearchResp struct {
	Query   string                `json:"query"`
	Photos  []photo.Response      `json:"photos"`
	Albums  []collection.ListResp `json:"albums"`
	Uploads []collection.ListResp `json:"uploads"`
}

// Query is a JSON type for filter query on photos
type Query struct {
	UploadFrom time.Time `json:"upload_from"`
	UploadTo   time.Time `json:"upload_to"`
	TakenFrom  time.Time `json:"taken_from"`
	TakenTo    time.Time `json:"taken_to"`
	WidthFrom  int       `json:"width_from"`
	WidthTo    int       `json:"width_to"`
	HeightFrom int       `json:"height_from"`
	HeightTo   int       `json:"height_to"`
	SizeFrom   int       `json:"size_from"`
	SizeTo     int       `json:"size_to"`
	Name       string    `json:"name"`
	OrderBy    string    `json:"order_by"`
}

// Result is a JSON type for query results on photos
type Result struct {
	Query  Query            `json:"query"`
	Photos []photo.Response `json:"photos"`
}
