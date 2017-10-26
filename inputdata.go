package main

type InputData struct {
	PubKeys    []CurvePoint    `json:"pubkeys"`
	Signatures []RingSignature `json:"signatures"`
}
