package annotate

import (
	sharedcli "ksef/cmd/ksef/commands/invoices/shared"

	"github.com/spf13/cobra"
)

var AnnotationsCommand = &cobra.Command{
	Use:   "annotate",
	Short: "zarządzaj adnotacjami faktur",
}

func init() {
	AnnotationsCommand.AddCommand(commentCmd, annotateExcludeCmd, annotateVat50Cmd, annotateFixedAssetCmd, listCmd, clearCmd)
}

var commentCmd = &cobra.Command{
	Use:     "comment [faktura.xml]",
	Short:   "dodaje adnotację (komentarz) do pozycji faktury",
	RunE:    runComment,
	Args:    cobra.ExactArgs(1),
	PreRunE: initAnnotationRule,
}

func init() {
	flagSet := commentCmd.Flags()
	sharedcli.ItemSelectorFlags(flagSet)
	flagSet.String("text", "", "treść komentarza")
	_ = commentCmd.MarkFlagRequired("comment")
}

func runComment(cmd *cobra.Command, args []string) error {
	commentText, err := cmd.Flags().GetString("text")
	if err != nil {
		return err
	}

	global, err := cmd.Flags().GetBool("global")
	if err != nil {
		return err
	}

	rules, err := buildAnnotationRule(cmd, commentText)
	if err != nil {
		return err
	}

	return ctx.AnnotationsMgr.AddItemRules(ctx.Invoice, rules, global)
}
