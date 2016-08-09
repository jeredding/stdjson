package config

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"time"
)

type StdjsonConfig struct {
	Stdout *StreamConfig `json:"stdout"`
	Stderr *StreamConfig `json:"stderr"`
}

type StreamConfig struct {
	Rewriter RewriterConfig `json:"rewriter"`
}

type RewriterConfig struct {
	GrokRewriter *GrokRewriterConfig `json:"grok"`
}

type GrokRewriterConfig struct {
	AddTimestamp    *string                     `json:"add_timestamp,omitempty"`
	NamedOnly       *bool                       `json:"named_only,omitempty"`
	ExtraPatterns   ExtraPatternsConfig         `json:"extra_patterns"`
	MatchPatterns   MatchPatternsConfig         `json:"match_patterns"`
	RecursiveFields []GrokRewriterRecurseConfig `json:"recursive_fields"`
	WhitelistFields []string                    `json:"whitelist_fields"`
	BlacklistFields []string                    `json:"blacklist_fields"`
	DefaultFields   map[string]interface{}      `json:"default_fields"`
	Multiline       MultilineConfig             `json:"multiline"`
}

type GrokRewriterRecurseConfig struct {
	Field     string
	Delimiter *string
	Trim      string
	Patterns  []string
}

type MultilineConfig struct {
	Delimiter           string   `json:"delimiter"`
	StripLastDelimiter  bool     `json:"strip_last_delimiter"`
	PrefixContinuations []string `json:"prefix_continuations"`
	Timeout             *string  `json:"timeout"`
}

func (mc *MultilineConfig) MustTimeoutDuration() *time.Duration {
	if mc.Timeout == nil {
		return nil
	}

	d, err := time.ParseDuration(*mc.Timeout)
	if err != nil {
		panic(err)
	}

	return &d
}

type ExtraPatternsConfig []ExtraPatternConfig

type MatchPatternsConfig []string

type ExtraPatternConfig struct {
	Inline *InlinePatternConfig `json:"inline"`
}

type InlinePatternConfig struct {
	Name    string `json:"name"`
	Pattern string `json:"pattern"`
}

func LoadConfigFromFile(file string) (*StdjsonConfig, error) {
	config := &StdjsonConfig{}

	configData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(configData, config); err != nil {
		return nil, err
	}

	return config, nil
}
