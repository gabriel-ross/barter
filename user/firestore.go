package user

import (
	"context"
	"fmt"
	"reflect"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/barter/model"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (svc *Service) create(ctx context.Context, data model.User) (_ model.User, err error) {
	data.ID = svc.db.Collection("users").NewDoc().ID
	_, err = svc.db.Collection("users").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return model.User{}, err
	}
	return data, nil
}

func (svc *Service) list(ctx context.Context, offset, limit int) (_ []model.User, err error) {
	resp := []model.User{}
	iter := svc.db.Collection("users").OrderBy("id", firestore.Asc).StartAt(offset).Limit(limit).Documents(ctx)
	for {
		dsnap, err := iter.Next()
		if err == iterator.Done {
			break
		}

		var m model.User
		dsnap.DataTo(&m)
		resp = append(resp, m)
	}
	return resp, nil
}

func (svc *Service) count(ctx context.Context) (_ int, err error) {
	docs, err := svc.db.Collection("users").Documents(ctx).GetAll()
	if err != nil {
		return 0, err
	}
	return len(docs), nil
}

func (svc *Service) read(ctx context.Context, id string) (_ model.User, err error) {
	dsnap, err := svc.db.Collection("users").Doc(id).Get(ctx)
	if err != nil {
		return model.User{}, err
	}

	var m model.User
	dsnap.DataTo(&m)
	return m, nil
}

func (svc *Service) exists(ctx context.Context, id string) (_ bool, err error) {
	_, err = svc.db.Collection("users").Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (svc *Service) set(ctx context.Context, id string, data model.User) (_ model.User, err error) {
	data.ID = id
	_, err = svc.db.Collection("users").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return model.User{}, err
	}
	return data, nil
}

func (svc *Service) updateNonZero(ctx context.Context, id string, data model.User) (_ model.User, err error) {

	foo := reflect.TypeOf(data)
	for i := 0; i < foo.NumField(); i++ {
		path := foo.Field(i).Tag.Get("firestore")
		fmt.Println(path)
	}

	// TODO: build firestore update slice
	// TODO: consider how this might work with nested objects

	return data, nil
}

func (svc *Service) delete(ctx context.Context, id string) (err error) {
	_, err = svc.db.Collection("users").Doc(id).Delete(ctx)
	return err
}
