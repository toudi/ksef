---
linkTitle: Autoryzacja
title: Autoryzacja
---

Autoryzacja w KSeF w wersji 2.0 została zmieniona i opiera się o certyfikaty oraz podpisy cyfrowe. Sam KSeF umożliwia wygenerowanie takiego certyfikatu, tym niemniej pierwsze logowanie należy przeprowadzić przy użyciu profilu zaufanego - tak aby umożliwić systemowi KSeF pobranie danych potrzebnych do wystawienia certyfikatu.

## Środowisko testowe

{{< callout type="important" >}}
Ważna informacja dotycząca środowiska testowego.

Na środowisku przedprodukcyjnym (demo) oraz produkcji system KSeF sam sprawdza powiązanie numeru PESEL z NIP (w przypadku osób prowadzących JDG). Na środowisku testowym musisz najpierw utworzyć powiązanie. Służy do tego poniższa komenda:

```
./ksef -t auth bind-nip -n 1111111111 -p 22222222222
```

Jest to o tyle istotne, że dopóki nie powiążesz numeru NIP z numerem PESEL, autoryzacja **NIE** będzie działać - nawet pomimo zastosowania poprawnego podpisu kryptograficznego.

Potraktuj to jako "krok 0" czyli czynność którą musisz wykonać zanim przystąpisz do dalszych czynności
{{< /callout >}}

### Certyfikat samopodpisany (self-signed)

Środowisko testowe KSeF umożliwia skorzystanie z certyfikatów samo-podpisanych (self-signed) dzięki czemu możesz pominąć ścieżkę z profilem zaufanym. Jest to z całą pewnością wygodne rozwiązanie, ale z drugiej strony środowisko KSeF wystawia też własne certyfikaty więc użyteczność tej metody jest krótkotrwała (tj. mam tu na myśli to, że ta metoda pozwala nam jedynie ominąć jednorazową autoryzację za pomocą profilu zaufanego)

Aby wygenerować samo-podpisany certyfikat użyj poniższej komendy:

```
./ksef -t certs gen-self-signed -p 22222222222
```

## Generowanie certyfikatów KSeF

Przy pierwszym uruchomieniu programu zastosuj poniższy scenariusz.

1. Wygenerowanie wyzwania autoryzacyjnego

   ```
   ./ksef auth init -n 1111111111
   ```

   W wyniku wywołania komendy utworzony zostanie plik `AuthTokenRequest.xml` który należy podpisać za pomocą profilu zaufanego. W tym celu udaj się na poniższą stronę: [podpis XAdES](https://moj.gov.pl/nforms/signer/upload?xFormsAppName=SIGNER&xadesPdf=true)

   ![](/images/mobywatel-sign-xades.png)

1. Podpisz plik a następnie zapisz go na dysku
1. Użyj modułu autoryzacyjnego aby zalogować się do systemu i wygenerować token sesyjny

   ```
   ./ksef auth login AuthTokenRequest.signed.xml
   ```

1. Użyj modułu certyfikatów aby wygenerować żądanie wydania certyfikatu

   ```
   ./ksef certs prepare-csr -n 1111111111 -a -o
   ```

   {{< callout type="info" >}}
   Znaczenie flag
   |flaga|flaga (skrót)|opis|
   |---|---|---|
   |`--auth`|`-a`|certyfikat służący do autoryzacji|
   |`--offline`|`-o`|certyfikat służący do podpisywania faktur wystawionych w trybie offline|
   {{< /callout >}}

1. Użyj metody synchronizacyjnej aby pobrać certyfikaty z KSeF

   ```
   ./ksef certs sync-csr -n 1111111111
   ```

   {{< callout type="important" >}}
   Wystawianie certyfikatów jest procesem **ASYNCHRONICZNYM**.

   W praktyce oznacza to, że powyższą komendę (`sync-csr`) należy wywołać dwukrotnie. Pierwsze wywołanie spowoduje wysłanie żądania certyfikatu i zapisanie numeru referencyjnego do bazy certyfikatów. Drugie wywołanie spowoduje pobranie wystawionych certyfikatów z KSeF

   Dodatkowa uwaga - pamiętaj o tym, że na środowisku testowym obowiązują dość restrykcyjne limity występowania o certyfikat dla JDG - **5 na 30 dni**. To istotne jeśli planujesz unieważniać certyfikaty w aplikacji użytkownika.

   {{< /callout >}}

## Każde kolejne uruchomienie

Podczas każdego kolejnego uruchomienia program będzie postępował według następującego scenariusza:

1. Spróbuje pobrać tokeny sesyjne z systemowego centrum kluczy (keyring)
1. Jeśli próba się powiedzie - wówczas klucze zostaną użyte i zyskasz kilka sekund
1. Jeśli próba się nie powiedzie (np. tokeny wygasły) - wówczas program użyje certyfikatów aby wygenerować nowe tokeny
