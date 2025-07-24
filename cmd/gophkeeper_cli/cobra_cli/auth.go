package cobra_cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/GophKeeper/internal/client"
	"github.com/GophKeeper/internal/handlers/user"
	"github.com/GophKeeper/internal/storage/bbolt"
)

// authCmd represents the auth command
func NewAuthCmd(gophKeeperClient *client.ClientImpl, bboltService *bbolt.ServiceImpl) *cobra.Command {
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

			res, errAuthUser := gophKeeperClient.AuthUser(ctx, &user.AuthUserRequest{
				Login:    login,
				Password: password,
			})
			if errAuthUser != nil {
				return fmt.Errorf("Ошибка при авторизации: %v", errAuthUser)
			}

			errSaveUser := bboltService.SaveUserData(bbolt.UserData{
				Login: login,
				JWT:   res.JWT,
			})
			if errSaveUser != nil {
				return fmt.Errorf("Ошибка при сохранении информации о пользователе: %v", errSaveUser)
			}

			return nil
		},
	}

	return authCmd
}
