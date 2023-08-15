package main

import (
	"backend/common/integrations/auth"
	"backend/common/server"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"

	"github.com/gofrs/uuid"
)

const FileSizeLimit = 20 * 1024 * 1024 // 20MB

var ImageTypeRegex = regexp.MustCompile(`^image/(png|jpe?g|gif)$`)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		server.ErrorResponse(w, http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("forum-token")
	if err != nil {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}

	userId := auth.CheckSession(cookie.Value)
	if userId == -1 {
		server.ErrorResponse(w, http.StatusUnauthorized)
		return
	}

	if r.ContentLength > FileSizeLimit {
		fmt.Println("return error")
		http.Error(w, fmt.Sprintf("Attachment is too big, max size is %vMB", FileSizeLimit/1024/1024), http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, FileSizeLimit)
	err = r.ParseMultipartForm(FileSizeLimit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Attachment is too big, max size is %vMB", FileSizeLimit/1024/1024), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		server.ErrorResponse(w, http.StatusBadRequest)
		return
	}

	contentType := header.Header.Get("Content-Type")

	if !ImageTypeRegex.MatchString(contentType) {
		http.Error(w, "Invalid file type, only images are supported", http.StatusBadRequest)
		return
	}

	defer file.Close()

	UUID, err := uuid.NewV4()
	if err != nil {
		server.ErrorResponse(w, http.StatusInternalServerError)
		return
	}

	fileName := UUID.String()

	dst, err := os.Create(path.Join(*attachmentsDir, fileName))
	if err != nil {
		server.ErrorResponse(w, http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		server.ErrorResponse(w, http.StatusInternalServerError)
		return
	}

	server.SendObject(w, "/api/attachments/"+fileName)
}
