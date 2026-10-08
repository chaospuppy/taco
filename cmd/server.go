package cmd

import (
	"errors"

	"github.com/rodneyxr/taco/internal/taco"
	"github.com/spf13/cobra"
)

type serverOptions struct {
	port   int
	listen string
}

// NewServerCommand builds the CLI web command
func NewServerCommand() *cobra.Command {
	opts := &serverOptions{}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run taco as a service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.run(cmd)
		},
	}

	flags := cmd.Flags()
	flags.IntVarP(&opts.port, "port", "p", 8080, "port to listen on")
	flags.StringVar(&opts.listen, "listen", "0.0.0.0", "listen address for server")

	return cmd
}

func (opts *serverOptions) run(cmd *cobra.Command) error {
	cfgPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	cfg, err := taco.LoadConfig(cfgPath, cmd.Flags().Changed("config"))
	if err != nil {
		return err
	}

	if len(cfg.Teams) == 0 {
		return errors.New("no teams found; add a teams section to taco.yaml")
	}

	return nil
}
