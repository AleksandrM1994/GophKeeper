package cobra_cli

import (
	"bufio"
	"crypto/aes"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/term"

	"github.com/GophKeeper/internal/client"
	"github.com/GophKeeper/internal/storage/sqlite"
	"github.com/GophKeeper/internal/utils"
	api "github.com/GophKeeper/pkg/api"
)

// authCmd represents the auth command
func NewAuthCmd(lg *zap.SugaredLogger, gophKeeperClient *client.ClientImpl, sqliteService *sqlite.ServiceImpl) *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Авторизация пользователя в системе",
		Long: `Авторизация пользователя в системе обязательна для выполнения команд, 
которые предоставляют доступ к возможностям системы.
Нужно ввести логин и пароль, программа запрашивает их у пользователя поочереди.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			reader := bufio.NewReader(os.Stdin)

			// Логин
			fmt.Print("Введите логин: ")
			login, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("ошибка при вводе логина: %w", err)
			}
			login = strings.TrimSpace(login)

			// Пароль
			var password string
			if term.IsTerminal(syscall.Stdin) {
				fmt.Print("Введите пароль: ")
				bytePwd, err := term.ReadPassword(syscall.Stdin)
				fmt.Println()
				if err != nil {
					return fmt.Errorf("ошибка при вводе пароля: %w", err)
				}
				password = string(bytePwd)
			} else {
				// НЕ-TTY: читаем как обычную строку
				fmt.Print("Введите пароль (будет видно на экране): ")
				pwdLine, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("ошибка при вводе пароля: %w", err)
				}
				password = strings.TrimSpace(pwdLine)
			}

			lg.Infof("Пользователь ввел: логин=%s", login)

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
