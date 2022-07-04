package controller

import (
	"fmt"
	"net/http"

	"github.com/OhanaFS/ohana/service"
	"github.com/gorilla/mux"
)

type UploadController struct {
	service service.UploadService
}

func RegisterUpload(r *mux.Router, service service.UploadService) {
	s := &UploadController{service}
	// Upload route
	r.HandleFunc("/v1/upload", s.uploadHandler)
}

func (s *UploadController) uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET LOST"))
	if r.Method == "GET" {
		fmt.Print(w, "upload", nil)
	} else {
		// Just for now
		var max_upload_size int64 = 10 * 1024 * 1024
		r.ParseMultipartForm(max_upload_size)

		// Get handler for filename, size and headers
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		// Create a new file
		message, err := s.service.CreateFile(file, handler.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Print("Message: ", message)
	}
}
