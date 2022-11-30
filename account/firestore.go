package account

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/barter/model"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (svc *Service) create(ctx context.Context, data model.Account) (_ model.Account, err error) {
	data.ID = svc.db.Collection("accounts").NewDoc().ID
	data.CreationTimestamp = time.Now()
	_, err = svc.db.Collection("accounts").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return model.Account{}, err
	}
	return data, nil
}

func (svc *Service) list(ctx context.Context, offset, limit int) (_ []model.Account, err error) {
	resp := []model.Account{}
	iter := svc.db.Collection("accounts").OrderBy("id", firestore.Asc).StartAt(offset).Limit(limit).Documents(ctx)
	for {
		dsnap, err := iter.Next()
		if err == iterator.Done {
			break
		}

		var m model.Account
		dsnap.DataTo(&m)
		resp = append(resp, m)
	}
	return resp, nil
}

func (svc *Service) count(ctx context.Context) (_ int, err error) {
	docs, err := svc.db.Collection("accounts").Documents(ctx).GetAll()
	if err != nil {
		return 0, err
	}
	return len(docs), nil
}

func (svc *Service) read(ctx context.Context, id string) (_ model.Account, err error) {
	dsnap, err := svc.db.Collection("accounts").Doc(id).Get(ctx)
	if err != nil {
		return model.Account{}, err
	}

	var m model.Account
	dsnap.DataTo(&m)
	return m, nil
}

func (svc *Service) exists(ctx context.Context, id string) (_ bool, err error) {
	_, err = svc.db.Collection("accounts").Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (svc *Service) update(ctx context.Context, id string, data model.Account) (_ model.Account, err error) {
	data.ID = id
	_, err = svc.db.Collection("accounts").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return model.Account{}, err
	}
	return data, nil
}

func (svc *Service) updateNonZero(ctx context.Context, id string, data model.Account) (_ model.Account, err error) {
	return data, nil
}

func (svc *Service) delete(ctx context.Context, id string) (err error) {
	_, err = svc.db.Collection("accounts").Doc(id).Delete(ctx)
	return err
}
