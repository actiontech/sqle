package dashboard

import "context"

// RedDotDetector defines a generic function signature for checking red dots.
// It only takes context and userID to keep the interface clean.
type RedDotDetector func(ctx context.Context, userID string) (bool, error)

var detectors []RedDotDetector

// RegisterDetector allows different modules/versions to plug in their logic.
func RegisterDetector(d RedDotDetector) {
	detectors = append(detectors, d)
}

// GetDashboardRedDotV2 iterates through all registered detectors.
func GetDashboardRedDotV2(ctx context.Context, userID string) (bool, error) {
	for _, detect := range detectors {
		has, err := detect(ctx, userID)
		if err != nil {
			return false, err
		}
		if has {
			return true, nil
		}
	}
	return false, nil
}
