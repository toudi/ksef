#set text(font: "CMU Sans Serif")
#set page(flipped: true, margin: 1.5cm)
#let annotations = yaml("annotations.yaml")

#let light-gray = color.mix((white, 70%), (gray, 30%))
#let table-border = rgb("666675")


#let cell(annotation, key) = {
  let empty = table.cell([---], align: horizon + center)
  if key in annotation {
    ([#{ annotation.at(key) }],)
  } else {
    (empty,)
  }
}

#let invoice-annotations(annotations) = {
  for annotation in annotations {
    ([#{ annotation.seller }],)
    ([#{ annotation.invoice }],)
    cell(annotation, "item-name")
    cell(annotation, "notes")
  }
}

#show table.cell.where(y: 0): cell => { align(horizon + center, text(cell, weight: "bold")) }

#table(
  fill: (_, y) => if y == 0 { rgb("#f0f0f0") },
  stroke: 0.5pt + table-border,
  columns: (1fr, auto, 1fr, 1fr),
  [Sprzedawca], [Faktura], [Pozycja], [Adnotacje],
  ..invoice-annotations(annotations.annotations),
)

#align(bottom, grid(
  columns: (1fr, 1fr),
  align: (left, right),
  [#{ annotations.metadata.report-date }], [Sporządzono w programie #{ annotations.metadata.generator }],
))
