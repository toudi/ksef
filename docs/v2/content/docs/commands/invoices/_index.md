---
linkTitle: invoices
---

# `invoices`

komendy do obsługiwania stanowej bazy faktur

Program trzyma wszystkie niezbędne dane potrzebne do wysyłki i synchronizacji faktur w katalogu `data`. Jego struktura jest następująca:

```
data
`-- ksef-test.mf.gov.pl
    |-- 1111111111
    |   `-- 2025
    |       |-- 12
    |       |   |-- registry.yaml
    |       |   |-- upload-sessions.yaml
    |       |   |-- upo
    |       |   |   |-- 20251218-EU-AABBCCDDEE-0011223344-55.pdf
    |       |   |   `-- 20251218-EU-AABBCCDDEE-0011223344-55.xml
    |       |   `-- wystawione
    |       |       `-- 0001-fv-12-34-01.xml
    |       `-- invoices.yaml
    |-- 2222222222
    |   `-- 2025
    |       |-- 12
    |       |   |-- otrzymane
    |       |   |   |-- 0001-moja-firma-abc-fv-12-34-01.pdf
    |       |   |   `-- 0001-moja-firma-abc-fv-12-34-01.xml
    |       |   |-- registry.yaml
    |       |   |-- upload-sessions.yaml
    |       |   `-- wystawione
    |       |       |-- 0001-fv-12-34-45.xml
    |       |       |-- 0002-fv-12-34-46.xml
    |       |       |-- 0003-fv-12-34-47.xml
    |       |       |-- 0004-fv-12-34-49.xml
    |       |       |-- 0005-fv-12-34-50.xml
    |       |       |-- 0006-fv-12-34-51.xml
    |       |       |-- 0007-fv-12-34-52.xml
    |       |       |-- 0008-fv-12-34-53.xml
    |       |       `-- 0009-fv-12-34-54.xml
    |       `-- invoices.yaml
```

Innymi słowy, powyższą strukturę można opisać następująco:

```
data/<bramka>/<nip>/<rok>
```

w środku katalogu z numerem roku znajdziemy plik `invoices.yaml` który stanowi rejestr faktur z całego roku oraz katalogi z miesiącami. w każdym katalogu miesiąca znajdziemy:

- plik `registry.yaml` który stanowi rejestr faktur z danego miesiąca
- plik `upload-sessions.yaml` który jest rejestrem wysyłek
- katalogi `wystawione`, `otrzymane`, `platnika` oraz `strony-upowaznionej` które odpowiadają odpowiednim typom podmiotu
