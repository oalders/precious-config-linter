package audit

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	if err := os.Chdir(dir); err != nil {
		panic(err)
	}
}

func TestPaths(t *testing.T) {
	config, err := PreciousConfig("testdata/precious.toml")
	assert.NoError(t, err)
	success, err := Paths(config)
	assert.NoError(t, err)
	assert.True(t, success)
}

func TestPatternsOkFailure(t *testing.T) {
	ignoreConfig, err := PpathConfig(".ppath.toml")
	assert.NoError(t, err)

	seen := make(matchCache)
	paths := []string{"foo", "bar", "go.mod", "**/*.go"}
	ok, err := patternsOk(seen, ignoreConfig, "golangci-lint", "include", paths)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(
		t,
		matchCache{"bar": false, "foo": false, "go.mod": true, "**/*.go": true},
		seen,
	)
}

func TestPatternsOkSuccess(t *testing.T) {
	ignoreConfig, err := PpathConfig(".ppath.toml")
	assert.NoError(t, err)

	seen := make(matchCache)
	paths := []string{"go.mod", "**/*.go"}
	ok, err := patternsOk(seen, ignoreConfig, "golangci-lint", "include", paths)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(
		t,
		matchCache{"go.mod": true, "**/*.go": true},
		seen,
	)
}

func TestPatternIgnored(t *testing.T) {
	ignoreConfig, err := PpathConfig(".ppath.toml")
	assert.NoError(t, err)

	assert.True(
		t,
		patternIgnored(ignoreConfig, "omegasort-gitignore", "this-does-not-exist.txt"),
	)

	assert.False(
		t,
		patternIgnored(ignoreConfig, "omegasort-gitignore", "this-does-not-exist"),
	)
}
