package utils

import (
	"fmt"
	"os"
)

const (
	STORAGE = "storage"           // where database will live
	MEDIA   = "web/static/media"  // all post media
	IMAGES  = "web/static/images" // all user profile images
)

/*
* CreateStorageFolder/0 - creates the folder where the database will live
 */
func CreatStorageFolder() error {
	// check if the folder exists
	_, err := os.Stat(STORAGE)
	if os.IsNotExist(err) {
		// the folder doesn't exist, so create it
		err := os.MkdirAll(STORAGE, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create '%s' directory: %w", STORAGE, err)
		}
	}
	// folder exists, do nothing, return nil
	return nil
}

/*
* CreateMediaFolder/0 - creates the folder for media files
 */
func CreatMediaFolder() error {
	// check if the folder exists
	_, err := os.Stat(MEDIA)
	if os.IsNotExist(err) {
		// the folder doesn't exist, so create it
		err := os.MkdirAll(MEDIA, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create '%s' directory: %w", MEDIA, err)
		}
	}
	// folder exists, do nothing, return nil
	return nil
}

/*
* CreateImagesFolder - creates the folder to hold all the user profile images folder
 */
func CreatImagesFolder() error {
	// check if the folder exists
	_, err := os.Stat(IMAGES)
	if os.IsNotExist(err) {
		// the folder doesn't exist, so create it
		err := os.MkdirAll(IMAGES, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create '%s' directory: %w", IMAGES, err)
		}
	}
	// folder exists, do nothing, return nil
	return nil
}
