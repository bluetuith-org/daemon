package endpoints

import (
	"context"
	"errors"
	"net/http"

	"github.com/bluetuith-org/api-native/api/bluetooth"
	"github.com/bluetuith-org/api-native/api/eventbus"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/puzpuzpuz/xsync/v3"
)

type authRequestEvent struct {
	ID            int64  `json:"auth_id,omitempty" doc:"The ID of the authorization request."`
	ReplyRequired bool   "json:\"reply_required,omitempty\" doc:\"If this parameter is set to 'true', use the `/auth/{auth_id}/{reply}` endpoint to respond to this request, otherwise ignore.\""
	AuthType      string `json:"auth_type,omitempty" enum:"pairing,transfer" doc:"The type of the authorization request."`

	PairingParams  *authPairingEvent  "json:\"pairing_params,omitempty\" doc:\"The parameters of the `pairing` authorization request.\""
	TransferParams *authTransferEvent "json:\"transfer_params,omitempty\" doc:\"The parameters of the `transfer` authorization request.\""
}

type authPairingEvent struct {
	PairingType string `json:"pairing_type,omitempty" enum:"display-pincode,display-passkey,confirm-passkey,authorize-pairing,authorize-service" doc:"The type of the pairing authorization request."`

	Address     bluetooth.MacAddress `json:"address,omitempty" doc:"The address of the device."`
	Pincode     string               `json:"pincode,omitempty" doc:"The provided pincode value."`
	Passkey     uint32               `json:"passkey,omitempty" doc:"The provided passkey value."`
	Entered     uint16               `json:"entered,omitempty" doc:"The entered passkey value."`
	ServiceUUID *uuid.UUID           `json:"uuid,omitempty" doc:"The service profile UUID."`
}

type authTransferEvent struct {
	Path           string                     `json:"path,omitempty" doc:"The path to the file."`
	FileProperties bluetooth.FileTransferData `json:"file_properties,omitempty" doc:"The properties of the file."`
}

type authEventReply struct {
	reply  bool
	reason string
}

type authEventID uint

const authEvent = authEventID(100)

var requests = xsync.NewMapOf[int64, chan authEventReply]()

func authEndpoint(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "auth",
		Method:      http.MethodGet,
		Path:        "/auth/{auth_id}/{reply}",
		Summary:     "Authorization",
		Description: "This endpoint enables responses to authorization requests, like device pairing or receiving file transfers.",
	}, func(_ context.Context, input *struct {
		ID     int64  "path:\"auth_id\" doc:\"The authorization ID provided by the `auth` event.\""
		Reply  string `path:"reply" json:"reply,omitempty" enum:"yes,no" doc:"The reply to an authorization request."`
		Reason string "query:\"reason\" json:\"reason,omitempty\" doc:\"An optional user-specified reason if the reply is `no`.\""
	}) (*struct{}, error) {
		if input.ID <= 0 {
			return nil, errors.New("Invalid authorization ID.")
		}

		ch, ok := requests.LoadAndDelete(input.ID)
		if !ok {
			return nil, errors.New("Authorization ID not found.")
		}

		ch <- authEventReply{input.Reply == "yes", input.Reason}

		return nil, nil
	})
}

type authorizer struct {
	id *xsync.Counter
}

func NewAuthorizer() *authorizer {
	return &authorizer{id: xsync.NewCounter()}
}

func (a *authorizer) AuthorizeTransfer(timeout bluetooth.AuthTimeout, path string, props bluetooth.FileTransferData) error {
	return a.sendAndWait(timeout, authRequestEvent{
		AuthType:      "transfer",
		ReplyRequired: true,
		TransferParams: &authTransferEvent{
			Path:           path,
			FileProperties: props,
		},
	})
}

func (a *authorizer) DisplayPinCode(timeout bluetooth.AuthTimeout, address bluetooth.MacAddress, pincode string) error {
	a.send(authRequestEvent{
		AuthType:      "pairing",
		ReplyRequired: false,
		PairingParams: &authPairingEvent{
			PairingType: "display-pincode",
			Address:     address,
			Pincode:     pincode,
		},
	})

	return nil
}

func (a *authorizer) DisplayPasskey(timeout bluetooth.AuthTimeout, address bluetooth.MacAddress, passkey uint32, entered uint16) error {
	a.send(authRequestEvent{
		AuthType:      "pairing",
		ReplyRequired: false,
		PairingParams: &authPairingEvent{
			PairingType: "display-passkey",
			Address:     address,
			Passkey:     passkey,
			Entered:     entered,
		},
	})

	return nil
}

func (a *authorizer) ConfirmPasskey(timeout bluetooth.AuthTimeout, address bluetooth.MacAddress, passkey uint32) error {
	return a.sendAndWait(timeout, authRequestEvent{
		AuthType:      "pairing",
		ReplyRequired: true,
		PairingParams: &authPairingEvent{
			PairingType: "confirm-passkey",
			Address:     address,
			Passkey:     passkey,
		},
	})
}

func (a *authorizer) AuthorizePairing(timeout bluetooth.AuthTimeout, address bluetooth.MacAddress) error {
	return a.sendAndWait(timeout, authRequestEvent{
		AuthType:      "pairing",
		ReplyRequired: true,
		PairingParams: &authPairingEvent{
			PairingType: "authorize-pairing",
			Address:     address,
		},
	})
}

func (a *authorizer) AuthorizeService(timeout bluetooth.AuthTimeout, address bluetooth.MacAddress, uuid uuid.UUID) error {
	return a.sendAndWait(timeout, authRequestEvent{
		AuthType:      "pairing",
		ReplyRequired: true,
		PairingParams: &authPairingEvent{
			PairingType: "authorize-service",
			Address:     address,
			ServiceUUID: &uuid,
		},
	})
}

func (a *authorizer) send(data authRequestEvent) int64 {
	a.id.Inc()
	data.ID = a.id.Value()

	eventbus.Publish(authEvent, data)

	return data.ID
}

func (a *authorizer) sendAndWait(timeout bluetooth.AuthTimeout, data authRequestEvent) error {
	var reply authEventReply

	ch := make(chan authEventReply, 1)
	requests.Store(a.send(data), ch)
	select {
	case <-timeout.Done():
	case reply = <-ch:
	}

	if reply.reply {
		return nil
	}

	return reply
}

func (i authEventID) String() string {
	return "auth"
}

func (i authEventID) Value() uint {
	return uint(i)
}

func (a authEventReply) Error() string {
	if a.reply {
		return ""
	}

	reason := "The authorization request was not accepted."
	if a.reason != "" {
		reason = a.reason
	}

	return reason
}
