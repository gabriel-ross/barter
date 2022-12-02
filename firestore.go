package barter

import "cloud.google.com/go/firestore"

// Filter operators
var (
	Eq  operator = "=="
	Ne           = "!="
	Lt           = "<"
	Leq          = "<="
	Gt           = ">"
	Geq          = ">="
)

type QueryOption func(firestore.Query) firestore.Query

func WithFilter(key string, op operator, val interface{}) func(firestore.Query) firestore.Query {
	return func(q firestore.Query) firestore.Query {
		return q.Where(key, string(op), val)
	}
}

func WithOrder(key string, dir firestore.Direction) func(firestore.Query) firestore.Query {
	return func(q firestore.Query) firestore.Query {
		return q.OrderBy(key, dir)
	}
}

func WithOffset(offset int) func(firestore.Query) firestore.Query {
	return func(q firestore.Query) firestore.Query {
		return q.StartAt(offset)
	}
}

func WithLimit(limit int) func(firestore.Query) firestore.Query {
	return func(q firestore.Query) firestore.Query {
		return q.Limit(limit)
	}
}

type operator string

func IsZeroValue(val interface{}) bool {
	switch t := val.(type) {
	case int:
		return t == 0
	case float64:
		return t == 0.0
	case string:
		return t == ""
	case []any:
		return len(t) == 0
	case map[any]any:
		return len(t) == 0
	}
	return false
}
