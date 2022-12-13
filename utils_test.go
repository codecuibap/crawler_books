package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	s, err := LoadConfig("nxbkimdong.com.vn")
	require.NoError(t, err)
	require.NotNil(t, s)
	require.Equal(t, s.Site, "nxbkimdong.com.vn")
}
