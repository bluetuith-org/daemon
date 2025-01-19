package endpoints

import (
	"github.com/bluetuith-org/api-native/api/bluetooth"
	"github.com/danielgtaylor/huma/v2"
)

type AddressInput struct {
	Address bluetooth.MacAddress
	Input   string `json:"address" path:"address" doc:"The Bluetooth MAC address."`
}

func (a *AddressInput) Resolve(_ huma.Context) []error {
	mac, err := bluetooth.ParseMAC(a.Input)
	if err != nil {
		return []error{&huma.ErrorDetail{
			Message:  err.Error(),
			Location: "address",
			Value:    a.Input,
		}}
	}

	a.Address = mac

	return nil
}
