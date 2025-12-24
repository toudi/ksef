package certificates

import (
	"fmt"
	"ksef/internal/certsdb"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "wyświetla listę dostępnych certyfikatów",
	RunE:  listCerts,
}

func init() {
	CertificatesCommand.AddCommand(listCommand)
}

func listCerts(cmd *cobra.Command, _ []string) error {
	vip := viper.GetViper()

	certDB, err := certsdb.OpenOrCreate(vip)
	if err != nil {
		return err
	}

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "id"},
			{Align: simpletable.AlignCenter, Text: "środowisko"},
			{Align: simpletable.AlignCenter, Text: "funkcja"},
			{Align: simpletable.AlignCenter, Text: "nip"},
			{Align: simpletable.AlignCenter, Text: "profil"},
			{Align: simpletable.AlignCenter, Text: "samopodpisany"},
		},
	}

	for _, cert := range certDB.Certs() {
		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{
				Text: cert.UID,
			},
			{
				Text: string(cert.Environment),
			},
			{
				Text: cert.UsageAsString(),
			},
			{
				Text: cert.NIP,
			},
			{
				Text: cert.ProfileName,
			},
			{
				Text: fmt.Sprintf("%t", cert.SelfSigned),
			},
		})
	}

	table.SetStyle(simpletable.StyleDefault)
	table.Println()
	return nil
}
