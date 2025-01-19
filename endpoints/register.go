package endpoints

import (
	"github.com/bluetuith-org/api-native/api/bluetooth"
	"github.com/danielgtaylor/huma/v2"
)

func Register(api huma.API, session bluetooth.Session) {
	rootEndpoints(api, session)
	adapterEndpoints(api, session)
	deviceEndpoints(api, session)
	obexEndpoints(api, session)
	networkEndpoints(api, session)
	mediaPlayerEndpoints(api, session)
}
