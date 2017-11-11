package core

import (
	"io"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"

	"github.com/golang/protobuf/proto"
)

// ConfigLoader is an utility to load V2Ray config from external source.
type ConfigLoader func(input io.Reader) (*Config, error)

var configLoaderCache = make(map[ConfigFormat]ConfigLoader)

// RegisterConfigLoader add a new ConfigLoader.
func RegisterConfigLoader(format ConfigFormat, loader ConfigLoader) error {
	configLoaderCache[format] = loader
	return nil
}

// LoadConfig loads config with given format from given source.
func LoadConfig(format ConfigFormat, input io.Reader) (*Config, error) {
	loader, found := configLoaderCache[format]
	if !found {
		return nil, newError(ConfigFormat_name[int32(format)], " is not loadable.")
	}
	return loader(input)
}

func loadProtobufConfig(input io.Reader) (*Config, error) {
	config := new(Config)
	data, err := buf.ReadAllToBytes(input)
	if err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

func init() {
	common.Must(RegisterConfigLoader(ConfigFormat_Protobuf, loadProtobufConfig))
}
