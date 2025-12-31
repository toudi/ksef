#import "../../common/xml-utils.typ": contains, extract

#let invoice-header(invoice) = {
  let header_items = ()
  header_items.push([Data wystawienia])
  header_items.push([#{ extract(invoice, "Fa.P_1") }])

  if contains(invoice, "Fa.P_6") {
    header_items.push([Data sprzedaży])
    header_items.push([#{ extract(invoice, "Fa.P_6") }])
  }
  grid(
    columns: 1,
    align: horizon,
    [= Faktura VAT],
    v(0.2em),
    [=== #{ extract(invoice, "Fa.P_2") }],
    v(3em),
    [#{
      table(
        columns: 2,
        stroke: 0em,
        gutter: -0.5em,
        ..header_items.flatten(),
      )
    }],
  )
}
