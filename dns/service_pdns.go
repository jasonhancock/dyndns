package dns

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// ServicePDNS handles updating DNS records.
type ServicePDNS struct {
	db *sqlx.DB
}

// NewService iniitalizes a new Service.
func NewServicePDNS(db *sqlx.DB) *ServicePDNS {
	s := &ServicePDNS{
		db: db,
	}

	return s
}

// DNS handles dns requests.
func (s *ServicePDNS) DNS(ctx context.Context, req Request) error {
	const query = `UPDATE records SET content=? WHERE name=?`

	_, err := s.db.ExecContext(ctx, query, req.Value, req.Name)
	return err
}
