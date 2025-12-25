package archive

import "os"

func (a *Archive) Cleanup() (err error) {
	return os.RemoveAll(a.outputDir)
}
