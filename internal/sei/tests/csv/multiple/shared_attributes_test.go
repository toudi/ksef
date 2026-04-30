package multiple_test

import (
	"bytes"
	"ksef/internal/sei"
	"ksef/internal/sei/tests/recorder"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

const inputCSV = `
"sekcja","shared-invoice-attributes","Faktura.Fa.Platnosc"
"FormaPlatnosci"
"6"

"sekcja","shared-invoice-attributes","Faktura.Fa.Platnosc.RachunekBankowy"
"NrRB","SWIFT"
"11223344","SWIFT1"
"22334455","SWIFT2"

"Sekcja","Faktura.Podmiot1.DaneIdentyfikacyjne"
"NIP","Nazwa",
1112223344,"Wytwórnia programów bezwartościowych"

"Sekcja","Faktura.Podmiot1.Adres"
"KodKraju","AdresL1","AdresL2"
"PL","ul. Kwietniowa 3/14","01-234 Mielno"

"-- ignore --","pierwsza faktura rozpoczyna się od sekcji faktura.fa",

"Sekcja","Faktura.Fa",
"KodWaluty","P_1","P_1M","P_2","P_6","RodzajFaktury",
"PLN",2026-04-01,"Gżegżółkowo","FV 2026/1",2026-04-15,"VAT",

"Sekcja","Faktura.Fa.Adnotacje",
"P_16","P_17","P_18","P_18A","Zwolnienie.P_19","Zwolnienie.P_19A","Zwolnienie.P_19N","NoweSrodkiTransportu.P_22N","P_23","PMarzy.P_PMarzyN",
"2","2","2","2","1","Art.43 ust.1 pkt.37 Ustawy o podatku od towarów i usług","","1","2","1",

"-- ignore --","pierwsza faktura jest nieopłacona, ale wartość '6' z sekcji 'shared-invoice-attributes' powinna się skopiować"
"Sekcja","Faktura.Fa.Platnosc"
"Zaplacono","DataZaplaty","TerminPlatnosci.Termin"
"","","2026-04-15"

"Sekcja","Faktura.Podmiot2.DaneIdentyfikacyjne",
"BrakID","Nazwa",
1,"Nabywca 1",

"Sekcja","Faktura.Podmiot2.Adres",
"KodKraju","AdresL1","AdresL2",
"PL","ul. Starościńska 1","82-200 Malbork",

"Sekcja","Faktura.Fa.FaWiersze.FaWiersz",
"P_7","P_8A","P_8B","P_9A","P_12",
"Fikcyjna usługa","szt",195,20,"zw",

"-- ignore --","kolejna faktura rozpoczyna się od sekcji Faktura.Fa"

"Sekcja","Faktura.Fa",
"KodWaluty","P_1","P_1M","P_2","P_6","RodzajFaktury",
"PLN",2026-04-01,"Gżegżółkowo","FV 2026/2",2026-04-15,"VAT",

"-- ignore --","druga faktura jest opłacona"
"Sekcja","Faktura.Fa.Platnosc"
"Zaplacono","DataZaplaty"
"1","2026-04-17"

"Sekcja","Faktura.Fa.Adnotacje",
"P_16","P_17","P_18","P_18A","Zwolnienie.P_19","Zwolnienie.P_19A","Zwolnienie.P_19N","NoweSrodkiTransportu.P_22N","P_23","PMarzy.P_PMarzyN",
"2","2","2","2","1","Art.43 ust.1 pkt.37 Ustawy o podatku od towarów i usług","","1","2","1",

"Sekcja","Faktura.Podmiot2.DaneIdentyfikacyjne",
"BrakID","Nazwa",
1,"Zygmunt III Waza",

"Sekcja","Faktura.Podmiot2.Adres",
"KodKraju","AdresL1","AdresL2",
"PL","plac Zamkowy","00-001 Warszawa",

"Sekcja","Faktura.Fa.FaWiersze.FaWiersz",
"P_7","P_8A","P_8B","P_9A","P_12",
"Konsultacja w sprawie karuzeli VAT","szt",142.00,12,"zw"
`

const expectedInvoice1 = `
<?xml version="1.0" encoding="utf-8"?>
<Faktura xmlns="http://crd.gov.pl/wzor/2025/06/25/13775/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <Naglowek>
    <KodFormularza kodSystemowy="FA (3)" wersjaSchemy="1-0E">FA</KodFormularza>
    <WariantFormularza>3</WariantFormularza>
    <DataWytworzeniaFa>2026-04-14T16:00:00Z</DataWytworzeniaFa>
    <SystemInfo>WSI Pegasus</SystemInfo>
  </Naglowek>
  <Podmiot1>
    <DaneIdentyfikacyjne>
      <NIP>1112223344</NIP>
      <Nazwa>Wytwórnia programów bezwartościowych</Nazwa>
    </DaneIdentyfikacyjne>
    <Adres>
      <KodKraju>PL</KodKraju>
      <AdresL1>ul. Kwietniowa 3/14</AdresL1>
      <AdresL2>01-234 Mielno</AdresL2>
    </Adres>
  </Podmiot1>
  <Podmiot2>
    <DaneIdentyfikacyjne>
      <BrakID>1</BrakID>
      <Nazwa>Nabywca 1</Nazwa>
    </DaneIdentyfikacyjne>
    <Adres>
      <KodKraju>PL</KodKraju>
      <AdresL1>ul. Starościńska 1</AdresL1>
      <AdresL2>82-200 Malbork</AdresL2>
    </Adres>
    <JST>2</JST>
    <GV>2</GV>
  </Podmiot2>
  <Fa>
    <KodWaluty>PLN</KodWaluty>
    <P_1>2026-04-01</P_1>
    <P_1M>Gżegżółkowo</P_1M>
    <P_2>FV 2026/1</P_2>
    <P_6>2026-04-15</P_6>
    <P_15>3900</P_15>
    <Adnotacje>
      <P_16>2</P_16>
      <P_17>2</P_17>
      <P_18>2</P_18>
      <P_18A>2</P_18A>
      <Zwolnienie>
        <P_19>1</P_19>
        <P_19A>Art.43 ust.1 pkt.37 Ustawy o podatku od towarów i usług</P_19A>
      </Zwolnienie>
      <NoweSrodkiTransportu>
        <P_22N>1</P_22N>
      </NoweSrodkiTransportu>
      <P_23>2</P_23>
      <PMarzy>
        <P_PMarzyN>1</P_PMarzyN>
      </PMarzy>
    </Adnotacje>
    <RodzajFaktury>VAT</RodzajFaktury>
    <FaWiersz>
      <NrWierszaFa>1</NrWierszaFa>
      <P_7>Fikcyjna usługa</P_7>
      <P_8A>szt</P_8A>
      <P_8B>195</P_8B>
      <P_9A>20</P_9A>
      <P_11>3900</P_11>
      <P_12>zw</P_12>
    </FaWiersz>
    <Platnosc>
      <TerminPlatnosci>
        <Termin>2026-04-15</Termin>
      </TerminPlatnosci>
      <FormaPlatnosci>6</FormaPlatnosci>
      <RachunekBankowy>
        <NrRB>11223344</NrRB>
        <SWIFT>SWIFT1</SWIFT>
      </RachunekBankowy>
      <RachunekBankowy>
        <NrRB>22334455</NrRB>
        <SWIFT>SWIFT2</SWIFT>
      </RachunekBankowy>
    </Platnosc>
  </Fa>
</Faktura>
`

const expectedInvoice2 = `
<?xml version="1.0" encoding="utf-8"?>
<Faktura xmlns="http://crd.gov.pl/wzor/2025/06/25/13775/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <Naglowek>
    <KodFormularza kodSystemowy="FA (3)" wersjaSchemy="1-0E">FA</KodFormularza>
    <WariantFormularza>3</WariantFormularza>
    <DataWytworzeniaFa>2026-04-14T16:00:01Z</DataWytworzeniaFa>
    <SystemInfo>WSI Pegasus</SystemInfo>
  </Naglowek>
  <Podmiot1>
    <DaneIdentyfikacyjne>
      <NIP>1112223344</NIP>
      <Nazwa>Wytwórnia programów bezwartościowych</Nazwa>
    </DaneIdentyfikacyjne>
    <Adres>
      <KodKraju>PL</KodKraju>
      <AdresL1>ul. Kwietniowa 3/14</AdresL1>
      <AdresL2>01-234 Mielno</AdresL2>
    </Adres>
  </Podmiot1>
  <Podmiot2>
    <DaneIdentyfikacyjne>
      <BrakID>1</BrakID>
      <Nazwa>Zygmunt III Waza</Nazwa>
    </DaneIdentyfikacyjne>
    <Adres>
      <KodKraju>PL</KodKraju>
      <AdresL1>plac Zamkowy</AdresL1>
      <AdresL2>00-001 Warszawa</AdresL2>
    </Adres>
    <JST>2</JST>
    <GV>2</GV>
  </Podmiot2>
  <Fa>
    <KodWaluty>PLN</KodWaluty>
    <P_1>2026-04-01</P_1>
    <P_1M>Gżegżółkowo</P_1M>
    <P_2>FV 2026/2</P_2>
    <P_6>2026-04-15</P_6>
    <P_15>1704</P_15>
    <Adnotacje>
      <P_16>2</P_16>
      <P_17>2</P_17>
      <P_18>2</P_18>
      <P_18A>2</P_18A>
      <Zwolnienie>
        <P_19>1</P_19>
        <P_19A>Art.43 ust.1 pkt.37 Ustawy o podatku od towarów i usług</P_19A>
      </Zwolnienie>
      <NoweSrodkiTransportu>
        <P_22N>1</P_22N>
      </NoweSrodkiTransportu>
      <P_23>2</P_23>
      <PMarzy>
        <P_PMarzyN>1</P_PMarzyN>
      </PMarzy>
    </Adnotacje>
    <RodzajFaktury>VAT</RodzajFaktury>
    <FaWiersz>
      <NrWierszaFa>1</NrWierszaFa>
      <P_7>Konsultacja w sprawie karuzeli VAT</P_7>
      <P_8A>szt</P_8A>
      <P_8B>142</P_8B>
      <P_9A>12</P_9A>
      <P_11>1704</P_11>
      <P_12>zw</P_12>
    </FaWiersz>
    <Platnosc>
      <Zaplacono>1</Zaplacono>
      <DataZaplaty>2026-04-17</DataZaplaty>
      <FormaPlatnosci>6</FormaPlatnosci>
      <RachunekBankowy>
        <NrRB>11223344</NrRB>
        <SWIFT>SWIFT1</SWIFT>
      </RachunekBankowy>
      <RachunekBankowy>
        <NrRB>22334455</NrRB>
        <SWIFT>SWIFT2</SWIFT>
      </RachunekBankowy>
    </Platnosc>
  </Fa>
</Faktura>
`

func TestSupportForSharedInvoiceAttributes(t *testing.T) {
	vip := viper.New()
	vip.Set("generator", "fa-3_1.0")
	vip.Set("csv.delimiter", ",")
	recorder := recorder.NewRecorder()
	sei, err := sei.SEI_Init(vip, sei.WithInvoiceReadyFunc(recorder.Ready))
	require.NoError(t, err)

	reader := bytes.NewReader([]byte(inputCSV))

	err = sei.ProcessReader(reader, "dummy.csv")
	require.NoError(t, err)
	require.Len(t, recorder.Invoices, 2)
	require.Equal(t, strings.TrimLeft(expectedInvoice1, " \n"), recorder.XMLInvoices[0])
	require.Equal(t, strings.TrimLeft(expectedInvoice2, " \n"), recorder.XMLInvoices[1])
}
