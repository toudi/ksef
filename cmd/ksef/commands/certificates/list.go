package certificates

import (
	"fmt"
	"ksef/internal/certsdb"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/samber/lo"
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
	defer certDB.Save()

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "id"},
			{Align: simpletable.AlignCenter, Text: "środowisko"},
			{Align: simpletable.AlignCenter, Text: "funkcja"},
			{Align: simpletable.AlignCenter, Text: "nip"},
			{Align: simpletable.AlignCenter, Text: "profil"},
			{Align: simpletable.AlignCenter, Text: "samopodpisany"},
			{Align: simpletable.AlignCenter, Text: "podmiot"},
		},
	}

	for _, cert := range certDB.Certs() {
		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{
				Text: cert.UID,
			},
			{
				Text: string(cert.EnvironmentId),
			},
			{
				Text: cert.UsageAsString(),
			},
			{
				Text: strings.Join(cert.NIP, ", "),
			},
			{
				Text: cert.ProfileName,
			},
			{
				Text: fmt.Sprintf("%t", cert.SelfSigned),
			},
			{
				Text: lo.FromPtrOr(cert.CN, "-"),
			},
		})
	}

	table.SetStyle(simpletable.StyleDefault)
	table.Println()
	return nil
}
