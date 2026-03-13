// Package config holds application configuration.
// To change the threat level for the demo, edit the ThreatLevel constant below
// and push to trigger the CI/CD pipeline.
package config

import (
	"log"
	"os"
	"strings"
)

// ThreatLevel is the demo's — change this single line and push.
// Valid values: "GREEN", "AMBER", "RED", "BLACK"
const ThreatLevel = "AMBER"

// Level represents a validated threat level.
type Level string

const (
	LevelGreen Level = "GREEN"
	LevelAmber Level = "AMBER"
	LevelRed   Level = "RED"
	LevelBlack Level = "BLACK"
)

// TailwindClasses holds the Tailwind CSS classes for a given threat level.
type TailwindClasses struct {
	Text       string
	Border     string
	Background string
	Pulse      bool
}

// levelClasses maps each threat level to its Tailwind CSS class set.
var levelClasses = map[Level]TailwindClasses{
	LevelGreen: {
		Text:       "text-emerald-400",
		Border:     "border-emerald-500",
		Background: "bg-emerald-900/20",
		Pulse:      false,
	},
	LevelAmber: {
		Text:       "text-amber-400",
		Border:     "border-amber-500",
		Background: "bg-amber-900/20",
		Pulse:      false,
	},
	LevelRed: {
		Text:       "text-red-400",
		Border:     "border-red-500",
		Background: "bg-red-900/20",
		Pulse:      true,
	},
	LevelBlack: {
		Text:       "text-white",
		Border:     "border-white",
		Background: "bg-black",
		Pulse:      true,
	},
}

// levelNumeric maps each threat level to its Prometheus gauge value.
var levelNumeric = map[Level]float64{
	LevelGreen: 0,
	LevelAmber: 1,
	LevelRed:   2,
	LevelBlack: 3,
}

// Config holds the loaded application configuration.
type Config struct {
	ThreatLevel Level
	Port        string
}

// Load reads configuration from environment variables, falling back to
// compiled-in constants. Invalid values are logged and replaced with defaults.
func Load() *Config {
	level := resolveLevel()
	port := resolvePort()

	return &Config{
		ThreatLevel: level,
		Port:        port,
	}
}

// resolveLevel reads THREAT_LEVEL env var, validates it, and falls back to
// the compiled-in ThreatLevel constant.
func resolveLevel() Level {
	raw := strings.ToUpper(strings.TrimSpace(os.Getenv("THREAT_LEVEL")))
	if raw == "" {
		// Use compiled-in constant
		raw = strings.ToUpper(ThreatLevel)
	}

	l := Level(raw)
	if _, ok := levelClasses[l]; !ok {
		log.Printf("WARNING: unrecognised THREAT_LEVEL %q — defaulting to GREEN", raw)
		return LevelGreen
	}
	return l
}

// resolvePort reads the PORT env var, defaulting to "8080".
func resolvePort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}

// Classes returns the Tailwind CSS class set for this config's threat level.
func (c *Config) Classes() TailwindClasses {
	return levelClasses[c.ThreatLevel]
}

// NumericLevel returns the numeric Prometheus gauge value for this config's threat level.
func (c *Config) NumericLevel() float64 {
	return levelNumeric[c.ThreatLevel]
}
