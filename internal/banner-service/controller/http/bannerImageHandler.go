package controller

import (
	"encoding/json"
	"io"
	"net/http"
	entity "retarget/pkg/entity"

	"github.com/gorilla/mux"
)

func (c *BannerController) UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	const maxFileSize int64 = 10 << 20

	if r.ContentLength > (maxFileSize) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Unsupported file size(max size 10MB): size your file is too large"))
		return
	}

	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid request"+err.Error()))
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Image not found in request"))
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to read image"))
		return
	}

	fileType := http.DetectContentType(buf)
	allowedTypes := map[string]struct{}{
		"image/jpeg": {},
		"image/jpg":  {},
		"image/png":  {},
		"image/gif":  {},
	}
	if _, ok := allowedTypes[fileType]; !ok {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Unsupported file type: upload only .png, .jpg, .jpeg or .gif files"))
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to seek image"))
		return
	}

	link, err := c.ImageUsecase.UploadBannerImage(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to upload image"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.NewResponse(false, link))
}

func (c *BannerController) DownloadImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := vars["image_id"]

	object, err := c.ImageUsecase.DownloadBannerImage(imageID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Image not found"))
		return
	}
	defer func() {
		if object != nil {
			object.Close()
		}
	}()

	buf := make([]byte, 512)
	_, err = object.Read(buf)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to read avatar: "+err.Error()))
		return
	}

	fileType := http.DetectContentType(buf)
	switch fileType {
	case "image/jpeg", "image/jpg":
		w.Header().Set("Content-Type", "image/jpeg")
	case "image/png":
		w.Header().Set("Content-Type", "image/png")
	case "image/gif":
		w.Header().Set("Content-Type", "image/gif")
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Unsupported file type"))
		return
	}

	_, err = object.Seek(0, io.SeekStart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to seek image"))
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=image")
	_, err = io.Copy(w, object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to download image"))
		return
	}
}
