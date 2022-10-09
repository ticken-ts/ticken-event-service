package api

import (
	"ticken-event-service/infra"
)

type Controller interface {
	Setup(router infra.Router)
}

type Middleware interface {
	Setup(router infra.Router)
}
