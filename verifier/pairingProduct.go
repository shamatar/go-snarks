package verifier

// uses only fast Cloudflare implementation

// TODO use pools for reduced GC

import (
	"errors"
	"math/big"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

// G1 is an abstract cyclic group. The zero value is suitable for use as the
// output of an operation, but cannot be used as an input.
type G1 = bn256.G1

// G2 is an abstract cyclic group. The zero value is suitable for use as the
// output of an operation, but cannot be used as an input.
type G2 = bn256.G2

// PairingCheck calculates the Optimal Ate pairing for a set of points.
func PairingCheck(a []*G1, b []*G2) bool {
	return bn256.PairingCheck(a, b)
}

// AddG1 parses raw data and does a G1 (small group) addition
// Expects 64 + 64 bytes of data
func AddG1(data []byte) ([]byte, error) {
	if len(data) != 128 {
		return nil, errors.New("Data should be 128 bytes long")
	}

	a := new(G1)
	_, err := a.Unmarshal(data[:64])

	if err != nil {
		return nil, err
	}

	b := new(G1)
	_, err = b.Unmarshal(data[64:])

	if err != nil {
		return nil, err
	}

	result := AddG1Parsed(a, b)
	return result.Marshal(), nil
}

// AddG1 parses raw data and does a G1 (small group) addition
// Expects 64 + 64 bytes of data
func AddG1Parsed(a, b *G1) *G1 {
	result := new(G1)
	result.Add(a, b)
	return result
}

// MulG1 parses raw data and does a G1 (small group) multiplication
// Expected 32 + 64 bytes of data
func MulG1(data []byte) ([]byte, error) {
	if len(data) != 96 {
		return nil, errors.New("Data should be 96 bytes long")
	}
	point := new(G1)
	_, err := point.Unmarshal(data[:64])
	if err != nil {
		return nil, err
	}
	scalar := new(big.Int).SetBytes(data[64:])
	result := MulG1Parsed(scalar, point)
	return result.Marshal(), nil
}

// MulG1Parsed parses does a G1 (small group) multiplication
func MulG1Parsed(scalar *big.Int, point *G1) *G1 {
	resultPoint := new(G1)
	resultPoint.ScalarMult(point, scalar)
	return resultPoint
}
