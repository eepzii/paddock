package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func New() (*FileManager, error) {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	configDirPath := filepath.Join(configPath, APP_NAME)
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		return nil, err
	}
	configFilePath := filepath.Join(configDirPath, CONFIG)

	cachedPath, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	cachedDirPath := filepath.Join(cachedPath, APP_NAME, BROWSER_PROFILE)
	if err := os.MkdirAll(cachedDirPath, 0755); err != nil {
		return nil, err
	}

	return &FileManager{configFilePath: configFilePath, browserProfileDirPath: cachedDirPath}, nil
}

func (fm *FileManager) SaveConfig(config Config) error {
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(fm.configFilePath, configData, 0644); err != nil {
		return err
	}
	return nil
}

func (fm *FileManager) LoadConfig() (Config, error) {
	var config Config

	configData, err := os.ReadFile(fm.configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return config, err
	}

	if err := json.Unmarshal(configData, &config); err != nil {
		return config, err
	}

	return config, nil
}

func (fm *FileManager) Reset() error {
	if err := fm.SaveConfig(Config{}); err != nil {
		return err
	}
	if err := tryRemoveAll(fm.browserProfileDirPath); err != nil {
		return err
	}
	return nil
}

func (fm *FileManager) ConfigFilePath() string {
	return fm.configFilePath
}

func (fm *FileManager) BrowserProfileDirPath() string {
	return fm.browserProfileDirPath
}
