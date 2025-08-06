package cobra_cli

import (
	"bufio"
	"crypto/aes"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/GophKeeper/internal/client"
	"github.com/GophKeeper/internal/storage/sqlite"
	"github.com/GophKeeper/internal/utils"
	api "github.com/GophKeeper/pkg/api"
)

// authCmd represents the auth command
func NewAuthCmd(gophKeeperClient *client.ClientImpl, sqliteService *sqlite.ServiceImpl) *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "авторизация пользователя в системе",
		Long: `авторизация пользователя в системе обязательна для выполнения команд, 
которые предоставляют доступ к возможностям системы
нужно ввести логин и пароль, программа запрашивает их у пользователя поочереди`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Шаг 1: спрашиваем логин
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Введите логин: ")
			login, err := reader.ReadString('\n')
			if err != nil {
				os.Exit(1)
			}
			// Убираем символ переноса строки
			login = login[:len(login)-1]

			// Шаг 2: спрашиваем пароль (в невидимом режиме)
			fmt.Print("Введите пароль: ")
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println() // перевод строки после ввода пароля
			if err != nil {
				fmt.Fprintf(os.Stderr, "Ошибка при вводе пароля: %v\n", err)
				os.Exit(1)
			}
			password := string(bytePassword)

			// Дальше вы можете использовать login и password
			fmt.Printf("Вы ввели:\n  логин: %s\n  пароль: %s\n", login, password)

			res, errAuthUser := gophKeeperClient.AuthUser(ctx, &api.AuthUserRequest{
				Login:    login,
				Password: password,
			})
			if errAuthUser != nil {
				return fmt.Errorf("ошибка при авторизации: %v", errAuthUser)
			}

			userData, errGetUserData := sqliteService.GetUserData(login)
			if errGetUserData != nil {
				return fmt.Errorf("ошибка при получении информации о пользователе: %v", errGetUserData)
			}

			if userData == nil {
				key, errGenerateKey := utils.GenerateKey(2 * aes.BlockSize)
				if errGenerateKey != nil {
					return fmt.Errorf("ошибка при генерации ключа: %v", errGenerateKey)
				}

				fmt.Println(key, len(key))

				userData = &sqlite.UserData{
					Login: login,
					Key:   key,
				}
				if errGenerateKey != nil {
					return fmt.Errorf("ошибка при генерации ключа: %v", errGenerateKey)
				}
			}
			userData.JWT = res.Jwt

			errSaveUser := sqliteService.SaveUserData(userData)
			if errSaveUser != nil {
				return fmt.Errorf("ошибка при сохранении информации о пользователе: %v", errSaveUser)
			}

			fmt.Println("Пользователь успешно авторизован!")

			return nil
		},
	}

	return authCmd
}
