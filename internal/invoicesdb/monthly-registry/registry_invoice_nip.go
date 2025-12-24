package monthlyregistry

import (
	"path/filepath"
	"strings"
)

const nipIndex = -3 // 3rd index from the end

func (r *Registry) GetNIP() (nip string) {
	// thanks to the fact that the directories are always following the same
	// structure, we can retrieve NIP quickly:
	//
	// data/<gateway>/<nip>/<year>/<month>
	//                ^^^^
	pathParts := strings.Split(r.dir, string(filepath.Separator))

	return pathParts[len(pathParts)+nipIndex]
}
