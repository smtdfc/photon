package photon


import (
	"os"
)

func (f *File) SaveToDisk(path string) error {
	fullPath := path + "/" + f.Filename
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(f.Content)
	if err != nil {
		return err
	}

	return nil
}