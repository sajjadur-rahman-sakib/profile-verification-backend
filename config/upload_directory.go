package config

import "os"

func UploadDirectory() {
	configuration := GetConfig()
	uploadDir := configuration.UploadDirectory
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}
}
