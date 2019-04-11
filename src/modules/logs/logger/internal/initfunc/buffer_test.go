// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package initfunc

import (
	"testing"

	"logs/assert"
	"logs/logger/writers"
)

func TestBuffer(t *testing.T) {
	a := assert.New(t)
	args := map[string]string{}

	w, err := Buffer(args)
	a.Error(err).Nil(w)

	args["size"] = "5"
	w, err = Buffer(args)
	a.NotError(err).NotNil(w)
	_, ok := w.(*writers.Buffer)
	a.True(ok)

	// 无法解析的 size 参数
	args["size"] = "5l"
	w, err = Buffer(args)
	a.Error(err).Nil(w)
}
