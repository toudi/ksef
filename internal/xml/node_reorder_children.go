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
		if children[i].Name == "FaWiersz" && children[j].Name == "FaWiersz" {
			s1 := children[i].Children[0].Value
			s2 := children[j].Children[0].Value

			if s1 < s2 {
				return true
			}

			// if s1 == s2 then let's arbitrarily position rows that are marked with "StanPrzed"
			// before these which do not have this flag
			if s1 == s2 {
				before1 := findChildOr(children[i].Children, "StanPrzed", "0")
				before2 := findChildOr(children[j].Children, "StanPrzed", "0")

				return before1 == "1" && before2 == "0"
			}
		}
		return pathOrdering[children[i].Name] < pathOrdering[children[j].Name]
	})

	for _, child := range children {
		if err = child.sortChildrenRecurse(path+"."+child.Name, ordering); err != nil {
			return fmt.Errorf("error sorting %s: %v", path, err)
		}
	}

	return nil
}

func findChildOr(children []*Node, nodeName string, defaultValue string) string {
	for _, child := range children {
		if child.Name == nodeName {
			return child.Value
		}
	}

	return defaultValue
}
