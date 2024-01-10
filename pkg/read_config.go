package covplots

import (
	"os"
	"io"
	"encoding/json"
)

type InputSet struct {
	Paths []string `json:"paths"`
	Name string `json:"name"`
	Functions []string `json:"functions"`
	FunctionArgs []any `json: "functionargs"`
	Extra any `json: "extra"`
}

type UltimateConfig struct {
	InputSets []InputSet `json:"inputsets"`
	Chrlens string `json:"chrlens"`
	Outpre string `json:"outpre"`
	Ylim []float64 `json:"ylim"`
	Plotfunc string `json: "plotfunc"`
	PlotfuncArgs any `json: "plotfuncargs"`
	Fullchr bool `json: "fullchr"`
	NoParent bool `json: "noparent"`
	ManualChrs []string `json: "manualchrs"`
	ManualChrsBedPath string `json: "manualchrsbedpath"]`
}

func ReadUltimateConfig(r io.Reader) ([]UltimateConfig, error) {
	cfgbytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var cfg []UltimateConfig
	err = json.Unmarshal(cfgbytes, &cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func GetUltimateConfig(path string) ([]UltimateConfig, error) {
	if path == "" {
		return ReadUltimateConfig(os.Stdin)
	}
	r, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	return ReadUltimateConfig(r)
}
