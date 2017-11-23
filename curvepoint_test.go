package main

import (
	"testing"
)


func TestCurvepointGenerate (t *testing.T) {
	Ap, As, err := generateKeyPair()
    if err != nil {
        t.Fatal(err)
    }
    if ! Ap.IsOnCurve() {
    	t.Fatal("Generated invalid public key")
    }
    if ! isValidSecretKey(As) {
    	t.Fatal("Generated invalid secret key")
    }
} 
