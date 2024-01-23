package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type UserType = int
type SubscriptionType = int

const (
	unknownUserType UserType = iota
	lead
	customer
	churned
	lostLead
)

const (
	unknownSubscriptionType SubscriptionType = iota
	basic
	premium
	exclusive
)

type UserAddRequest struct {
	UserType       UserType
	Email          string
	SubType        SubscriptionType
	PaymentDetails PaymentDetails
}

type UserModifyRequest struct {
	ID             string
	UserType       UserType
	Email          string
	SubType        SubscriptionType
	PaymentDetails PaymentDetails
}

type User struct {
	ID             string
	PaymentDetails PaymentDetails
}

type PaymentDetails struct {
	stripeTokenID string
}

type LeadRequest struct {
	email string
}

type UserManager interface {
	AddUser(ctx context.Context, request UserAddRequest) (User, error)
	ModifyUser(ctx context.Context, request UserModifyRequest) (User, error)
}

type UserHandler interface {
	IsUserSubscriptionActive(ctx context.Context, userID string) bool
}

type UserActiveResponse struct {
	IsActive bool
}

func router(u UserHandler) {
	m := mux.NewRouter()

	m.HandleFunc("/user/{id}/active", func(w http.ResponseWriter, r *http.Request) {
		// check auth, etc

		uID := mux.Vars(r)["userID"]
		if uID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		isActive := u.IsUserSubscriptionActive(r.Context(), uID)

		b, err := json.Marshal(UserActiveResponse{IsActive: isActive})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(b)
	}).Methods(http.MethodGet)
}
