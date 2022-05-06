package importer

import (
	"fmt"

	"github.com/joao3101/mailchimp-data-importer/internal/config"
)

func Importer(conf *config.AppConfig) error {
	fmt.Println(conf)
	return nil
}
