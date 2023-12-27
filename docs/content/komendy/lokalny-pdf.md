# Wizualiacja PDF na podstawie własnego szablonu

::: warning
Ta funkcjonalność jest cały czas rozwijana - jestem otwarty na wszelkie sugestie
:::

Niestety, w związku z tym, że eFaktura jest plikiem XML mamy pod górkę już na starcie. Teoretycznie, aby skonwertować plik XML na HTML wystarczy użyć transformacji XSLT. Tyle tylko, że plik faktury zawiera przestrzenie nazw co powoduje, że zapytania XPath stają się niezbyt przyjemne i czytelne.

Zamiast tego wpadłem na pomysł aby do wizualizacji służyła mikroaplikacja javascriptowa która na wejściu przyjmie surowy plik XML, następnie przetworzy go na obiekt javascriptowy, wyświetli w HTML a następnie wydrukuje za pomocą puppetteer

::: warning
Zdaję sobie sprawę, że jest to wybór **niemerytoryczny** i że nie każdemu może przypaść do gustu. Dlatego zakładam, że silników drukujących mogłoby być kilka i konfigurowane byłyby w pliku konfiguracyjnym a ten który zaimplementowałem możnaby traktować jedynie jako wersję "referencyjną"
:::

::: danger
Przykładowa aplikacja uruchamia przeglądarkę chrome / chromium z parametrem `--disable-web-security` aby załadować lokalne pliki z dysku. Opcją numer dwa byłoby serwowanie plików przez program kliencki na jakimś losowym porcie - póki co po prostu nie chce mi się tego implementować. Inną opcją jest użycie programu gotenberg ale wymaga to działającego dockera
:::

## Rozumiem zagrożenie, co mam zrobić?

Zerknij na katalog `examples/local-pdf-printout`. Zawiera on dwa katalogi:

| katalog           | znaczenie                                                       |
| ----------------- | --------------------------------------------------------------- |
| invoice-rendering | Zawiera **szablon aplikacji** wyświetlającej fakturę oraz QRKod |
| printing          | Zawiera projekt oraz przykład skryptu drukującego               |

### Instalacja bibliotek

```shell
cd invoice-rendering
npm install
cd ../printing
npm install
```

### Zmiana szablonu

Referencyjna aplikacja ma zaimplementowany szablon który na 99% nie będzie spełniać Twoich oczekiwań ale da Ci pogląd na to w jaki sposób można odnosić się do danych zawartych w fakturze, jak renderować QRKod etc etc.

Mikroaplikacja jest napisana w Vue.js. Ponownie - jest to wybór **niemerytoryczny** i przyczyna jest bardzo prosta - kod wyjściowy generowany przez Vue.js bardziej przypada mi do gustu niż inne silniki.

W mojej opinii nie ma to jednak **najmniejszego** znaczenia, ponieważ - ponownie - mój sposób jest jedynie przykładowy i równie dobrze możesz napisać aplikację w react / angularze albo svelte. Grunt, żeby działała w identyczny sposób.

Oto kilka najistotniejszych plików wraz z krótkim omówieniem

| plik                    | znaczenie                                                                                                                                                                                          |
| ----------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| index.html              | Podstawowy szablon. Zwróć uwagę na specyficzne wartości wpisane do sekcji `head > meta`. Na ich podstawie aplikacja będzie w stanie zwizualizować fakturę                                          |
| src/annotate-invoice.js | Tworzenie anotacji (dodatkowych pomocniczych informacji do wyświetlenia). Zawiera funkcję która podmienia pola z XML na mnemoniki a także oblicza częściowe i całkowite kwoty netto / brutto / VAT |
| src/components/\*.vue   | pliki komponentów do renderowania faktury. Jeśli kiedykolwiek robiłeś coś w vue to z pewnością dasz radę je ogarnąć                                                                                |

### Sekcja `head > meta`

Jest to dość istotna kwestia związana z szablonem. Wygląda ona następująco:

```html
<meta name="invoice:qrcode" content="__qrcode_url__" />
<meta name="invoice:seiRefNo" content="__invoice_sei_ref_no__" />
<meta name="invoice" content="__invoice_base64__" />
```

Do sekcji `meta` mój program przekazuje następujące dane

| zmienna                  | znaczenie                                                   |
| ------------------------ | ----------------------------------------------------------- |
| `__qrcode_url__`         | Link do weryfikacji faktury który znajduje się na kodzie QR |
| `__invoice_sei_ref_no__` | Numer faktury który widnieje w KSeF                         |
| `__invoice_base64__`     | Treść pliku XML z fakturą zakodowana w base64               |

::: info
Jeśli testujesz szablon, użyj następującej komendy:

```shell
npm run dev
```

Wówczas uruchomi się serwer na porcie 5173. Oczywiście jeśli zostawisz domyślne wartości meta to szablon się wysypie. Pamiętaj więc, żeby do testów wstawić tam własne wartości. O ile link QR i numer faktury z KSeF nie mają większego znaczenia o tyle jeśli nie wpiszesz faktycznego base64 z XML'a z fakturą do `meta[name=invoice]` to .. niczego nie potestujesz :)
:::

### Budowanie wynikowego szablonu

Jeśli dostosowałeś szablon do swoich potrzeb to zbuduj jego finalną wersję:

```shell
npm run build -- --outDir=../template
```

::: warning
Podczas wizualizacji program tworzy plik `render.html` który wysyła do procesu drukującego. Teoretycznie więc jeśli będziesz mieć sporo wydruków, warto rozważyć trzymanie szablonu w ramdysku żeby nie niszczyć sobie dysku :-)
:::

## Jak powiązać to z klientem?

```text
./ksef render-pdf
Usage of render-pdf:
  -i string
    	plik XML do wizualizacji
  -m	stwórz katalog, jeśli wskazany do zapisu nie istnieje
  -o string
    	ścieżka do zapisu PDF (domyślnie katalog pliku statusu + {nrRef}.pdf)
  -p string
    	ścieżka do pliku rejestru
```

Klient działa w następujący sposób:

1. Oblicza sumę sha256 pliku XML faktury
2. Znajduje metadane KSeF faktury w pliku rejestru
3. Koduje plik do base64
4. Bierze wskazany przez Ciebie szablon HTML a następnie wypełnia pola w `meta` odpowiednimi wartościami
5. Tworzy plik render.html (oryginalny index.html pozostaje bez zmian)
6. Renderuje PDF

Więcej o konfiguracji silników do renderowania PDF czytaj w rozdziale [Konfiguracja](/content/konfiguracja)

### Przykłady wywołań:

::: warning
Warunkiem **koniecznym** jest wskazanie lokalizacji pliku konfiguracyjnego poprzez przełącznik `-c`
:::

#### Zapisanie PDF'a w katalogu `katalog` z domyślną nazwą `{numerFakturyZKSeF}.pdf`

```shell
./ksef -c config.yaml render-pdf -i katalog/invoice.xml -o katalog -p katalog/registry.yaml
```

#### Zapisanie PDF'a w katalogu z własną nazwą

```shell
./ksef -c config.yaml render-pdf -i katalog/invoice.xml -o katalog/mojafaktura.pdf -p katalog/registry.yaml
```

#### Zapisanie PDF'a w innym katalogu z domyślną nazwą

```shell
./ksef -c config.yaml render-pdf -i katalog/invoice.xml -o inny-katalog -m -p katalog/registry.yaml
```

::: warning
zwróć uwagę na flagę `-m` - w przeciwnym wypadku jeśli `inny-katalog` nie istnieje, program zwróci błąd
:::

#### Zapisanie PDF'a w innym katalogu z własną nazwą

```shell
./ksef -c config.yaml render-pdf -i katalog/invoice.xml -o inny-katalog/mojafaktura.pdf -m -p katalog/registry.yaml
```

::: warning
zwróć uwagę na flagę `-m` - w przeciwnym wypadku jeśli `inny-katalog` nie istnieje, program zwróci błąd
:::
