package avatar

import (
	"encoding/json"
	"io"
	"net/http"
	entity "retarget/pkg/entity"
)

func (c *AvatarController) UploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	}

	user, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
	}
	userID := user.UserID

	if r.ContentLength > (10 << 20) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Unsupported file size(max size 10MB): size your file is too large"))
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid request"))
		return
	}

	file, _, err := r.FormFile("avatar")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Avatar not found in request"))
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to read avatar"))
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
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to seek avatar"))
		return
	}

	err = c.avatarUsecase.UploadAvatar(userID, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to upload avatar"))
		c.logger.Errorw("Failed to upload avatar", "err", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.NewResponse(false, "Got and Saved"))
}
