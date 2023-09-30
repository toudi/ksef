kompilacja programu:

go build ksef/cmd/main.go

w celu rekompilacji schematów:

go run parse_schemas.go

wówczas program sparsuje schematy z katalogu "schemas" i wygeneruje odpowiednie struktury w katalogu "generators"

Jeśli zastanawiasz się po kiego grzyba jest ten generowany kod spieszę odpowiedzieć, że niestety ministerstwo używa typu sequence a on wymusza aby elementy w drzewie występowały w określonej kolejności (sic!) miałem więc do wyboru albo zaimplementować struktury w ten sposób, żeby ręcznie wklepać je do kodu w golang albo zaimplementować je w sposób ogólniejszy aby to użytkownik programu wypełniał te pola / atrybuty które wie, że potrzebuje a program na podstawie sparsowania schemy posortuje atrybuty według kolejności i XML przejdzie walidację.

## generowanie faktur:

./ksef generate -d ';' -f faktury.csv -o katalog-docelowy [-t]

(parametr -t uzywa klucza publicznego bramki testowej do generowania metadanych)

## podpisywanie faktur.

w katalogu docelowym program stworzy pliki:

- metadata.xml
- metadata.zip
- metadata.zip.aes

( oraz pliki z fakturami )

plik metadata.xml nalezy podpisać - uzywajac karty kryptograficznej lub profiu zaufanego. podpisany cyfrowo plik nalezy zapisac na komputerze i
uzyc kolejnej komendy:

## wysylka faktur

./ksef upload [-t] -f metadata.xml.signed

## pobieranie upo

./ksef status metadata.xml.signed.ref

jesli status przetworzenia faktur bedzie poprawny, program pobierze upo i wygeneruje plik:

- metadata.html (jesli nie znajdzie w sciezce programu wkhtmltopdf)
- metadata.pdf (jesli znajdzie wkhtmltopdf)

UPO jest gotowe do druku.
