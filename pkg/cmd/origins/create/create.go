package create

import (
  "strconv"
	"context"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"

	msg "github.com/aziontech/azion-cli/messages/domains"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
  sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/spf13/cobra"
)

type Fields struct {
	ApplicationID        int64
	Name                 string 
	OriginType           string 
	Addresses            []string
	OriginProtocolPolicy string 
	HostHeader           string 
	OriginPath           string 
	HmacAuthentication   string
	HmacRegionName       string 
	HmacAccessKey        string 
	HmacSecretKey        string 
	Path                 string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.DomainsCreateUsage,
		Short:         msg.DomainsCreateShortDescription,
		Long:          msg.DomainsCreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli origins create --application-id 1673635839 --name "ffcafe222sdsdffdf" --addresses "httpbin.org" --host-header "asdf.safe" --origin-type "single_origin" --origin-protocol-policy "http" --origin-path "/requests" --hmac-authentication "false"
        $ azioncli origins create --application-id 1673635839 --name "drink coffe" --addresses "asdfg.asd" --host-header "host"
        $ azioncli origins create --application-id 1673635839 --in "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.CreateOriginsRequest{}
			if cmd.Flags().Changed("in") {
				var (
					file *os.File
					err  error
				)
				if fields.Path == "-" {
					file = os.Stdin
				} else {
					file, err = os.Open(fields.Path)
					if err != nil {
						return fmt.Errorf("%w: %s", utils.ErrorOpeningFile, fields.Path)
					}
				}
				err = cmdutil.UnmarshallJsonFromReader(file, &request)
        fmt.Println("err ", err)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("name") || 
          !cmd.Flags().Changed("addresses") || !cmd.Flags().Changed("host-header") {  // flags requireds
					return msg.ErrorMandatoryCreateFlags
				}

        request.SetName(fields.Name)
        request.SetAddresses(prepareAddresses(fields.Addresses))
        request.SetHostHeader(fields.HostHeader)
        if cmd.Flags().Changed("origin-type") {
          request.SetOriginType(fields.OriginType)
        }
        if cmd.Flags().Changed("origin-protocol-policy") {
          request.SetOriginProtocolPolicy(fields.OriginProtocolPolicy)
        }
        if cmd.Flags().Changed("origin-path") {
          request.SetOriginPath(fields.OriginPath)
        }
        if cmd.Flags().Changed("hmac-authentication") { 
					hmacAuth, err := strconv.ParseBool(fields.HmacAuthentication)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorActiveFlag, fields.HmacAuthentication)
					}
          request.SetHmacAuthentication(hmacAuth)
        }
        if cmd.Flags().Changed("hmac-region-name") { 
          request.SetHmacRegionName(fields.HmacRegionName)
        }
        if cmd.Flags().Changed("hmac-access-key") {
          request.SetHmacAccessKey(fields.HmacAccessKey)
        }
        if cmd.Flags().Changed("hmac-secret-key") {
          request.SetHmacSecretKey(fields.HmacSecretKey)
        }
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.CreateOrigins(context.Background(), fields.ApplicationID, &request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateDomain.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, msg.DomainsCreateOutputSuccess, response.GetOriginId())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&fields.ApplicationID, "application-id", "a", 0, msg.DomainsCreateFlagIsActive)
	flags.StringVar(&fields.Name, "name", "", msg.DomainsCreateFlagName)
	flags.StringVar(&fields.OriginType, "origin-type", "", msg.DomainsCreateFlagName)
	flags.StringSliceVar(&fields.Addresses, "addresses", []string{}, msg.DomainsCreateFlagCnames)
	flags.StringVar(&fields.OriginProtocolPolicy, "origin-protocol-policy", "", msg.DomainsCreateFlagCnameAccessOnly)
	flags.StringVar(&fields.HostHeader, "host-header", "", msg.DomainsCreateFlagDigitalCertificateId)
	flags.StringVar(&fields.OriginPath, "origin-path", "", msg.DomainsCreateFlagEdgeApplicationId)
	flags.StringVar(&fields.HmacAuthentication, "hmac-authentication", "", msg.DomainsCreateFlagIsActive)
  flags.StringVar(&fields.HmacRegionName, "hmac-region-name", "", msg.DomainsCreateFlagIsActive)
	flags.StringVar(&fields.HmacAccessKey, "hmac-access-key", "", msg.DomainsCreateFlagIsActive)
	flags.StringVar(&fields.HmacSecretKey, "hmac-secret-key", "", msg.DomainsCreateFlagIsActive)
	flags.StringVar(&fields.Path, "in", "", msg.DomainsCreateFlagIn)
	flags.BoolP("help", "h", false, msg.DomainsCreateHelpFlag)
	return cmd
}

func prepareAddresses(addrs []string) (addresses []sdk.CreateOriginsRequestAddresses) {
  var addr sdk.CreateOriginsRequestAddresses
  for _, v := range addrs {
    addr.Address = v 
    addresses = append(addresses, addr)
  }
  return
}
