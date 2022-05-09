package main

import (
	"time"

	"github.com/joao3101/mailchimp-data-importer/internal/config"
	"github.com/joao3101/mailchimp-data-importer/internal/sync"
	"github.com/rs/zerolog/log"
)

func main() {
	duration := time.Duration(config.Config.Ticker.Timer)
	tickerChan := time.NewTicker(time.Second * duration).C
	log.Info().Msg("Starting app")
	for {
		newSync, err := sync.NewSync(&config.Config)
		if err != nil {
			log.Panic().Msgf("failed to create a new sync: %v", err)
		}
		if err := newSync.Sync(); err != nil {
			log.Panic().Msgf("failed to run app: %v", err)
		}
		<-tickerChan
	}
}
