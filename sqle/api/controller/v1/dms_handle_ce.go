//go:build !enterprise
// +build !enterprise

package v1

import (
	"context"
)

func (h BeforeArchiveProject) Handle(ctx context.Context, currentUserId string, dataResourceId string) error {
	return nil
}
func (h AfterDeleteProject) Handle(ctx context.Context, currentUserId string, dataResourceId string) error {
	return nil
}
func (h BeforeDeleteProject) Handle(ctx context.Context, currentUserId string, dataResourceId string) error {
	return nil
}
func (h AfterCreateProject) Handle(ctx context.Context, currentUserId string, dataResourceId string) error {
	return nil
}
