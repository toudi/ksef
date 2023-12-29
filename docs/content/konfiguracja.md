# Konfiguracja programu

## Konfiguracja loggerów

```yaml
logging:
  nazwa-loggera: poziom
```

Więcej o dostępnych loggerach i poziomach przeczytasz w rozdziale [Logowanie](/content/logowanie)

## Konfiguracja silników lokalnego wydruku

Masz do dyspozycji dwie możliwości wyrenderowania PDF. W katalogu z przykładami znajdziesz referencyjną mikroaplikację renderującą fakturę. Więcej na jej temat dowiesz się z rozdziału [Wizualizacja PDF](/content/komendy/lokalny-pdf)

### lokalny puppeteer

W katalogu z programem w podkatalogu `examples/local-pdf-printout` znajdziesz katalog `printing` który zawiera potrzebny skrypt `pdf.js` oraz `package.json` który instaluje pakiet `puppeteer`

::: danger
W przypadku wybrania silnika `pupetter` wymagane jest uruchomienie przeglądarki chrome / chromium w trybie headless z parametrem `--allow-file-access-from-files` aby załadować lokalne pliki z dysku. Teoretycznie i tak masz pełną kontrolę nad plikami które znajdują się w aplikacji, tym niemniej miej na względzie oczywiste konsekwencje związane z bezpieczeństwem tego rozwiązania
:::

```yaml
pdf-renderer:
  engine: puppeteer
  node_bin: /usr/bin/node
  browser_bin: /usr/bin/brave-browser
  template_path: /sciezka/do/pliku/index.html
  rendering_script: /sciezka/do/pliku/pdf.js
```

::: warning
Jeśli korzystasz z systemu MacOS wówczas ścieżki mogą być podobne do poniższych:

```yaml
pdf-renderer:
  node_bin: /opt/homebrew/opt/node@20/bin/node
  browser_bin: "/Applications/Brave Browser.app/Contents/MacOS/Brave Browser"
```

:::

### gotenberg

https://gotenberg.dev/

Projekt gotenberg to skonteneryzowane środowisko do konwersji różnych formatów do PDF, w tym HTML. Pod spodem używa on również projektu pupetteer, z tą jednak różnicą iż uruchamia się go poprzez dockera. Ma to swoje wady i zalety. Zaletą jest fakt że nawet jeśli do szablonu renderującego wkradnie się złośliwy kod to zostanie on ograniczony do dockera. Kolejną zaletą jest z pewnością mniejsza ilość zależności po stronie systemu użytkownika. Ale .. musisz dysponować hostem na którym uruchomisz dockera z obrazem gotenberg.

::: warning
Warunkiem **koniecznym** dostosowania Twojego szablonu do współpracy z gotenberg jest to aby wszystkie pliki aplikacji znajdowały się w tym samym katalogu. To właśnie dlatego w dostarczonej przykładowej aplikacji znajdziesz poniższy fragment (w pliku `vite.config.js`):

```js
export default defineConfig({
  build: {
    // required for gotenberg
    assetsDir: '',
  },
```

Opcja ta powoduje wrzucenie wszystkich plików wynikowych (html / CSS / JS) do wspólnego katalogu.
:::

Przykładowa konfiguracja wygląda tak:

```yaml
pdf-renderer:
  engine: gotenberg
  host: http://127.0.0.1:3000/
  template_path: /sciezka/do/pliku/index.html
```

::: info
Podczas procesu wydruku, program odczyta wszystkie pliki z katalogu w którym znajduje się `index.html` oraz prześle je do gotenberga. Dzięki temu, jeśli masz jakieś obrazki lub inne rzeczy to po prostu upewnij się że znajdują się w tym samym katalogu i ścieżki do nich również nie wykraczają poza ścieżkę. Jest to ograniczenie samego projektu gotenberg i warto o nim pamiętać
:::

::: warning
Zwróć uwagę aby podczas uruchamiania gotenberga przez dockera przekazać mu flagę `--chromium-allow-file-access-from-files`.

np:

```shell
docker run -d --rm -p 3000:3000 gotenberg/gotenberg:7 gotenberg --chromium-allow-file-access-from-files
```

W przeciwnym wypadku, gotenberg nie załaduje zasobów i zwróci pustego PDF'a.
:::
