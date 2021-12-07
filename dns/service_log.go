package dns

import (
	"context"

	"github.com/jasonhancock/go-logger"
)

// LoggingService is a middleware service for logging requests and errors.
type LoggingService struct {
	svc SVC
	log *logger.L
}

// NewLoggingService initializes a logging service.
func NewLoggingService(svc SVC, l *logger.L) *LoggingService {
	return &LoggingService{
		svc: svc,
		log: l,
	}
}

// DNS logs info about the request when errors occur.
func (s *LoggingService) DNS(ctx context.Context, req Request) error {
	err := s.svc.DNS(ctx, req)
	if err != nil {
		s.log.Err(
			"dns_error",
			"error", err.Error(),
			"name", req.Name,
			"value", req.Value,
		)
	}

	return err
}
