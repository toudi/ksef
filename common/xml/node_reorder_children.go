package xml

import (
	"fmt"
	"sort"
)

func (node *Node) ApplyOrdering(ordering map[string]map[string]int) error {
	return node.sortChildrenRecurse(node.Name, ordering)
}

func (node *Node) sortChildrenRecurse(path string, ordering map[string]map[string]int) error {
	var err error
	children := node.Children

	if children == nil {
		return nil
	}

	pathOrdering, exists := ordering[path]
	if !exists {
		fmt.Printf("[WARNING] cannot determing ordering for children within %s; no-op", path)
		return nil
	}

	sort.Slice(children, func(i, j int) bool {
		return pathOrdering[children[i].Name] < pathOrdering[children[j].Name]
	})

	for _, child := range children {
		if err = child.sortChildrenRecurse(path+"."+child.Name, ordering); err != nil {
			return fmt.Errorf("error sorting %s: %v", path, err)
		}
	}

	return nil
}
