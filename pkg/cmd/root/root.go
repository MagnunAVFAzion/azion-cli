package root

import (
	"fmt"
	"net/http"
	"time"

	msg "github.com/aziontech/azion-cli/messages/root"
	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/constants"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	tokenFlag  string
	configFlag string
)

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	version := version.BinVersion

	rootCmd := &cobra.Command{
		Use:     msg.RootUsage,
		Short:   color.New(color.Bold).Sprint(fmt.Sprintf(msg.RootDescription, version)),
		Version: version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := doPreCommandCheck(cmd, f, PreCmd{
				config: configFlag,
				token:  tokenFlag,
			})

			if err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceErrors: true, // Silence errors, so the help message won't be shown on flag error
		SilenceUsage:  true, // Silence usage on error
	}

	rootCmd.SetIn(f.IOStreams.In)
	rootCmd.SetOut(f.IOStreams.Out)
	rootCmd.SetErr(f.IOStreams.Err)

	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		rootHelpFunc(f, cmd, args)
	})

	//Global flags
	rootCmd.PersistentFlags().StringVarP(&tokenFlag, "token", "t", "", msg.RootTokenFlag)
	rootCmd.PersistentFlags().StringVarP(&configFlag, "config", "c", "", msg.RootConfigFlag)

	//other flags
	rootCmd.Flags().BoolP("help", "h", false, msg.RootHelpFlag)

	//set template for -v flag
	rootCmd.SetVersionTemplate(color.New(color.Bold).Sprint("Azion CLI " + version + "\n")) // TODO: Change to version.BinVersion once 1.0 is released

	return rootCmd
}

func Execute() {
	streams := iostreams.System()
	httpClient := &http.Client{
		Timeout: 10 * time.Second, // TODO: Configure this somewhere
	}

	// TODO: Ignoring errors since the file might not exist, maybe warn the user?
	tok, _ := token.ReadFromDisk()
	viper.SetEnvPrefix("AZIONCLI")
	viper.AutomaticEnv()
	viper.SetDefault("token", tok)
	viper.SetDefault("api_url", constants.ApiURL)
	viper.SetDefault("storage_url", constants.StorageApiURL)

	factory := &cmdutil.Factory{
		HttpClient: httpClient,
		IOStreams:  streams,
		Config:     viper.GetViper(),
	}

	cmd := NewRootCmd(factory)

	cobra.CheckErr(cmd.Execute())
}
