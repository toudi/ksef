---
title: Instalacja programu
---

Istnieje kilka metod instalacji

## Użycie prekompilowanych wersji binarnych

Aby pobrać binarki skompilowane przez CI/CD githuba, udaj się na stronę releaseów githuba: [wydania](https://github.com/toudi/ksef/releases)

Istnieje kilka binarek, w zależności od systemu operacyjnego i architektury procesora. I tak:

| System operacyjny | Architektura              | Binarka                  |
| ----------------- | ------------------------- | ------------------------ |
| Linux             | x86-64                    | ksef-linux-x86_64.tar.gz |
| Linux             | arm64                     | ksef-linux-arm64.tar.gz  |
| Mac OS            | x86-64                    | ksef-mac-x86_64.tar.gz   |
| Mac OS            | apple sillicon - m1/m2/.. | ksef-mac-arm64.tar.gz    |
| Windows           | x86-64                    | ksef-windows-x86_64.zip  |

## Własnoręczna kompilacja

Aby skompilować program własnoręcznie potrzebujesz kompilatora języka go. W tym celu pobierz go ze strony https://go.dev/dl/

Następnie wykonaj kompilację

```shell
go mod tidy
go build -o ksef cmd/ksef/main.go
```

Możesz w tym celu użyć przygotowanego pliku `Makefile`:

```shell
make
```
