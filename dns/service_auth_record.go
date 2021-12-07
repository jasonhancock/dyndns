package dns

import (
	"context"

	"github.com/pkg/errors"
)

var errNotAllowed = errors.New("address not allowed")

// AuthRecordService is a middleware service for ensuring a record is on our approved list of records to manage.
type AuthRecordService struct {
	svc     SVC
	allowed map[string]struct{}
}

// NewAuthRecordService initializes an AuthRecordServict.
func NewAuthRecordService(svc SVC, allowed []string) *AuthRecordService {
	s := &AuthRecordService{
		svc:     svc,
		allowed: make(map[string]struct{}, len(allowed)),
	}

	for _, v := range allowed {
		s.allowed[v] = struct{}{}
	}

	return s
}

// DNS logs info about the request when errors occur.
func (s *AuthRecordService) DNS(ctx context.Context, req Request) error {
	if _, ok := s.allowed[req.Name]; !ok {
		return errNotAllowed
	}

	return s.svc.DNS(ctx, req)
}
