package service

import (
	"context"
	"errors"
	"incomster/backend/store"
	"incomster/config"
	"incomster/core"
	errs "incomster/pkg/apperrors"
	"incomster/pkg/jwt"
	"time"
)

type AccountService struct {
	user      store.IUserStore
	session   store.ISessionStore
	tokenizer *jwt.Tokenizer
	config    *config.Config
}

func NewAccountService(session store.ISessionStore, user store.IUserStore, tokenizer *jwt.Tokenizer, config *config.Config) *AccountService {
	return &AccountService{session: session, user: user, tokenizer: tokenizer, config: config}
}

func (s *AccountService) Register(ctx context.Context, input *core.UserCreateInput) (*core.Session, error) {
	user, err := s.user.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenizer.Generate(user.ID, user.Role)
	if err != nil {
		return nil, failedToGenerateToken
	}

	session, err := s.session.Create(ctx, &core.SessionCreateInput{
		UserID: user.ID,
		Token:  token,
	})
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AccountService) Login(ctx context.Context, input *core.UserLoginInput) (*core.Session, error) {
	user, err := s.user.Get(ctx, &core.UserGetInput{Username: &input.Username})
	if errors.Is(err, errs.ErrorUserNotFound) {
		return nil, invalidCredential
	}
	if err != nil {
		return nil, failedToRetrieveUser
	}

	session, err := s.session.Get(ctx, &core.SessionGetInput{UserID: user.ID})
	if err != nil {
		return s.handleFirstLogin(ctx, user)
	}

	return s.handleSecondLogin(ctx, session)
}

func (s *AccountService) Logout(ctx context.Context, userId int) error {
	_, err := s.session.Delete(ctx, &core.SessionGetInput{UserID: userId})
	if err != nil {
		return err
	}
	return nil
}

func (s *AccountService) handleFirstLogin(ctx context.Context, input *core.User) (*core.Session, error) {
	// Генерируем новый токен
	token, err := s.tokenizer.Generate(input.ID, input.Role)
	if err != nil {
		return nil, failedToGenerateToken
	}

	// Добавляем новую сессию
	session, err := s.session.Create(ctx, &core.SessionCreateInput{UserID: input.ID, Token: token})
	if err != nil {
		return nil, failedToCreateSession
	}

	return session, nil
}

func (s *AccountService) handleSecondLogin(ctx context.Context, input *core.Session) (*core.Session, error) {
	// Парсим данные токена
	claims, err := s.tokenizer.Parse(input.Token)
	if err != nil {
		return nil, failedToValidateToken
	}

	// Если данные валидны на текущий момент - возвращаем ошибку некорректных данных авторизации.
	// Просто чтоб пользователь не знал что введенные данные валидные данные.
	if claims.IsValidAt(time.Now()) {
		return nil, invalidCredential
	}

	// Генерируем новый токен
	token, err := s.tokenizer.Generate(claims.UserID, claims.Role)
	if err != nil {
		return nil, failedToGenerateToken
	}

	// Обновляем данные сессии
	session, err := s.session.Update(ctx, &core.SessionUpdateInput{Id: input.ID, UserID: input.UserID, Token: token})
	if err != nil {
		return nil, failedToUpdateToken
	}

	return session, nil
}
