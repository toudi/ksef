package xml

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/samber/lo"
)

func (node *Node) DumpToWriter(file io.Writer, indent int) error {
	prefix := strings.Repeat("  ", indent)
	tailPrefix := ""
	var err error

	fmt.Fprintf(file, "%s<%s", prefix, node.Name)

	if node.Attribs != nil {
		// unfortunetely we have to rely on this hack since iterating over maps is
		// non-deterministic in go by design. and we need a deterministic order in
		// order to achieve the same checksum
		var attribNames = lo.Keys(node.Attribs)
		slices.Sort(attribNames)

		var attribs []string

		for _, name := range attribNames {
			attribs = append(attribs, fmt.Sprintf("%s=\"%s\"", name, node.Attribs[name]))
		}
		fmt.Fprintf(file, " %s", strings.Join(attribs, " "))
	}
	fmt.Fprintf(file, ">")
	if node.Value != "" {
		fmt.Fprintf(file, "%s", node.Value)
	}
	if node.Children != nil {
		fmt.Fprintf(file, "\n")
		tailPrefix = prefix

		for _, child := range node.Children {
			if err = child.DumpToWriter(file, indent+1); err != nil {
				return fmt.Errorf("error dumping child: %v", err)
			}
		}
	}
	fmt.Fprintf(file, "%s</%s>\n", tailPrefix, node.Name)

	return nil
}
