// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package sqlite

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseTypeRejectsIntOverflow(t *testing.T) {
	tooLarge := strconv.Itoa(math.MaxInt) + "1"

	_, err := ParseType("varchar(" + tooLarge + ")")
	require.EqualError(t, err, `parse size "`+tooLarge+`"`)

	_, err = ParseType("decimal(" + tooLarge + ",1)")
	require.EqualError(t, err, `parse precision "`+tooLarge+`"`)

	_, err = ParseType("decimal(10," + tooLarge + ")")
	require.EqualError(t, err, `parse scale "10"`)
}
