# Diagnostyka / Logowanie

Aby ułatwić diagnozowanie problemów z programem, przewidziałem logowanie na kilku poziomach. Oznacza to, że możesz włączyć loggery selektywnie, w zależności od tego gdzie program się wywala / działa w sposób niepożądany.

## Jak włączyć logowanie

Stwórz plik konfiguracyjny z poniższą składnią

```yaml
logging:
  nazwa-loggera: poziom
```

dostępne poziomy: `info`, `debug`, `error`

A następnie wywołaj program z przełącznikiem `-c`

::: info

```shell
./ksef -c config.yaml -log - ...
```

:::

::: warning
Przełącznik -log jest bardzo ważny ponieważ odpowiada za przekierowanie wyjścia logów. Domyślnie jest ono wyłączone, nawet jeśli ustawisz poziomy logowania. Przykłady wywołań:

przekierowanie logów do pliku

```shell
./ksef -c config.yaml -log plik-wyjscia.txt ...
```

przekierowanie logów na stdout

```shell
./ksef -c config.yaml -log - ...
```

:::

## Dostępne loggery

| logger           | znaczenie                                   |
| ---------------- | ------------------------------------------- |
| main             | główny logger programu                      |
| interactive      | logger sesji interaktywnej                  |
| interactive.http | logger sesji interaktywnej (zapytania HTTP) |
| batch            | logger sesji wsadowej                       |
| batch.http       | logger sesji wsadowej (zapytania HTTP)      |
| download         | logger pobierania faktur                    |
| download.http    | logger pobierania faktur (zapytania HTTP)   |
| upload           | logger wysyłki faktur                       |
| upload.http      | logger wysyłki faktur (zapytania HTTP)      |
| upo              | logger pobierania UPO                       |
| upo.http         | logger pobierania UPO (zapytania HTTP)      |

## Tryb verbose

Program obsługuje również tryb `verbose` (flaga `-v`) natomiast przestrzegam przed jej stosowaniem - generuje bardzo obszerne wyjście. Flaga `verbose` oznacza przełączenie wszystkich dostępnych loggerów na poziom `DEBUG` (tak, żebyś nie musiał robić tego ręcznie)

::: info
Przykład wywołania:

```shell
./ksef -v -log - ...
```
