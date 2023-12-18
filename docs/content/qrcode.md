# Kod QR na wizualizacji faktury

Stosunkowo niedawno (tj. około wersji 1.5 API) wizualizacje PDF oferowane przez ministerstwo zaczęły zawierać kod QR służący do weryfikacji, czy faktura istnieje w zasobach KSeF. Po krótkiej analizie treści obrazka stwierdziłem, że zawiera ona sumę `sha256` dokumentu źródłowego faktury (co jest zrozumiałe) oraz jej numer referencyjny w KSeF (co jest niezbyt szczęśliwe). Link zapisany w kodzie jest postaci następującej:

```text
https://{środowisko}/web/verify/{numerFakturyWKSeF}/{sumaKontrolna}
```

gdzie `sumaKontrolna` to skrót `sha256` (w postaci bajtów) zaenkodowany przez `base64`

Oznacza to, że jedyną możliwością wygenerowania kodu QR (przynajmniej na chwilę obecną) jest wysłanie faktury do KSeF ponieważ inaczej nie otrzymamy jej numeru referencyjnego. To przeczy wcześniejszym założeniom ministerstwa finansów jakoby kody QR można było generować off-line tj. w przypadku niedostępności KSeF.

Tak czy siak, program wygeneruje dla Ciebie link który możesz przepuścić przez dowolną bibiotekę do obsługi kodów QR i wygenerować obrazek i zapisze go w pliku rejestru:

```yaml
invoices:
  - referenceNumber: FV 00/11/22 - TEST QR
    ksefReferenceNumber: 1111111111-22222222-XXXXXXXXXXXX-ZZ
    qrcode-url: https://ksef-test.mf.gov.pl/web/verify/1111111111-22222222-XXXXXXXXXXXX-ZZ/VPuBCK3cwQvsnprWQSwWSglJvokGtUH%2FQCsPyUPiXK0%3D
```

Oprócz tego, program wygeneruje przykładowy kod QR i zapisze go jako plik `{numerFakturyKSeF}.svg`. Możesz użyć go jeśli nie posiadasz żadnej biblioteki do generowania kodów QR.

Zauważysz zapewne, że pewne znaki są kodowane w postaci `urlsafe` i wydaje się to być działanie zamierzone ze strony ministerstwa finansów.
