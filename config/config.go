package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"strings"
)

func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.UnmarshalStrict(content, cfg)
	if err != nil {
		return nil, err
	}
	setConfigDefaults(cfg)
	setRuleDefaults(cfg)
	return cfg, nil
}

func setConfigDefaults(config *Config) (error) {
	for i := 0; i < len(config.WhitelistObjectNames); i++ {
		var regex *regexp.Regexp
		var raw_regex string
		var err error = nil
		if strings.Contains(config.WhitelistObjectNames[i], "*"){
			raw_regex = strings.Replace(config.WhitelistObjectNames[i], "*", ".+", 1)
			regex, err = regexp.Compile(raw_regex)
			if err != nil {
				return err
			}
			config.WhitelistObjectNames = append(config.WhitelistObjectNames[:i], config.WhitelistObjectNames[i+1:]...)
		}
		config.WhitelistObjectNamesRegexp = append(config.WhitelistObjectNamesRegexp, regex)
	}
	for i := 0; i < len(config.BlacklistObjectNames); i++ {
		var regex *regexp.Regexp
		var raw_regex string
		var err error = nil
		if strings.Contains(config.BlacklistObjectNames[i], "*"){
			raw_regex = strings.Replace(config.BlacklistObjectNames[i], "*", ".+", 1)
			regex, err = regexp.Compile(raw_regex)
			if err != nil {
				return err
			}
			config.BlacklistObjectNames = append(config.BlacklistObjectNames[:i], config.BlacklistObjectNames[i+1:]...)
		}
		config.BlacklistObjectNamesRegexp = append(config.BlacklistObjectNamesRegexp, regex)
	}
	return nil
}

func setRuleDefaults(config *Config) (error){
	for k, v := range config.Rules {
		if v.ValueFactor == 0.0 {
			config.Rules[k].ValueFactor = 1.0
		}
		compiled, err := regexp.Compile(v.Pattern)
		if err != nil {
			return err
		}
		config.Rules[k].PatternRegexp = compiled
	}
	return nil
}

type Config struct {
	StartDelaySeconds    int    `yaml:"startDelaySeconds,omitempty"`
	HostPort             string `yaml:"hostPort,omitempty"`
	JmxUrl               string `yaml:"jmxUrl,omitempty"`
	Username             string `yaml:"username,omitempty"`
	Password             string `yaml:"password,omitempty"`
	Ssl                  bool   `yaml:"ssl,omitempty"`
	LowercaseOutputName  bool   `yaml:"lowercaseOutputName,omitempty"`
	LowercaseOutputLabelNames bool `yaml:"lowercaseOutputLabelNames,omitempty"`
	WhitelistObjectNames []string `yaml:"whitelistObjectNames,omitempty"`
	WhitelistObjectNamesRegexp []*regexp.Regexp
	BlacklistObjectNames []string `yaml:"blacklistObjectNames,omitempty"`
	BlacklistObjectNamesRegexp []*regexp.Regexp
	Rules                []Rule   `yaml:"rules,omitempty"`
	LastUpdate           int64    `yaml:"lastUpdate,omitempty""`
}

type Rule struct {
	Pattern string `yaml:"pattern,omitempty"`
	PatternRegexp *regexp.Regexp
	Name string `yaml:"name,omitempty"`
	Type string `yaml:"type,omitempty"`
	Value string `yaml:"value,omitempty"`
	ValueFactor float64 `yaml:"valueFactor,omitempty"`
	Help string `yaml:"help,omitempty"`
	AttrNameSnakeCase bool `yaml:"attrNameSnakeCase,omitempty"`
	Labels map[string]string `yaml:"labels,omitempty"`
}

