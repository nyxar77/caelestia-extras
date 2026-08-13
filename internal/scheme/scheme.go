package scheme

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

var hex = regexp.MustCompile(`^[0-9A-Fa-f]{6}$`)

type Scheme struct {
	Colours struct {
		Primary string `json:"primary"`
	} `json:"colours"`
	Mode string `json:"mode"`
}

func Read(path string) (Scheme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Scheme{}, fmt.Errorf("read Caelestia scheme: %w", err)
	}
	var scheme Scheme
	if err := json.Unmarshal(data, &scheme); err != nil {
		return Scheme{}, fmt.Errorf("parse Caelestia scheme: %w", err)
	}
	if !hex.MatchString(scheme.Colours.Primary) {
		return Scheme{}, fmt.Errorf("scheme primary colour must be a six-digit hex value")
	}
	if scheme.Mode != "dark" && scheme.Mode != "light" {
		return Scheme{}, fmt.Errorf("scheme mode must be dark or light")
	}
	return scheme, nil
}
