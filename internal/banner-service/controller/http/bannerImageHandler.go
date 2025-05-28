package controller

import (
	"io"
	"net/http"
	response "retarget/pkg/entity"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
)

func (c *BannerController) UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	const maxFileSize int64 = 10 << 20

	if r.ContentLength > (maxFileSize) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Unsupported file size(max size 10MB): size your file is too large"))
		resp := response.NewResponse(true, "Unsupported file size(max size 10MB): size your file is too large")
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid request"+err.Error()))
		resp := response.NewResponse(true, "Invalid request"+err.Error())
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Image not found in request"))
		resp := response.NewResponse(true, "Image not found in request")
		easyjson.MarshalToWriter(&resp, w)
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to read image"))
		resp := response.NewResponse(true, "Failed to read image")
		easyjson.MarshalToWriter(&resp, w)
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
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Unsupported file type: upload only .png, .jpg, .jpeg or .gif files"))
		resp := response.NewResponse(true, "Unsupported file type: upload only .png, .jpg, .jpeg or .gif files")
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to seek image"))
		resp := response.NewResponse(true, "Failed to seek image")
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	link, err := c.ImageUsecase.UploadBannerImage(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to upload image"))
		resp := response.NewResponse(true, "Failed to upload image")
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	// json.NewEncoder(w).Encode(entity.NewResponse(false, link))
	resp := response.NewResponse(false, link)
	easyjson.MarshalToWriter(&resp, w)
}

func (c *BannerController) DownloadImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageID := vars["image_id"]

	object, err := c.ImageUsecase.DownloadBannerImage(imageID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Image not found"))
		resp := response.NewResponse(true, "Image not found")
		easyjson.MarshalToWriter(&resp, w)
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
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to read image: "+err.Error()))
		resp := response.NewResponse(true, "Failed to read image: "+err.Error())
		easyjson.MarshalToWriter(&resp, w)
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
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Unsupported file type"))
		resp := response.NewResponse(true, "Unsupported file type")
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	_, err = object.Seek(0, io.SeekStart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to seek image"))
		resp := response.NewResponse(true, "Failed to seek image")
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=image")
	_, err = io.Copy(w, object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to download image"))
		resp := response.NewResponse(true, "Failed to download image")
		easyjson.MarshalToWriter(&resp, w)
		return
	}
}
