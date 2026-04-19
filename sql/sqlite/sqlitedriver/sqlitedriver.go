// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Package sqlitedriver registers modernc.org/sqlite under both "sqlite" (its
// default name) and "sqlite3" so existing callers of sql.Open("sqlite3", ...)
// keep working after the swap away from the CGO driver mattn/go-sqlite3.
package sqlitedriver

import (
	"database/sql"

	sqlite "modernc.org/sqlite"
)

func init() {
	sql.Register("sqlite3", &sqlite.Driver{})
}
