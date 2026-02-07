package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
)

const (
	defaultPrefix       = "ET_"
	defaultDelimiter    = "."
	defaultSeparator    = "__"
	defaultYamlFilePath = "config.yml"
)

var c Config

type Option struct {
	Prefix       string
	Delimiter    string
	Separator    string
	YamlFilePath string
	CallbackEnv  func(string) string
}

// our environment variables must prefix with `ET_` (Expense Tracker)
// for nested .env should use `__` aka: ET__DB__HOST.
func defaultCallbackEnv(source string) string {
	base := strings.ToLower(strings.TrimPrefix(source, defaultPrefix))
	return strings.ReplaceAll(base, defaultSeparator, defaultDelimiter)
}

func init() {
	k := koanf.New(defaultDelimiter)

	// load default configuration from Default function
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default config: %s", err)
	}

	// load configuration from yaml file
	if err := k.Load(file.Provider(defaultYamlFilePath), yaml.Parser()); err != nil {
		log.Printf("error loading config from `config.yml` file: %s", err)
	}

	// load from environment variable
	if err := k.Load(env.Provider(defaultPrefix, defaultDelimiter, defaultCallbackEnv), nil); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &c); err != nil {
		log.Fatalf("error unmarshaling config: %s", err)
	}
}

// C returns the loaded configuration
func C() Config {
	return c
}

// New creates a new configuration with custom options
func New(opt Option) Config {
	k := koanf.New(opt.Delimiter)

	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default config: %s", err)
	}

	if err := k.Load(file.Provider(opt.YamlFilePath), yaml.Parser()); err != nil {
		log.Printf("error loading config from `%s` file: %s", opt.YamlFilePath, err)
	}

	if err := k.Load(env.Provider(opt.Prefix, opt.Delimiter, opt.CallbackEnv), nil); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &c); err != nil {
		log.Fatalf("error unmarshaling config: %s", err)
	}

	return c
}