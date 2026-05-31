---
linkTitle: dump
---

# `dump`

Generowanie podsumowania miesiąca

Kiedy zakończysz oznaczanie pozycji faktur flagami (komenda `annotate`) i jesteś w gotowości aby wysłać raport do księgowych - skorzystaj z komendy `dump`. Zbiera ona wszystkie faktury do jednej paczki i dodatkowo generuje PDF z adnotacjami.

Komenda przewiduje dwa tryby pliku wyjściowego:

### `dump` – archiwum ZIP

```bash
./ksef invoices dump [rok] [miesiąc]
```

(domyślnie raport jest generowany dla poprzedniego miesiąca)

Komenda tworzy archiwum ZIP, które zawiera:

- pliki PDF wszystkich faktur z danego miesiąca (faktury kwalifikujące się do
  JPK – wystawione i otrzymane, które nie są w całości wyłączone)
- opcjonalnie pliki XML (flaga `--xml` / `-x`)
- plik **`accountant-notes.pdf`** – jeśli w danym miesiącu istnieją faktury
  z adnotacjami

Opóźnione wywołanie dla stycznia 2026:

```bash
./ksef invoices dump 2026 01
```

Z flagą `-x` (dołącz pliki XML):

```bash
./ksef invoices dump 2026 01 --xml
```

Niestandardowa ścieżka pliku wyjściowego:

```bash
./ksef invoices dump 2026 01 -o archiwum-styczen.zip
```

### `dump pdf` – scalony plik PDF

```bash
./ksef invoices dump pdf [rok] [miesiąc]
```

Komenda łączy wszystkie odnalezione pliki PDF faktur w jeden plik PDF. Jeśli
w danym miesiącu istnieją faktury z adnotacjami, na początku dokumentu
zostaje umieszczona strona **accountant-notes** (tabela z pozycjami i
uwagami). Kolejność w scalonym pliku:

1. **accountant-notes.pdf** (jeśli są adnotacje)
2. **faktury** (w kolejności z rejestru)

```bash
./ksef invoices dump pdf 2026 01 -o wszystkie-faktury-styczen.pdf
```

Domyślna nazwa pliku: `invoices-merged-{rok}-{miesiąc}.pdf`.

---

## Konfiguracja szablonu Typst

Raport accountant-notes jest generowany za pomocą
[Typst](https://typst.app/) — systemu składu dokumentów. Szablon domyślnie
znajduje się w:

```
examples/local-pdf-printout/typst/annotation/accountant-notes.typ
```

Konfiguracja (w `config.yaml` lub przez zmienne środowiskowe):

| Klucz konfiguracji             | Opis                                        | Domyślna wartość                                                                        |
|--------------------------------|---------------------------------------------|-----------------------------------------------------------------------------------------|
| `monthly-dump.typst-template-path` | Ścieżka do pliku szablonu Typst            | `examples/local-pdf-printout/typst/annotation/accountant-notes.typ`                     |
| `monthly-dump.generator`           | Nazwa generatora widoczna w stopce raportu | `WSI Pegasus`                                                                           |

Przykład konfiguracji w `config.yaml`:

```yaml
monthly-dump:
  typst-template-path: /ścieżka/do/mojego-szablonu.typ
  generator: "Moja Firma Sp. z o.o."
```
