---
weight: 100
title: "Instalacja"
toc: true
---

Aby pobrać najnowszą wersję klienta KSeF masz dwie opcje:

## Użycie wersji binarnych

Przejdź na stronę wydań projektu [ https://github.com/toudi/ksef/releases ] i pobierz wersję odpowiednią dla Twojego systemu operacyjnego. Pamiętaj aby zweryfikować sumę kontrolną pobranego pliku.

## Kompilacja programu ze źródeł

W tym celu musisz najpierw zainstalować kompilator języka golang ze strony https://go.dev/dl/

Następnie sklonuj repozytorium lub pobierz jego archiwum z Github'a: https://github.com/toudi/ksef

w katalogu z repozytorium wydaj następującą komendę:

```shell
go build -o ksef cmd/ksef/main.go
```
