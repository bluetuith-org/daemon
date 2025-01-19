package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bluetuith-org/api-native/api/bluetooth"
	"github.com/bluetuith-org/api-native/platform"
	"github.com/bluetuith-org/daemon/endpoints"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/google/uuid"
)

func main() {
	session := platform.Session()
	_, err := session.Start(authHandler{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = session.Stop()
	}()

	// Create a new router & API.
	router := http.NewServeMux()
	api := humago.New(router, huma.DefaultConfig("My API", "1.0.0"))
	endpoints.Register(api, session)

	// Start the server!
	http.ListenAndServe("127.0.0.1:8888", router)
}

type authHandler struct{}

func (authHandler) AuthorizeTransfer(ctx context.Context, path string, props bluetooth.FileTransferData) (_ error) {
	return nil
}
func (authHandler) DisplayPinCode(ctx context.Context, address bluetooth.MacAddress, pincode string) (_ error) {
	return nil
}
func (authHandler) DisplayPasskey(ctx context.Context, address bluetooth.MacAddress, passkey uint32, entered uint16) (_ error) {
	return nil
}
func (authHandler) ConfirmPasskey(ctx context.Context, address bluetooth.MacAddress, passkey uint32) (_ error) {
	return nil
}
func (authHandler) AuthorizePairing(ctx context.Context, address bluetooth.MacAddress) (_ error) {
	return nil
}
func (authHandler) AuthorizeService(ctx context.Context, address bluetooth.MacAddress, uuid uuid.UUID) (_ error) {
	return nil
}
