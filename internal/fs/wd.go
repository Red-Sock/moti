package fs

import (
	"os"

	"github.com/rs/zerolog/log"
)

func GetWdOrDie() string {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("could not get working directory")
	}

	return workDir
}
