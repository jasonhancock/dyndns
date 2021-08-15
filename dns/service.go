package dns

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

var errNotAllowed = errors.New("address not allowed")

// Service handles updating DNS records.
type Service struct {
	allowed map[string]struct{}
	db      *sqlx.DB
}

// NewService iniitalizes a new Service.
func NewService(allowedAddrs []string, db *sqlx.DB) *Service {
	s := &Service{
		allowed: make(map[string]struct{}, len(allowedAddrs)),
		db:      db,
	}

	for _, v := range allowedAddrs {
		s.allowed[v] = struct{}{}
	}

	return s
}

// DNS handles dns requests.
func (s *Service) DNS(ctx context.Context, req Request) error {
	if _, ok := s.allowed[req.Name]; !ok {
		return errNotAllowed
	}

	const query = `UPDATE records SET content=? WHERE name=?`

	_, err := s.db.ExecContext(ctx, query, req.Value, req.Name)
	return err
}
