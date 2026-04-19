// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Package sqlitedriver is the Atlas fork's SQLite driver shim. It registers
// modernc.org/sqlite (a pure-Go SQLite) under both "sqlite" (modernc's native
// name) and "sqlite3" (the legacy name mattn/go-sqlite3 used and every Atlas
// call site still passes to sql.Open), and it translates mattn-style DSN
// query params (e.g. "_fk=1", "_busy_timeout=5000") to the PRAGMA statements
// modernc needs to hear, so pre-existing Atlas URLs keep working after the
// swap away from the CGO driver.
package sqlitedriver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"
	"strings"

	sqlite "modernc.org/sqlite"
)

// mattnPragma maps mattn/go-sqlite3's underscore-prefixed DSN query params
// (and their aliases) to the PRAGMA name modernc should issue on connection
// open. The value from the DSN is passed through unchanged — SQLite PRAGMAs
// accept mattn's full value vocabulary (0/1, on/off, true/false, etc.).
var mattnPragma = map[string]string{
	"_foreign_keys":             "foreign_keys",
	"_fk":                       "foreign_keys",
	"_defer_foreign_keys":       "defer_foreign_keys",
	"_defer_fk":                 "defer_foreign_keys",
	"_busy_timeout":             "busy_timeout",
	"_timeout":                  "busy_timeout",
	"_journal_mode":             "journal_mode",
	"_journal":                  "journal_mode",
	"_synchronous":              "synchronous",
	"_sync":                     "synchronous",
	"_locking_mode":             "locking_mode",
	"_locking":                  "locking_mode",
	"_case_sensitive_like":      "case_sensitive_like",
	"_cslike":                   "case_sensitive_like",
	"_recursive_triggers":       "recursive_triggers",
	"_rt":                       "recursive_triggers",
	"_auto_vacuum":              "auto_vacuum",
	"_vacuum":                   "auto_vacuum",
	"_secure_delete":            "secure_delete",
	"_query_only":               "query_only",
	"_writable_schema":          "writable_schema",
	"_ignore_check_constraints": "ignore_check_constraints",
}

func init() {
	// modernc's package-level RegisterConnectionHook attaches to its internal
	// singleton Driver (registered as "sqlite"). We register our own Driver
	// instance under "sqlite3", so the hook has to be attached per-instance.
	drv := &sqlite.Driver{}
	drv.RegisterConnectionHook(applyMattnDSNParams)
	sql.Register("sqlite3", drv)
}

// applyMattnDSNParams issues a PRAGMA statement for each mattn-style query
// param found in dsn. Unknown params are ignored (modernc will have already
// surfaced anything it cares about). Non-PRAGMA mattn params (_loc, _mutex,
// _txlock) aren't translated — no Atlas code uses them.
func applyMattnDSNParams(conn sqlite.ExecQuerierContext, dsn string) error {
	qIdx := strings.Index(dsn, "?")
	if qIdx < 0 {
		return nil
	}
	params, err := url.ParseQuery(dsn[qIdx+1:])
	if err != nil {
		return nil
	}
	ctx := context.Background()
	for k, vs := range params {
		pragma, ok := mattnPragma[k]
		if !ok || len(vs) == 0 {
			continue
		}
		stmt := fmt.Sprintf("PRAGMA %s = %s", pragma, vs[len(vs)-1])
		if _, err := conn.ExecContext(ctx, stmt, []driver.NamedValue{}); err != nil {
			return fmt.Errorf("sqlitedriver: apply %s: %w", pragma, err)
		}
	}
	return nil
}
