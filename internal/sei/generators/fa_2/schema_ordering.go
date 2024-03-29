// this file is generated by go generate. Please do not modify it manually!
package fa_2

var FA_2ChildrenOrder = map[string]map[string]int{
	"Faktura.PodmiotUpowazniony.AdresKoresp": {"KodKraju": 0,"AdresL1": 1,"AdresL2": 2,"GLN": 3},
	"TAdres": {"KodKraju": 0,"AdresL1": 1,"AdresL2": 2,"GLN": 3},
	"Faktura.Podmiot1": {"PrefiksPodatnika": 0,"NrEORI": 1,"DaneIdentyfikacyjne": 2,"Adres": 3,"AdresKoresp": 4,"DaneKontaktowe": 5,"StatusInfoPodatnika": 6},
	"Faktura.Fa.Podmiot1K": {"PrefiksPodatnika": 0,"DaneIdentyfikacyjne": 1,"Adres": 2},
	"Faktura.PodmiotUpowazniony.Adres": {"KodKraju": 0,"AdresL1": 1,"AdresL2": 2,"GLN": 3},
	"Faktura.Fa.Rozliczenie": {"Obciazenia": 0,"SumaObciazen": 1,"Odliczenia": 2,"SumaOdliczen": 3,"DoZaplaty": 4,"DoRozliczenia": 5},
	"Faktura.Fa.Platnosc.ZaplataCzesciowa": {"KwotaZaplatyCzesciowej": 0,"DataZaplatyCzesciowej": 1},
	"Faktura.Fa.Platnosc.TerminPlatnosci": {"Termin": 0,"TerminOpis": 1},
	"Faktura.Fa.Zamowienie.ZamowienieWiersz": {"NrWierszaZam": 0,"UU_IDZ": 1,"P_7Z": 2,"IndeksZ": 3,"GTINZ": 4,"PKWiUZ": 5,"CNZ": 6,"PKOBZ": 7,"P_8AZ": 8,"P_8BZ": 9,"P_9AZ": 10,"P_11NettoZ": 11,"P_11VatZ": 12,"P_12Z": 13,"P_12Z_XII": 14,"P_12Z_Zal_15": 15,"GTUZ": 16,"ProceduraZ": 17,"KwotaAkcyzyZ": 18,"StanPrzedZ": 19},
	"TPodmiot2": {"NIP": 0,"KodUE": 1,"NrVatUE": 2,"KodKraju": 3,"NrID": 4,"BrakID": 5,"Nazwa": 6},
	"Faktura.Podmiot1.DaneKontaktowe": {"Email": 0,"Telefon": 1},
	"Faktura.Podmiot3.DaneKontaktowe": {"Email": 0,"Telefon": 1},
	"Faktura.Fa.Adnotacje.PMarzy": {"P_PMarzy": 0,"P_PMarzy_2": 1,"P_PMarzy_3_1": 2,"P_PMarzy_3_2": 3,"P_PMarzy_3_3": 4,"P_PMarzyN": 5},
	"Faktura.Podmiot1.DaneIdentyfikacyjne": {"NIP": 0,"Nazwa": 1},
	"Faktura.PodmiotUpowazniony.DaneKontaktowe": {"EmailPU": 0,"TelefonPU": 1},
	"Faktura.Fa.WarunkiTransakcji.Transport": {"RodzajTransportu": 0,"TransportInny": 1,"OpisInnegoTransportu": 2,"Przewoznik": 3,"NrZleceniaTransportu": 4,"OpisLadunku": 5,"LadunekInny": 6,"OpisInnegoLadunku": 7,"JednostkaOpakowania": 8,"DataGodzRozpTransportu": 9,"DataGodzZakTransportu": 10,"WysylkaZ": 11,"WysylkaPrzez": 12,"WysylkaDo": 13},
	"Faktura.Podmiot2.DaneIdentyfikacyjne": {"NIP": 0,"KodUE": 1,"NrVatUE": 2,"KodKraju": 3,"NrID": 4,"BrakID": 5,"Nazwa": 6},
	"Faktura.Fa.Podmiot2K.Adres": {"KodKraju": 0,"AdresL1": 1,"AdresL2": 2,"GLN": 3},
	"Faktura.Podmiot2.DaneKontaktowe": {"Email": 0,"Telefon": 1},
	"Faktura.Fa.Zamowienie": {"WartoscZamowienia": 0,"ZamowienieWiersz": 1},
	"Faktura.Stopka.Informacje": {"StopkaFaktury": 0},
	"Faktura.Fa.DodatkowyOpis": {"NrWiersza": 0,"Klucz": 1,"Wartosc": 2},
	"Faktura.Fa.Adnotacje.Zwolnienie": {"P_19": 0,"P_19A": 1,"P_19B": 2,"P_19C": 3,"P_19N": 4},
	"Faktura.Fa.Platnosc": {"Zaplacono": 0,"DataZaplaty": 1,"ZnacznikZaplatyCzesciowej": 2,"ZaplataCzesciowa": 3,"TerminPlatnosci": 4,"FormaPlatnosci": 5,"PlatnoscInna": 6,"OpisPlatnosci": 7,"RachunekBankowy": 8,"RachunekBankowyFaktora": 9,"Skonto": 10},
	"Faktura.Fa.Podmiot1K.DaneIdentyfikacyjne": {"NIP": 0,"Nazwa": 1},
	"Faktura.Podmiot3.AdresKoresp": {"KodKraju": 0,"AdresL1": 1,"AdresL2": 2,"GLN": 3},
	"Faktura.Fa.Platnosc.RachunekBankowy": {"NrRB": 0,"SWIFT": 1,"RachunekWlasnyBanku": 2,"NazwaBanku": 3,"OpisRachunku": 4},
	"Faktura.Podmiot1.Adres": {"KodKraju": 0,"AdresL1": 1,"AdresL2": 2,"GLN": 3},
	"Faktura.PodmiotUpowazniony": {"NrEORI": 0,"DaneIdentyfikacyjne": 1,"Adres": 2,"AdresKoresp": 3,"DaneKontaktowe": 4,"RolaPU": 5},
	"Faktura.Fa.DaneFaKorygowanej": {"DataWystFaKorygowanej": 0,"NrFaKorygowanej": 1,"NrKSeF": 2,"NrKSeFFaKorygowanej": 3,"NrKSeFN": 4},
	"Faktura.Fa.Podmiot2K": {"DaneIdentyfikacyjne": 0,"Adres": 1,"IDNabywcy": 2},
	"Faktura.Fa.WarunkiTransakcji.Transport.Przewoznik": {"DaneIdentyfikacyjne": 0,"AdresPrzewoznika": 1},
	"schema": {"Faktura": 0},
	"Faktura.Podmiot3": {"IDNabywcy": 0,"NrEORI": 1,"DaneIdentyfikacyjne": 2,"Adres": 3,"AdresKoresp": 4,"DaneKontaktowe": 5,"Rola": 6,"RolaInna": 7,"OpisRoli": 8,"Udzial": 9,"NrKlienta": 10},
	"Faktura.Fa.Adnotacje.NoweSrodkiTransportu.NowySrodekTransportu": {"P_22A": 0,"P_NrWierszaNST": 1,"P_22BMK": 2,"P_22BMD": 3,"P_22BK": 4,"P_22BNR": 5,"P_22BRP": 6,"P_22B": 7,"P_22B1": 8,"P_22B2": 9,"P_22B3": 10,"P_22B4": 11,"P_22BT": 12,"P_22C": 13,"P_22C1": 14,"P_22D": 15,"P_22D1": 16},
	"Faktura.Fa.Rozliczenie.Odliczenia": {"Kwota": 0,"Powod": 1},
	"Faktura.Fa.WarunkiTransakcji.Transport.WysylkaPrzez": {"KodKraju": 0,"AdresL1": 1,"AdresL2": 2,"GLN": 3},
	"Faktura.Podmiot2.AdresKoresp": {"KodKraju": 0,"AdresL1": 1,"AdresL2": 2,"GLN": 3},
	"TPodmiot1": {"NIP": 0,"Nazwa": 1},
	"Faktura.Fa.ZaliczkaCzesciowa": {"P_6Z": 0,"P_15Z": 1,"KursWalutyZW": 2},
	"Faktura.Fa.FaWiersz": {"NrWierszaFa": 0,"UU_ID": 1,"P_6A": 2,"P_7": 3,"Indeks": 4,"GTIN": 5,"PKWiU": 6,"CN": 7,"PKOB": 8,"P_8A": 9,"P_8B": 10,"P_9A": 11,"P_9B": 12,"P_10": 13,"P_11": 14,"P_11A": 15,"P_11Vat": 16,"P_12": 17,"P_12_XII": 18,"P_12_Zal_15": 19,"KwotaAkcyzy": 20,"GTU": 21,"Procedura": 22,"KursWaluty": 23,"StanPrzed": 24},
	"Faktura.Stopka": {"Informacje": 0,"Rejestry": 1},
	"Faktura.Podmiot3.DaneIdentyfikacyjne": {"NIP": 0,"IDWew": 1,"KodUE": 2,"NrVatUE": 3,"KodKraju": 4,"NrID": 5,"BrakID": 6,"Nazwa": 7},
	"Faktura.Fa.Podmiot1K.Adres": {"KodKraju": 0,"AdresL1": 1,"AdresL2": 2,"GLN": 3},
	"Faktura.Fa.WarunkiTransakcji.Zamowienia": {"DataZamowienia": 0,"NrZamowienia": 1},
	"Faktura.Fa.WarunkiTransakcji.Transport.WysylkaDo": {"KodKraju": 0,"AdresL1": 1,"AdresL2": 2,"GLN": 3},
	"Faktura.Fa.Platnosc.RachunekBankowyFaktora": {"NrRB": 0,"SWIFT": 1,"RachunekWlasnyBanku": 2,"NazwaBanku": 3,"OpisRachunku": 4},
	"Faktura.Fa.WarunkiTransakcji.Transport.Przewoznik.AdresPrzewoznika": {"KodKraju": 0,"AdresL1": 1,"AdresL2": 2,"GLN": 3},
	"Faktura.Podmiot2": {"NrEORI": 0,"DaneIdentyfikacyjne": 1,"Adres": 2,"AdresKoresp": 3,"DaneKontaktowe": 4,"NrKlienta": 5,"IDNabywcy": 6},
	"Faktura.Stopka.Rejestry": {"PelnaNazwa": 0,"KRS": 1,"REGON": 2,"BDO": 3},
	"Faktura.Naglowek": {"KodFormularza": 0,"WariantFormularza": 1,"DataWytworzeniaFa": 2,"SystemInfo": 3},
	"Faktura.Fa.Adnotacje.NoweSrodkiTransportu": {"P_22": 0,"P_42_5": 1,"NowySrodekTransportu": 2,"P_22N": 3},
	"Faktura.Podmiot3.Adres": {"KodKraju": 0,"AdresL1": 1,"AdresL2": 2,"GLN": 3},
	"Faktura.PodmiotUpowazniony.DaneIdentyfikacyjne": {"NIP": 0,"Nazwa": 1},
	"TNaglowek": {"KodFormularza": 0,"WariantFormularza": 1,"DataWytworzeniaFa": 2,"SystemInfo": 3},
	"TPodmiot3": {"NIP": 0,"IDWew": 1,"KodUE": 2,"NrVatUE": 3,"KodKraju": 4,"NrID": 5,"BrakID": 6,"Nazwa": 7},
	"Faktura": {"Naglowek": 0,"Podmiot1": 1,"Podmiot2": 2,"Podmiot3": 3,"PodmiotUpowazniony": 4,"Fa": 5,"Stopka": 6},
	"Faktura.Fa.Adnotacje": {"P_16": 0,"P_17": 1,"P_18": 2,"P_18A": 3,"Zwolnienie": 4,"NoweSrodkiTransportu": 5,"P_23": 6,"PMarzy": 7},
	"Faktura.Fa": {"KodWaluty": 0,"P_1": 1,"P_1M": 2,"P_2": 3,"WZ": 4,"P_6": 5,"OkresFa": 6,"P_13_1": 7,"P_14_1": 8,"P_14_1W": 9,"P_13_2": 10,"P_14_2": 11,"P_14_2W": 12,"P_13_3": 13,"P_14_3": 14,"P_14_3W": 15,"P_13_4": 16,"P_14_4": 17,"P_14_4W": 18,"P_13_5": 19,"P_14_5": 20,"P_13_6_1": 21,"P_13_6_2": 22,"P_13_6_3": 23,"P_13_7": 24,"P_13_8": 25,"P_13_9": 26,"P_13_10": 27,"P_13_11": 28,"P_15": 29,"KursWalutyZ": 30,"Adnotacje": 31,"RodzajFaktury": 32,"PrzyczynaKorekty": 33,"TypKorekty": 34,"DaneFaKorygowanej": 35,"OkresFaKorygowanej": 36,"NrFaKorygowany": 37,"Podmiot1K": 38,"Podmiot2K": 39,"P_15ZK": 40,"KursWalutyZK": 41,"ZaliczkaCzesciowa": 42,"FP": 43,"TP": 44,"DodatkowyOpis": 45,"FakturaZaliczkowa": 46,"ZwrotAkcyzy": 47,"FaWiersz": 48,"Rozliczenie": 49,"Platnosc": 50,"WarunkiTransakcji": 51,"Zamowienie": 52},
	"Faktura.Fa.WarunkiTransakcji.Umowy": {"DataUmowy": 0,"NrUmowy": 1},
	"Faktura.Fa.Platnosc.Skonto": {"WarunkiSkonta": 0,"WysokoscSkonta": 1},
	"Faktura.Fa.WarunkiTransakcji": {"Umowy": 0,"Zamowienia": 1,"NrPartiiTowaru": 2,"WarunkiDostawy": 3,"KursUmowny": 4,"WalutaUmowna": 5,"Transport": 6,"PodmiotPosredniczacy": 7},
	"TKluczWartosc": {"NrWiersza": 0,"Klucz": 1,"Wartosc": 2},
	"TRachunekBankowy": {"NrRB": 0,"SWIFT": 1,"RachunekWlasnyBanku": 2,"NazwaBanku": 3,"OpisRachunku": 4},
	"Faktura.Fa.OkresFa": {"P_6_Od": 0,"P_6_Do": 1},
	"Faktura.Fa.FakturaZaliczkowa": {"NrKSeFZN": 0,"NrFaZaliczkowej": 1,"NrKSeFFaZaliczkowej": 2},
	"Faktura.Fa.WarunkiTransakcji.Transport.WysylkaZ": {"KodKraju": 0,"AdresL1": 1,"AdresL2": 2,"GLN": 3},
	"Faktura.Fa.Rozliczenie.Obciazenia": {"Kwota": 0,"Powod": 1},
	"Faktura.Fa.Podmiot2K.DaneIdentyfikacyjne": {"NIP": 0,"KodUE": 1,"NrVatUE": 2,"KodKraju": 3,"NrID": 4,"BrakID": 5,"Nazwa": 6},
	"Faktura.Fa.WarunkiTransakcji.Transport.Przewoznik.DaneIdentyfikacyjne": {"NIP": 0,"KodUE": 1,"NrVatUE": 2,"KodKraju": 3,"NrID": 4,"BrakID": 5,"Nazwa": 6},
	"Faktura.Podmiot2.Adres": {"KodKraju": 0,"AdresL1": 1,"AdresL2": 2,"GLN": 3},
}