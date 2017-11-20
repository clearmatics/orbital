package main

type inputData struct {
	AliceToBob *StealthSession  `json:"alice2bob"`
	BobToAlice *StealthSession  `json:"bob2alice"`
	PubKeys    []CurvePoint    `json:"pubkeys"`
	Signatures []RingSignature `json:"signatures"`
}
