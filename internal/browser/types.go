package browser

import (
	"github.com/eepzii/paddock/internal/storage"
)

type Browser struct {
	commandPath string
	profilePath string
	headless    bool
	brand       Brand
}

type Config struct {
	FileManager       *storage.FileManager
	CustomBrowserPath string
	Headless          bool
}
