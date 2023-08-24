//go:build !enterprise
// +build !enterprise

package v1

import (
	"context"
)

func (h BeforeArvhiveNamespace) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
	return nil
}
func (h AfterDeleteNamespace) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
	return nil
}
func (h BeforeDeleteNamespace) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
	return nil
}
func (h AfterCreateNamespace) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
	return nil
}
