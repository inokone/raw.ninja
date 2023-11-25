package image

// Writer is an interface for changing images (RAW or processed).
type Writer interface {
	Store(id string, raw []byte, thumbnail []byte) error

	Delete(id string) error
}

// Loader is an interface for loading images (RAW or processed).
type Loader interface {
	LoadThumbnail(id string) ([]byte, error)

	LoadImage(id string) ([]byte, error)
}

// Storer is an interface for types that can store images (RAW or processed).
type Storer interface {
	Writer

	Loader
}
