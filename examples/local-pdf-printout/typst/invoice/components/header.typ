#import "../../common/xml-utils.typ": extract

#let invoice-header(invoice) = {
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
        [Data wystawienia], [#{ extract(invoice, "Fa.P_1") }],
        [Data sprzedaży], [#{ extract(invoice, "Fa.P_6") }],
      )
    }],
  )
}
