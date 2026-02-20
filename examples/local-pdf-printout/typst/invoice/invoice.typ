#let meta = yaml("meta.yaml")
#let invoice = xml("invoice.xml").first()

#import "../common/xml-utils.typ": *

#import "./components/qr-codes.typ": qr-codes
#import "./components/header.typ": invoice-header
#import "./settings/page.typ": default-header-footer, header-or-footer, page-footer, page-header

#set text(font: "CMU Sans Serif")
#let page-settings = meta.at("page", default: (:))
#let header-content = header-or-footer(page-settings.at("header", default: (:)))
#let footer-content = header-or-footer(page-settings.at("footer", default: (:)))

#let top-margin = 1cm
#if header-content != default-header-footer {
  top-margin = 2cm
}

#set page(
  margin: (x: 1cm, bottom: 1cm, top: top-margin),
  header: page-header(header-content),
  footer: page-footer(footer-content),
)

#grid(
  columns: (3fr, 2fr),
  align: horizon,
  [#{ qr-codes(meta) }], [#align(right, { invoice-header(invoice) })],
)

#import "./components/participants.typ": participants

#{ participants(invoice) }

#if extract(invoice, "Fa.RodzajFaktury") == "KOR" {
  import "./components/items-delta.typ": delta-summary, items-delta

  [== Przed korektą]
  { items-delta(invoice, true) }
  [== Po korekcie]
  { items-delta(invoice, false) }
  [== Wartość korekty]
  delta-summary(invoice)
} else {
  import "./components/items.typ": items, vat-summary

  v(2em)

  { items(invoice) }

  v(1em)

  align(end, { vat-summary(invoice) })
}
