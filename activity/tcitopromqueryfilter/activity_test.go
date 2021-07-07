/*
* BSD 3-Clause License
* Copyright © 2020. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package tcitoprommetrics

import (
	"testing"

	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestPromMetrics(t *testing.T) {

	mf := mapper.NewFactory(resolve.GetBasicResolver())
	iCtx := test.NewActivityInitContext(nil, mf)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())

	input := &Input{
		Query: "a=b,c=d,e=f",
	}
	tc.SetInputObject(input)

	act.Eval(tc)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.NotEmpty(t, output.Filter)
	assert.Equal(t, "b", output.Filter["a"], "First object not retrieded properly")

}
