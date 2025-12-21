#import "../../common/xml-utils.typ": children, extract
#import "../../common/colors.typ": light-gray, table-border

#let address-data(p) = {
  let data = (v(0.5em),)
  data.push([#{ extract(p, "Adres.AdresL1") }])
  data.push([#{ extract(p, "Adres.AdresL2") }])
  data.push([NIP: #{ extract(p, "DaneIdentyfikacyjne.NIP") }])

  return data
}

#let participant(p, role) = {
  block(
    stroke: 0.5pt + table-border,
    table(
      columns: 1fr,
      row-gutter: -4pt,
      fill: (_, y) => if y == 0 { light-gray },
      stroke: 0pt,
      [#{ role }],
      [#text({ extract(p, "DaneIdentyfikacyjne.Nazwa") }, weight: "bold")],
      ..address-data(p),
    ),
  )
}

#let participants(invoice) = {
  let columns = (1fr, 1fr)
  let podmiot1 = extract(invoice, "Podmiot1")
  let podmiot2 = extract(invoice, "Podmiot2")
  let podmiot3 = children(invoice, "Podmiot3")

  if podmiot3.len() > 0 {
    podmiot3 = podmiot3.first()
    columns = (1fr, 1fr, 1fr)
  }

  grid(
    columns: columns,
    column-gutter: 3pt,
    [#{ participant(podmiot1, "SPRZEDAWCA") }],
    [#{ participant(podmiot2, "NABYWCA") }],
    if podmiot3.len() > 0 {
      [#{ participant(podmiot3, "ODBIORCA") }]
    }
  )
}
