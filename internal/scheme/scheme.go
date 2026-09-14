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
		Values  string `json:"-"`
	} `json:"colours"`
	Mode string `json:"mode"`
}

func (s *Scheme) UnmarshalJSON(data []byte) error {
	type rawScheme struct {
		Colours map[string]string `json:"colours"`
		Mode    string            `json:"mode"`
	}
	var raw rawScheme
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	encoded, err := json.Marshal(raw.Colours)
	if err != nil {
		return err
	}
	s.Colours.Values = string(encoded)
	s.Colours.Primary = raw.Colours["primary"]
	s.Mode = raw.Mode
	return nil
}

func (s Scheme) Color(name string) string {
	if s.Colours.Values == "" {
		return ""
	}
	var colours map[string]string
	if err := json.Unmarshal([]byte(s.Colours.Values), &colours); err != nil {
		return ""
	}
	return colours[name]
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
