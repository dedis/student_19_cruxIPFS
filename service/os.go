package service

import "os"

// default file mode for files that the program writes to the system
const defaultFileMode os.FileMode = 0777

// CreateEmptyDir create an empty directory at the given path
func CreateEmptyDir(path string) error {
	// remove existing dir if any
	os.RemoveAll(path)
	// create the empty directorys
	err := os.Mkdir(path, defaultFileMode)
	return err
}
