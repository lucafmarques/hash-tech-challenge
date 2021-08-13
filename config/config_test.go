package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	assert.IsType(t, Config{}, config, "Failed asserting config type")
	assert.Equal(t, DEVELOPMENT, config.Service.Environment, "Failed asserting default Service.Environment in default Config")
}

func TestLoadFromEnv(t *testing.T) {
	config := Config{}
	env := "TEST"
	path := "config.yaml"
	os.Setenv(env, path)

	err := config.LoadFromEnv(env, path)

	assert.Nil(t, err, "Failed asserting nil error")
}

func TestLoadFromEnvNoEnv(t *testing.T) {
	config := Config{}
	path := "config.yaml"

	err := config.LoadFromEnv("INEXISTENT_ENV", path)

	assert.Nil(t, err, "Failed asserting nil error")
}

func TestLoadFromEnvPathError(t *testing.T) {
	config := Config{}
	path := "fake_path_that_shouldnt_exist.go"

	err := config.LoadFromEnv("INEXISTENT_ENV", path)

	assert.Error(t, err, "Failed asserting error reading config file")
}

func TestLoadFromEnvWrongConfig(t *testing.T) {
	config := Config{}
	path := "temp_test_config.yaml"
	f, _ := ioutil.TempFile("", path)
	f.Write([]byte(`
	service:
		core: "Not actual Core structure"
	`))
	defer os.Remove(path)

	err := config.LoadFromEnv("INEXISTENT_ENV", f.Name())
	assert.Error(t, err, "Failed asserting error loading config")
}
