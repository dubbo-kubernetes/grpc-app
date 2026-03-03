package util

import (
	"fmt"
	"google.golang.org/grpc/status"
	"regexp"
	"strings"
)

func FirstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// formatGRPCError formats gRPC errors similar to grpcurl output
func FormatGRPCError(err error, index int32, total int32) string {
	if err == nil {
		return ""
	}

	// Extract gRPC status code and message using status package
	code := "Unknown"
	message := err.Error()

	// Try to extract gRPC status
	if st, ok := status.FromError(err); ok {
		code = st.Code().String()
		message = st.Message()
	} else {
		// Fallback: try to extract from error string
		if strings.Contains(message, "code = ") {
			// Extract code like "code = Unavailable"
			codeMatch := regexp.MustCompile(`code = (\w+)`).FindStringSubmatch(message)
			if len(codeMatch) > 1 {
				code = codeMatch[1]
			}
			// Extract message after "desc = "
			descMatch := regexp.MustCompile(`desc = "?([^"]+)"?`).FindStringSubmatch(message)
			if len(descMatch) > 1 {
				message = descMatch[1]
			} else {
				// If no desc, try to extract message after code
				parts := strings.SplitN(message, "desc = ", 2)
				if len(parts) > 1 {
					message = strings.Trim(parts[1], `"`)
				}
			}
		}
	}

	// Format similar to grpcurl (single line format)
	if total == 1 {
		return fmt.Sprintf("ERROR:\nCode: %s\nMessage: %s", code, message)
	}
	return fmt.Sprintf("[%d] Error: rpc error: code = %s desc = %s", index, code, message)
}
