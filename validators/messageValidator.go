package validators

import (
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

func MessageValidator(jsonStr string) (*models.Order, error) {
	var order models.Order

	// Декодирование JSON в структуру
	if err := json.Unmarshal([]byte(jsonStr), &order); err != nil {
		return nil, fmt.Errorf("ошибка декодирования JSON: %v", err)
	}

	// Создание валидатора
	validate := validator.New()

	// Валидация структуры
	if err := validate.Struct(order); err != nil {
		return nil, fmt.Errorf("ошибка валидации структуры: %v", err)
	}

	return &order, nil
}
