package service

import (
	"context"
	"incomster/backend/store"
	"incomster/core"
	errs "incomster/pkg/apperrors"
	"incomster/pkg/ctxutil"
	"incomster/pkg/jwt"
	"time"
)

type SecurityService struct {
	session   store.ISessionStore
	tokenizer *jwt.Tokenizer
}

func NewSecurityService(session store.ISessionStore, tokenizer *jwt.Tokenizer) *SecurityService {
	return &SecurityService{tokenizer: tokenizer, session: session}
}

func (s *SecurityService) HandleBearerAuth(ctx context.Context, operation, token string) (context.Context, error) {
	// TODO: handle user role here?
	// Надо ли тут проверять роли, ссылаясь на аргумент <operation>?
	// Возвращать ошибку если пользователь с ролью "user" пытается вызвать операцию для которой нужна роль "admin"
	//------------------------------------------------------------------------------------------------------------------
	// UPD :
	// Концепция изменилась, для админки сделаю отдельный API или эндпоинты, если, блядь, успею.
	// Соответственно нет нужды проверять роли здесь.
	//------------------------------------------------------------------------------------------------------------------
	// Но я все ровно не уверен что это был бы корректный подход проверять роли тут, хотя и достаточно удобный.

	// Проверка наличия токена в базе авторизированных сессий
	session, err := s.session.Get(ctx, &core.SessionGetInput{Token: token})
	if err != nil {
		return nil, errs.Unauthorized("invalid token", err)
	}

	// Получение данных токена ранее авторизованной сессии
	claims, err := s.tokenizer.Parse(session.Token)
	if err != nil {
		return ctx, errs.Unauthorized("invalid token", err)
	}

	// Проверка срока пригодности токена
	if !claims.IsValidAt(time.Now()) {
		return ctx, errs.Unauthorized("token expired")
	}

	// Добавление данных в контекст
	ctx = ctxutil.WithUserId(ctx, claims.UserID)
	return ctx, nil
}
