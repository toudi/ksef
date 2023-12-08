package xml

import (
	"fmt"
	"strings"
)

func (node *Node) GetChild(name string) (*Node, error) {
	for _, child := range node.Children {
		if child.Name == name {
			return child, nil
		}
	}

	return nil, fmt.Errorf("cannot locate child")
}

func (node *Node) GetOrCreateChild(name string, isArray bool) (*Node, bool) {
	if node.Children == nil {
		node.Children = make([]*Node, 0)
	}

	for _, child := range node.Children {
		if child.Name == name && !isArray {
			return child, false
		}
	}

	newNode := &Node{Name: name}
	node.Children = append(node.Children, newNode)
	return newNode, true
}

func (node *Node) CreateChild(path string, isArray bool) (*Node, bool) {
	var created bool = false
	var target *Node = node

	pathParts := strings.Split(path, ".")
	// numParts := len(pathParts)

	// fmt.Printf("path parts: %v\n", pathParts)
	// targetName := pathParts[numParts-1]

	for _, nodeName := range pathParts {
		// fmt.Printf("target is %s\n", target.Name)
		// if i < numParts-1 {
		// fmt.Printf("i = %d; nodeName=%s\n", i, nodeName)
		target, created = target.GetOrCreateChild(nodeName, isArray)
		// fmt.Printf("set target to %s\n", target.Name)
		// }
	}

	// target, created = target.appendChild(targetName, isArray)

	return target, created
}

func (node *Node) LocateNode(path string) (*Node, error) {
	// fmt.Printf("locateNode %s\n", path)
	var target *Node = node
	var found bool

	pathParts := strings.Split(path, ".")
	numParts := len(pathParts)

	for i, part := range pathParts {
		if part == target.Name {
			// locate next target
			if i < numParts-1 {
				found = false
				for _, child := range target.Children {
					if child.Name == pathParts[i+1] {
						target = child
						found = true
						break
					}
				}
				if !found {
					return nil, fmt.Errorf("could not locate node")
				}
			}
		}
	}

	return target, nil

}

func getNodeAndAttrib(name string) (string, string) {
	_path_attrib := strings.Split(name, "#")
	if len(_path_attrib) == 2 {
		return _path_attrib[0], _path_attrib[1]
	}

	return name, ""
}

func (node *Node) SetValue(path string, value string) {
	pathParts := strings.Split(path, ".")

	var nodeName string
	var attribName string
	var target *Node = node

	// fmt.Printf("path=%s; value=%s\n", path, value)

	for _, part := range pathParts {
		if target.Name == part {
			continue
		}

		nodeName, attribName = getNodeAndAttrib(part)
		if nodeName == "" {
			// we've reached the end of the path and we have an attribute name.
			break
		}

		target, _ = target.GetOrCreateChild(nodeName, false)
	}

	if attribName != "" {
		if target.Attribs == nil {
			target.Attribs = make(map[string]string)
		}
		target.Attribs[attribName] = value
	} else {
		target.Value = value
	}
}

func (node *Node) SetValuesFromMap(data map[string]string) {
	for key, value := range data {
		node.SetValue(key, value)
	}
}

func (node *Node) SetData(path string, data map[string]string) {
	for key, value := range data {
		node.SetValue(path+"."+key, value)
	}
}

func (node *Node) ValueOf(path string) (string, error) {
	child, err := node.GetChild(path)
	if err != nil {
		return "", err
	}

	return child.Value, nil
}

func (node *Node) DeleteChild(name string) {
	var index int = -1

	for i, child := range node.Children {
		if child.Name == name {
			index = i
			break
		}
	}

	if node.Children != nil && index > 0 {
		node.Children = append(node.Children[:index], node.Children[index+1:]...)
	}
}
