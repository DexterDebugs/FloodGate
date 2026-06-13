package config
//lowercase - private, uppercase - public

import (
	"os"	//For readfile
	"gopkg.in/yaml.v3"	//for Unmarshal
)

type Route struct {
	Path string		`yaml:"path"`	//Field Names
	Target string	`yaml:"target"`
}

type RateLimit struct {
	Algorithm string	`yaml:"algorithm"`
	Limit int	`yaml:"limit"`
	WindowSeconds int	`yaml:"window_seconds"`
}

type Config struct {		//list of routes
	Routes []Route `yaml:"routes"`	//a slice (list) of routes above
	APIKeys []string	`yaml:"api_keys"`
	RateLimit RateLimit		`yaml:"rate_limit"`
}

func Load(path string)	(*Config, error){
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{} 	//Creating a new empty config from struct, & takes its address, so cfg is a *Config pointer
	err = yaml.Unmarshal(data, cfg)	//err, not error. error is a Go built-in type name — don't reassign it.
	if err != nil{
		return nil, err
	}

	return cfg, err
}

/*
What Unmarshal actually does under the hood:
When you call:
goyaml.Unmarshal(data, cfg)

The YAML library does this:

Parses the byte stream as YAML syntax (finds keys, values, lists, indentation)
Looks at the struct you passed in (cfg points to a Config)
Uses reflection — a Go feature that lets code inspect other code's types at runtime — to look at Config's fields and their YAML tags
Sees Config has a field Routes tagged yaml:"routes", finds routes in the YAML, matches them up
For each item in the YAML list, sees the field is []Route, recursively does the same for each Route's Path and Target fields
Writes the parsed values directly into the struct's memory (which is why you passed a pointer — Unmarshal needs to mutate your variable)

The YAML tags you wrote (`yaml:"path"`) are the bridge between YAML keys (lowercase) and Go fields (capitalized).
 Without those tags, the library wouldn't know path in YAML maps to Path in Go.
 */