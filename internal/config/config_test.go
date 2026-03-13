package config

import (
	"os"
	"testing"
)

func TestLoad_DefaultsToCompiledConstant(t *testing.T) {
	os.Unsetenv("THREAT_LEVEL")
	cfg := Load()
	if cfg.ThreatLevel != LevelGreen {
		t.Errorf("expected GREEN (compiled constant), got %q", cfg.ThreatLevel)
	}
}

func TestLoad_EnvVarOverridesConstant(t *testing.T) {
	tests := []struct {
		envVal string
		want   Level
	}{
		{"GREEN", LevelGreen},
		{"AMBER", LevelAmber},
		{"RED", LevelRed},
		{"BLACK", LevelBlack},
		{"green", LevelGreen}, // case-insensitive
		{"amber", LevelAmber}, // case-insensitive
	}

	for _, tt := range tests {
		t.Run(tt.envVal, func(t *testing.T) {
			t.Setenv("THREAT_LEVEL", tt.envVal)
			cfg := Load()
			if cfg.ThreatLevel != tt.want {
				t.Errorf("THREAT_LEVEL=%q: expected %q, got %q", tt.envVal, tt.want, cfg.ThreatLevel)
			}
		})
	}
}

func TestLoad_InvalidLevelFallsBackToGreen(t *testing.T) {
	t.Setenv("THREAT_LEVEL", "ULTRAVIOLET")
	cfg := Load()
	if cfg.ThreatLevel != LevelGreen {
		t.Errorf("expected GREEN fallback for invalid level, got %q", cfg.ThreatLevel)
	}
}

func TestLoad_DefaultPort(t *testing.T) {
	os.Unsetenv("PORT")
	cfg := Load()
	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %q", cfg.Port)
	}
}

func TestLoad_CustomPort(t *testing.T) {
	t.Setenv("PORT", "9090")
	cfg := Load()
	if cfg.Port != "9090" {
		t.Errorf("expected port 9090, got %q", cfg.Port)
	}
}

func TestConfig_NumericLevel(t *testing.T) {
	tests := []struct {
		level Level
		want  float64
	}{
		{LevelGreen, 0},
		{LevelAmber, 1},
		{LevelRed, 2},
		{LevelBlack, 3},
	}
	for _, tt := range tests {
		t.Run(string(tt.level), func(t *testing.T) {
			cfg := &Config{ThreatLevel: tt.level}
			if got := cfg.NumericLevel(); got != tt.want {
				t.Errorf("NumericLevel() for %q = %v, want %v", tt.level, got, tt.want)
			}
		})
	}
}

func TestConfig_Classes_ContainsTailwindClasses(t *testing.T) {
	tests := []struct {
		level         Level
		wantTextClass string
	}{
		{LevelGreen, "text-emerald-400"},
		{LevelAmber, "text-amber-400"},
		{LevelRed, "text-red-400"},
		{LevelBlack, "text-white"},
	}
	for _, tt := range tests {
		t.Run(string(tt.level), func(t *testing.T) {
			cfg := &Config{ThreatLevel: tt.level}
			classes := cfg.Classes()
			if classes.Text != tt.wantTextClass {
				t.Errorf("Classes().Text for %q = %q, want %q", tt.level, classes.Text, tt.wantTextClass)
			}
		})
	}
}

func TestConfig_Classes_PulseOnlyForRedAndBlack(t *testing.T) {
	noPulse := []Level{LevelGreen, LevelAmber}
	pulse := []Level{LevelRed, LevelBlack}

	for _, l := range noPulse {
		cfg := &Config{ThreatLevel: l}
		if cfg.Classes().Pulse {
			t.Errorf("expected no pulse for %q", l)
		}
	}
	for _, l := range pulse {
		cfg := &Config{ThreatLevel: l}
		if !cfg.Classes().Pulse {
			t.Errorf("expected pulse for %q", l)
		}
	}
}
