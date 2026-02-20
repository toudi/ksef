#import "../../common/colors.typ": light-gray, table-border
#import "../../common/xml-utils.typ": children, contains, extract
#import "@preview/zero:0.5.0": num, set-group
#import "./items/item-price.typ": add-percent-if-numeric-rate, item-price
#import "./items/total-by-vat-rate.typ": aggregate-by-vat-rate-numbers, empty-row-keys, vat-aggregate-total-id
#set-group(separator: ",", threshold: 4)

#let invoice-items-rows-delta(rows) = {
  return rows
    .filter(e => "tag" in e)
    .enumerate()
    .map(((index, row)) => {
      let row-keys = row.children.filter(e => "tag" in e).map(e => e.tag)
      if row-keys in empty-row-keys {
        (
          [#{ index + 1 }],
          table.cell(text([--- (usunięcie pozycji) ---]), colspan: 9, align: center),
        )
      } else {
        let item-amounts = item-price(row)
        (
          [#{ index + 1 }],
          [#{ extract(row, "P_7") }],
          [#{ if contains(row, "P_8A") { extract(row, "P_8A") } }],
          [#{ extract(row, "P_8B") }],
          [#{ text(num(str(item-amounts.net), digits: 2, math: false)) }],
          [#{ item-amounts.gross }],
          [#{ add-percent-if-numeric-rate(extract(row, "P_12")) }],
          [#{ text(num(str(item-amounts.amount.net), digits: 2, math: false)) }],
          [#{ text(num(str(item-amounts.amount.vat), digits: 2, math: false)) }],
          [#{ text(num(str(item-amounts.amount.gross), digits: 2, math: false)) }],
        )
      }
    })
    .flatten()
}

#let items-delta(invoice, before) = {
  set text(size: 8pt)
  let items = children(children(invoice, "Fa").first(), "FaWiersz").filter(e => "tag" in e)

  if before {
    items = items.filter(
      e => e
        .children
        .filter(c => "tag" in c)
        .map(
          c => c.tag,
        )
        .contains("StanPrzed"),
    )
  } else {
    items = items.filter(
      e => {
        let children-keys = e
          .children
          .filter(c => "tag" in c)
          .map(
            c => c.tag,
          )
        let _foo = "StanPrzed" not in children-keys
        return _foo
      },
    )
  }
  show table.cell.where(y: 0): cell => { align(horizon + center, text(cell, weight: "bold")) }
  show table.cell: cell => {
    if (cell.x in (3, 4, 5, 7, 8, 9)) { align(cell, right) } else if (cell.x in (0, 2, 6)) {
      align(cell, horizon + center)
    } else { cell }
  }

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

    ..invoice-items-rows-delta(items),
  )
}

#let totals-summary(items) = {
  let vat-rate-buckets = aggregate-by-vat-rate-numbers(items)

  for (vat-rate, subtotal) in vat-rate-buckets {
    if vat-rate == vat-aggregate-total-id {
      continue
    }

    (
      [#{ vat-rate }],
      [#{ text(num(str(subtotal.at("net")), digits: 2, math: false)) }],
      [#{ text(num(str(subtotal.at("vat")), digits: 2, math: false)) }],
      [#{ text(num(str(subtotal.at("gross")), digits: 2, math: false)) }],
    )
  }
  let total = vat-rate-buckets.at(vat-aggregate-total-id)
  (
    [*RAZEM*],
    [#{ text(num(str(total.at("net")), digits: 2, math: false), weight: "bold") }],
    [#{ text(num(str(total.at("vat")), digits: 2, math: false), weight: "bold") }],
    [#{ text(num(str(total.at("gross")), digits: 2, math: false), weight: "bold") }],
  )
}

#let delta-summary(invoice) = {
  set text(size: 8pt)
  let items = children(children(invoice, "Fa").first(), "FaWiersz")

  show table.cell.where(y: 0): cell => { align(horizon + center, text(cell, weight: "bold")) }

  align(end, table(
    columns: (auto, auto, auto, auto),
    stroke: 0.5pt + table-border,
    fill: (_, y) => if y == 0 { light-gray },
    [], [Różnica (netto)], [Różnica (VAT)], [Różnica (brutto)],
    ..totals-summary(items),
  ))
}
