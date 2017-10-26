// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main

import (
	"bytes"
	"fmt"
	"math/big"
)

type RingSignature struct {
	Tau    CurvePoint `json:"tau"`
	Ctlist []*big.Int `json:"ctlist"`
}

func (s RingSignature) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("tau: %s\n", s.Tau))
	buffer.WriteString(fmt.Sprintf("ctlist: [\n"))
	for i := 0; i < len(s.Ctlist); i++ {
		buffer.WriteString(fmt.Sprintf("\t%s\n", s.Ctlist[i]))
	}
	buffer.WriteString(fmt.Sprintf("]\n"))

	return buffer.String()
}
