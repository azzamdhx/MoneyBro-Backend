package graph

import (
	"github.com/azzamdhx/moneybro/backend/internal/services"
)

type Resolver struct {
	Services *services.Services
}

func NewResolver(svc *services.Services) *Resolver {
	return &Resolver{
		Services: svc,
	}
}
