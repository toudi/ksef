#import "../../common/xml-utils.typ": children, contains, extract
#import "../../common/colors.typ": light-gray, table-border

#let light-table-border = 0.5pt + table-border;

#let address-data(p) = {
  let data = ()
  let ident = extract(p, "DaneIdentyfikacyjne")

  if contains(p, "Adres.AdresL1") {
    data.push([#{ extract(p, "Adres.AdresL1") }])
  }
  if contains(p, "Adres.AdresL2") {
    data.push([#{ extract(p, "Adres.AdresL2") }])
  }

  let nip = ""
  if contains(ident, "NIP") {
    nip = extract(ident, "NIP")
  } else if contains(ident, "NrVatUE") {
    nip = extract(ident, "NrVatUE")
  } else if contains(ident, "NrID") {
    nip = extract(ident, "NrID")
  }
  data.push([NIP: #{ nip }])

  return grid.cell(
    stroke: (top: none, left: light-table-border, bottom: light-table-border, right: light-table-border),
    inset: 4pt,
    align: bottom,
    grid(..data.flatten()),
  )
}

#let participant(p, role) = {
  grid.cell(
    stroke: (left: light-table-border, top: light-table-border, right: light-table-border, bottom: none),
    grid(
      columns: 1fr,
      inset: 4pt,
      fill: (_, y) => if y == 0 { light-gray },

      [#{ role }],
      [#text({ extract(p, "DaneIdentyfikacyjne.Nazwa") }, weight: "bold")],
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
    grid(
      columns: columns,
      column-gutter: 3pt,
      [#{ participant(podmiot1, "SPRZEDAWCA") }],
      [#{ participant(podmiot2, "NABYWCA") }],
      if podmiot3.len() > 0 {
        [#{ participant(podmiot3, "ODBIORCA") }]
      }
    ),
    // this looks really weird, but I simply could not force
    // typst to align the address data to the bottom of the cell
    // so instead of wasting more time I decided to add the second
    // grid just for that.
    grid(
      columns: columns,
      column-gutter: 3pt,
      [#{ address-data(podmiot1) }],
      [#{ address-data(podmiot2) }],
      if podmiot3.len() > 0 {
        [#{ address-data(podmiot3) }]
      }
    ),
  )
}
