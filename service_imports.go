package nacelle

import "github.com/go-nacelle/service"

type ServiceContainer = service.Container

var (
	NewServiceContainer        = service.New
	NewServiceContainerOverlay = service.NewOverlay
)
