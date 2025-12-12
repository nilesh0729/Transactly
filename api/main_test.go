package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	Anuskh "github.com/nilesh0729/Transactly/db/Result"
	"github.com/nilesh0729/Transactly/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store Anuskh.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(store, config)
	require.NoError(t, err)

	return server

}
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())

}
