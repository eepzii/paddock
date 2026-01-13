package browser

import (
	"github.com/eepzii/paddock/internal/storage"
)

type Browser struct {
	commandPath string
	profilePath string
	headless    bool
	brand       Brand
	proxy       Proxy
}

type Config struct {
	FileManager       *storage.FileManager
	CustomBrowserPath string
	Headless          bool
	Proxy             Proxy
}

type Proxy struct {
	Address  string
	User     string
	Password string
}
