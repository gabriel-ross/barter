package transaction

import (
	"context"

	"github.com/gabriel-ross/barter"
	"github.com/gabriel-ross/barter/model"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (svc *Service) create(ctx context.Context, data model.Transaction) (_ model.Transaction, err error) {
	data.ID = svc.db.Collection("transactions").NewDoc().ID
	_, err = svc.db.Collection("transactions").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return model.Transaction{}, err
	}

	return data, nil
}

func (svc *Service) list(ctx context.Context, options ...barter.QueryOption) (_ []model.Transaction, err error) {
	resp := []model.Transaction{}
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

		var m model.Transaction
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

func (svc *Service) read(ctx context.Context, id string) (_ model.Transaction, err error) {
	dsnap, err := svc.db.Collection("transactions").Doc(id).Get(ctx)
	if err != nil {
		return model.Transaction{}, err
	}

	var m model.Transaction
	dsnap.DataTo(&m)
	return m, nil
}

func (svc *Service) exists(ctx context.Context, id string) (_ bool, err error) {
	_, err = svc.db.Collection("transactions").Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (svc *Service) update(ctx context.Context, id string, data model.Transaction) (_ model.Transaction, err error) {
	data.ID = id
	_, err = svc.db.Collection("transactions").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return model.Transaction{}, err
	}
	return data, nil
}

func (svc *Service) updateNonZero(ctx context.Context, id string, data model.Transaction) (_ model.Transaction, err error) {
	return data, nil
}

func (svc *Service) delete(ctx context.Context, id string) (err error) {
	_, err = svc.db.Collection("transactions").Doc(id).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
