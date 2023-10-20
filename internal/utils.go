package internal

import "regexp"

// The default Compute Engine service account is always PROJECT_NUMBER-compute@developer.gserviceaccount.com
func IsDefaultComputeEngineSA(name string) bool {
	regionPattern := regexp.MustCompile(`^[0-9]+-compute@developer\.gserviceaccount\.com$`)
	return regionPattern.MatchString(name)
}

// Check if a variable name or value may be sensitive
// name: check the presence of a keyword in the name (case insensitive)
// value: check the presence of a private key pattern
func MayBeSensitive(name string, value string) bool {
	privateKeyPattern := regexp.MustCompile(`-----BEGIN (RSA|EC|DSA|GPP|OPENSSH) PRIVATE KEY-----`)
	namePattern := regexp.MustCompile(`(key|secret|password|token|private|credential|auth|certificate|cert|pem|ssl|tls|ssh|rsa|ecdsa|dsa|gpp)(?i)`)

	return namePattern.MatchString(name) || privateKeyPattern.MatchString(value)
}
