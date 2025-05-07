package avatar

import (
	"encoding/json"
	"io"
	"net/http"
	entity "retarget/pkg/entity"
)

func (c *AvatarController) DownloadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(entity.СtxKeyRequestID{}).(string)
	if r.Method != http.MethodGet {
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

	object, err := c.avatarUsecase.DownloadAvatar(userID, requestID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Avatar not found"))
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
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to seek avatar"))
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=avatar")
	_, err = io.Copy(w, object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to download avatar"))
		return
	}
}
