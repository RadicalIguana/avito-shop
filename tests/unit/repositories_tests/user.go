package repositories_test

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/RadicalIguana/avito-shop/internal/repositories"
// 	"github.com/pashagolub/pgxmock"
// 	"github.com/stretchr/testify/assert"
// )

// func TestGetUserCoins(t *testing.T) {
// 	mockDB, err := pgxmock.NewConn()
// 	if err != nil {
// 		t.Fatalf("Не удалось создать мок базы данных: %v", err)
// 	}
// 	defer mockDB.Close()

// 	r := repositories.NewUserInfoRepository(mockDB)
// 	ctx := context.Background()
// 	userID := 1
// 	t.Run("успешное получение количества монет", func(t *testing.T) {
// 		rows := pgxmock.NewRows([]string{"coins"}).AddRow(100)
// 		mockDB.ExpectQuery("SELECT coins FROM users WHERE id = $1").WithArgs(userID).WillReturnRows(rows)

// 		coins, err := r.GetUserCoins(ctx, userID)
// 		assert.NoError(t, err)
// 		assert.Equal(t, 100, coins)
// 	})

// 	t.Run("ошибка запроса", func(t *testing.T) {
// 		mockDB.ExpectQuery("SELECT coins FROM users WHERE id = $1").WithArgs(userID).WillReturnError(errors.New("db error"))

// 		coins, err := r.GetUserCoins(ctx, userID)
// 		assert.Error(t, err)
// 		assert.Equal(t, 0, coins)
// 	})

// 	t.Run("отсутствие пользователя", func(t *testing.T) {
// 		rows := pgxmock.NewRows([]string{"coins"}) // пустой результат
// 		mockDB.ExpectQuery("SELECT coins FROM users WHERE id = $1").WithArgs(userID).WillReturnRows(rows)

// 		coins, err := r.GetUserCoins(ctx, userID)
// 		assert.Error(t, err)
// 		assert.Equal(t, 0, coins)
// 	})
// }
