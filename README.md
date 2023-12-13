# Dokumentacja programu

Dokumentację znajdziesz w katalogu `docs`. Masz kilka możliwości jej odczytania:

## Lokalne przebudowanie

```shell
cd docs
npm ci
npm run docs:dev
```

Dokumentacja zostanie udostępniona na porcie 5173 ( http://localhost:5173 )

## Plik PDF

Przy każdym buildzie programu generowany jest PDF z dokumentacją. Można go znaleźć w środku archiwum oraz w artefaktach buildu

Możesz też własnoręcznie wygenerować plik PDF używając następujących komend:

```shell
cd docs
npm ci
npm run export-pdf
```

(w katalogu `docs` pojawi się plik `ksef-dokumentacja-uzytkownika.pdf`)
