package ports

import (
	"context"

	"github.com/nicodelara/microblogging-uala/internal/timeline/domain"
)

// TimelineService define la interfaz para el servicio de timeline
type TimelineService interface {
	// GetTimeline obtiene el timeline de un usuario
	GetTimeline(ctx context.Context, username string, offset, limit int) (*domain.Timeline, error)
}

// TimelineUseCase define la interfaz para los casos de uso del timeline
type TimelineUseCase interface {
	// GetUserTimeline obtiene el timeline de un usuario
	GetUserTimeline(ctx context.Context, username string, offset, limit int) (*domain.Timeline, error)
}
