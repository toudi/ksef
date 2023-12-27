# Konfiguracja programu

## Konfiguracja loggerów

```yaml
logging:
  nazwa-loggera: poziom
```

Więcej o dostępnych loggerach i poziomach przeczytasz w rozdziale [Logowanie](/content/logowanie)

## Konfiguracja silników lokalnego wydruku

```yaml
pdf-renderer:
  engine: puppeteer
  node_bin: /usr/bin/node
  browser_bin: /usr/bin/brave-browser
  template_path: /sciezka/do/pliku/index.html
  rendering_script: /sciezka/do/pliku/pdf.js
```

::: warning
Wymagane jest posiadanie przeglądarki opartej o silnik Chromium oraz zainstalowanie pakietu `puppeteer-core`.
:::

::: warning
Jeśli korzystasz z systemu MacOS wówczas ścieżki mogą być podobne do poniższych:

```yaml
pdf-renderer:
  node_bin: /opt/homebrew/opt/node@20/bin/node
  browser_bin: "/Applications/Brave Browser.app/Contents/MacOS/Brave Browser"
```

:::
