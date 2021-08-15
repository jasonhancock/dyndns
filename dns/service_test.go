package dns

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

// TODO: this test doesn't do any setup....you need a mysql server running with
// the pdns database already created in it.

func TestServiceDNS(t *testing.T) {
	conn, err := connectDB("127.0.0.1", "3306", "pdns", "root", "password")
	require.NoError(t, err)
	defer conn.Close()

	t.Run("normal", func(t *testing.T) {
		svc := NewService([]string{"home.jasonhancock.com"}, conn)

		require.NoError(t, svc.DNS(context.Background(), Request{
			Name:  "home.jasonhancock.com",
			Value: "192.168.22.33",
		}))
	})

	t.Run("not allowed", func(t *testing.T) {
		svc := NewService([]string{}, conn)

		require.Equal(t, errNotAllowed, svc.DNS(context.Background(), Request{
			Name:  "home.jasonhancock.com",
			Value: "192.168.22.33",
		}))
	})
}

func connectDB(dbHost, dbPort, dbName, dbUser, dbPass string) (*sqlx.DB, error) {
	c := mysql.NewConfig()

	c.User = dbUser
	c.Passwd = dbPass
	c.DBName = dbName
	c.Net = "tcp"
	c.Addr = fmt.Sprintf("%s:%s", dbHost, dbPort)

	// Require MultiStatements in order to be able to apply migrations.
	c.MultiStatements = true

	// Automatically parse into time.Time
	c.ParseTime = true

	db, err := sqlx.Open("mysql", c.FormatDSN())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
