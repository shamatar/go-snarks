package verifier

import (
	"bytes"
	"errors"
	"math/big"
	"math/rand"
	"strings"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
)

// verify a Pinocchio snark
// 1 pairing equation, 9 group elements of a proof

// VerifyingKey lists precomputed parameters
// and points for multiplication of witness
type VerifyingKey struct {
	A          *G2
	B          *G1
	C          *G2
	gamma      *G2
	gammaBeta1 *G1
	gammaBeta2 *G2
	Z          *G2
	IC         []*G1 // set of G1 point to multiply a witness on
}

// Proof lists a number of aparameters
// and arbitrary number of inputs
// note that where Verifying key has G2 group, Proof has G1 group
// and vice versa
type Proof struct {
	A  *G1
	Ap *G1
	B  *G2
	Bp *G1
	C  *G1
	Cp *G1
	K  *G1
	H  *G1
}

// Witness is type alias for a set of scalars
// Each scalar is in fact integer mod q, where q is BN256 G1 group order
type Witness = []*big.Int

func padBigInt(bn *big.Int) ([]byte, error) {
	marshalled := bn.Bytes()
	length := len(marshalled)
	if length > 32 {
		return nil, errors.New("Integer is too large")
	} else if length == 32 {
		return marshalled, nil
	} else {
		padding := make([]byte, 32-length)
		return append(padding, marshalled...), nil
	}
}

func base16bi(s string) (*big.Int, error) {
	val := strings.TrimPrefix(s, "0x")
	n, success := new(big.Int).SetString(val, 16)
	if !success {
		return nil, errors.New("Can not parse string to number")
	}
	return n, nil
}

func base10bi(s string) (*big.Int, error) {
	n, success := new(big.Int).SetString(s, 10)
	if !success {
		return nil, errors.New("Can not parse string to number")
	}
	return n, nil
}

func NewG1(xCoord, yCoord *big.Int) (*G1, error) {
	point := new(G1)
	xMarshalled, _ := padBigInt(xCoord)
	yMarshalled, _ := padBigInt(yCoord)
	_, err := point.Unmarshal(append(xMarshalled, yMarshalled...))
	if err != nil {
		return nil, err
	}
	return point, nil
}

func NewG1FromStrings(xCoord, yCoord string, radix int) (*G1, error) {
	if radix == 10 {
		x, err := base10bi(xCoord)
		if err != nil {
			return nil, err
		}
		y, err := base10bi(yCoord)
		if err != nil {
			return nil, err
		}
		return NewG1(x, y)
	} else if radix == 16 {
		x, err := base16bi(xCoord)
		if err != nil {
			return nil, err
		}
		y, err := base16bi(yCoord)
		if err != nil {
			return nil, err
		}
		return NewG1(x, y)
	}
	return nil, errors.New("Work with radix 16 and 10 only")
}

func NewG2(aCoords, bCoords [2]*big.Int) (*G2, error) {
	point := new(G2)
	aMarshalled, _ := padBigInt(aCoords[0])
	bMarshalled, _ := padBigInt(aCoords[1])
	cMarshalled, _ := padBigInt(bCoords[0])
	dMarshalled, _ := padBigInt(bCoords[1])
	marshalled := append(aMarshalled, bMarshalled...)
	marshalled = append(marshalled, cMarshalled...)
	marshalled = append(marshalled, dMarshalled...)
	_, err := point.Unmarshal(marshalled)
	if err != nil {
		return nil, err
	}
	return point, nil
}

func NewG2FromStrings(aCoords, bCoords [2]string, radix int) (*G2, error) {
	if radix == 10 {
		a, err := base10bi(aCoords[0])
		if err != nil {
			return nil, err
		}
		b, err := base10bi(aCoords[1])
		if err != nil {
			return nil, err
		}
		c, err := base10bi(bCoords[0])
		if err != nil {
			return nil, err
		}
		d, err := base10bi(bCoords[1])
		if err != nil {
			return nil, err
		}
		return NewG2([2]*big.Int{a, b}, [2]*big.Int{c, d})
	} else if radix == 16 {
		a, err := base16bi(aCoords[0])
		if err != nil {
			return nil, err
		}
		b, err := base16bi(aCoords[1])
		if err != nil {
			return nil, err
		}
		c, err := base16bi(bCoords[0])
		if err != nil {
			return nil, err
		}
		d, err := base16bi(bCoords[1])
		if err != nil {
			return nil, err
		}
		return NewG2([2]*big.Int{a, b}, [2]*big.Int{c, d})
	}
	return nil, errors.New("Work with radix 16 and 10 only")
}

// GetG1Base gets G1 generator
func GetG1Base() *G1 {
	generator, _ := NewG1FromStrings("1", "2", 10)
	return generator
}

// GetG2Base gets G2 generator
func GetG2Base() *G2 {
	aCoords := [2]string{
		"11559732032986387107991004021392285783925812861821192530917403151452391805634",
		"10857046999023057135944570762232829481370756359578518086990519993285655852781",
	}
	bCoords := [2]string{
		"4082367875863433681332203403145435568316851327593401208105741076214120093531",
		"8495653923123431417604973247489272438418190587263600148770280649306958101930",
	}
	generator, _ := NewG2FromStrings(aCoords, bCoords, 10)
	return generator
}

// naiveSplitVerification computes A LOT of pairing
// follows the ZoKrates logic for verification in smart-contracts
// where randomness is not available
func naiveSplitVerification(witness Witness, proof *Proof, vk *VerifyingKey) error {
	if len(witness)+1 != len(vk.IC) {
		return errors.New("Invalid length of the witness")
	}
	G2Base := GetG2Base()
	witnessAccululator := new(G1).Set(vk.IC[0])
	temp := new(G1)
	for i, w := range witness {
		temp.ScalarMult(vk.IC[i+1], w)
		witnessAccululator.Add(temp, witnessAccululator)
	}
	// e(proof.A, vk.A) == e(-proof.Ap, G2)
	success := PairingCheck([]*G1{proof.A, temp.Neg(proof.Ap)}, []*G2{vk.A, G2Base})
	if !success {
		return errors.New("First pairing equation has failed")
	}
	// e(vk.B, proof.B) + e(-proof.Bp, G2)
	success = PairingCheck([]*G1{vk.B, temp.Neg(proof.Bp)}, []*G2{proof.B, G2Base})
	if !success {
		return errors.New("Second pairing equation has failed")
	}
	// e(proof.C, vk.C) + e(-proof.Cp, G2)
	success = PairingCheck([]*G1{proof.C, temp.Neg(proof.Cp)}, []*G2{vk.C, G2Base})
	if !success {
		return errors.New("Third pairing equation has failed")
	}

	// e(proof.K, vk.gamma) + e(- witnessAccumulator - proof.A - proof.C, vk.gammaBeta2) + e(-vk.gammaBeta1, proof.B)
	t := new(G1)
	t.Add(witnessAccululator, proof.A)
	t.Add(t, proof.C)
	t.Neg(t)

	success = PairingCheck([]*G1{proof.K, t, temp.Neg(vk.gammaBeta1)}, []*G2{vk.gamma, vk.gammaBeta2, proof.B})
	if !success {
		return errors.New("Fourth pairing equation has failed")
	}

	// e(witnessAccumulator + proof.A, proof.B) + e(- proof.H, vk.Z) + e(-proof.C, G2)
	u := new(G1)
	u.Add(witnessAccululator, proof.A)

	success = PairingCheck([]*G1{u, temp.Neg(proof.H), new(G1).Neg(proof.C)}, []*G2{proof.B, vk.Z, G2Base})
	if !success {
		return errors.New("Fifth pairing equation has failed")
	}
	return nil
}

// agregatedVerification takes some entropy and tried to run ONE pairing check
// if sum of pairing with arbitrary coefficients holds than most likely each of those holds
func agregatedVerification(witness Witness, proof *Proof, vk *VerifyingKey) error {
	if len(witness)+1 != len(vk.IC) {
		return errors.New("Invalid length of the witness")
	}
	G2Base := GetG2Base()
	witnessAccululator := new(G1).Set(vk.IC[0])
	temp := new(G1)

	for i, w := range witness {
		temp.ScalarMult(vk.IC[i+1], w)
		witnessAccululator.Add(temp, witnessAccululator)
	}
	// grab some entopy for pairing checks
	entropy := make([]*big.Int, 5)
	for i := range entropy {
		slice := make([]byte, 32)
		rand.Read(slice)
		entropy[i] = new(big.Int).SetBytes(slice)
	}

	pairings := make([]*GT, 5)

	// e(proof.A, vk.A) + e(-proof.Ap, G2)
	pair := bn256.Miller(proof.A, vk.A)
	pair = pair.Add(pair, bn256.Miller(temp.Neg(proof.Ap), G2Base))
	pairings[0] = pair
	// emptyGT := new(GT)
	// emptyGT.Set(pair)
	// emptyGT.Finalize()
	// if bytes.Compare(emptyGT.Marshal(), IdentityBytes) != 0 {
	// 	return errors.New("Pairing check has failed")
	// }

	// e(vk.B, proof.B) + e(-proof.Bp, G2)
	pair = bn256.Miller(vk.B, proof.B)
	pair = pair.Add(pair, bn256.Miller(temp.Neg(proof.Bp), G2Base))
	pairings[1] = pair
	// emptyGT.Set(pair)
	// emptyGT.Finalize()
	// if bytes.Compare(emptyGT.Marshal(), IdentityBytes) != 0 {
	// 	return errors.New("Pairing check has failed")
	// }

	// e(proof.C, vk.C) + e(-proof.Cp, G2)
	pair = bn256.Miller(proof.C, vk.C)
	pair = pair.Add(pair, bn256.Miller(temp.Neg(proof.Cp), G2Base))
	pairings[2] = pair
	// emptyGT.Set(pair)
	// emptyGT.Finalize()
	// if bytes.Compare(emptyGT.Marshal(), IdentityBytes) != 0 {
	// 	return errors.New("Pairing check has failed")
	// }

	// e(proof.K, vk.gamma) + e(- witnessAccumulator - proof.A - proof.C, vk.gammaBeta2) + e(-vk.gammaBeta1, proof.B)
	t := new(G1)
	t.Add(witnessAccululator, proof.A)
	t.Add(t, proof.C)
	t.Neg(t)
	pair = bn256.Miller(proof.K, vk.gamma)
	pair = pair.Add(pair, bn256.Miller(t, vk.gammaBeta2))
	pair = pair.Add(pair, bn256.Miller(temp.Neg(vk.gammaBeta1), proof.B))
	pairings[3] = pair
	// emptyGT.Set(pair)
	// emptyGT.Finalize()
	// if bytes.Compare(emptyGT.Marshal(), IdentityBytes) != 0 {
	// 	return errors.New("Pairing check has failed")
	// }

	// e(witnessAccumulator + proof.A, proof.B) + e(- proof.H, vk.Z) + e(-proof.C, G2)
	u := new(G1)
	u.Add(witnessAccululator, proof.A)

	pair = bn256.Miller(u, proof.B)
	pair = pair.Add(pair, bn256.Miller(temp.Neg(proof.H), vk.Z))
	pair = pair.Add(pair, bn256.Miller(new(G1).Neg(proof.C), G2Base))
	pairings[4] = pair
	// emptyGT.Set(pair)
	// emptyGT.Finalize()
	// if bytes.Compare(emptyGT.Marshal(), IdentityBytes) != 0 {
	// 	return errors.New("Pairing check has failed")
	// }

	linearCombination := new(GT).Set(pairings[0])
	for i := 1; i < len(pairings); i++ {
		linearCombination.Add(linearCombination, pairings[i].ScalarMult(pairings[i], entropy[i]))
		// emptyGT.Set(linearCombination)
		// emptyGT.Finalize()
		// if bytes.Compare(emptyGT.Marshal(), IdentityBytes) != 0 {
		// 	return errors.New("Pairing check has failed")
		// }
	}
	linearCombination.Finalize()
	if bytes.Compare(linearCombination.Marshal(), IdentityBytes) != 0 {
		return errors.New("Pairing check has failed")
	}
	return nil
}
