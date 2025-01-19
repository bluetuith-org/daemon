package endpoints

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bluetuith-org/api-native/api/bluetooth"
	"github.com/danielgtaylor/huma/v2"
)

func adapterEndpoints(api huma.API, session bluetooth.Session) {
	devicesEndpoint(api, session)
	statesEndpoint(api, session)
	adapterPropertiesEndpoint(api, session)
}

func devicesEndpoint(api huma.API, session bluetooth.Session) {
	type AdapterDevicesOutput struct {
		Body []bluetooth.DeviceData
	}

	huma.Register(api, huma.Operation{
		OperationID: "adapter-devices",
		Method:      http.MethodGet,
		Path:        "/adapter/{address}/devices",
		Summary:     "Devices",
		Description: "This endpoint fetches the devices associated with an adapter.",
		Tags:        []string{"Adapter"},
	}, func(_ context.Context, input *struct {
		AddressInput
	}) (*AdapterDevicesOutput, error) {
		adapterCall := session.Adapter(input.Address)

		devices, err := adapterCall.Devices()

		return &AdapterDevicesOutput{devices}, err
	})
}

func adapterPropertiesEndpoint(api huma.API, session bluetooth.Session) {
	type AdapterPropertiesOutput struct {
		Body bluetooth.AdapterData
	}

	huma.Register(api, huma.Operation{
		OperationID: "adapter-properties",
		Method:      http.MethodGet,
		Path:        "/adapter/{address}/properties",
		Summary:     "Properties",
		Description: "This endpoint fetches the properties of an adapter.",
		Tags:        []string{"Adapter"},
	}, func(_ context.Context, input *struct {
		AddressInput
	}) (*AdapterPropertiesOutput, error) {
		adapterCall := session.Adapter(input.Address)

		properties, perr := adapterCall.Properties()
		if perr != nil {
			return nil, perr
		}

		return &AdapterPropertiesOutput{properties}, nil
	})
}

func statesEndpoint(api huma.API, session bluetooth.Session) {
	type AdapterStatesInput struct {
		Powered      string `query:"powered" enum:"enable,disable" doc:"Set the adapter's powered state."`
		Pairable     string `query:"pairable" enum:"enable,disable" doc:"Set the adapter's pairable state."`
		Discoverable string `query:"discoverable" enum:"enable,disable" doc:"Set the adapter's discoverable state."`
		Discovery    string `query:"discovery" enum:"enable,disable" doc:"Toggle the adapter's device discovery mode."`
	}

	type AdapterStatesOutput struct {
		Body struct {
			PoweredState      string `json:"powered,omitempty" enum:"enabled,disabled" doc:"The adapter's powered state."`
			PairableState     string `json:"pairable,omitempty" enum:"enabled,disabled" doc:"The adapter's pairable state."`
			DiscoverableState string `json:"discoverable,omitempty" enum:"enabled,disabled" doc:"The adapter's discoverable state."`
			DiscoveryState    string `json:"discovery,omitempty" enum:"enabled,disabled" doc:"The adapter's device discovery mode state."`
		}
	}

	huma.Register(api, huma.Operation{
		OperationID: "adapter-states",
		Method:      http.MethodGet,
		Path:        "/adapter/{address}/states",
		Summary:     "States",
		Description: "This endpoint, when called by itself, fetches the different states (powered, pairable, discoverable and device discovery) of an adapter. Use the **query parameters** to `enable` or `disable` each state. Note that when **discovery** is **enabled**, all discovered devices will be published to the `/event` stream, with the ***event-name*** as *'device'*, and with ***event-action*** as *'added'*.",
		Tags:        []string{"Adapter"},
	}, func(_ context.Context, input *struct {
		AdapterStatesInput
		AddressInput
	}) (*AdapterStatesOutput, error) {
		states := &AdapterStatesOutput{}
		adapterCall := session.Adapter(input.Address)

		inputs := []struct {
			InputToCheck      string
			EnableFunc        func() error
			DisableFunc       func() error
			SetStatesProperty func(string)
		}{
			{
				InputToCheck: input.Discovery,
				EnableFunc:   adapterCall.StartDiscovery,
				DisableFunc:  adapterCall.StopDiscovery,
				SetStatesProperty: func(toggle string) {
					states.Body.DiscoveryState = toggle
				},
			},
			{
				InputToCheck: input.Discoverable,
				EnableFunc: func() error {
					return adapterCall.SetDiscoverableState(true)
				},
				DisableFunc: func() error {
					return adapterCall.SetDiscoverableState(false)
				},
				SetStatesProperty: func(toggle string) {
					states.Body.DiscoverableState = toggle
				},
			},
			{
				InputToCheck: input.Pairable,
				EnableFunc: func() error {
					return adapterCall.SetPairableState(true)
				},
				DisableFunc: func() error {
					return adapterCall.SetPairableState(false)
				},
				SetStatesProperty: func(toggle string) {
					states.Body.PairableState = toggle
				},
			},
			{
				InputToCheck: input.Powered,
				EnableFunc: func() error {
					return adapterCall.SetPoweredState(true)
				},
				DisableFunc: func() error {
					return adapterCall.SetPoweredState(false)
				},
				SetStatesProperty: func(toggle string) {
					states.Body.PoweredState = toggle
				},
			},
		}

		var errs error
		var emptyInputs int

		for _, in := range inputs {
			var err error
			var state string

			switch in.InputToCheck {
			case "enable":
				err = in.EnableFunc()
				state = "enabled"
			case "disable":
				err = in.DisableFunc()
				state = "disabled"
			case "":
				emptyInputs++
				continue
			}
			if err != nil {
				if errs == nil {
					errs = err
				} else {
					errs = fmt.Errorf("%w, %w", errs, err)
				}
			}

			in.SetStatesProperty(state)
		}
		if errs != nil {
			return nil, errs
		}

		if emptyInputs == len(inputs) {
			properties, perr := adapterCall.Properties()
			if perr != nil {
				return nil, perr
			}

			states.Body.DiscoverableState = toggleStr(properties.Discoverable)
			states.Body.PairableState = toggleStr(properties.Pairable)
			states.Body.DiscoveryState = toggleStr(properties.Discovering)
			states.Body.PoweredState = toggleStr(properties.Powered)
		}

		return states, nil
	})
}

func toggleStr(val bool) string {
	if val {
		return "enabled"
	}

	return "disabled"
}
