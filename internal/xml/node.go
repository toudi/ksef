package xml

type Node struct {
	Name      string
	Value     string
	Namespace string
	Attribs   map[string]string
	Children  []*Node
}
