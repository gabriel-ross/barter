package account

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/barter"
	"github.com/gabriel-ross/barter/model"
	"google.golang.org/api/iterator"
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

func (svc *Service) list(ctx context.Context, options ...barter.QueryOption) (_ []model.Account, err error) {
	resp := []model.Account{}
	query := svc.db.Collection("accounts").Query
	for _, option := range options {
		query = option(query)
	}
	iter := query.Documents(ctx)
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

func (svc *Service) count(ctx context.Context, options ...barter.QueryOption) (_ int, err error) {
	query := svc.db.Collection("accounts").Query
	for _, option := range options {
		query = option(query)
	}
	docs, err := query.Documents(ctx).GetAll()
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

func (svc *Service) set(ctx context.Context, id string, data model.Account) (_ model.Account, err error) {
	data.ID = id
	_, err = svc.db.Collection("accounts").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return model.Account{}, err
	}
	return data, nil
}

func (svc *Service) updateNonZero(ctx context.Context, id string, data model.Account) (_ model.Account, err error) {
	updates := []firestore.Update{}
	if data.Owner != "" {
		updates = append(updates, firestore.Update{
			Path:  "owner",
			Value: data.Owner,
		})
	}
	for key, val := range data.Funds {
		updates = append(updates, firestore.Update{
			Path:  "funds." + key,
			Value: firestore.Increment(val),
		})
	}
	if data.Reputation > 0 {
		updates = append(updates, firestore.Update{
			Path:  "reputation",
			Value: data.Reputation,
		})
	}

	_, err = svc.db.Collection("accounts").Doc(id).Update(ctx, updates)
	if err != nil {
		return model.Account{}, err
	}

	return data, nil
}

func (svc *Service) delete(ctx context.Context, id string) (err error) {
	_, err = svc.db.Collection("accounts").Doc(id).Delete(ctx)
	return err
}
