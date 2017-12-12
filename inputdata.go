package main

type inputData struct {
	AliceToBob *StealthSession `json:"alice2bob"`
	BobToAlice *StealthSession `json:"bob2alice"`
	Message    []byte          `json:"message"`
	PubKeys    []CurvePoint    `json:"ring"`
	Signatures []RingSignature `json:"signatures"`
}
