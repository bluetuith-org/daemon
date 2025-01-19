package endpoints

import (
	"context"
	"net/http"

	bluetooth "github.com/bluetuith-org/api-native/api/bluetooth"
	"github.com/bluetuith-org/api-native/api/eventbus"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
)

type eventPublisher struct {
	sender sse.Sender
}

func (e *eventPublisher) Publish(id uint, name string, data any) {
	e.sender(sse.Message{
		ID:    int(id),
		Data:  data,
		Retry: 0,
	})
}

func eventsEndpoint(api huma.API) {
	newPublisher := func(sender sse.Sender) (*eventPublisher, eventbus.EventSubscriber) {
		eh := eventbus.NilHandler()
		return &eventPublisher{sender}, eh
	}

	sse.Register(api, huma.Operation{
		OperationID: "events",
		Method:      http.MethodGet,
		Path:        "/events",
		Summary:     "Events",
	}, map[string]any{
		"auth":         authRequestEvent{},
		"adapter":      bluetooth.AdapterEvent(),
		"error":        bluetooth.ErrorEvent(),
		"device":       bluetooth.DeviceEvent(),
		"mediaplayer":  bluetooth.MediaEvent(),
		"filetransfer": bluetooth.FileTransferEvent(),
	}, func(ctx context.Context, input *struct{}, send sse.Sender) {
		eventbus.RegisterEventHandlers(newPublisher(send))
		defer eventbus.DisableEvents()

		<-ctx.Done()
	})
}
