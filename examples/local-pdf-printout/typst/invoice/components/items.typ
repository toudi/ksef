#import "../../common/colors.typ": light-gray, table-border
#import "../../common/xml-utils.typ": children, contains, extract
#import "./items/item-price.typ": add-percent-if-numeric-rate, item-price
#import "./items/total-amounts.typ": calc-total-amounts
#import "./items/total-by-vat-rate.typ": aggregage-by-vat-rate
#import "@preview/zero:0.5.0": num, set-group
#set-group(separator: ",", threshold: 4)

#let invoice-items-rows(rows) = {
  set-group(separator: ",", threshold: 4)

  let result = ()
  for (index, row) in rows.filter(e => "tag" in e).enumerate() {
    let item_row = ()
    let item-amounts = item-price(row)
    item_row.push([#{ index + 1 }])
    item_row.push([#{ extract(row, "P_7") }])
    item_row.push([#{ if contains(row, "P_8A") { extract(row, "P_8A") } }])
    item_row.push([#{ extract(row, "P_8B") }])
    item_row.push([#{ text(num(str(item-amounts.net), digits: 2, math: false)) }])
    item_row.push([#{ item-amounts.gross }])
    item_row.push([#{ add-percent-if-numeric-rate(extract(row, "P_12")) }])
    item_row.push([#{ text(num(str(item-amounts.amount.net), digits: 2, math: false)) }])
    item_row.push([#{ text(num(str(item-amounts.amount.vat), digits: 2, math: false)) }])
    item_row.push([#{ text(num(str(item-amounts.amount.gross), digits: 2, math: false)) }])
    result.push(item_row)
  }

  return result.flatten()
}


#let items(invoice) = {
  set text(size: 8pt)
  let items = children(children(invoice, "Fa").first(), "FaWiersz")
  show table.cell.where(y: 0): cell => { align(horizon + center, text(cell, weight: "bold")) }
  show table.cell: cell => {
    if (cell.x in (3, 4, 5, 7, 8, 9)) { align(cell, right) } else if (cell.x in (0, 2, 6)) {
      align(cell, horizon + center)
    } else { cell }
  }

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
    table.cell(text([RAZEM], weight: "bold"), colspan: 7, align: left),
    [#{ text(num(str(total-amounts.net), digits: 2, math: false), weight: "bold") }],
    [#{ text(num(str(total-amounts.vat), digits: 2, math: false), weight: "bold") }],
    [#{ text(num(str(total-amounts.gross), digits: 2, math: false), weight: "bold") }],
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
