#let default-header-footer = (
  left: [#{ "" }],
  center: [#{ align(center, [#{ "" }]) }],
  right: [#{ align(end, [#{ "" }]) }],
)

#let header-or-footer(config) = {
  let content = default-header-footer

  if config != none {
    if "left" in config {
      content.insert("left", [#{ config.left }])
    }
    if "center" in config {
      content.insert("center", [#{ align(center, config.center) }])
    }
    if "right" in config {
      content.insert("right", [#{ align(end, text(config.right)) }])
    }
  }

  return content
}

#let page-header(content) = {
  if content != default-header-footer {
    grid(
      grid(
        columns: (1fr, 1fr, 1fr),
        content.left, content.center, content.right,
      ),
      [#v(0.5em)],
      [#line(length: 100%)],
    )
  }
}

#let page-footer(content) = {
  if content != default-header-footer {
    grid(
      [#line(length: 100%)],
      [#v(0.5em)],
      grid(
        columns: (1fr, 1fr, 1fr),
        content.left, content.center, content.right,
      ),
    )
  }
}
