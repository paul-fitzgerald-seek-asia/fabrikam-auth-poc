package screening

import "log"

type ScreeningRouteConfig interface {
	ScreeningService() Service
	Logger() *log.Logger
}
