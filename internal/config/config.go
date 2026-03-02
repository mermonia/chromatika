package config

import (
	_ "embed"
)

type Config struct {
	DMLuminosityThreshold float32 `toml:"dark_mode_luminosity_threshold"`

	DMminimumNeutralLuminosity float32 `toml:"dark_mode_minimum_neutral_luminosity"`
	DMmaximumNeutralLuminosity float32 `toml:"dark_mode_maximum_neutral_luminosity"`
	LMminimumNeutralLuminosity float32 `toml:"light_mode_minimum_neutral_luminosity"`
	LMmaximumNeutralLuminosity float32 `toml:"light_mode_maximum_neutral_luminosity"`

	DMdarkestChromaProportion   float32 `toml:"dark_mode_darkest_chroma_proportion"`
	DMbrightestChromaProportion float32 `toml:"dark_mode_brightest_chroma_proportion"`
	LMdarkestChromaProportion   float32 `toml:"light_mode_darkest_chroma_proportion"`
	LMbrightestChromaProportion float32 `toml:"light_mode_brightest_chroma_proportion"`

	DMbgLuminosityDelta float32 `toml:"dark_mode_bg_luminosity_deltas"`
	LMbgLuminosityDelta float32 `toml:"light_mode_bg_luminosity_deltas"`

	DMpreferredAccentLuminosity float32 `toml:"dark_mode_preferred_accent_luminosity"`
	LMpreferredAccentLuminosity float32 `toml:"light_mode_preferred_accent_luminosity"`
	LuminosityBias              float32 `toml:"luminosity_bias"`

	DMpreferredAccentChroma float32 `toml:"dark_mode_preferred_accent_chroma"`
	LMpreferredAccentChroma float32 `toml:"light_mode_preferred_accent_chroma"`
	ChromaBias              float32 `toml:"chroma_bias"`

	AnalogousDegreeShift float32 `toml:"analogous_degree_shift"`
	TriadicDegreeShift   float32 `toml:"triadic_degree_shift"`

	PrimaryStepSize float32 `toml:"primary_step_size"`
}
