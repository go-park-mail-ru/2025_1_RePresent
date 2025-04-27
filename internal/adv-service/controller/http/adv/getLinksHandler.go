package adv

import (
	"encoding/json"
	"fmt"
	"net/http"
	entity "retarget/pkg/entity"
	response "retarget/pkg/entity"
)

func (c *AdvController) GetLinksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	}
	// requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	links, err := c.advUsecase.GetLinks(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponseWithBody(true, err.Error(), nil))
		return
	}

	type LinkInfo struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	}

	linksResponse := make([]LinkInfo, len(links))
	for i, link := range links {
		linksResponse[i] = LinkInfo{
			URL:    fmt.Sprintf("http://re-target.ru/api/v1/adv/iframe/%s", link.TextLink),
			Height: link.Height,
			Width:  link.Width,
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.NewResponseWithBody(false, "Links retrieved successfully",
		map[string]interface{}{
			"links": linksResponse,
		},
	))
}
