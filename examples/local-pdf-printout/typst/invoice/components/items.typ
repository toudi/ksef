#import "../../common/colors.typ": light-gray, table-border
#import "../../common/xml-utils.typ": children, extract
#import "./items/item-price.typ": item-price
#import "./items/total-amounts.typ": calc-total-amounts
#import "./items/total-by-vat-rate.typ": aggregage-by-vat-rate

#let invoice-items-rows(rows) = {
  let result = ()
  for (index, row) in rows.filter(e => "tag" in e).enumerate() {
    let item_row = ()
    let item-amounts = item-price(row)
    item_row.push([#{ index + 1 }])
    item_row.push([#{ extract(row, "P_7") }])
    item_row.push([#{ extract(row, "P_8A") }])
    item_row.push([#{ extract(row, "P_8B") }])
    item_row.push([#{ item-amounts.net }])
    item_row.push([#{ item-amounts.gross }])
    item_row.push([#{ extract(row, "P_12") }])
    item_row.push([#{ item-amounts.amount.net }])
    item_row.push([#{ item-amounts.amount.vat }])
    item_row.push([#{ item-amounts.amount.gross }])
    result.push(item_row)
  }

  return result.flatten()
}


#let items(invoice) = {
  set text(size: 8pt)
  let items = children(children(invoice, "Fa").first(), "FaWiersz")
  show table.cell.where(y: 0): cell => { align(horizon + center, text(cell, weight: "bold")) }
  show table.cell: cell => { if (cell.x in (4, 5, 7, 8, 9)) { align(cell, right) } else { cell } }

  let total-amounts = calc-total-amounts(items)

  table(
    columns: (auto, 1fr, 10%, auto, 10%, 10%, auto, 10%, 10%, 10%),
    fill: (_, y) => if y == 0 { light-gray },
    stroke: 0.5pt + table-border,
    [L.p.],
    [Nazwa],
    [Jednostka],
    [Ilość],
    [Cena \ jedn.\ netto],
    [Cena \ jedn.\ brutto],
    [Stawka VAT],
    [Wartość \ sprzedaży \ netto],
    [Wartość \ VAT],
    [Wartość \ sprzedaży \ brutto],
    ..invoice-items-rows(items),
    table.cell(text([RAZEM], weight: "bold"), colspan: 7),
    [#{ text([#{ total-amounts.net }], weight: "bold") }],
    [#{ text([#{ total-amounts.vat }], weight: "bold") }],
    [#{ text([#{ total-amounts.gross }], weight: "bold") }],
  )
}


#let vat-summary(invoice) = {
  set text(size: 8pt)
  [#{ text([Podsumowanie stawek podatku], size: 10pt) }]
  let items = children(children(invoice, "Fa").first(), "FaWiersz")
  table(
    columns: (auto, auto, auto, auto),
    fill: (_, y) => if y == 0 { light-gray },
    stroke: 0.5pt + table-border,
    [Stawka], [Kwota netto], [Kwota podatku], [Kwota brutto],
    ..aggregage-by-vat-rate(items),
  )
}
