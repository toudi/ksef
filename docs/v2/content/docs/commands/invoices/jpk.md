---
linkTitle: jpk
weight: 10
---

# `jpk`

Generowanie deklaracji VAT w formacie JPK_V7M na podstawie zarejestrowanych faktur.

Działanie komendy jest dość oczywiste - program iteruje po wszystkich fakturach (sprzedażowych i kosztowych), oblicza VAT naliczony oraz należny i generuje deklarację JPK. Oczywiście musisz mieć na względzie to, że program działa jedynie na tych fakturach, które są zarejestrowane w KSeF. Jeśli masz inne faktury (np. zagraniczne, które nie podlegają rejestracji w KSeF) musisz wyedytować deklarację. Mozesz to zrobić korzystając z edytora przygotowanego przez ministerstwo finansów: https://e-mikrofirma.mf.gov.pl/jpk-form/read-file

---

## Rozliczenie nadwyżki VAT

Jeśli z raportu wynika, że kwota VAT naliczonego przekracza VAT należny (tj. występuje nadwyżka podatku naliczonego nad należnym), JPK_V7M przewiduje trzy tryby rozliczenia tej nadwyżki. Program domyślnie wybiera tryb przeniesienia na następny okres rozliczeniowy (tj. pole `P_62`). Jest możliwość wskazania innego zachowania poprzez linię komend lub plik konfiguracyjny. Poniżej opiszę jak to zrobić

### Alternatywne tryby rozliczenia

Dostępne są dwa alternatywne tryby, które można wybrać – **ale tylko jeden z nich na raz**, ponieważ się wzajemnie wykluczają.
Tu chciałbym zaznaczyć, że być może według instrukcji ministerstwa finansów te opcje (tj. zwrot na rachunek oraz przeksięgowanie na poczet przyszłych zobowiązań) wcale nie muszą się wzajemnie wykluczać, tym niemniej instrukcja jest tak okrutnie niezrozumiała że nie podejmuję się jej interpretacji. Być może podoła temu KIS albo polonista? W razie wątpliwości odsyłam do strony https://podatki.gov.pl/podatki-firmowe/jednolity-plik-kontrolny/jpk_vat-z-deklaracja/pytania-i-odpowiedzi (sekcja "Wypełnianie JPK VAT z deklaracją")

#### Zwrot na rachunek (refund)

Nadwyżka VAT może zostać zwrócona na wskazany rachunek bankowy. Dostępne są następujące warianty trybów zwrotu:

| Tryb zwrotu | Opis | Pole formularza |
|-------------|------| --- |
| `15-bank` | Zwrot w terminie 15 dni na rachunek bankowy | `P_540` |
| `25-vat` | Zwrot w terminie 25 dni na rachunek VAT | `P_55` |
| `25-bank` | Zwrot w terminie 25 dni na rachunek bankowy | `P_56` |
| `40-bank` | Zwrot w terminie 40 dni na rachunek bankowy | `P_560` |
| `180-bank` | Zwrot w terminie 180 dni na rachunek bankowy | `P_58` |

#### Przeksięgowanie na inne zobowiązanie podatkowe (offset-tax)

Nadwyżka VAT może zostać przeksięgowana na poczet przyszłych zobowiązań podatkowych innego rodzaju. Wymaga podania kodu zobowiązania. Wypełnione zostaną następujące pola:

- `P_59` (wpisana zostanie wartość `1`)
- `P_60` (wpisana zostanie kwota nadwyżki)
- `P_61` (wpisana zostanie wskazana nazwa zobowiązania)

Niestety tu zaczynają się schody. Według instrukcji ministerstwa finansów (cytuję):

> Aby zaliczyć nadwyżkę podatku naliczonego nad należnym (względnie części zwrotu) na poczet przyszłych zobowiązań podatkowych należy w pozycji P_54 wskazać część lub całość kwoty wykazanej w pozycji P_53. W jednym z pól P_55 – P_58 należy wybrać przysługujący termin zwrotu (nawet jeśli całość zwrotu zaliczana jest na poczet innych zobowiązań podatkowych) przez wpisanie wartości „1”. W polu P_59 należy podać wartość „1”, w P_60 podać wysokość zwrotu do zaliczenia na poczet przyszłych zobowiązań podatkowych. Następnie w P_61 opisać rodzaj przyszłego zobowiązania podatkowego. Pole P_62, czyli wysokość nadwyżki podatku naliczonego nad należnym do przeniesienia na następny okres rozliczeniowy, jest opisana jako różnica między polami P_53 i P_54. Zatem kwota w P_60 zawiera się w kwocie P_54.

Niestety dla mnie jest to kompletny bełkot. Aktualnie program wpisuje całość kwoty do zwrotu w pole `P_54`, Wartość `1` w pole `P_55` (co wskazuje konto VAT jako konto do zwrotu) a następnie wypełnia pola `P_59`, `P_60`, `P_61`. **Nie mam pojęcia czy jest to poprawne działanie**. Tak wypełniony plik JPK przechodzi walidację programem `xmllint`

---

## Konfiguracja

Tryb rozliczenia nadwyżki VAT można skonfigurować na dwa sposoby: przez flagi programu lub przez plik ustawień podmiotu.

### Flagi linii poleceń

Flagi mają pierwszeństwo nad ustawieniami z pliku:

| Flaga | Opis |
|-------|------|
| `--jpk.surplus.refund <tryb>` | Zwrot nadwyżki na rachunek – np. `25-vat`, `15-bank`, `40-bank` |
| `--jpk.surplus.offset-tax <kod>` | Przeksięgowanie na inne zobowiązanie podatkowe |

Przykład – zwrot w terminie 25 dni na rachunek VAT:

```bash
./ksef invoices jpk --jpk.surplus.refund 25-vat
```

Przykład – przeksięgowanie na inne zobowiązanie:

```bash
./ksef invoices jpk --jpk.surplus.offset-tax PPE
```

### Plik ustawień podmiotu

Ustawienia można również zapisać w pliku konfiguracyjnym podmiotu. Plik ten znajduje się w katalogu danych podmiotu i ma postać:

```yaml
jpk:
  surplus:
    # carry-over: true
    # refund: 25-vat
    offset-tax: PPE
```

{{< callout type="warning" >}}
Opcje `refund` oraz `offset-tax` wykluczają się wzajemnie – można wybrać tylko jedną z nich. Pole `carry-over` jest domyślne i można je jawnie ustawić - jeśli bardzo chcesz.
{{< /callout >}}

