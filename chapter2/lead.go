package main

import "context"

// Updated code
type Lead struct {
	id string
}

type LeadCreator interface {
	CreateLead(ctx context.Context, request LeadRequest) (Lead, error)
}
