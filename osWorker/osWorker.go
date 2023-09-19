package osWorker

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

type OsWorker struct {
	Params map[string]any
}

func (ow *OsWorker) LsCommand(dir string) (fs []string, ds []string, err error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("Reading directory %v is failed. Error: %w.", dir, err)
	}

	for _, file := range files {
		if file.IsDir() {
			ds = append(ds, file.Name())
		} else {
			fs = append(fs, file.Name())
		}
	}
	return
}

func (ow *OsWorker) CreateDirCommand(dir string, name string) (err error) {
	_, err = os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Reading directory %v is failed. Error: %w.", dir, err)
	}

	err = os.Mkdir(name, 0755)
	if err != nil {
		err = fmt.Errorf("Creating directory %v is failed. Error: %w.", name, err)
	}
	return
}

func (ow *OsWorker) DeleteCommand(dir string, name string) (err error) {
	fi, err := os.Stat(name)
	if err != nil {
		return fmt.Errorf("Reading element %v is failed. Error: %w.", name, err)
	}

	if fi.IsDir() {
		err = os.RemoveAll(name)
		if err != nil {
			return fmt.Errorf("Deleting element %v is failed. Error: %w.", name, err)
		}
	} else {
		err = os.Remove(name)
		if err != nil {
			return fmt.Errorf("Deleting element %v is failed. Error: %w.", name, err)
		}
	}

	return
}

func (ow *OsWorker) RenameCommand(dir string, oldName string, newName string) (err error) {
	_, err = os.Stat(oldName)
	if err != nil {
		return fmt.Errorf("Reading element %v is failed. Error: %w.", oldName, err)
	}

	err = os.Rename(oldName, newName)
	if err != nil {
		return fmt.Errorf("Renaming element %v is failed. Error: %w.", newName, err)
	}

	return
}

func (ow *OsWorker) FileDowlCommand(path string) (fileBytes []byte, err error) {
	fileBytes, err = os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Reading file %v is failed. Error: %w.", path, err)
	}

	return
}

func (ow *OsWorker) FileLoadCommand(file *multipart.File, path string) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Creating uploaded file %v is failed. Error: %w.", path, err)
	}
	defer f.Close()
	_, err = io.Copy(f, *file)
	if err != nil {
		return fmt.Errorf("Creating uploaded file %v is failed. Error: %w.", path, err)
	}
	return
}
