package cobra_cli

import (
	"bufio"
	"crypto/aes"
	"fmt"
	"os"
	"strings"

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
				return fmt.Errorf("ошибка при вводе логина: %v", err)
			}
			login = strings.TrimSpace(login)

			// Шаг 2: спрашиваем пароль (в невидимом режиме)
			fmt.Print("Введите пароль: ")
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return fmt.Errorf("ошибка при вводе пароля: %v", err)
			}
			fmt.Println() // перевод строки после ввода пароля
			password := string(bytePassword)

			fmt.Printf("Вы ввели:\n  логин: %s\n  пароль: %s\n", login, password)

			// Шаг 3: авторизация пользователя
			res, errAuthUser := gophKeeperClient.AuthUser(ctx, &api.AuthUserRequest{
				Login:    login,
				Password: password,
			})
			if errAuthUser != nil {
				return fmt.Errorf("ошибка при авторизации: %v", errAuthUser)
			}

			// Шаг 4: получение информации о пользователе из SQLite
			userData, errGetUserData := sqliteService.GetUserData(ctx, login)
			if errGetUserData != nil {
				return fmt.Errorf("ошибка при получении информации о пользователе: %v", errGetUserData)
			}

			if userData == nil {
				// Генерация ключа
				key, errGenerateKey := utils.GenerateKey(2 * aes.BlockSize)
				if errGenerateKey != nil {
					return fmt.Errorf("ошибка при генерации ключа: %v", errGenerateKey)
				}

				userData = &sqlite.UserData{
					Login: login,
					Key:   key,
				}
			}
			// Шаг 5: установка JWT
			userData.JWT = res.Jwt

			fmt.Printf("данные для сохранения в БД: %w", userData)

			// Шаг 6: сохранение данных пользователя
			errSaveUser := sqliteService.SaveUserData(ctx, userData)
			if errSaveUser != nil {
				return fmt.Errorf("ошибка при сохранении информации о пользователе: %v", errSaveUser)
			}

			fmt.Println("Пользователь успешно авторизован!")
			return nil
		},
	}

	return authCmd
}
