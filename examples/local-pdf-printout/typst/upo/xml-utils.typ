#let extract(elem, path) = {
  if path == "" {
    if elem.len() == 1 {
      return elem.first()
    }
    return elem
  }

  /// we need to recurse
  /// let's find out if elem is a dict containing children
  /// or is it already an array of children
  let children
  if type(elem) == array {
    children = elem.filter(e => "tag" in e)
  } else {
    children = elem.children.filter(e => "tag" in e)
  }

  let path_elements = path.split(".")
  let child_name = path_elements.remove(0)

  let new_path = path_elements.join(".", default: "")

  return extract(children.filter(e => e.tag == child_name).first().children, new_path)
}

#let children(node, tag) = {
  return node.children.filter(e => "tag" in e and e.tag == tag)
}
