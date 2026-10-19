package fs

import "github.com/jaevor/go-nanoid"

var uidFactory func() string

const maxUIDGenerationAttempts = 10

func init() {
	uidFactory = nanoid.MustCustomASCII("abcdefgh0123456789", 5)
}

func (f *FSCertsStorage) generateUID() (string, error) {
	collision := true
	var uid string
	var generationAttempt int
	for collision {
		uid = uidFactory()
		_, collision = f.uidIndex[uid]
		generationAttempt += 1
		if generationAttempt > maxUIDGenerationAttempts {
			return "", ErrUIDGenerationAttemptsExceeded
		}
	}

	return uid, nil
}
