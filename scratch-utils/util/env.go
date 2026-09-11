// Copyright (c) 2026 Tigera, Inc. All rights reserved.

package util

import "os"

// EnvOr returns the value of key, or def when key is unset OR set to the empty
// string. Empty-means-absent is deliberate: a workflow that declares an env var
// it has no value for should get the default, not "".
//
// Where the difference matters -- LocalSecret and --put-env, which must tell a
// deliberately empty secret from a missing one -- use os.LookupEnv directly.
func EnvOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
