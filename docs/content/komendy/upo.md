# Pobieranie UPO

Z oczywistych względów, aby pobrać UPO musisz najpierw wysłać faktury do KSeF. Program stworzy wtedy plik statusu wraz z numerami przesłanych faktur.

Aby pobrać UPO z ministerstwa, skorzystaj z komendy `status`

```shell
./ksef status
Usage of status:
  -o string
    	ścieżka do zapisu UPO (domyślnie katalog pliku statusu + {nrRef}.pdf)
  -p string
    	ścieżka do pliku statusu
  -xml
    	zapis UPO jako plik XML
```

na przykład:

```shell
./ksef status -p sciezka-do-status.yaml
```
