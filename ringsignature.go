// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main

import (
	"math/big"
	"encoding/json"
)

// A RingSignature is represented as a curve point and the signature data itself
type RingSignature struct {
	Tau    CurvePoint `json:"tau"`
	Ctlist []*big.Int `json:"ctlist"`
}

// MarshalJSON converts a RingSignature to a JSON representation
func (rs *RingSignature) MarshalJSON() ([]byte, error) {
	// XXX: go has no easy way of doing this without iterating
	ctlist := make([]*hexBig, len(rs.Ctlist))
	for i, v := range rs.Ctlist {
		ctlist[i] = (*hexBig)(v)
	}

	return json.Marshal(&struct {
		Tau    CurvePoint `json:"tau"`
		Ctlist []*hexBig `json:"ctlist"`
	}{
		Tau: rs.Tau,
		Ctlist: ctlist,
	})
}

// UnmarshalJSON converts a JSON representation to a RingSignature struct 
func (rs *RingSignature) UnmarshalJSON(data []byte) error {
	var aux struct {
		Tau    CurvePoint `json:"tau"`
		Ctlist []*hexBig `json:"ctlist"`
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	ctlist := make([]*big.Int, len(aux.Ctlist))
	for i, v := range aux.Ctlist {
		ctlist[i] = (*big.Int)(v)
	}
	rs.Ctlist = ctlist
	rs.Tau = aux.Tau
	return nil
}
