package cobra_cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/GophKeeper/internal/client"
	"github.com/GophKeeper/internal/handlers/private_data"
)

var (
	textMode bool
	fileMode bool
	authMode bool
	bankMode bool
)

// authCmd represents the auth command
func NewSaveCmd(gophKeeperClient *client.ClientImpl) *cobra.Command {
	saveCmd := &cobra.Command{
		Use:   "save",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("save called")
			ctx := cmd.Context()

			modes := 0
			for _, b := range []bool{textMode, fileMode, authMode, bankMode} {
				if b {
					modes++
				}
			}
			if modes != 1 {
				return errors.New("нужно указать ровно один режим: --text или --file или --auth или --bank")
			}

			reader := bufio.NewReader(os.Stdin)
			text, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			text = strings.TrimRight(text, "\r\n")

			fmt.Println("text: ", text)

			err = gophKeeperClient.SavePrivateData(ctx, &private_data.SavePrivateDataRequest{
				Data: []byte("{\"data\":\"123\"}"),
			})
			if err != nil {
				return fmt.Errorf("save private data: %w", err)
			}

			return nil
		},
	}
	return saveCmd
}
