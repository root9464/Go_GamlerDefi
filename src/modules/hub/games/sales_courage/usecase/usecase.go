package sales_courafe_usecase

import (
	"context"

	sales_courage_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/sales_courage/entity"
)

type SalesCourageUsecase struct {
}

func (u *SalesCourageUsecase) StartGame() {}

func (u *SalesCourageUsecase) RollDices() []int { return nil }

func (u *SalesCourageUsecase) ProcessMove() {}

func (u *SalesCourageUsecase) IssueDeck(ctx context.Context, userID string) (*sales_courage_models.Deck, error) {
	// ведущий выбирает колоду, а так же игррока, которому он предоставляет колоду
	return nil, nil
}

func (u *SalesCourageUsecase) AddMoney(ctx context.Context, userID string, amount int) error {
	// ведущий добавляет монеты игроку
	// если введет отрицательное число, то игрок теряет монеты
	//
	// по мимо того, что у пользователя высвечивается колода, админ получает сообщение, что пользователь выбрал карту
	return nil
}

func (u *SalesCourageUsecase) GetUserHand(ctx context.Context, userID string) (*sales_courage_models.Deck, error) {
	return nil, nil
}
