package dns

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/stretchr/testify/require"
)

func TestDetermineZone(t *testing.T) {
	tests := []struct {
		desc     string
		zones    []types.HostedZone
		domain   string
		expected string
		err      error
	}{
		{
			"simple",
			[]types.HostedZone{
				{Id: aws.String("xyz123"), Name: aws.String("foo.com.")},
				{Id: aws.String("abc123"), Name: aws.String("example.com.")},
			},
			"foo.example.com",
			"abc123",
			nil,
		},
		{
			"not found",
			[]types.HostedZone{
				{Id: aws.String("abc123"), Name: aws.String("example.com.")},
			},
			"foo.com",
			"",
			errZoneNotFound,
		},
		{
			"more specific domain found",
			[]types.HostedZone{
				{Id: aws.String("abc123"), Name: aws.String("example.com.")},
				{Id: aws.String("xyz123"), Name: aws.String("foo.example.com.")},
			},
			"bar.foo.example.com",
			"xyz123",
			nil,
		},
		{
			"more specific domain found first",
			[]types.HostedZone{
				{Id: aws.String("xyz123"), Name: aws.String("foo.example.com.")},
				{Id: aws.String("abc123"), Name: aws.String("example.com.")},
			},
			"bar.foo.example.com",
			"xyz123",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			id, err := determineZoneID(tt.zones, tt.domain)
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, id)
		})
	}
}
