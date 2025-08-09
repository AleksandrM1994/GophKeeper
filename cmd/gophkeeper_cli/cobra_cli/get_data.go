package cobra_cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/GophKeeper/internal/storage/sqlite"
	"github.com/GophKeeper/internal/utils"
)

// NewGetCmd создает команду для получения сохраненных данных по логину
func NewGetCmd(lg *zap.SugaredLogger, sqliteService *sqlite.ServiceImpl) *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Получить сохраненные данные пользователя по его логину",
		Long: `Команда позволяет пользователю ввести свой логин и получить список всех 
сохранённых приватных данных, связанных с этим логином.`,
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

			// Шаг 2: Получаем данные пользователя
			userData, errGetUserData := sqliteService.GetUserData(ctx, login)
			if errGetUserData != nil {
				return fmt.Errorf("ошибка при получении информации о пользователе: %v", errGetUserData)
			}

			if userData == nil {
				fmt.Println("Данные пользователя не найдены.")
				return nil
			}

			fmt.Printf("логин: %s\n", login)

			// Шаг 3: Получаем все приватные данные этого пользователя
			privateDataList, err := sqliteService.GetPrivateData(ctx, login)
			if err != nil {
				return fmt.Errorf("ошибка при получении приватных данных: %v", err)
			}

			// Шаг 4: Выводим результаты
			if len(privateDataList) == 0 {
				fmt.Println("У пользователя нет сохранённых данных.")
				return nil
			}

			fmt.Printf("Найдено %d записей для пользователя %s:\n", len(privateDataList), login)
			for i, pd := range privateDataList {
				decryptData, errDecrypt := utils.Decrypt(pd.Data, pd.Nonce, userData.Key)
				if errDecrypt != nil {
					lg.Errorf("Ошибка при расшифровке данных: %v", errDecrypt)
					fmt.Printf("%d. ID: %s | Тип: %s | Ошибка: %v\n", i+1, pd.ID, pd.Type, errDecrypt)
					continue
				}
				fmt.Printf("%d. ID: %s | Тип: %s | Данные: %s\n", i+1, pd.ID, pd.Type, string(decryptData))
			}

			return nil
		},
	}

	return getCmd
}
