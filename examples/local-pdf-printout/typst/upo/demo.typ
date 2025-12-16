#let upo = xml("upo.xml").first()
#import "xml-utils.typ": *

#let extract-auth-context(elem) = {
  let auth_type_child = elem
    .filter(e => "tag" in e and e.tag == "IdKontekstu")
    .first()
    .children
    .filter(e => "tag" in e)
    .first()

  return (
    ctx_type: auth_type_child.tag,
    value: auth_type_child.children.first(),
  )
}

#set text(font: "CMU Sans Serif")
#set par(spacing: 0.4em)
#set page(
  paper: "a4",
  flipped: true,
  margin: (x: 1cm, bottom: 1cm),
  header: [
    #grid(
      columns: (1fr, auto),
      [#align(horizon, [*Krajowy system #text(fill: red, "e")-Faktur*])],
      grid(
        align: right,
        row-gutter: 0.4em,
        columns: 1,
        [Nazwa pełna podmiotu, któremu doręczono dokument elektroniczny: *Ministerstwo Finansów*],
        [Informacja o dokumencie: *Dokument został zarejestrowany w systemie teleinformatycznym Ministerstwa Finansów*],
      ),
    )
    #line(length: 100%)
  ],
)

#set heading(numbering: "1.")
#let light-gray = color.mix((white, 70%), (gray, 30%))
#set table(stroke: 0.2pt)

#text(size: 1.2em, [*Urzędowe poświadczenie odbioru dokumentu elektronicznego KSeF*])
#v(2em)

#let opisPotwierdzenia = extract(upo, "OpisPotwierdzenia")
#let uwierzytelnienie = extract(upo, "Uwierzytelnienie").filter(e => "tag" in e)
#let authContext = extract-auth-context(uwierzytelnienie)

#table(
  columns: 2,
  fill: (x, _) => if x == 0 { light-gray },
  stroke: 0.1pt,

  [Numer referencyjny sesji], [#{ extract(upo, "NumerReferencyjnySesji") }],

  [Strona dokumentu UPO ], [#{ extract(opisPotwierdzenia, "Strona") }],
  [Całkowita liczba stron dokumentu UPO ], [#{ extract(opisPotwierdzenia, "LiczbaStron") }],
  [Zakres dokumentów od ], [#{ extract(opisPotwierdzenia, "ZakresDokumentowOd") }],
  [Zakres dokumentów do ], [#{ extract(opisPotwierdzenia, "ZakresDokumentowDo") }],
  [Całkowita liczba dokumentów ], [#{ extract(opisPotwierdzenia, "CalkowitaLiczbaDokumentow") }],
  [Typ kontekstu ], [#{ upper(authContext.at("ctx_type")) }],
  [Identyfikator kontekstu uwierzytelnienia ], [#{ authContext.at("value") }],
  [Skrót dokumentu uwierzytelniającego ], [#{ extract(uwierzytelnienie, "SkrotDokumentuUwierzytelniajacego") }],
  [Nazwa pliku XSD struktury logicznej dotycząca przesłanego dokumentu ],
  [#{ extract(upo, "NazwaStrukturyLogicznej") }],

  [Kod formularza przedłoonego dokumentu elektronicznego ], [#{ extract(upo, "KodFormularza") }],
)

#v(1fr)

#let invoices = children(upo, "Dokument")

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

#set text(size: 0.8em)
#show table.cell.where(x: 7): cell => { if cell.y > 0 { text(7pt)[#align(horizon, cell)] } else { cell } }
#show table.cell.where(x: 4): cell => { align(horizon + center, cell) }
#show table.cell: cell => { if cell.y == 0 { text(cell, weight: "bold") } else { cell } }
#show table.cell: cell => {
  if (cell.x == 5 or cell.x == 6) and cell.y > 0 { format-dt(cell.body) } else { cell }
}
#show table.cell: cell => align(horizon + center, cell)
#table(
  columns: (auto, 26%, auto, auto, auto, auto, auto, 1fr),
  fill: (_, y) => if y == 0 { light-gray },
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
