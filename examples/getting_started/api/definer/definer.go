package definer

import (
	"context"
	"fmt"

	"github.com/justtrackio/gosoline/examples/getting_started/api/handler"
	"github.com/justtrackio/gosoline/pkg/apiserver"
	"github.com/justtrackio/gosoline/pkg/cfg"
	"github.com/justtrackio/gosoline/pkg/log"
)

func ApiDefiner(ctx context.Context, config cfg.Config, logger log.Logger) (*apiserver.Definitions, error) {
	definitions := &apiserver.Definitions{}

	euroHandler, err := handler.NewEuroHandler(ctx, config, logger)
	if err != nil {
		return nil, fmt.Errorf("can not create euroHandler: %w", err)
	}

	euroAtDateHandler, err := handler.NewEuroAtDateHandler(ctx, config, logger)
	if err != nil {
		return nil, fmt.Errorf("can not create euroAtDateHandler: %w", err)
	}

	definitions.GET("/euro/:amount/:currency", apiserver.CreateHandler(euroHandler))
	definitions.GET("/euro-at-date/:amount/:currency/:date", apiserver.CreateHandler(euroAtDateHandler))

	return definitions, nil
}
