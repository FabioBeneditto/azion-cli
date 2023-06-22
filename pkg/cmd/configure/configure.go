package configure

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/messages/configure"
	msg "github.com/aziontech/azion-cli/messages/configure"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var configureToken string

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	// configureCmd represents the configure command
	configureCmd := &cobra.Command{
		Use:   msg.ConfigureUsage,
		Short: msg.ConfigureShortDescription,
		Long:  configure.ConfigureLongDescription,
		Example: heredoc.Doc(`
		$ azioncli configure --help
		$ azioncli configure --token azion123456abcdefg789asas1011hijklmn

        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := token.New(&token.Config{
				Client: f.HttpClient,
				Out:    f.IOStreams.Out,
			})
			if err != nil {
				return fmt.Errorf("%s: %w", utils.ErrorTokenManager, err)
			}

			if configureToken == "" {
				return utils.ErrorTokenNotProvided
			}

			valid, err := t.Validate(&configureToken)
			if err != nil {
				return err
			}

			if !valid {
				return utils.ErrorInvalidToken
			}

			if err := t.Save(); err != nil {
				return err
			}

			return nil
		},
	}

	configureCmd.SetIn(f.IOStreams.In)
	configureCmd.SetOut(f.IOStreams.Out)
	configureCmd.SetErr(f.IOStreams.Err)

	configureCmd.Flags().StringVarP(&configureToken, "token", "t", "", msg.ConfigureFlagToken)
	configureCmd.Flags().BoolP("help", "h", false, msg.ConfigureHelpFlag)

	return configureCmd
}
