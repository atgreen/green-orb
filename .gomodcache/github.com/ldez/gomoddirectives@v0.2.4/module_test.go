package gomoddirectives

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetModuleFile(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(wd)
	}()

	err = os.Chdir("./testdata/a/")
	require.NoError(t, err)

	file, err := GetModuleFile()
	require.NoError(t, err)

	assert.Equal(t, "github.com/ldez/gomoddirectives/testdata/a", file.Module.Mod.Path)
}

func TestGetModuleFile_here(t *testing.T) {
	file, err := GetModuleFile()
	require.NoError(t, err)

	assert.Equal(t, "github.com/ldez/gomoddirectives", file.Module.Mod.Path)
}
