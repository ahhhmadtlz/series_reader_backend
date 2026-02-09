package config

import (
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
)

const (
	defaultPrefix       = "SERIES_"
	defaultDelimiter    = "."
	defaultSeparator    = "__"
	defaultYamlFilePath = "config.yml"
)

type Option struct {
	Prefix       string
	Delimiter    string
	Separator    string
	YamlFilePath string
	CallbackEnv  func(string) string
}

func defaultCallbackEnv(source string) string {
	base := strings.ToLower(strings.TrimPrefix(source, defaultPrefix))
	return strings.ReplaceAll(base, defaultSeparator, defaultDelimiter)
}

func DefaultOption() Option {
	return Option{
		Prefix:       defaultPrefix,
		Delimiter:    defaultDelimiter,
		Separator:    defaultSeparator,
		YamlFilePath: defaultYamlFilePath,
		CallbackEnv:  defaultCallbackEnv,
	}
}

func Load(opt Option) (*Config, error) {
	k := koanf.New(opt.Delimiter)

	// 1. defaults (from struct)
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		return nil, err
	}

	// 2. yaml file (optional, but recommended)
	if err := k.Load(file.Provider(opt.YamlFilePath), yaml.Parser()); err != nil {
		return nil, err
	}

	// 3. environment variables
	if err := k.Load(
		env.Provider(opt.Prefix, opt.Delimiter, opt.CallbackEnv),
		nil,
	); err != nil {
		return nil, err
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
