package config

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configFile := `
	dark_mode_luminosity_threshold = 0.1

	dark_mode_minimum_neutral_luminosity = 0.2
	dark_mode_maximum_neutral_luminosity = 0.3
	light_mode_minimum_neutral_luminosity = 0.4
	light_mode_maximum_neutral_luminosity = 0.5

	dark_mode_darkest_chroma_proportion = 1
	dark_mode_brightest_chroma_proportion = 10
	light_mode_darkest_chroma_proportion = 10
	light_mode_brightest_chroma_proportion = 1

	dark_mode_bg_luminosity_deltas = 2
	light_mode_bg_luminosity_deltas = 3

	dark_mode_preferred_accent_luminosity = 40
	light_mode_preferred_accent_luminosity = 50
	luminosity_bias = 0.5

	dark_mode_preferred_accent_chroma = 10
	light_mode_preferred_accent_chroma = 20
	chroma_bias = 0.7

	analogous_degree_shift = 30
	triadic_degree_shift = 120

	primary_step_size = 30
	`

	baseConfig := &Config{
		DMLuminosityThreshold: 0.1,
		
		DMminimumNeutralLuminosity: 0.2,
		DMmaximumNeutralLuminosity: 0.3,
		LMminimumNeutralLuminosity: 0.4,
		LMmaximumNeutralLuminosity: 0.5,

		DMdarkestChromaProportion: 1,
		DMbrightestChromaProportion: 10,
		LMdarkestChromaProportion: 10,
		LMbrightestChromaProportion: 1,

		DMbgLuminosityDelta: 2,
		LMbgLuminosityDelta: 3,

		DMpreferredAccentLuminosity: 40,
		LMpreferredAccentLuminosity: 50,
		LuminosityBias: 0.5,

		DMpreferredAccentChroma: 10,
		LMpreferredAccentChroma: 20,
		ChromaBias: 0.7,

		AnalogousDegreeShift: 30,
		TriadicDegreeShift: 120,
		PrimaryStepSize: 30,
	}

	tempDir := t.TempDir()
	file, err := os.CreateTemp(tempDir, "testConfig*")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err := file.Write([]byte(configFile)); err != nil {
		log.Fatal(err)
	}

	c, err := LoadConfig(file.Name())
	if err != nil {
		log.Fatal(err)
	}

	if !reflect.DeepEqual(c, baseConfig) {
		t.Fatal("read config does not have the expected values")
	}
}







