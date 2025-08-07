package cobra_cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/GophKeeper/internal/storage/sqlite"
	"github.com/GophKeeper/internal/utils"
)

// NewGetCmd создает команду для получения сохраненных данных по логину
func NewGetCmd(sqliteService *sqlite.ServiceImpl) *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Получить сохраненные данные пользователя по его логину",
		Long: `Команда позволяет пользователю ввести свой логин и получить список всех 
сохранённых приватных данных, связанных с этим логином.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// спрашиваем логин
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Введите логин: ")
			login, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("ошибка при вводе логина: %v", err)
			}
			login = strings.TrimSpace(login)

			// Получаем данные пользователя
			userData, errGetUserData := sqliteService.GetUserData(ctx, login)
			if errGetUserData != nil {
				return fmt.Errorf("ошибка при получении информации о пользователе: %v", errGetUserData)
			}

			if userData == nil {
				fmt.Println("Данные пользователя не найдены.")
				return nil
			}

			fmt.Printf("логин: %s\n", login)

			// Получаем все приватные данные этого пользователя
			privateDataList, err := sqliteService.GetPrivateData(ctx, login)
			if err != nil {
				return fmt.Errorf("ошибка при получении приватных данных: %v", err)
			}

			// Выводим результаты
			fmt.Printf("Найдено %d записей для пользователя %s:\n", len(privateDataList), login)
			for i, pd := range privateDataList {
				decryptData, errDecrypt := utils.Decrypt(pd.Data, pd.Nonce, userData.Key)
				if errDecrypt != nil {
					fmt.Println("Ошибка при расшифровке данных:", errDecrypt)
					return nil
				}
				fmt.Printf("%d. ID: %s | Тип: %s | Данные: %s\n", i+1, pd.ID, pd.Type, string(decryptData))
			}

			return nil
		},
	}

	return getCmd
}
