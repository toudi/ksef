---
linkTitle: import
---

# `import`

Komenda `import` służy do importu certyfikatów wygenerowanych przez MCU (Moduł Certyfikatów i Uprawnień)

Linki do środowisk:

| Środowisko | Link                               |
| ---------- | ---------------------------------- |
| Produkcja  | https://mcu.mf.gov.pl/web/login    |
| Demo       | https://web2tr-ksef.mf.gov.pl/web/ |

Poniżej przedstawiam samouczek jak wygenerować certyfikat na środowisku demo

1. Zaloguj się do aplikacji

   ![](/images/auth-sign/sign-entry.png)

1. Wrowadź swój numer NIP

   ![](/images/auth-sign/sign-context.png)

1. Kliknij "Podpisz"

   ![](/images/mcu-import/trusted-profile-sign.png)

1. Wybierz metodę autoryzacji. Ja wybrałem profil zaufany bo akurat w momencie tworzenia tego rozdziału dokumentacji nie działała mi aplikacja eDO App :-)

   ![](/images/mcu-import/choose-auth-method.png)

1. W aplikacji MCU kliknij na menu "Certyfikaty" a następnie "Wnioskuj o certyfikat"

   ![](/images/mcu-import/certificates-menu.png)

1. Wypełnij formularz. Zwróć uwagę na nazwę certyfikatu - wartość tego pola stawnowić będzie podstawę nazwy pliku (bez rozszerzenia) klucza prywatnego oraz certyfikatu

   ![](/images/mcu-import/generate-pkey.png)

   {{< callout type="warning" >}}
   Zapisz hasło w osobnym pliku tekstowym jeśli nie ufasz swojej pamięci krótkotrwałej. Będzie ono potrzebne w kolejnym kroku
   {{< /callout >}}

1. Wybierz przeznaczenie certyfikatu. Zacznijmy od autoryzacji

   ![](/images/mcu-import/cert-usage-authentication.png)

1. Kliknij "Wyślij wniosek o wydanie certyfikatu" a następnie utwórz kolejny wniosek
1. Pamiętaj o wybraniu innej nazwy (ja wybrałem "offline-demo")

   ![](/images/mcu-import/generate-pkey-offline.png)

1. Wybierz kolejne przeznaczenie certyfikatu - tym razem "Podpis linku"

   ![](/images/mcu-import/cert-usage-offline.png)

1. Przejdź do pozycji "Lista certyfikatów" i skopiuj numery seryjne oraz pobierz pliki certyfikatów

   ![](/images/mcu-import/certificate.png)

# Import certyfikatu do programu

## Odszyfrowanie

{{< callout type="warning" >}}
Pliki klucza prywatnego są zabezpieczone hasłem. Niestety, pomimo najszczerszych chęci nie byłem w stanie znaleźć biblioteki w golang która byłaby w stanie je odszyfrować. Jako, że czas mnie gonił postanowiłem iść na ciężki kompromis i po prostu dodać odszyfrowane klucze. Jeśli wiesz w jaki sposób odszyfrować klucze "w locie" w golang - będę wdzięczny za podpowiedź :)
{{< /callout >}}

```
openssl pkey -in authentication.key -passin file:pass.txt -out authentication-decrypted.pem
```

gdzie `pass.txt` to ścieżka do pliku z zapisanym hasłem certyfikatu

## Dodanie odszyfrowanego klucza oraz certyfikatu

{{< gateway-flags >}}

Obsługiwane flagi

| flaga           | znaczenie                                                                        |
| --------------- | -------------------------------------------------------------------------------- |
| `-n`            | Numer NIP podmiotu                                                               |
| `-p`            | ścieżka do odszyfrowanego pliku klucza prywatnego                                |
| `--certificate` | ścieżka do pliku certyfikatu w formacie PEM                                      |
| `--serial`      | numer seryjny certyfikatu (do uzyskania na stronie MCU)                          |
| `--usage`       | przeznaczenie certyfikatu. Dopuszczalne wartości: `Authentication` lub `Offline` |

```
./ksef --demo-gateway certs import -n 1111111111 -p authentication-decrypted.pem --certificate authentication.crt --serial AABBCCDDEEFFGGHH --usage Authentication
```
