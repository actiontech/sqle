//go:build !enterprise
// +build !enterprise

package v1

import (
	"context"
)

func (h BeforeArvhiveProject) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
	return nil
}
func (h AfterDeleteProject) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
	return nil
}
func (h BeforeDeleteProject) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
	return nil
}
func (h AfterCreateProject) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
	return nil
}
