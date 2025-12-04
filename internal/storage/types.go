package storage

import "time"

type FileManager struct {
	configFilePath        string
	browserProfileDirPath string
}

type Config struct {
	Email             string    `json:"email"`
	SubscriptionToken string    `json:"subscriptionToken"`
	TokenExpiration   time.Time `json:"tokenExpiration"`
}
