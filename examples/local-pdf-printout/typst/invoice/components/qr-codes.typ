#import "@preview/tiaoma:0.3.0"

#let qr-codes(meta) = {
  let urls = meta.at("invoice").at("qr-codes")
  let columns = 1fr
  let offline = "offline" in urls
  let invoice-label = ""
  let font-size-label = 9pt
  if offline {
    columns = (1fr, 1fr)
    invoice-label = "OFFLINE"
  } else {
    invoice-label = meta.at("invoice").at("ksef-ref-no")
    font-size-label = 11pt
  }
  grid(
    columns: columns,
    grid(
      gutter: 0.5em,
      align: center,
      [#{
        tiaoma.qrcode(
          { urls.at("invoice") },
          width: 4cm,
        )
      }],
      [#text({ invoice-label }, size: font-size-label)],
    ),
    if "offline" in urls {
      grid(
        gutter: 0.5em,
        align: center,
        [#{
          tiaoma.qrcode(
            { urls.at("offline") },
            width: 4cm,
          )
        }],
        [#text([CERTYFIKAT], size: font-size-label)],
      )
    },
  )
}
