# Pobieranie UPO

Z oczywistych względów, aby pobrać UPO musisz najpierw wysłać faktury do KSeF. Program stworzy wtedy plik statusu wraz z numerami przesłanych faktur.

Aby pobrać UPO z ministerstwa, skorzystaj z komendy `status`

```shell
./ksef status
Usage of status:
  -m	stwórz katalog, jeśli wskazany do zapisu nie istnieje
  -o string
    	ścieżka do zapisu UPO (domyślnie katalog pliku rejestru + {nrRef}.pdf)
  -p string
    	ścieżka do pliku rejestru
  -xml
    	zapis UPO jako plik XML
```

na przykład:

```shell
./ksef status -p sciezka-do-registry.yaml
```

Jeśli chcesz zapisać UPO w innym katalogu niż katalog rejestru:

```shell
./ksef status -p sciezka-do-registry.yaml -o inny-katalog
```

::: warning
Katalog `inny-katalog` musi istnieć na dysku. W przeciwnym wypadku program zgłosi błąd: ```wskazany katalog nie istnieje a nie użyłeś opcji `-m` ```. Jeśli chcesz, aby program stworzył katalog wyjściowy, wywołaj program w ten sposób:

```shell
./ksef status -p sciezka-do-registry.yaml -o nieistniejacy/katalog -m
```
:::

Jeśli chcesz ręcznie wskazać nazwę UPO:

```shell
./ksef status -p sciezka-do-registry.yaml -o nieistniejacy/katalog/upo.pdf -m
```