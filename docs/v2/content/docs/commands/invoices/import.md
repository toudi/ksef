---
linkTitle: import
---

# `import`

Importuj faktury z plików CSV, XLSX, YAML lub faktury w formacie XML wygenerowane przez inny program do bazy faktur.

## Tryb "offline"

KSeF przewiduje zadeklarowanie przesyłanej faktury jako wystawionej w trybie "Offline". Aby to zrobić, użyj flagi `--offline` przy imporcie:

```
./ksef invoices import --offline [ ... ]
```

## importowanie z plików CSV / XLSX / YAML

Tą czynność opisałem dość szczegółowo w opisie komendy [`generate`](/docs/commands/generate). Jedyna różnica jest taka, że tutaj zaimportowane faktury znajdą się w hierarchi katalogu `data`.

### tryb automatycznego wystawiania korekt

Program przewiduje umożliwienie automatycznego wystawiania korekt przesyłanych faktur. Jest to zdecydowanie niebezpieczna funkcjonalność (domyślnie wyłączona) a jej włączenie jest zabezpieczone flagami, tym niemniej przewidziałem ją przez wzgląd na programy księgowe które nie są w stanie robić tego samoczynnie. Program wystawia korektę tylko i wyłącznie w jednym przypadku - kiedy zauważy, że próbujesz zaimportować po raz kolejny taką samą fakturę (tj. o takim samym numerze) ale jej zawartość (tj. pozycje) się różnią. Wówczas program wygeneruje fakturę korygującą.

Dodatkowo, jeśli w danym roku chociaż raz użyjesz flagi `--auto-correction`, wówczas program konsekwentnie będzie zapisywał dane oryginalnych faktur w rejestrze więc spodziewaj się, że jego rozmiar "spuchnie".

{{< callout type="warning" >}}
Program **NIE** wygeneruje faktury korygującej w sytuacji w której zmienił się nabywca faktury. Według instrukcji KSeF (i chyba zdrowego rozsądku) w takiej sytuacji należy najpierw wyzerować oryginalną fakturę (tj. wygenerować fakturę korygującą będącą odwrotnością oryginalnej tak aby wartość netto/brutto była równa zero) a następnie wygenerować zupełnie nową fakturę dla nowego nabywcy. O ile numeracja faktur korygujących może być inna, o tyle mój program nie jest w stanie przewidzieć numeru nowej faktury
{{< /callout >}}

Feature flagi:

| flaga                         | Opis                                                      |
| ----------------------------- | --------------------------------------------------------- |
| `--auto-correction`           | Włącz tryb automatycznego wystawiania faktur korygujących |
| `--auto-correction.numbering` | Schemat numeracji dla faktur korygujących                 |

Schemat numeracji

Dostępne placeholdery:

| Placeholder        | Opis                                               |
| ------------------ | -------------------------------------------------- |
| `{year}`           | numer roku                                         |
| `{month}`          | numer miesiąca                                     |
| `{count}`          | numer kolejny faktury korygującej w danym roku     |
| `{currMonthCount}` | numer kolejny faktury korygującej w danym miesiącu |

## importowanie faktur w formacie XML (wytworzonych przez inny program)

Jeśli generujesz faktury w formacie XML w odrębnym programie i używasz mojego programu jedynie do wysyłki, musisz je najpierw zaimportować do struktury katalogowej - tak aby program obliczył sumy kontrolne i przygotował faktury do wysyłki

{{< gateway-flags >}}

Przykłady wywołań

```
./ksef import plik-faktury.xml
./ksef import katalog/*.xml
```

Program ustali numer NIP na podstawie samego pliku XML więc nie musisz go podawać
