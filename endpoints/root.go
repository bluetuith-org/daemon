package endpoints

import (
	"context"
	"net/http"

	"github.com/bluetuith-org/api-native/api/bluetooth"
	"github.com/danielgtaylor/huma/v2"
)

func rootEndpoints(api huma.API, session bluetooth.Session) {
	eventsEndpoint(api)
	authEndpoint(api)

	adaptersEndpoint(api, session)
}

func adaptersEndpoint(api huma.API, session bluetooth.Session) {
	type AdaptersOutput struct {
		Body []bluetooth.AdapterData
	}

	huma.Register(api, huma.Operation{
		OperationID: "adapters",
		Method:      http.MethodGet,
		Path:        "/adapters",
		Summary:     "Adapters",
		Description: "This endpoint fetches all available adapters.",
	}, func(_ context.Context, input *struct{}) (*AdaptersOutput, error) {
		return &AdaptersOutput{session.Adapters()}, nil
	})
}
