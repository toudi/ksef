---
linkTitle: annotate
weight: 6
---

# `annotate`

Zarządzanie adnotacjami faktur zakupowych – opisywanie pozycji faktur flagami i
komentarzami

---

## Cel i powód istnienia

Podczas pracy z fakturami i przekazywaniem ich do księgowości napotkałem na następujące sytuacje:

- Faktura której pozycja jest niezwykle enigmatyczna. Przykładowo - zakup stacji dokującej gdzie opisem jest jej numer modelu. Dla księgowych jest to spory problem, bo nie wiedzą co to w istocie jest
- Faktura której nie chcę wrzucać w koszty
- Faktura gdzie dozwolone jest odliczenie tylko 50% VAT (akurat w takim przypadku podeujrzewam, że księgowi wiedzą o wiele lepiej ode mnie jak rozpoznawać takie przypadki)

Dotychczas (tj. przed wprowadzeniem KSeF) do adnotacji wystarczało zapisanie małego komentarza na wydruku. Jednak kiedy faktury są elektroniczne a KSeF w ogóle nie przewiduje takiej funkcjonalności - postanowiłem dorobić taką funkcjonalność do programu.

Adnotacje są przechowywane w bazie faktur i mogą być lokalne (dotyczą
konkretnej faktury) albo globalne (ustawione dla danego sprzedawcy – stosowane
automatycznie do wszystkich faktur od tego kontrahenta).

---

## Związek z generowaniem JPK_V7M

Podczas generowania JPK program odczytuje rejestr faktur, a dla każdej pozycji
zakupu sprawdza, czy istnieje dla niej adnotacja. W zależności od flagi:

| Flaga               | Wpływ na JPK                                                                                                                                                 |
|---------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `exclude`           | Pozycja jest **pomijana** – nie trafia do JPK. Dotyczy np. zakupów prywatnych.                                                                              |
| `vat-50-percent`    | Do JPK trafia tylko **połowa** VAT-u. Dotyczy zakupów pojazdów, paliwa itp.                                        |
| `fixed-asset`       | Pozycja jest klasyfikowana jako **środek trwały**. Kwoty bazowe i VAT trafiają do pól **K_40/K_41** zamiast K_42/K_43.                                     |
| `comment`           | Nie wpływa na wartości JPK – komentarz jest widoczny przy wydruku **accountant-notes** (raport dla księgowego).                                             |

---

## Podstawowe użycie

Wszystkie komendy z grupy `annotate` oczekują ścieżki do pliku XML faktury
zakupowej (i - co mam nadzieję jest logiczne - **tylko** zakupowej. W przypadku próby wywołania komendy `annotate` na fakturze wystawionej program nie pozwoli na to i wyświetli komunikat błędu):

```bash
./ksef invoices annotate <komenda> <faktura.xml> [flagi]
```

W każdej komendzie możesz wskazać pozycje faktury za pomocą flagi
`--items` (lista numerów pozycji – od 1) lub znaku `*` (wszystkie pozycje).

---

## `comment` – dodawanie komentarza tekstowego

```bash
./ksef invoices annotate comment \
  --items 1 \
  --text "Myszka do laptopa" \
  data/.../otrzymane/0001-moja-firma-fv-xyz.xml
```

{{< callout type="info" >}}
Flaga `--text` jest wymagana.
{{< /callout >}}

---

## `exclude` – wyłączenie pozycji z JPK

Oznacza pozycje jako wyłączone z raportu JPK (zakup prywatny, niepodlegający
odliczeniu). 

```bash
./ksef invoices annotate exclude \
  --items 2 \
  data/.../otrzymane/0001-moja-firma-fv-xyz.xml
```

---

## `vat-50` – ograniczenie odliczenia VAT do 50%

Przydatne przy zakupie paliwa do pojazdów lub leasingu tychże, gdzie ustawa ogranicza odliczenie
VAT do 50%. 

```bash
./ksef invoices annotate vat-50 \
  --items 1 \
  data/.../otrzymane/0001-moja-firma-fv-xyz.xml
```

---

## `fixed-asset` – oznaczenie środka trwałego

Pozycja oznaczona jako środek trwały trafia do odrębnych pól deklaracji VAT
(K_40 – podstawa, K_41 – VAT zamiast domyślnych K_42/K_43).

```bash
./ksef invoices annotate fixed-asset \
  --items 1 \
  data/.../otrzymane/0001-moja-firma-fv-xyz.xml
```

---

## `list` – podgląd adnotacji

Wyświetla tabelę z pozycjami faktury oraz ich adnotacjami.

```bash
./ksef invoices annotate list data/.../otrzymane/0001-moja-firma-fv-xyz.xml
```

Przykładowe wyjście:

```
+--------+-----------------------------+------------+--------------------------------+
| wiersz |            nazwa            | stawka VAT |             uwagi              |
+--------+-----------------------------+------------+--------------------------------+
| 1      | USB-C Multifunction Adapter | 23         | Stacja dokująca USB do laptopa |
+--------+-----------------------------+------------+--------------------------------+
```

---

## `clear` – usunięcie wszystkich adnotacji

Usuwa wszystkie adnotacje z danej faktury:

```bash
./ksef invoices annotate clear data/.../otrzymane/0001-moja-firma-fv-xyz.xml
```

---

## Adnotacje globalne i lokalne

Domyślnie adnotacje są zapisywane lokalnie – w rejestrze miesięcznym (`registry.yaml`). Można
je jednak zapisać globalnie, w ustawieniach podmiotu, za pomocą flagi `--global`
(`-g`). Adnotacje globalne są wiązane z NIP-em sprzedawcy i automatycznie
stosowane do wszystkich faktur od tego kontrahenta. Funkcję tę przewidziałem na wypadek leasingów bądź np. tankowania u tego samego sprzedawcy aby nie trzeba było co miesiąc powielać identycznych flag.

Przykład – globalne wyłączenie pozycji „Kawa" od kontrahenta o NIP 1111111111:

```bash
./ksef invoices annotate exclude \
  --items 3 \
  --global \
  data/.../otrzymane/0001-moja-firma-fv-xyz.xml
```

Od tej pory każda faktura od tego sprzedawcy, która zawiera pozycję pasującą
do hasha (nazwa, indeks, GTIN lub PKWiU), automatycznie otrzyma adnotację
`exclude`.

---

## Generowanie podsumowania miesiąca

Adnotacje są zbierane i prezentowane w formie raportu **accountant-notes**
(pdf z tabelą pozycji i ich adnotacji). Raport ten jest generowany przy okazji
dwóch komend z grupy `dump`:

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
2. **faktury po kolei** (w kolejności z rejestru)

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

### Szablon `accountant-notes.typ`

Poniżej znajduje się domyślny szablon Typst. Tworzy on tabelę z czterema
kolumnami: sprzedawca, faktura, pozycja, adnotacje. Dane są odczytywane
z pliku `annotations.yaml`, który program generuje automatycznie w
katalogu tymczasowym.

```typst
#set text(font: "CMU Sans Serif")
#set page(flipped: true, margin: 1.5cm)
#let annotations = yaml("annotations.yaml")

#let light-gray = color.mix((white, 70%), (gray, 30%))
#let table-border = rgb("666675")

#let cell(annotation, key) = {
  let empty = table.cell([---], align: horizon + center)
  if key in annotation {
    ([#{ annotation.at(key) }],)
  } else {
    (empty,)
  }
}

#let invoice-annotations(annotations) = {
  for annotation in annotations {
    ([#{ annotation.seller }],)
    ([#{ annotation.invoice }],)
    cell(annotation, "item-name")
    cell(annotation, "notes")
  }
}

#show table.cell.where(y: 0): cell => { align(horizon + center, text(cell, weight: "bold")) }

#table(
  fill: (_, y) => if y == 0 { rgb("#f0f0f0") },
  stroke: 0.5pt + table-border,
  columns: (1fr, auto, 1fr, 1fr),
  [Sprzedawca], [Faktura], [Pozycja], [Adnotacje],
  ..invoice-annotations(annotations.annotations),
)

#align(bottom, grid(
  columns: (1fr, 1fr),
  align: (left, right),
  [#{ annotations.metadata.report-date }], [Sporządzono w programie #{ annotations.metadata.generator }],
))
```

Szablon oczekuje w swoim katalogu pliku `annotations.yaml` o następującej
strukturze (generowanej automatycznie przez program):

```yaml
metadata:
  report-date: "2026-01-31"
  generator: "WSI Pegasus"
annotations:
  - seller: "Dostawca Sp. z o.o."
    invoice: "FV 01/2026"
    item: "1"
    item-name: "Materiały biurowe"
    notes: "50% VAT, środek trwały"
  - seller: "Inny Dostawca"
    invoice: "FV 02/2026"
    item: "2"
    item-name: "Paliwo"
    notes: "50% VAT"
```

---

## Przepływ pracy – przykład

1. Importujesz fakturę zakupową.
2. Przeglądasz pozycje za pomocą `annotate list`.
3. Oznaczasz wybrane pozycje flagami:

   ```bash
   ./ksef invoices annotate exclude --items 3 faktura.xml
   ./ksef invoices annotate vat-50  --items 1 faktura.xml
   ./ksef invoices annotate comment --items 2 --text "Wyjaśnić z dostawcą" faktura.xml
   ```

4. Pod koniec miesiąca generujesz JPK — program uwzględnia adnotacje
   (pomija wyłączone pozycje, stosuje 50% VAT, odpowiednio klasyfikuje
   środki trwałe).
5. Generujesz podsumowanie dla księgowego:

   ```bash
   ./ksef invoices dump 2026 01
   ```

   lub scalony PDF:

   ```bash
   ./ksef invoices dump pdf 2026 01
   ```

6. Archiwum (lub scalony PDF) wysyłasz do biura rachunkowego wraz z raportem
   accountant-notes.

---

## Dostępne flagi

### Wspólne dla `comment`, `exclude`, `vat-50`, `fixed-asset`

| Flaga                | Skrót | Opis                                           | Domyślna wartość |
|----------------------|-------|------------------------------------------------|------------------|
| `--items`            |       | Numery pozycji (od 1) lub `*` (wszystkie)      |                  |
| `--global`           | `-g`  | Zapisz adnotację w ustawieniach podmiotu (globalna) | `false`     |

### `comment`

| Flaga                | Opis                 | Wymagana |
|----------------------|----------------------|----------|
| `--text`             | Treść komentarza     | tak      |

### `dump`

| Flaga                | Skrót | Opis                                           | Domyślna wartość          |
|----------------------|-------|------------------------------------------------|---------------------------|
| `-n, --nip`          |       | Numer NIP podmiotu                             |                           |
| `-o, --output`       |       | Ścieżka pliku wyjściowego                      | `invoices-dump-{r}-{m}.zip` |
| `-x, --xml`          |       | Uwzględnij pliki XML w archiwum                | `false`                   |

### `dump pdf`

| Flaga                | Skrót | Opis                                           | Domyślna wartość              |
|----------------------|-------|------------------------------------------------|-------------------------------|
| `-n, --nip`          |       | Numer NIP podmiotu                             |                               |
| `-o, --output`       |       | Ścieżka pliku wyjściowego                      | `invoices-merged-{r}-{m}.pdf` |
