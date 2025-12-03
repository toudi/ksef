---
linkTitle: sign
---

# `sign`

Komenda służąca do podpisywania pliku wyzwania za pomocą certyfikatu autoryzacyjnego

Przykładowe wywołanie:

```sh
./ksef auth sign -f plik-wyzwania.xml -o podpisany-plik.xml
```

Możesz użyć tej komendy aby zalogować się do aplikacji użytkownika KSeF. Zademonstruję to na przykładzie środowiska testowego:

1. Wejdź na stronę aplikacji użytkownika: https://web2te-ksef.mf.gov.pl/web/
1. Wybierz opcję "Uwierzytelnij się w Krajowym systemie e-Faktur"

   ![](/images/auth-sign/sign-entry.png)

1. Wybierz opcję "Zaloguj się certyfikatem kwalifikowanym"

   ![](/images/auth-sign/sign-choice.png)

1. Wprowadź swój numer NIP

   ![](/images/auth-sign/sign-context.png)

1. Na pytanie "Czy Twój certyfikat zawiera numer PESEL lub ( ... ) wybierz odpowiedź "Tak"

   ![](/images/auth-sign/sign-pesel-question.png)

1. Kliknij przycisk "Pobierz żądanie autoryzacyjne" i zapisz plik XML wygenerowany przez system KSeF. Jest to plik wyzwania

   ![](/images/auth-sign/sign-get-challenge-file.png)

1. Użyj komendy `auth sign` aby podpisać plik wyzwania. Program sam odczyta numer NIP i wybierze odpowiedni certyfikat

   ```sh
   ./ksef auth sign -f plik-wyzwania.xml
   ```

   {{< gateway-flags >}}

1. Kliknij przycisk "Dodaj plik" i wskaż podpisany plik wyzwania

   ![](/images/auth-sign/sign-upload-signed-challenge.png)
   ![](/images/auth-sign/sign-final.png)

1. To wszystko. Po kliknięciu w przycisk "Dalej" zostaniesz przeniesiony do aplikacji użytkownika KSeF.
