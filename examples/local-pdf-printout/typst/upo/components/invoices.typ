#import "../../common/xml-utils.typ": children, extract
#import "../../common/colors.typ": light-gray, table-border

#let table-rows(invoices) = {
  let rows = ()

  for (index, invoice) in invoices.enumerate() {
    let row = ()
    row.push([#{ index + 1 }])
    row.push([#{ extract(invoice, "NumerKSeFDokumentu") }])
    row.push([#{ extract(invoice, "NumerFaktury") }])
    row.push([#{ extract(invoice, "NipSprzedawcy") }])
    row.push([#{ extract(invoice, "DataWystawieniaFaktury") }])
    row.push([#{ extract(invoice, "DataPrzeslaniaDokumentu") }])
    row.push([#{ extract(invoice, "DataNadaniaNumeruKSeF") }])
    row.push([#{ extract(invoice, "SkrotDokumentu") }])
    rows.push(row)
  }

  return rows
}

#let format-dt(dt) = {
  let datetime_fragment = dt.text.slice(0, 19)
  let timezone_fragment = dt.text.slice(-6)
  let parsed_dt = toml(bytes("date = " + datetime_fragment)).date
  return table.cell(
    grid(
      columns: 1,
      row-gutter: -4pt,
      inset: 5pt,
      text(9pt)[#{ parsed_dt.display("[year]-[month]-[day]") }],
      text(9pt)[#{ parsed_dt.display("[hour]:[minute]:[second]") } #{ timezone_fragment }],
    ),
  )
}

#let invoices-table(invoices) = {
  set text(size: 0.8em)
  show table.cell.where(x: 7): cell => { if cell.y > 0 { text(7pt)[#align(horizon, cell)] } else { cell } }
  show table.cell.where(x: 4): cell => { align(horizon + center, cell) }
  show table.cell: cell => { if cell.y == 0 { text(cell, weight: "bold") } else { cell } }
  show table.cell: cell => {
    if (cell.x == 5 or cell.x == 6) and cell.y > 0 { format-dt(cell.body) } else { cell }
  }
  show table.cell: cell => align(horizon + center, cell)
  table(
    columns: (auto, 26%, auto, auto, auto, auto, auto, 1fr),
    fill: (_, y) => if y == 0 { light-gray },
    stroke: 0.5pt + table-border,
    table.header(
      [L.p.],
      [Numer identyfikujący \ fakturę w KSeF],
      [Numer faktury],
      [NIP sprzedawcy],
      [Data wystawienia \ faktury],
      [Data przesłania \ do KSeF],
      [Data nadania \ numeru KSeF],
      [Wartość funkcji\  skrótu złożonego \ dokumentu],
    ),

    ..table-rows(invoices).flatten(),
  )
}
