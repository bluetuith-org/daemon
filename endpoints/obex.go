package endpoints

import (
	"context"
	"net/http"

	"github.com/bluetuith-org/api-native/api/bluetooth"
	"github.com/danielgtaylor/huma/v2"
	"github.com/puzpuzpuz/xsync/v3"
)

var transferMap = xsync.NewMapOf[bluetooth.MacAddress, context.CancelFunc]() 

func obexEndpoints(api huma.API, session bluetooth.Session) {
}

func sendFileEndpoint(api huma.API, session bluetooth.Session) {
	type SendFilesInput struct {
		FilePaths []string `json:"file_paths" required:"true" doc:"The full paths of files to be sent."`
	}

	huma.Register(api, huma.Operation{
		OperationID: "device-send-file",
		Method:      http.MethodPost,
		Path:        "/device/{address}/send_file",
		Summary:     "Send files (Object Push)",
		Description: "This endpoint attempts to send files to a device.",
		Tags:        []string{"Device"},
	}, func(_ context.Context, input *struct {
		AddressInput
		SendFilesInput
	}) (*struct{}, error) {
		obexCall := session.Obex(input.Address)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if err := obexCall.FileTransfer().CreateSession(ctx); err != nil {
			return nil, err
		}

		return nil, nil
	})
}
