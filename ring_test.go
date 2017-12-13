// Copyright (c) 2017 Clearmatics Technologies Ltd

// SPDX-License-Identifier: LGPL-3.0+

package main

import (
	"testing"
)

func generateRing(i int) Ring {
	ring := &Ring{}
	ring.Generate(i)
	return *ring
}

func TestRingGenerate(t *testing.T) {
	i := 1
	ring := generateRing(i)

	expected := i
	actual := len(ring.PubKeys)
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

	actual = len(ring.PrivKeys)
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}

}

func TestRingSignature(t *testing.T) {
	i := 1
	r := generateRing(i)
	message := []byte("foobarbaz")

	_, err := r.Signature(r.PrivKeys[0], message, 0)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkRingSignature(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		i := 1
		r := generateRing(i)
		message := []byte("foobarbaz")

		b.StartTimer()
		_, err := r.Signature(r.PrivKeys[0], message, 0)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestRingSignatures(t *testing.T) {
	i := 4
	r := generateRing(i)
	s := []byte("foobarbaz")

	sigs, err := r.Signatures(s)
	if err != nil {
		t.Fatal(err)
	}

	expected := i
	actual := len(sigs)
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func BenchmarkRingSignatures(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		i := 4
		r := generateRing(i)
		message := []byte("foobarbaz")

		b.StartTimer()
		_, err := r.Signatures(message)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestVerifySignature(t *testing.T) {
	i := 1
	r := generateRing(i)
	message := []byte("foobarbaz")

	sig, err := r.Signature(r.PrivKeys[0], message, 0)
	if err != nil {
		t.Fatal(err)
	}

	expected := true
	actual := r.VerifySignature(message, *sig)
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func BenchmarkVerifySignature(b *testing.B) {

	for n := 0; n < b.N; n++ {
		// Stop clock
		b.StopTimer()
		i := 1
		r := generateRing(i)
		message := []byte("foobarbaz")

		sig, err := r.Signature(r.PrivKeys[0], message, 0)
		if err != nil {
			b.Fatal(err)
		}

		expected := true

		// Start clock
		b.StartTimer()

		actual := r.VerifySignature(message, *sig)
		if actual != expected {
			b.Errorf("Expected %v but got %v", expected, actual)
		}
	}
}

func TestVerifySignatures(t *testing.T) {
	i := 4
	r := generateRing(i)
	message := []byte("foobarbaz")
	sigs, err := r.Signatures(message)
	if err != nil {
		t.Fatal(err)
	}

	for _, sig := range sigs {
		expected := true
		actual := r.VerifySignature(message, sig)
		if actual != expected {
			t.Errorf("Expected %v but got %v", expected, actual)
		}
	}
}

func BenchmarkVerifySignatures(b *testing.B) {

	for n := 0; n < b.N; n++ {
		// Stop clock
		b.StopTimer()
		i := 4
		r := generateRing(i)
		message := []byte("foobarbaz")
		sigs, err := r.Signatures(message)
		if err != nil {
			b.Fatal(err)
		}

		// Start clock
		b.StartTimer()
		for _, sig := range sigs {
			expected := true
			actual := r.VerifySignature(message, sig)
			if actual != expected {
				b.Errorf("Expected %v but got %v", expected, actual)
			}
		}
	}
}

func TestVerifySignatureBad(t *testing.T) {
	i := 1
	r := generateRing(i)
	message := []byte("foobarbaz")

	sig, err := r.Signature(r.PrivKeys[0], message, 0)
	if err != nil {
		t.Fatal(err)
	}
	expected := false
	actual := r.VerifySignature([]byte("badmessage"), *sig)
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestPubKeyIndex(t *testing.T) {
	i := 3
	r := generateRing(i)
	actual := r.PubKeyIndex(r.PubKeys[0])
	expected := 0
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestPubKeyIndexNotFound(t *testing.T) {
	i := 3
	r := generateRing(i)
	c := CurvePoint{}
	actual := r.PubKeyIndex(c)
	expected := -1
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}
