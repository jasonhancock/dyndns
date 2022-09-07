package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	cl "github.com/jasonhancock/cobra-logger"
	"github.com/jasonhancock/dyndns/dns"
	phttp "github.com/jasonhancock/dyndns/http"
	"github.com/jasonhancock/dyndns/version"
	"github.com/jasonhancock/go-env"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

const (
	defaultServerHTTPAddr = ":6061"
)

func NewCmd(wg *sync.WaitGroup, info version.Info) *cobra.Command {
	var (
		httpAddr string
		allowed  string
		logConf  *cl.Config
		//dbConf   *DBConfig
	)

	cmd := &cobra.Command{
		Use:          "server",
		Short:        "Starts the server.",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			l := logConf.Logger(os.Stdout)

			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			// construct AWS Route53 client
			c, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
			if err != nil {
				return err
			}

			/*
				db, err := dbConf.Connect()
				if err != nil {
					return errors.Wrap(err, "connecting to database")
				}
			*/

			routerHTTP := mux.NewRouter()
			v1HTTP := mux.NewRouter()
			routerHTTP.PathPrefix("/v1/").Handler(http.StripPrefix("/v1", v1HTTP))

			{
				var svc dns.SVC
				svc = dns.NewServiceRoute53(route53.NewFromConfig(c))
				svc = dns.NewAuthRecordService(svc, strings.Split(strings.TrimSpace(allowed), ","))
				svc = dns.NewLoggingService(svc, l.New("dns"))
				dns.NewHTTPServer(v1HTTP, svc)
			}

			// start the HTTP server
			{
				httpLogger := l.New("http_servers")

				// start up the http server
				if err := phttp.NewHTTPServer(ctx, httpLogger, wg, routerHTTP, httpAddr); err != nil {
					return err
				}
			}

			wg.Wait()
			return nil
		},
	}

	logConf = cl.NewConfig(cmd)
	//dbConf = NewDBConfig(cmd)

	cmd.Flags().StringVar(
		&httpAddr,
		"http-addr",
		env.String("ADDR", defaultServerHTTPAddr),
		"Interface:port to bind the HTTP interface to.",
	)

	cmd.Flags().StringVar(
		&allowed,
		"allowed-addrs",
		os.Getenv("ALLOWED_ADDRS"),
		"comma delimited list of addresses that are allowed to be managed by dyndns",
	)

	return cmd
}

type DBConfig struct {
	Host string
	Port int
	Name string
	User string
	Pass string
}

func NewDBConfig(cmd *cobra.Command) *DBConfig {
	c := &DBConfig{}

	cmd.Flags().StringVar(
		&c.Host,
		"db-host",
		env.String("DB_HOST", "127.0.0.1"),
		"Database host",
	)

	cmd.Flags().IntVar(
		&c.Port,
		"db-port",
		env.Int("DB_PORT", 3306),
		"Database port",
	)

	cmd.Flags().StringVar(
		&c.Name,
		"db-name",
		os.Getenv("DB_NAME"),
		"Database name",
	)

	cmd.Flags().StringVar(
		&c.User,
		"db-user",
		os.Getenv("DB_USER"),
		"Database user",
	)

	cmd.Flags().StringVar(
		&c.Pass,
		"db-pass",
		os.Getenv("DB_PASS"),
		"Database password",
	)

	return c
}

// Logger gets the logger
func (cfg *DBConfig) Connect() (*sqlx.DB, error) {
	return connectDB(cfg.Host, strconv.Itoa(cfg.Port), cfg.Name, cfg.User, cfg.Pass)
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
