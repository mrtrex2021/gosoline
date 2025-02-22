//go:build fixtures
// +build fixtures

package fixtures

import (
	"context"
	"fmt"

	"github.com/justtrackio/gosoline/pkg/cfg"
	"github.com/justtrackio/gosoline/pkg/log"
)

type fixtureLoaderSettings struct {
	Enabled bool `cfg:"enabled" default:"false"`
}

type fixtureLoader struct {
	logger        log.Logger
	settings      *fixtureLoaderSettings
	writerFactory func(factory FixtureWriterFactory) (FixtureWriter, error)
}

func NewFixtureLoader(ctx context.Context, config cfg.Config, logger log.Logger) FixtureLoader {
	logger = logger.WithChannel("fixture_loader")

	settings := &fixtureLoaderSettings{}
	config.UnmarshalKey("fixtures", settings)

	return &fixtureLoader{
		logger:   logger,
		settings: settings,
		writerFactory: func(factory FixtureWriterFactory) (FixtureWriter, error) {
			return factory(ctx, config, logger)
		},
	}
}

func (f *fixtureLoader) Load(ctx context.Context, fixtureSets []*FixtureSet) error {
	if !f.settings.Enabled {
		f.logger.Info("fixture loader is not enabled")
		return nil
	}

	for _, fs := range fixtureSets {
		if !fs.Enabled {
			f.logger.Info("skipping disabled fixture set")
			continue
		}

		if fs.Writer == nil {
			return fmt.Errorf("fixture set is missing a writer")
		}

		writer, err := f.writerFactory(fs.Writer)
		if err != nil {
			return fmt.Errorf("can not create writer: %w", err)
		}

		if fs.Purge {
			err := writer.Purge(ctx)
			if err != nil {
				return fmt.Errorf("error during purging of fixture set: %w", err)
			}
		}

		if err = writer.Write(ctx, fs); err != nil {
			return fmt.Errorf("error during loading of fixture set: %w", err)
		}
	}

	return nil
}
