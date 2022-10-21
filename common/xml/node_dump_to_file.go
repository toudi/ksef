package xml

import (
	"fmt"
	"os"
	"strings"
)

func (node *Node) DumpToFile(file *os.File, indent int) error {
	prefix := strings.Repeat("  ", indent)
	tailPrefix := ""
	var err error

	fmt.Fprintf(file, "%s<%s", prefix, node.Name)

	if node.Attribs != nil {
		var attribs = []string{}
		for name, value := range node.Attribs {
			attribs = append(attribs, fmt.Sprintf("%s=\"%s\"", name, value))
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
			if err = child.DumpToFile(file, indent+1); err != nil {
				return fmt.Errorf("error dumping child: %v", err)
			}
		}
	}
	fmt.Fprintf(file, "%s</%s>\n", tailPrefix, node.Name)

	return nil
}
