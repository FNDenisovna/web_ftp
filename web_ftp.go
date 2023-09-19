package main

import (
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"sync"
	osWorker "web_ftp/osWorker"
)

var mu sync.Mutex
var ow IOsWorker

const port = "8080"
const rootDir = `D:\`

type IOsWorker interface {
	LsCommand(dir string) (fs []string, ds []string, err error)
	CreateDirCommand(dir string, name string) error
	DeleteCommand(dir string, name string) error
	RenameCommand(dir string, oldName string, newName string) error
	FileDowlCommand(path string) (fileBytes []byte, err error)
	FileLoadCommand(file *multipart.File, path string) (err error)
}

func main() {
	mux := http.NewServeMux()
	log.Printf("Server started on localhost on port %v.\n", port)
	log.Printf("Let use it with url \"http://localhost:8080/ls\" in your browser.\n")

	ow = &osWorker.OsWorker{}
	//mux.HandleFunc("/", ls)
	mux.HandleFunc("/ls", ls)
	mux.HandleFunc("/create_dir", createDir)
	mux.HandleFunc("/delete", delete)
	mux.HandleFunc("/rename", rename)
	mux.HandleFunc("/command", command)
	mux.HandleFunc("/file", file)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// Обработчик, возвращающий содержимое директории.
func ls(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL.Path = %q\n", r.URL.Path)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: " + err.Error()))
		log.Printf("error: %s\n", err)
		return
	}

	customDir := r.FormValue("dir")
	log.Printf("parameter from url = %s\n", customDir)
	fullDir := filepath.Join(rootDir, customDir)
	log.Printf("customDir = %s\n", fullDir)

	mu.Lock()
	fs, ds, err := ow.LsCommand(fullDir)
	mu.Unlock()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: " + err.Error()))
		log.Printf("error: %s\n", err)
		return
	}

	w.Write([]byte(getListDirView(fs, ds, fullDir, customDir)))
}

// Обработчик нажатия кнопки и принятия решения о команде.
func command(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL.Path = %q\n", r.URL.Path)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: " + err.Error()))
		log.Printf("error: %s\n", err)
		return
	}

	_, fh, _ := r.FormFile("file")
	if fh != nil {
		uploadFile(w, r)
		return
	}

	name := r.FormValue("name")
	newName := r.FormValue("new_name")

	if newName != "" {
		if name != "" {
			rename(w, r)
		} else {
			createDir(w, r)
		}
	} else if name != "" {
		delete(w, r)
	} else {
		ls(w, r)
	}
}

// Обработчик, возвращающий содержимое директории.
func createDir(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL.Path = %q\n", r.URL.Path)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: " + err.Error()))
		log.Printf("error: %s\n", err)
		return
	}

	customDir := r.FormValue("dir")
	newName := r.FormValue("new_name")
	log.Printf("parameter from url = %s\n", customDir)
	fullDir := filepath.Join(rootDir, customDir)
	log.Printf("customDir = %s\n", fullDir)
	newDir := filepath.Join(fullDir, newName)
	log.Printf("newDir = %s\n", newDir)

	mu.Lock()
	err = ow.CreateDirCommand(fullDir, newDir)
	mu.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: " + err.Error()))
		log.Printf("error: %s\n", err)
		return
	}

	ls(w, r)
}

// Обработчик, удаляющий директорию или файл.
func delete(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL.Path = %q\n", r.URL.Path)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: " + err.Error()))
		log.Printf("error: %s\n", err)
		return
	}

	customDir := r.FormValue("dir")
	delName := r.FormValue("name")
	log.Printf("parameter from url = %s\n", customDir)
	fullDir := filepath.Join(rootDir, customDir)
	log.Printf("customDir = %s\n", fullDir)
	delDir := filepath.Join(fullDir, delName)
	log.Printf("delDir = %s\n", delDir)

	mu.Lock()
	err = ow.DeleteCommand(fullDir, delDir)
	mu.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: " + err.Error()))
		log.Printf("error: %s\n", err)
		return
	}

	ls(w, r)
}

// Обработчик, переименовывающий директорию или файл.
func rename(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL.Path = %q\n", r.URL.Path)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: " + err.Error()))
		log.Printf("error: %s\n", err)
		return
	}

	customDir := r.FormValue("dir")
	oldName := r.FormValue("name")
	newName := r.FormValue("new_name")
	log.Printf("parameter from url = %s\n", customDir)
	fullDir := filepath.Join(rootDir, customDir)
	log.Printf("customDir = %s\n", fullDir)
	oldDir := filepath.Join(fullDir, oldName)
	newDir := filepath.Join(fullDir, newName)
	log.Printf("oldDir = %s, newDir = %s\n", oldDir, newDir)

	mu.Lock()
	err = ow.RenameCommand(fullDir, oldDir, newDir)
	mu.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: " + err.Error()))
		log.Printf("error: %s\n", err)
		return
	}

	ls(w, r)
}

// Обработчик, передающий на скачку файл.
func file(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL.Path = %q\n", r.URL.Path)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: " + err.Error()))
		log.Printf("error: %s\n", err)
		return
	}

	customDir := r.FormValue("dir")
	name := r.FormValue("name")

	log.Printf("parameter from url = %s\n", customDir)
	fullDir := filepath.Join(rootDir, customDir)
	log.Printf("customDir = %s\n", fullDir)
	file := filepath.Join(fullDir, name)
	log.Printf("file = %s\n", file)

	mu.Lock()
	fileBytes, err := ow.FileDowlCommand(file)
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment")
	w.Write(fileBytes)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL.Path = %q\n", r.URL.Path)

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: " + err.Error()))
		log.Printf("error: %s\n", err)
		return
	}

	customDir := r.FormValue("dir")
	log.Printf("parameter from url = %s\n", customDir)
	fullDir := filepath.Join(rootDir, customDir)
	log.Printf("customDir = %s\n", fullDir)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: Uploading File from web is failed. " + err.Error()))
		log.Printf("Uploading File from web error: %s\n", err)
		return
	}
	defer file.Close()

	fileDest := filepath.Join(fullDir, fileHeader.Filename)
	log.Printf("fileDest = %s\n", fileDest)

	mu.Lock()
	err = ow.FileLoadCommand(&file, fileDest)
	mu.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: " + err.Error()))
		log.Printf("Creating file %v error: %s\n", fileDest, err)
		return
	}

	ls(w, r)
}
