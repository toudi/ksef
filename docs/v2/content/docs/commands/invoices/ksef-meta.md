---
linkTitle: ksef-meta
---

# `ksef-meta`

{{< gateway-flags >}}

Komenda zwraca metadane KSeF faktury, w tym numer referencyjny faktury, numer KSeF oraz kody QR.

Przykładowe wywołanie:

```
./ksef invoices ksef-meta 0001-fv-12-34-01
```

## Dostępne flagi

| Flaga | Opis | Domyślna wartość |
|-------|------|------------------|
| `-f, --format` | Format wyjścia (`json`, `yaml`, `toml`) | tekstowy |
| `-o, --output` | Ścieżka pliku wyjściowego (domyślnie: stdout) | stdout |
| `-n, --nip` | Numer NIP podmiotu | - |

## Format tekstowy (domyślny)

```
Numer faktury        : 0001-fv-12-34-01
Numer faktury w KSeF : 1112223344-20260101-0011AABBCC99-A2
Kody QR
  Weryfikacyjny      : https://ksef.mf.gov.pl/client-app/invoice/...
  Offline            : https://ksef.mf.gov.pl/certificate/...
```

## Format JSON

```
./ksef invoices ksef-meta 0001-fv-12-34-01 -f json
```

```json
{
  "ref-no": "0001-fv-12-34-01",
  "ksef-ref-no": "1112223344-20260101-0011AABBCC99-A2",
  "qrcodes": {
    "invoice": "https://ksef.mf.gov.pl/client-app/invoice/...",
    "offline": "https://ksef.mf.gov.pl/certificate/..."
  }
}
```

## Format YAML

```
./ksef invoices ksef-meta 0001-fv-12-34-01 -f yaml
```

```yaml
ref-no: 0001-fv-12-34-01
ksef-ref-no: 1112223344-20260101-0011AABBCC99-A2
qrcodes:
  invoice: https://ksef.mf.gov.pl/client-app/invoice/...
  offline: https://ksef.mf.gov.pl/certificate/...
```

## Format TOML

```
./ksef invoices ksef-meta 0001-fv-12-34-01 -f toml
```

```toml
ref-no = "0001-fv-12-34-01"
ksef-ref-no = "1112223344-20260101-0011AABBCC99-A2"

[qrcodes]
invoice = "https://ksef.mf.gov.pl/client-app/invoice/..."
offline = "https://ksef.mf.gov.pl/certificate/..."
```

## Zapis do pliku

Aby zapisać metadane do pliku, użyj flagi `--output`:

```
./ksef invoices ksef-meta 0001-fv-12-34-01 -f json -o metadata.json
```
