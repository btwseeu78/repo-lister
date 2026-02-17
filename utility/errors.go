package utility

import (
	"fmt"
	"strings"
)

// HandleRegistryError interprets common container registry errors and returns
// user-friendly messages with actionable hints.
//
// operation describes what was being attempted (e.g. "listing tags", "pulling image", "pushing image", "copying image").
// target is the image or repository reference involved.
func HandleRegistryError(err error, operation string, target string) error {
	if err == nil {
		return nil
	}
	msg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msg, "unauthorized") ||
		strings.Contains(msg, "authentication required") ||
		strings.Contains(msg, "denied"):
		return fmt.Errorf("authentication failed while %s '%s'. Please check your credentials or Kubernetes secret (try using -s \"my-docker-secret\"). Original error: %v", operation, target, err)

	case strings.Contains(msg, "name_unknown") ||
		strings.Contains(msg, "manifest unknown") ||
		strings.Contains(msg, "404"):
		return fmt.Errorf("repository '%s' not found while %s. Please verify the image name is correct. Original error: %v", target, operation, err)

	case strings.Contains(msg, "dial tcp") ||
		strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "i/o timeout"):
		return fmt.Errorf("registry unreachable while %s '%s'. Please check your network connection and registry URL. Original error: %v", operation, target, err)

	case strings.Contains(msg, "forbidden"):
		return fmt.Errorf("access denied while %s '%s'. You don't have permission to perform this operation. Original error: %v", operation, target, err)

	default:
		return fmt.Errorf("failed %s '%s': %w", operation, target, err)
	}
}
