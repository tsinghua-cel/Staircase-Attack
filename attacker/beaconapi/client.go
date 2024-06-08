package beaconapi

import (
	"context"
	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/rs/zerolog"
)

func NewClient(ctx context.Context, endpoint string) (eth2client.Service, error) {
	return http.New(ctx,
		// WithAddress supplies the address of the beacon node, as a URL.
		http.WithAddress(endpoint),
		// LogLevel supplies the level of logging to carry out.
		http.WithLogLevel(zerolog.WarnLevel),
	)
}
