package endpoints

import (
	"context"
	"net/http"

	"github.com/bluetuith-org/api-native/api/bluetooth"
	"github.com/danielgtaylor/huma/v2"
)

func mediaPlayerEndpoints(api huma.API, session bluetooth.Session) {
	mediaPlayerControlEndpoint(api, session)
	mediaPlayerPropertiesEndpoint(api, session)
}

func mediaPlayerPropertiesEndpoint(api huma.API, session bluetooth.Session) {
	type MediaPropertiesOutput struct {
		Body bluetooth.MediaData
	}

	huma.Register(api, huma.Operation{
		OperationID: "device-media-player-properties",
		Method:      http.MethodGet,
		Path:        "/device/{address}/media_player/properties",
		Summary:     "Properties",
		Description: "This endpoint sends a media control command to the device's media player, if available.",
		Tags:        []string{"Media Player"},
	}, func(_ context.Context, input *struct {
		AddressInput
	}) (*MediaPropertiesOutput, error) {
		var err error

		mediaCall := session.MediaPlayer(input.Address)
		properties, err := mediaCall.Properties()
		if err != nil {
			return nil, err
		}

		return &MediaPropertiesOutput{properties}, nil
	})
}

func mediaPlayerControlEndpoint(api huma.API, session bluetooth.Session) {
	type MediaControlInput struct {
		Control string `json:"control_type" path:"control_type" enum:"play,pause,next,previous,fast-forward,rewind,stop" doc:"The type of control command to send to the device's media player."`
	}

	huma.Register(api, huma.Operation{
		OperationID: "device-media-player-controls",
		Method:      http.MethodGet,
		Path:        "/device/{address}/media_player/control/{control_type}",
		Summary:     "Controls",
		Description: "This endpoint sends a media control command to the device's media player, if available.",
		Tags:        []string{"Media Player"},
	}, func(_ context.Context, input *struct {
		AddressInput
		MediaControlInput
	}) (*struct{}, error) {
		var err error

		mediaCall := session.MediaPlayer(input.Address)
		switch input.Control {
		case "play":
			err = mediaCall.Play()
		case "pause":
			err = mediaCall.Pause()
		case "next":
			err = mediaCall.Next()
		case "previous":
			err = mediaCall.Previous()
		case "fast-forward":
			err = mediaCall.FastForward()
		case "rewind":
			err = mediaCall.Rewind()
		case "stop":
			err = mediaCall.Stop()
		}

		return nil, err
	})
}
