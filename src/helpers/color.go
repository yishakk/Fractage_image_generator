package helpers

import (
	"errors"
	"image/color"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Generates a color with random rgb values (a=255).
func RandomColor() color.RGBA {
	randomColor := color.RGBA{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		255,
	}
	return randomColor
}

// Generates a color with random rgba values.
func RandomAlphaColor() color.RGBA {
	randomColor := color.RGBA{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
	}
	return randomColor
}

// Parses a given color value.
func ParseColor(txt string) (color.RGBA, error) {
	colorParsers := make(map[string]func(string) (color.RGBA, error), 2)
	colorParsers["#"] = ParseHexColor
	colorParsers["rgb"] = ParseRGBColor
	colorText := strings.Trim(strings.ToLower(txt), " ")

	for prefix, parser := range colorParsers {
		if strings.HasPrefix(colorText, prefix) {
			return parser(colorText)
		}
	}
	return ParseNameColor(txt)
}

// Parses a hexadecimal color value.
//
// Format of text: #[0-9a-f]{3, 12}
func ParseHexColor(txt string) (color.RGBA, error) {
	valid, err := regexp.Match("^#[0-9a-f]*$", []byte(txt))
	if err != nil || !valid {
		return color.RGBA{}, errors.New("Invalid hex color pattern")
	}
	valuesTxt := []rune(SubString(txt, 1, len(txt)))
	colorParts := make([]rune, 8)
	if len(valuesTxt) == 3 || len(valuesTxt) == 4 {
		colorParts[0] = valuesTxt[0]
		colorParts[1] = valuesTxt[0]
		colorParts[2] = valuesTxt[1]
		colorParts[3] = valuesTxt[1]
		colorParts[4] = valuesTxt[2]
		colorParts[5] = valuesTxt[2]
		colorParts[6] = 'f'
		colorParts[7] = 'f'
		if len(valuesTxt) == 4 {
			colorParts[6] = valuesTxt[3]
			colorParts[7] = valuesTxt[3]
		}
	} else if len(valuesTxt) == 3*2 || len(valuesTxt) == 4*2 {
		colorParts[0] = valuesTxt[0]
		colorParts[1] = valuesTxt[1]
		colorParts[2] = valuesTxt[2]
		colorParts[3] = valuesTxt[3]
		colorParts[4] = valuesTxt[4]
		colorParts[5] = valuesTxt[5]
		colorParts[6] = 'f'
		colorParts[7] = 'f'
		if len(valuesTxt) == 4*2 {
			colorParts[6] = valuesTxt[6]
			colorParts[7] = valuesTxt[7]
		}
	} else {
		return color.RGBA{}, errors.New("Invalid hex color pattern")
	}
	r, err := strconv.ParseInt(string(colorParts[0:2]), 16, 16)
	if err != nil || r < 0 || r > 255 {
		return color.RGBA{}, err
	}
	g, err := strconv.ParseInt(string(colorParts[2:4]), 16, 16)
	if err != nil || g < 0 || g > 255 {
		return color.RGBA{}, err
	}
	b, err := strconv.ParseInt(string(colorParts[4:6]), 16, 16)
	if err != nil || b < 0 || b > 255 {
		return color.RGBA{}, err
	}
	a, err := strconv.ParseInt(string(colorParts[6:8]), 16, 16)
	if err != nil || a < 0 || a > 255 {
		return color.RGBA{}, err
	}
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}, nil
}

// Parses an rgb(a)? color value.
//
// Format of text: rgb(\d+, \d+, \d+) or rgb(\d+, \d+, \d+)
func ParseRGBColor(txt string) (color.RGBA, error) {
	valid, err := regexp.Match("^rgba?\\([\\d, ]+\\)$", []byte(txt))
	if err != nil || !valid {
		return color.RGBA{}, errors.New("Invalid rgb color pattern")
	}
	start := 4
	if txt[start-1] == 'a' {
		start++
	}
	valuesTxt := SubString(txt, start, -1)
	colorParts, err := GetCSV(valuesTxt)
	if err != nil || !(len(colorParts) == 3 || len(colorParts) == 4) {
		return color.RGBA{}, errors.New("Invalid rgb color pattern")
	}
	r, g, b, a := 0, 0, 0, 255
	r, err = strconv.Atoi(colorParts[0])
	if err != nil || r < 0 || r > 255 {
		return color.RGBA{}, errors.New("Invalid rgb color pattern")
	}
	g, err = strconv.Atoi(colorParts[1])
	if err != nil || g < 0 || g > 255 {
		return color.RGBA{}, errors.New("Invalid rgb color pattern")
	}
	b, err = strconv.Atoi(colorParts[2])
	if err != nil || b < 0 || b > 255 {
		return color.RGBA{}, errors.New("Invalid rgb color pattern")
	}
	if len(colorParts) > 3 {
		a, err = strconv.Atoi(colorParts[3])
		if err != nil || a < 0 || a > 255 {
			return color.RGBA{}, errors.New("Invalid rgb color pattern")
		}
	}
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}, nil
}

// Returns the color value of a predetermined color that
// matches the given name.
func ParseNameColor(txt string) (color.RGBA, error) {
	file, err := os.ReadFile("src/data/colors.yaml")
	if err != nil {
		return color.RGBA{}, err
	}
	colors := make(map[string]string, 5)
	err = yaml.Unmarshal(file, colors)
	for name, value := range colors {
		if name == txt {
			return ParseColor(value)
		}
	}
	return color.RGBA{}, errors.New("Invalid color pattern")
}
