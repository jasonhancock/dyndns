package update

import (
	"net"

	"github.com/jasonhancock/dyndns/dns"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var (
		url      string
		username string
		password string
		iface    string
		name     string
	)

	cmd := &cobra.Command{
		Use:          "update",
		Short:        "Updates a dns record.",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return errors.New("required parameter not set: -name")
			}

			ip, err := v4Addr(iface)
			if err != nil {
				return err
			}

			d := dns.NewClient(username, password, url)
			return d.Set(dns.Request{Name: name, Value: ip})
		},
	}

	cmd.Flags().StringVar(
		&url,
		"url-addr",
		"https://dyndns.jasonhancock.com/v1/dns",
		"The dyndns endpoint",
	)

	cmd.Flags().StringVar(
		&username,
		"username",
		"",
		"The username for http basic auth",
	)

	cmd.Flags().StringVar(
		&password,
		"password",
		"",
		"The password for http basic auth",
	)

	cmd.Flags().StringVar(
		&iface,
		"iface",
		"eth0",
		"Which interface to query.",
	)

	cmd.Flags().StringVar(
		&name,
		"name",
		"",
		"The dns name to update",
	)

	return cmd
}

var errIpNotFound = errors.New("no ipv4 IP found")

// v4Addr returns the ipv4 address of the interface identified by name or an error
func v4Addr(name string) (string, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return "", err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	for _, a := range addrs {
		switch v := a.(type) {
		case *net.IPNet:
			if v.IP.To4() != nil {
				return v.IP.String(), nil
			}
		}
	}

	return "", errIpNotFound
}
