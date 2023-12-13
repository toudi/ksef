# Pobieranie wizualizacji PDF faktury

KSeF umożliwia pobranie wizualizacji PDF faktury (zapewne w celu wysłania jej kontrahentowi). Klient obsługuje dwa przypadki pobrania PDF'a z fakturą

```shell
./ksef download-pdf
Usage of download-pdf:
  -i string
    	numer faktury do pobrania
  -o string
    	ścieżka do zapisu PDF (domyślnie katalog pliku statusu + {nrRef}.pdf)
  -p string
    	ścieżka do pliku statusu
```

## Kiedy masz dostęp do źródłowego pliku XML

Jeśli wciąż masz na dysku źródłowy plik faktury wówczas generowanie PDF przebiegnie o wiele szybciej. Wydaj następującą komendę:

```shell
./ksef download-pdf -i sciezka-do-pliku.xml -p sciezka-do-pliku-status.yaml
```

Wówczas program odczyta numer faktury z XML'a, sprawdzi pod jakim numerem KSeF ta faktura została zarejestrowana a następnie pobierze PDF'a do wskazanego katalogu. Jeśli nie wskażesz żadnego katalogu wyjściowego, wówczas użyty zostanie katalog w którym znajduje się plik statusu.

## Kiedy masz dostęp do pliku statusu ale straciłeś źródłowy XML

```shell
./ksef download-pdf -i numer-referencyjny-faktury-w-ksef -p sciezka-do-pliku-status.yaml
```

Sytuacja jest gorsza ale nie beznadziejna. Program zainicjuje sesję interaktywną i pobierze źródłowy XML z KSeF a następnie użyje bramki ministerstwa do pobrania PDF. Szkopuł w tym, że aby pobrać XML trzeba chwilę poczekać na zakończenie sprawdzania uwierzytelnień. Jest to proces asynchroniczny i u mnie potrafiło to zająć grubo ponad 20 sekund.
