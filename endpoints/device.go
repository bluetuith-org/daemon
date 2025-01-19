package endpoints

import (
	"context"
	"net/http"

	"github.com/bluetuith-org/api-native/api/bluetooth"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

func deviceEndpoints(api huma.API, session bluetooth.Session) {
	connectEndpoint(api, session)
	disconnectEndpoint(api, session)
	pairEndpoint(api, session)
	removeEndpoint(api, session)
	devicePropertiesEndpoint(api, session)
}

func devicePropertiesEndpoint(api huma.API, session bluetooth.Session) {
	type DevicePropertiesOutput struct {
		Body bluetooth.DeviceData
	}

	huma.Register(api, huma.Operation{
		OperationID: "device-properties",
		Method:      http.MethodGet,
		Path:        "/device/{address}/properties",
		Summary:     "Properties",
		Description: "This endpoint fetches the properties of the device.",
		Tags:        []string{"Device"},
	}, func(_ context.Context, input *struct {
		AddressInput
	}) (*DevicePropertiesOutput, error) {
		deviceCall := session.Device(input.Address)

		properties, err := deviceCall.Properties()

		return &DevicePropertiesOutput{properties}, err
	})
}

func removeEndpoint(api huma.API, session bluetooth.Session) {
	huma.Register(api, huma.Operation{
		OperationID: "device-remove",
		Method:      http.MethodGet,
		Path:        "/device/{address}/remove",
		Summary:     "Remove",
		Description: "This endpoint removes a device from its associated adapter.",
		Tags:        []string{"Device"},
	}, func(_ context.Context, input *struct {
		AddressInput
	}) (*struct{}, error) {
		deviceCall := session.Device(input.Address)

		return nil, deviceCall.Remove()
	})
}

func pairEndpoint(api huma.API, session bluetooth.Session) {
	huma.Register(api, huma.Operation{
		OperationID: "device-pair",
		Method:      http.MethodGet,
		Path:        "/device/{address}/pair",
		Summary:     "Pairing",
		Description: "This endpoint starts a pairing process to an unpaired device in pairing mode. If the `cancel` parameter is specified, an ongoing pairing operation to the device, if it exists, will be stopped.",
		Tags:        []string{"Device"},
	}, func(_ context.Context, input *struct {
		AddressInput
		Cancel bool `query:"cancel" doc:"Specifies if an ongoing pairing operation to the device should be cancelled."`
	}) (*struct{}, error) {
		deviceCall := session.Device(input.Address)

		if input.Cancel {
			return nil, deviceCall.CancelPairing()
		}

		return nil, deviceCall.Pair()
	})
}

func connectEndpoint(api huma.API, session bluetooth.Session) {
	huma.Register(api, huma.Operation{
		OperationID: "device-connect",
		Method:      http.MethodGet,
		Path:        "/device/{address}/connect",
		Summary:     "Connection",
		Description: "This endpoint starts a connection process to a paired device. If a service profile UUID is specified, it will attempt to connect to it, otherwise a profile will be chosen and connected to automatically.",
		Tags:        []string{"Device"},
	}, func(_ context.Context, input *struct {
		AddressInput
		UUID uuid.UUID `query:"profile_uuid" format:"uuid" doc:"The Bluetooth service profile UUID."`
	}) (*struct{}, error) {
		deviceCall := session.Device(input.Address)

		if input.UUID != uuid.Nil {
			return nil, deviceCall.ConnectProfile(input.UUID)
		}

		return nil, deviceCall.Connect()
	})
}

func disconnectEndpoint(api huma.API, session bluetooth.Session) {
	huma.Register(api, huma.Operation{
		OperationID: "device-disconnect",
		Method:      http.MethodGet,
		Path:        "/device/{address}/disconnect",
		Summary:     "Disconnection",
		Description: "This endpoint starts a disconnection process from a paired device. If a service profile UUID is specified, it will attempt to disconnect from it.",
		Tags:        []string{"Device"},
	}, func(_ context.Context, input *struct {
		AddressInput
		UUID uuid.UUID `query:"profile_uuid" format:"uuid" doc:"The Bluetooth service profile UUID."`
	}) (*struct{}, error) {
		deviceCall := session.Device(input.Address)

		if input.UUID != uuid.Nil {
			return nil, deviceCall.DisconnectProfile(input.UUID)
		}

		return nil, deviceCall.Disconnect()
	})
}
