package main

import (
	"github.com/joao3101/mailchimp-data-importer/internal/config"
	"github.com/joao3101/mailchimp-data-importer/internal/importer"
	"github.com/rs/zerolog/log"
)

func main() {
	imp := importer.NewImporter(&config.Config)
	if err := imp.Import(&config.Config); err != nil {
		log.Panic().Msgf("failed to run app: %v", err)
	}
}
