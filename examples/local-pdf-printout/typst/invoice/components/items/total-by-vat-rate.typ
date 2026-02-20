#import "../../../common/xml-utils.typ": children, contains, extract
#import "./item-price.typ": item-price
#import "@preview/zero:0.5.0": num, set-group
#set-group(separator: ",", threshold: 4)

#let vat-rate-descriptions = (
  "np I": "np z wyłączeniem art. 100 ust 1 pkt 4 ustawy",
  "np II": "np na podstawie art. 100 ust. 1 pkt 4 ustawy",
)

#let vat-aggregate-total-id = "--total--"

#let empty-row-keys = (
  ("NrWierszaFa", "StanPrzed"),
  ("NrWierszaFa",),
)


// this function merely prepares the aggregates that are later used by the
// rendering functions
#let aggregate-by-vat-rate-numbers(items) = {
  let vat-rate-buckets = (
    (vat-aggregate-total-id): (
      net: decimal(0),
      gross: decimal(0),
      vat: decimal(0),
    ),
  )

  for item in items {
    let item-keys = item.children.filter(e => "tag" in e).map(e => e.tag)
    let is-before = contains(item, "StanPrzed")
    let empty-row = item-keys in empty-row-keys
    if empty-row {
      continue
    }
    let vat-rate = extract(item, "P_12")
    let amounts = item-price(item)
    // initialize the description to vat rate itself
    let vat-rate-description = vat-rate
    if vat-rate.match(regex("^[0-9]+$")) != none {
      // if it's a number, let's add the percentage sign and let it be
      // the description
      vat-rate-description = vat-rate + " %"
    } else {
      // it's not a number - let's check if we have some special
      // rendering for it
      if vat-rate in vat-rate-descriptions {
        vat-rate-description = vat-rate-descriptions.at(vat-rate)
      }
    }

    if vat-rate-description not in vat-rate-buckets {
      vat-rate-buckets.insert(vat-rate-description, (
        net: decimal(0),
        gross: decimal(0),
        vat: decimal(0),
      ))
    }
    if is-before {
      amounts.amount.net *= -1
      amounts.amount.gross *= -1
      amounts.amount.vat *= -1
    }
    if vat-rate == "8" {
      let amounts-clone = amounts
    }
    vat-rate-buckets.at(vat-rate-description).at("net") += amounts.amount.net
    vat-rate-buckets.at(vat-rate-description).at("gross") += amounts.amount.gross
    vat-rate-buckets.at(vat-rate-description).at("vat") += amounts.amount.vat

    vat-rate-buckets.at(vat-aggregate-total-id).at("net") += amounts.amount.net
    vat-rate-buckets.at(vat-aggregate-total-id).at("gross") += amounts.amount.gross
    vat-rate-buckets.at(vat-aggregate-total-id).at("vat") += amounts.amount.vat
  }

  return vat-rate-buckets
}

#let aggregage-by-vat-rate(items) = {
  let vat-rate-buckets = aggregate-by-vat-rate-numbers(items)
  let rows = vat-rate-buckets.pairs().sorted(key: it => (it.at(1).gross)).rev()

  let table-rows = ()
  for row in rows {
    if row.at(0) == vat-aggregate-total-id {
      continue
    }
    let table-row = ()
    table-row.push([#{ row.at(0) }])
    table-row.push([#{ text(num(str(row.at(1).net), digits: 2, math: false)) }])
    table-row.push([#{ text(num(str(row.at(1).vat), digits: 2, math: false)) }])
    table-row.push([#{ text(num(str(row.at(1).gross), digits: 2, math: false)) }])
    table-rows.push(table-row)
  }
  return table-rows.flatten()
}
