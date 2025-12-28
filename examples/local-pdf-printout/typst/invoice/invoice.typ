#let meta = yaml("meta.yaml")
#let invoice = xml("invoice.xml").first()

#import "../common/xml-utils.typ": *

#import "./components/qr-codes.typ": qr-codes
#import "./components/header.typ": invoice-header

#let header-or-footer(config) = {
  let content = (
    left: [#h(1fr)],
    center: [#h(1fr)],
    right: [#h(1fr)],
  )

  if config != none {
    if "left" in config {
      content.insert("left", [#{ config.left }])
    }
    if "center" in config {
      content.insert("center", [#{ config.center }])
    }
    if "right" in config {
      content.insert("right", [#{ config.right }])
    }
  }

  return content
}

#let page-header(meta) = {
  if "header" in meta {
    let content = header-or-footer(meta.at("header"))
    grid(
      grid(
        columns: 3,
        content.left, content.center, content.right,
      ),
      [#v(0.5em)],
      [#line(length: 100%)],
    )
  }
}

#let page-footer(meta) = {
  if "footer" in meta {
    let content = header-or-footer(meta.at("footer"))

    grid(
      [#line(length: 100%)],
      [#v(0.5em)],
      grid(
        columns: 3,
        content.left, content.center, content.right,
      ),
    )
  }
}

#set text(font: "CMU Sans Serif")
#set page(
  margin: (x: 1cm, bottom: 1cm, top: 2cm),
  header: page-header(meta),
  footer: page-footer(meta),
)

#grid(
  columns: (3fr, 2fr),
  align: horizon,
  [#{ qr-codes(meta) }], [#align(right, { invoice-header(invoice) })],
)

#import "./components/participants.typ": participants

#{ participants(invoice) }

#import "./components/items.typ": items, vat-summary

#v(2em)

#{ items(invoice) }

#v(1em)

#align(end, { vat-summary(invoice) })
