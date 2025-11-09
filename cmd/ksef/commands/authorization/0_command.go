package authorization

import (
	"github.com/spf13/cobra"
)

var AuthCommand = &cobra.Command{
	Use:   "auth",
	Short: "zarządzanie autoryzacją KSeF",
}
