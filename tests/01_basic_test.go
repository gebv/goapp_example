package tests

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test01_Basic(t *testing.T) {
	addr, err := url.Parse("http://" + *addressF + "/foo/bar")
	require.NoError(t, err)

	res, err := http.Get(addr.String())
	require.NoError(t, err)
	dat, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "Hi there, I love foo/bar!", string(dat))
}
