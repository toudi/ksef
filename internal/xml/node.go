package xml

type Node struct {
	Name     string
	Value    string
	Attribs  map[string]string
	Children []*Node
}
