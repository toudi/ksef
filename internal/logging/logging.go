package logging

import (
	"io"

	"github.com/rs/zerolog"
)

var SeiLogger zerolog.Logger

func InitLogging(output io.Writer) {
	SeiLogger = zerolog.New(output)
}
