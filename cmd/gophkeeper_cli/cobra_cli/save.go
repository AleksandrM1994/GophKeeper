package cobra_cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/GophKeeper/internal/client"
	"github.com/GophKeeper/internal/client/dto"
	"github.com/GophKeeper/internal/repository"
	"github.com/GophKeeper/internal/storage/sqlite"
	"github.com/GophKeeper/internal/utils"
)

var (
	textMode             bool
	fileMode             bool
	authMode             bool
	bankMode             bool
	login                string
	textData             string
	filePath             string
	authLogin            string
	authPass             string
	bankCardNumber       string
	bankCardPersonName   string
	bankCardActiveDateTo string
	bankCardCVC          string
)

func NewSaveCmd(gophKeeperClient *client.ClientImpl, sqliteService *sqlite.ServiceImpl) *cobra.Command {
	saveCmd := &cobra.Command{
		Use:   "save",
		Short: "Сохранить данные (текст, файл, аутентификацию или банковские данные)",
		Long: `Команда save позволяет сохранять разные типы данных:
- текст (--text)
- файл (--file)
- данные аутентификации (--auth)
- банковские данные (--bank)

Обязательно указать --login и ровно один из флагов.`,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			if login == "" {
				return errors.New("обязательно указать --user")
			}

			var data []byte
			var dataType repository.PrivateDataType

			switch {
			case textMode:
				dataType = repository.PrivateDataTypeText
				data = []byte(textData)
			case fileMode:
				dataType = repository.PrivateDataTypeFile
				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					return fmt.Errorf("не удалось прочитать файл: %w", err)
				}
				data = fileContent
			case authMode:
				dataType = repository.PrivateDataTypeAuth
				authData := fmt.Sprintf("login:%s|password:%s", authLogin, authPass)
				data = []byte(authData)
			case bankMode:
				dataType = repository.PrivateDataTypeBank
				bankData := fmt.Sprintf("number:%s|name:%s|date:%s|cvc:%s",
					bankCardNumber, bankCardPersonName, bankCardActiveDateTo, bankCardCVC)
				data = []byte(bankData)
			}

			userData, errGetUserData := sqliteService.GetUserData(ctx, login)
			if errGetUserData != nil {
				return fmt.Errorf("get user data: %w", errGetUserData)
			}

			fmt.Printf("Пользователь найден: %s\n", userData)

			encryptData, nonce, errEncrypt := utils.Encrypt(data, userData.Key)
			if errEncrypt != nil {
				return fmt.Errorf("encrypt data: %w", errEncrypt)
			}

			fmt.Println(nonce)

			err := gophKeeperClient.SavePrivateData(ctx, &dto.SavePrivateDataRequest{
				Type:  dataType,
				Data:  encryptData,
				JWT:   userData.JWT,
				Nonce: nonce,
				Login: login,
			})
			if err != nil {
				fmt.Println(err)
				return fmt.Errorf("save private data: %w", err)
			}

			fmt.Printf("Данные успешно сохранены под логином '%s' (тип: %s)\n", login, dataType)
			return nil
		},
	}

	// Регистрация флагов
	saveCmd.Flags().BoolVarP(&textMode, "text", "t", false, "Save text data")
	saveCmd.Flags().StringVarP(&textData, "text-data", "", "", "Text to save")

	saveCmd.Flags().BoolVarP(&fileMode, "file", "f", false, "Save file data")
	saveCmd.Flags().StringVarP(&filePath, "file-path", "", "", "Path to file to save")

	saveCmd.Flags().BoolVarP(&authMode, "auth", "a", false, "Save authentication data")
	saveCmd.Flags().StringVarP(&authLogin, "auth-login", "", "", "Authentication login")
	saveCmd.Flags().StringVarP(&authPass, "auth-pass", "", "", "Authentication password")

	saveCmd.Flags().BoolVarP(&bankMode, "bank", "b", false, "Save bank data")
	saveCmd.Flags().StringVarP(&bankCardNumber, "bank-number", "", "", "Bank card number")
	saveCmd.Flags().StringVarP(&bankCardPersonName, "bank-name", "", "", "Bank card number")
	saveCmd.Flags().StringVarP(&bankCardActiveDateTo, "bank-date", "", "", "Bank card number")
	saveCmd.Flags().StringVarP(&bankCardCVC, "bank-cvc", "", "", "Bank CVC code")

	saveCmd.Flags().StringVarP(&login, "login", "l", "", "Обязательный логин для данных")

	return saveCmd
}
