package controller

import (
	"github.com/gorilla/mux"
)

type LinkBuilder interface {
	BannerImageURL(imageID string) (string, error)
}

type linkBuilder struct {
	router *mux.Router
}

func NewLinkBuilder(router *mux.Router) LinkBuilder {
	return &linkBuilder{router: router}
}

func (b *linkBuilder) BannerImageURL(imageID string) (string, error) {
	url, err := b.router.Get("download_image").URL("image_id", imageID)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
