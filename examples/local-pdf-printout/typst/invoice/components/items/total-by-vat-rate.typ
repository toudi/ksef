#import "../../../common/xml-utils.typ": children, extract
#import "./item-price.typ": item-price

#let vat-rate-descriptions = (
  "np I": "np z wyłączeniem art. 100 ust 1 pkt 4 ustawy",
  "np II": "np na podstawie art. 100 ust. 1 pkt 4 ustawy",
)

#let aggregage-by-vat-rate(items) = {
  let vat-rate-buckets = (:)

  for item in items {
    let vat-rate = extract(item, "P_12")
    let amounts = item-price(item)
    let vat-rate-description = ""
    if vat-rate.contains(regex("[0-9]+")) {
      vat-rate-description = vat-rate + " %"
    } else {
      vat-rate-description = vat-rate-descriptions.at(vat-rate)
    }

    if vat-rate-description not in vat-rate-buckets {
      vat-rate-buckets.insert(vat-rate-description, (
        net: decimal(0),
        gross: decimal(0),
        vat: decimal(0),
      ))
    }
    vat-rate-buckets.at(vat-rate-description).at("net") += amounts.amount.net
    vat-rate-buckets.at(vat-rate-description).at("gross") += amounts.amount.gross
    vat-rate-buckets.at(vat-rate-description).at("vat") += amounts.amount.vat
  }

  let rows = vat-rate-buckets.pairs().sorted(key: it => (it.at(1).gross)).rev()

  let table-rows = ()
  for row in rows {
    let table-row = ()
    table-row.push([#{ row.at(0) }])
    table-row.push([#{ row.at(1).net }])
    table-row.push([#{ row.at(1).vat }])
    table-row.push([#{ row.at(1).gross }])
    table-rows.push(table-row)
  }
  return table-rows.flatten()
}
