package http

import (
	"context"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jasonhancock/go-logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPServer(t *testing.T) {
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(handle)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
	})

	addr, u := newRandomPort(t)
	var wg sync.WaitGroup
	require.NoError(t, NewHTTPServer(ctx, logger.Silence(), &wg, r, addr))

	// give the server time to spin up
	time.Sleep(1 * time.Second)

	req, err := http.NewRequest(http.MethodGet, u+"/foo", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, "hello world", string(b))

	cancel()
	wg.Wait()
}

// newRandomPort identifies a port on the localhost interface for use during tests
// and returns the string in host:port format as well as a url with an http scheme.
// It uses similar methodology to how the net/http/httptest server chooses a port.
func newRandomPort(t *testing.T) (string, string) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := l.Addr()
	l.Close()
	return addr.String(), "http://" + addr.String()
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello world"))
}
