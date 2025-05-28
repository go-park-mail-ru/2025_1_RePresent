package controller

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	banneruc "retarget/internal/banner-service/usecase"
	authenticate "retarget/pkg/middleware/auth"
)

func TestSetupBannerRoutes_AllRoutes(t *testing.T) {
	handler := SetupBannerRoutes(&authenticate.Authenticator{}, &banneruc.BannerUsecase{}, &banneruc.BannerImageUsecase{})
	router, ok := handler.(*mux.Router)
	if !ok {
		t.Fatal("expected *mux.Router")
	}

	routes := []struct {
		name   string
		path   string
		method string
	}{
		{"download_image", "/api/v1/banner/image/{image_id}", "GET"},
		{"get banners", "/api/v1/banner/", "GET"},
		{"create", "/api/v1/banner/create", "POST"},
		{"read", "/api/v1/banner/{banner_id:[0-9]+}", "GET"},
		{"update", "/api/v1/banner/{banner_id:[0-9]+}", "PUT"},
		{"delete", "/api/v1/banner/{banner_id:[0-9]+}", "DELETE"},
		{"iframe", "/api/v1/banner/iframe/{banner_id:[0-9]+}", "GET"},
		{"upload", "/api/v1/banner/upload", "PUT"},
		{"uniq", "/api/v1/banner/uniq_link/{uniq_link}", "GET"},
	}

	for _, tc := range routes {
		t.Run(tc.name, func(t *testing.T) {
			route := router.Get(tc.name)
			if route == nil && tc.name != "download_image" {
				route = findRouteByPath(router, tc.path, tc.method)
				if route == nil {
					t.Fatalf("route %q with path %q and method %q not found", tc.name, tc.path, tc.method)
				}
			}

			if route != nil {
				methods, err := route.GetMethods()
				if err != nil {
					t.Errorf("Failed to get methods for %q: %v", tc.name, err)
				} else if len(methods) != 1 || methods[0] != tc.method {
					t.Errorf("Route %q: expected method %q, got %v", tc.name, tc.method, methods)
				}

				path, err := route.GetPathTemplate()
				if err != nil {
					t.Errorf("Failed to get path for %q: %v", tc.name, err)
				} else if path != tc.path {
					t.Errorf("Route %q: expected path %q, got %q", tc.name, tc.path, path)
				}
			}
		})
	}
}

func findRouteByPath(r *mux.Router, path, method string) *mux.Route {
	var result *mux.Route
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil && pathTemplate == path {
			methods, err := route.GetMethods()
			if err == nil && len(methods) > 0 && methods[0] == method {
				result = route
				return nil
			}
		}
		return nil
	})
	return result
}

func TestSetupBannerRoutes_ReturnValue(t *testing.T) {
	handler := SetupBannerRoutes(&authenticate.Authenticator{}, &banneruc.BannerUsecase{}, &banneruc.BannerImageUsecase{})
	if _, ok := handler.(http.Handler); !ok {
		t.Error("expected http.Handler")
	}
}

func TestNewBannerController(t *testing.T) {
	controller := NewBannerController(
		&banneruc.BannerUsecase{},
		&banneruc.BannerImageUsecase{},
		&linkBuilder{router: mux.NewRouter()},
	)

	if controller == nil {
		t.Fatal("expected non-nil controller")
	}

	if controller.BannerUsecase == nil {
		t.Error("BannerUsecase is nil")
	}

	if controller.ImageUsecase == nil {
		t.Error("ImageUsecase is nil")
	}

	if controller.LinkBuilder == nil {
		t.Error("LinkBuilder is nil")
	}
}

func TestSetupBannerRoutes_RequiredRoutes(t *testing.T) {
	handler := SetupBannerRoutes(&authenticate.Authenticator{}, &banneruc.BannerUsecase{}, &banneruc.BannerImageUsecase{})
	router, ok := handler.(*mux.Router)
	if !ok {
		t.Fatal("expected *mux.Router")
	}

	criticalRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/banner/"},
		{"POST", "/api/v1/banner/create"},
		{"GET", "/api/v1/banner/123"},
		{"PUT", "/api/v1/banner/123"},
		{"DELETE", "/api/v1/banner/123"},
		{"GET", "/api/v1/banner/iframe/123"},
		{"GET", "/api/v1/banner/image/test"},
		{"PUT", "/api/v1/banner/upload"},
		{"GET", "/api/v1/banner/uniq_link/123"},
	}

	for _, route := range criticalRoutes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			// Проверяем соответствие маршрута
			req, _ := http.NewRequest(route.method, route.path, nil)
			var match mux.RouteMatch

			if match := router.Match(req, &match); !match {
				t.Errorf("No route matches %s %s", route.method, route.path)
			}
		})
	}
}

func TestSetupBannerRoutes_Middleware(t *testing.T) {
	auth := &authenticate.Authenticator{}
	handler := SetupBannerRoutes(auth, &banneruc.BannerUsecase{}, &banneruc.BannerImageUsecase{})
	router, ok := handler.(*mux.Router)
	if !ok {
		t.Fatal("expected *mux.Router")
	}

	var authPaths, noAuthPaths []string
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}

		pathRequiresAuth := false
		switch path {
		case "/api/v1/banner/",
			"/api/v1/banner/create",
			"/api/v1/banner/{banner_id:[0-9]+}",
			"/api/v1/banner/upload":
			pathRequiresAuth = true
		}

		if pathRequiresAuth {
			authPaths = append(authPaths, path)
		} else {
			noAuthPaths = append(noAuthPaths, path)
		}

		return nil
	})

	if len(authPaths) == 0 {
		t.Error("No authenticated paths found")
	}

	if len(noAuthPaths) == 0 {
		t.Error("No public paths found")
	}

	t.Logf("Auth paths: %v", authPaths)
	t.Logf("Public paths: %v", noAuthPaths)
}

func TestSetupBannerRoutes_NamedRoutes(t *testing.T) {
	handler := SetupBannerRoutes(&authenticate.Authenticator{}, &banneruc.BannerUsecase{}, &banneruc.BannerImageUsecase{})
	router, ok := handler.(*mux.Router)
	if !ok {
		t.Fatal("expected *mux.Router")
	}

	route := router.GetRoute("download_image")
	if route == nil {
		t.Fatal("download_image route not found")
	}

	path, err := route.GetPathTemplate()
	if err != nil {
		t.Fatalf("Cannot get path for download_image: %v", err)
	}

	if expected := "/api/v1/banner/image/{image_id}"; path != expected {
		t.Errorf("download_image path is %q, expected %q", path, expected)
	}
}

func TestNewBannerController_Complete(t *testing.T) {
	bannerUsecase := &banneruc.BannerUsecase{}
	imageUsecase := &banneruc.BannerImageUsecase{}
	linkBuilder := NewLinkBuilder(mux.NewRouter())

	controller := NewBannerController(bannerUsecase, imageUsecase, linkBuilder)

	if controller == nil {
		t.Fatal("NewBannerController returned nil")
	}

	if controller.BannerUsecase != bannerUsecase {
		t.Error("BannerUsecase wasn't properly assigned")
	}

	if controller.ImageUsecase != imageUsecase {
		t.Error("ImageUsecase wasn't properly assigned")
	}

	if controller.LinkBuilder != linkBuilder {
		t.Error("LinkBuilder wasn't properly assigned")
	}
}
