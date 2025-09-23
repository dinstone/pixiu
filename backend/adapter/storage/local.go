package storage

import (
	"os"
	"path"
)

// LocalStorage provides reading and writing application data to the user's
// configuration directory.
type LocalStorage struct {
	ConfPath string
}

// NewLocalStorage returns a localStore instance.
func NewLocalStorage(parentdir string, filename string) *LocalStorage {
	return &LocalStorage{
		ConfPath: path.Join(parentdir, filename),
	}
}

// Load reads the given file in the user's configuration directory and returns
// its contents.
func (l *LocalStorage) Load() ([]byte, error) {
	d, err := os.ReadFile(l.ConfPath)
	if err != nil {
		return nil, err
	}
	return d, err
}

// Store writes data to the user's configuration directory at the given
// filename.
func (l *LocalStorage) Store(data []byte) error {
	dir := path.Dir(l.ConfPath)
	if err := ensureDirExists(dir); err != nil {
		return err
	}
	if err := os.WriteFile(l.ConfPath, data, 0777); err != nil {
		return err
	}
	return nil
}

// ensureDirExists checks for the existence of the directory at the given path,
// which is created if it does not exist.
func ensureDirExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err = os.Mkdir(path, 0777); err != nil {
			return err
		}
	}
	return nil
}
