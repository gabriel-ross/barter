package auth

import (
	"context"

	"github.com/gabriel-ross/barter/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (svc *Service) createState(ctx context.Context, id string) (err error) {
	_, err = svc.db.Collection("states").Doc(id).Set(ctx, struct{}{})
	if err != nil {
		return err
	}
	return nil
}

func (svc *Service) read(ctx context.Context, id string) (err error) {
	_, err = svc.db.Collection("states").Doc(id).Get(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (svc *Service) createUser(ctx context.Context, m model.User) (_ model.User, err error) {
	_, err = svc.db.Collection("users").Doc(m.ID).Get(ctx)
	if status.Code(err) == codes.NotFound {
		_, err = svc.db.Collection("users").Doc(m.ID).Set(ctx, m)
		if err != nil {
			return model.User{}, err
		}
	} else if err != nil {
		return m, err
	}
	return model.User{}, nil
}
