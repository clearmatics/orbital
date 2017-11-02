package main

type inputData struct {
	PubKeys    []CurvePoint    `json:"pubkeys"`
	Signatures []RingSignature `json:"signatures"`
}
