// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package schema_test

import (
	"testing"

	"ariga.io/atlas/sql/schema"

	"github.com/stretchr/testify/require"
)

// testObject is a named realm-level object for testing object filtering.
type testObject struct {
	schema.Object
	typ  string
	name string
}

func (o *testObject) SpecType() string { return o.typ }
func (o *testObject) SpecName() string { return o.name }

func TestExcludeRealm(t *testing.T) {
	t.Run("NoPatterns", func(t *testing.T) {
		r := schema.NewRealm(schema.New("public"), schema.New("internal"))
		got, err := schema.ExcludeRealm(r, nil)
		require.NoError(t, err)
		require.Len(t, got.Schemas, 2)
	})

	t.Run("ExcludeSchemaByName", func(t *testing.T) {
		r := schema.NewRealm(schema.New("public"), schema.New("internal"))
		got, err := schema.ExcludeRealm(r, []string{"internal"})
		require.NoError(t, err)
		require.Len(t, got.Schemas, 1)
		require.Equal(t, "public", got.Schemas[0].Name)
	})

	t.Run("ExcludeSchemaGlob", func(t *testing.T) {
		r := schema.NewRealm(schema.New("public"), schema.New("internal"), schema.New("other"))
		got, err := schema.ExcludeRealm(r, []string{"*"})
		require.NoError(t, err)
		require.Empty(t, got.Schemas)
	})

	t.Run("ExcludeTableInSchema", func(t *testing.T) {
		r := schema.NewRealm(
			schema.New("public").AddTables(
				schema.NewTable("users"),
				schema.NewTable("pets"),
			),
		)
		got, err := schema.ExcludeRealm(r, []string{"public.users"})
		require.NoError(t, err)
		require.Len(t, got.Schemas, 1)
		require.Len(t, got.Schemas[0].Tables, 1)
		require.Equal(t, "pets", got.Schemas[0].Tables[0].Name)
	})

	t.Run("ExcludeTableGlob", func(t *testing.T) {
		r := schema.NewRealm(
			schema.New("public").AddTables(
				schema.NewTable("users"),
				schema.NewTable("user_roles"),
				schema.NewTable("pets"),
			),
		)
		got, err := schema.ExcludeRealm(r, []string{"public.user*"})
		require.NoError(t, err)
		require.Len(t, got.Schemas[0].Tables, 1)
		require.Equal(t, "pets", got.Schemas[0].Tables[0].Name)
	})

	t.Run("ExcludeColumnInTable", func(t *testing.T) {
		r := schema.NewRealm(
			schema.New("public").AddTables(
				schema.NewTable("users").AddColumns(
					schema.NewStringColumn("name", "text"),
					schema.NewStringColumn("password", "text"),
				),
			),
		)
		got, err := schema.ExcludeRealm(r, []string{"public.users.password"})
		require.NoError(t, err)
		require.Len(t, got.Schemas[0].Tables[0].Columns, 1)
		require.Equal(t, "name", got.Schemas[0].Tables[0].Columns[0].Name)
	})

	t.Run("ExcludeObjectByName", func(t *testing.T) {
		r := schema.NewRealm(schema.New("public")).AddObjects(
			&testObject{typ: "extension", name: "pgcrypto"},
			&testObject{typ: "extension", name: "fuzzystrmatch"},
		)
		got, err := schema.ExcludeRealm(r, []string{"pgcrypto"})
		require.NoError(t, err)
		require.Len(t, got.Objects, 1)
		require.Equal(t, "fuzzystrmatch", got.Objects[0].(schema.SpecTypeNamer).SpecName())
	})

	t.Run("ExcludeObjectByTypeSelector", func(t *testing.T) {
		r := schema.NewRealm(schema.New("public")).AddObjects(
			&testObject{typ: "extension", name: "pgcrypto"},
			&testObject{typ: "other", name: "something"},
		)
		got, err := schema.ExcludeRealm(r, []string{"*[type=extension]"})
		require.NoError(t, err)
		require.Len(t, got.Objects, 1)
		require.Equal(t, "something", got.Objects[0].(schema.SpecTypeNamer).SpecName())
	})
}

func TestExcludeSchema(t *testing.T) {
	t.Run("NoPatterns", func(t *testing.T) {
		r := schema.NewRealm(schema.New("public").AddTables(schema.NewTable("users"), schema.NewTable("pets")))
		got, err := schema.ExcludeSchema(r.Schemas[0], nil)
		require.NoError(t, err)
		require.Len(t, got.Tables, 2)
	})

	t.Run("ExcludeTable", func(t *testing.T) {
		r := schema.NewRealm(schema.New("public").AddTables(schema.NewTable("users"), schema.NewTable("pets")))
		got, err := schema.ExcludeSchema(r.Schemas[0], []string{"users"})
		require.NoError(t, err)
		require.Len(t, got.Tables, 1)
		require.Equal(t, "pets", got.Tables[0].Name)
	})

	t.Run("MissingRealm", func(t *testing.T) {
		s := schema.New("public")
		_, err := schema.ExcludeSchema(s, []string{"users"})
		require.ErrorContains(t, err, "missing realm")
	})
}
