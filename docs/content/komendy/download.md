# Pobieranie faktur

```text
./ksef download
Usage of download:
  -cost
    	Synchronizuj faktury kosztowe (Podmiot2)
  -d string
    	Katalog docelowy
  -income
    	Synchronizuj faktury przychodowe (Podmiot1)
  -nip string
    	Numer NIP podmiotu
  -refresh string
    	odświeża istniejący rejestr faktur według istniejącego pliku
  -start-date string
    	Data początkowa
  -subject3
    	Synchronizuj faktury podmiotu innego (Podmiot3)
  -subjectAuthorized
    	Synchronizuj faktury podmiotu upoważnionego (???)
  -t	użyj bramki testowej
```

Istnieje kilka przewidzianych trybów dla tej komendy.

## Pobieranie listy faktur

Najpierw wybierz rodzaj faktur który Cię interesuje:

| przełącznik          | znaczenie                                                                                                                           |
| -------------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| `-income`            | Faktury przychodowe, tj. wystawione przez Ciebie                                                                                    |
| `-cost`              | Faktury kosztowe                                                                                                                    |
| `-subject3`          | Faktury gdzie występujesz jako strona 3. Z tego co się orientuję jest to sytuacja taka kiedy jesteś zamawiającym ale nie płatnikiem |
| `-subjectAuthorized` | Nie mam zielonego pojęcia do czego to służy. Czyżby jakaś opcja dla księgowych żeby mogli pobierać faktury swoich klientów?         |

parametr `-start-date` określa datę początkową filtrowania faktur. Jest to o tyle istotne, że ta wartość zostanie zapisana do pliku rejestru (więcej o tym w sekcji poniżej). Data końcowa to zawsze czas bieżący

Przykładowe wywołanie:

```shell
./ksef download -t -nip 1111111111 -cost -d kosztowe/2023-12 -start-date 2023-12-01
```

## Synchronizowanie już pobranej listy faktur

Ta opcja przewidziana jest dla sytuacji w której chcesz kilka razy w miesiącu synchronizować faktury. Czyli przykładowo jeśli w kroku poprzednim stworzyłeś katalog `kosztowe/2023-12` to możesz teraz go odświeżyć / zsynchronizować:

```shell
./ksef download -refresh kosztowe/2023-12/registry.yaml
```

Jak widzisz, ilość parametrów jest znacznie mniejsza ponieważ wszystkie potrzebne dane znajdują się w pliku `registry.yaml`. Jeśli program zauważy, że jakaś faktura już znajduje się w rejestrze to ją pominie

::: info
KSeF udostępnia całkiem sporo danych w nagłówkach faktur więc serializuję je do pliku rejestru. Są tam takie informacje jak data wystawienia, rodzaj faktury, dane firmy itd itd.
:::

::: warning
Ta opcja **nie** pobiera źródłowych faktur w formie XML i/lub PDF - od tego jest komenda `download-pdf`. Przyczyna jest prozaiczna - faktur może być sporo i wolałem rozdzielić te komendy na dwie osobne
:::
