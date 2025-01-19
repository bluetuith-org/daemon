package endpoints

import (
	"net/http"

	ac "github.com/bluetuith-org/api-native/api/appcapability"
	bluetooth "github.com/bluetuith-org/api-native/api/bluetooth"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func Register(router *http.ServeMux, session bluetooth.Session, collection ac.Collection) huma.API {
	api := humago.New(router, huma.DefaultConfig("My API", "1.0.0"))
	api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		ctx.SetHeader("Retry-After", "10")
		next(ctx)
	})
	api.OpenAPI().Info = &huma.Info{
		Title:       "My API",
		Description: "# Description",
		Contact: &huma.Contact{
			Name:  "Contact",
			URL:   "https://contact.com",
			Email: "",
		},
		License: &huma.License{
			Name:       "MIT",
			Identifier: "",
			URL:        "",
		},
		Version: "",
	}


	rootEndpoints(api, session)
	adapterEndpoints(api, session)
	deviceEndpoints(api, session)

	if collection.Has(ac.CapabilitySendFile, ac.CapabilityReceiveFile) {
		obexEndpoints(api, session)
	}

	if collection.Has(ac.CapabilityNetwork) {
		networkEndpoints(api, session)
	}

	if collection.Has(ac.CapabilityMediaPlayer) {
		mediaPlayerEndpoints(api, session)
	}

	return api
}
