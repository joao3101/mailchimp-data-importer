package main

import (
	"time"

	"github.com/joao3101/mailchimp-data-importer/internal/config"
	"github.com/joao3101/mailchimp-data-importer/internal/importer"
	"github.com/rs/zerolog/log"
)

func main() {
	duration := time.Duration(config.Config.Ticker.Timer)
	tickerChan := time.NewTicker(time.Second * duration).C
	log.Info().Msg("Starting app")
	for {
		imp := importer.NewImporter(&config.Config)
		if err := imp.Import(&config.Config); err != nil {
			log.Panic().Msgf("failed to run app: %v", err)
		}
		<-tickerChan
	}
}
