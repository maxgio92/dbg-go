package stats

import (
	"github.com/falcosecurity/dbg-go/pkg/root"
	"github.com/falcosecurity/dbg-go/pkg/validate"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type fileStatter struct {
	root.Looper
}

func NewFileStatter() Statter {
	return &fileStatter{Looper: root.NewFsLooper(root.BuildConfigPath)}
}

func (f *fileStatter) Info() string {
	return "gathering stats for local config files"
}

func (f *fileStatter) GetDriverStats(opts root.Options) (driverStatsByDriverVersion, error) {
	driverStatsByVersion := make(driverStatsByDriverVersion)
	err := f.LoopFiltered(opts, "computing stats", "config", func(driverVersion, configPath string) error {
		dStats := driverStatsByVersion[driverVersion]
		err := getConfigStats(&dStats, configPath)
		driverStatsByVersion[driverVersion] = dStats
		return err
	})
	return driverStatsByVersion, err
}

func getConfigStats(dStats *driverStats, configPath string) error {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	var driverkitYaml validate.DriverkitYaml
	err = yaml.Unmarshal(configData, &driverkitYaml)
	if err != nil {
		return errors.WithMessagef(err, "config: %s", configPath)
	}

	slog.Debug("fetching stats", "parsedConfig", driverkitYaml)

	if driverkitYaml.Output.Probe != "" {
		dStats.NumProbes++
	}
	if driverkitYaml.Output.Module != "" {
		dStats.NumModules++
	}
	return nil
}
