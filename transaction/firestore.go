package transaction

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gabriel-ross/barter"
	"github.com/gabriel-ross/barter/model"
	"google.golang.org/api/iterator"
)

var zeroTime time.Time

func (svc *Service) create(ctx context.Context, data model.Transaction) (_ model.Transaction, err error) {
	data.ID = svc.db.Collection("transactions").NewDoc().ID
	_, err = svc.db.Collection("transactions").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return model.Transaction{}, err
	}

	payments := []firestore.Update{}
	credits := []firestore.Update{}
	for key, val := range data.Quantities {
		payments = append(payments, firestore.Update{
			Path:  "balances." + key,
			Value: firestore.Increment((-1) * val),
		})
		credits = append(payments, firestore.Update{
			Path:  "balances." + key,
			Value: firestore.Increment(val),
		})
	}
	_, err = svc.db.Collection("accounts").Doc(data.SenderAccountID).Update(ctx, payments)
	if err != nil {
		return model.Transaction{}, err
	}
	_, err = svc.db.Collection("accounts").Doc(data.RecipientAccountID).Update(ctx, credits)
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

func (svc *Service) set(ctx context.Context, id string, data model.Transaction) (_ model.Transaction, err error) {
	data.ID = id
	_, err = svc.db.Collection("transactions").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return model.Transaction{}, err
	}
	return data, nil
}

func (svc *Service) updateNonZero(ctx context.Context, id string, data model.Transaction) (_ model.Transaction, err error) {
	updates := []firestore.Update{}
	if len(data.Quantities) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "quantities",
			Value: data.Quantities,
		})
	}
	if data.SenderAccountID != "" {
		updates = append(updates, firestore.Update{
			Path:  "sender",
			Value: data.SenderAccountID,
		})
	}
	if data.RecipientAccountID != "" {
		updates = append(updates, firestore.Update{
			Path:  "sender",
			Value: data.RecipientAccountID,
		})
	}
	if data.Timestamp.After(zeroTime) {
		updates = append(updates, firestore.Update{
			Path:  "timestamp",
			Value: data.Timestamp,
		})
	}

	_, err = svc.db.Collection("transactions").Doc(id).Update(ctx, updates)
	if err != nil {
		return model.Transaction{}, err
	}

	return data, nil
}

func (svc *Service) delete(ctx context.Context, id string) (err error) {
	_, err = svc.db.Collection("transactions").Doc(id).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (svc *Service) deleteWithCascade(ctx context.Context, id string) (err error) {
	data, err := svc.read(ctx, id)
	if err != nil {
		return err
	}

	_, err = svc.db.Collection("transactions").Doc(id).Delete(ctx)
	if err != nil {
		return err
	}

	payments := []firestore.Update{}
	credits := []firestore.Update{}
	for key, val := range data.Quantities {
		payments = append(payments, firestore.Update{
			Path:  "balances." + key,
			Value: firestore.Increment((-1) * val),
		})
		credits = append(payments, firestore.Update{
			Path:  "balances." + key,
			Value: firestore.Increment(val),
		})
	}
	_, err = svc.db.Collection("accounts").Doc(data.SenderAccountID).Update(ctx, credits)
	if err != nil {
		return err
	}
	_, err = svc.db.Collection("accounts").Doc(data.RecipientAccountID).Update(ctx, payments)
	if err != nil {
		return err
	}

	return nil
}
