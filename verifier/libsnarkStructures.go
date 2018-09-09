package verifier

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func skipSpaces(r *bufio.Reader) {
	for {
		rune, _, err := r.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				return
			}
		}
		c := string(rune)
		if c != " " {
			r.UnreadRune()
			return
		}
	}
}

func consumeNewLine(r *bufio.Reader) {
	rune, _, err := r.ReadRune()
	if err != nil {
		if err.Error() == "EOF" {
			return
		}
	}
	c := string(rune)
	if c != "\n" {
		r.UnreadRune()
		return
	}
}

func ReadInt(r *bufio.Reader) (uint64, error) {
	x := ""
	for {
		rune, _, err := r.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return 0, err
		}
		c := string(rune)
		if c != " " && c != "\n" {
			x = x + c
		} else {
			r.UnreadRune()
			break
		}
	}
	i, err := strconv.ParseUint(x, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func ReadBigIntAsString(r *bufio.Reader) (string, error) {
	x := ""
	for {
		rune, _, err := r.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return "", err
		}
		c := string(rune)
		if c != " " && c != "\n" {
			x = x + c
		} else {
			r.UnreadRune()
			break
		}
	}
	return x, nil
}

func ReadG1(r *bufio.Reader) (*G1, error) {
	rune, _, err := r.ReadRune()
	if err != nil {
		return nil, err
	}
	str := string(rune)
	if str == "1" {
		return new(G1), nil
	} else if str != "0" {
		return nil, errors.New("Invalid encoding format")
	}

	skipSpaces(r)
	xCoord, err := ReadBigIntAsString(r)
	if err != nil {
		return nil, err
	}
	skipSpaces(r)
	yCoord, err := ReadBigIntAsString(r)
	if err != nil {
		return nil, err
	}
	consumeNewLine(r)
	if err != nil {
		return nil, err
	}
	fmt.Println("X = " + xCoord)
	fmt.Println("Y = " + yCoord)
	newPoint, err := NewG1FromStrings(xCoord, yCoord, 10)
	if err != nil {
		return nil, err
	}
	return newPoint, nil
}

func ReadG2(r *bufio.Reader) (*G2, error) {
	rune, _, err := r.ReadRune()
	if err != nil {
		return nil, err
	}
	str := string(rune)
	if str == "1" {
		return new(G2), nil
	} else if str != "0" {
		return nil, errors.New("Invalid encoding format")
	}
	skipSpaces(r)
	a, err := ReadBigIntAsString(r)
	if err != nil {
		return nil, err
	}
	skipSpaces(r)
	b, err := ReadBigIntAsString(r)
	if err != nil {
		return nil, err
	}
	skipSpaces(r)
	c, err := ReadBigIntAsString(r)
	if err != nil {
		return nil, err
	}
	skipSpaces(r)
	d, err := ReadBigIntAsString(r)
	if err != nil {
		return nil, err
	}
	consumeNewLine(r)
	fmt.Println("A = " + a)
	fmt.Println("B = " + b)
	fmt.Println("C = " + c)
	fmt.Println("D = " + d)
	newPoint, err := NewG2FromStrings([2]string{a, b}, [2]string{c, d}, 10)
	// if err != nil {
	// 	return nil, err
	// }
	return newPoint, nil
}

type SparseVector struct {
	first *G1
	rest  []*G1
}

type LibsnarkVerifyingKey struct {
	A          *G2
	B          *G1
	C          *G2
	Gamma      *G2
	GammaBeta1 *G1
	GammaBeta2 *G2
	Z          *G2
	IC         *SparseVector
}

// func ParseFromBytes(b []byte) {
// 	r := bytes.NewReader(b)
// 	reader := bufio.NewReader(r)
// 	for {
// 		_, err := ReadG1(reader)
// 		if err != nil {
// 			return
// 		}
// 	}
// }

func (sv *SparseVector) ParseFromReader(r *bufio.Reader) error {
	first, err := ReadG1(r)
	if err != nil {
		return err
	}
	ReadInt(r)
	consumeNewLine(r)
	indSize, err := ReadInt(r)
	if err != nil {
		return err
	}
	consumeNewLine(r)
	for i := 0; i < int(indSize); i++ {
		_, err := ReadInt(r)
		if err != nil {
			return err
		}
		consumeNewLine(r)
	}
	valuesSize, err := ReadInt(r)
	if err != nil {
		return err
	}
	consumeNewLine(r)
	points := make([]*G1, valuesSize)
	fmt.Println("Verification has inputs = " + strconv.FormatUint(valuesSize, 10))
	for i := 0; i < int(valuesSize); i++ {
		point, err := ReadG1(r)
		if err != nil {
			return err
		}
		consumeNewLine(r)
		points[i] = point
	}
	sv.first = first
	sv.rest = points
	return nil
}

func (vk *LibsnarkVerifyingKey) ParseFromFile(filename string) error {
	r, err := os.Open(filename)
	defer r.Close()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(r)
	return vk.ParseFromReader(reader)
}

func (vk *LibsnarkVerifyingKey) ParseFromReader(reader *bufio.Reader) error {
	A, err := ReadG2(reader)
	if err != nil {
		return err
	}
	consumeNewLine(reader)

	B, err := ReadG1(reader)
	if err != nil {
		return err
	}
	consumeNewLine(reader)

	C, err := ReadG2(reader)
	if err != nil {
		return err
	}
	consumeNewLine(reader)

	Gamma, err := ReadG2(reader)
	if err != nil {
		return err
	}
	consumeNewLine(reader)

	GammaBeta1, err := ReadG1(reader)
	if err != nil {
		return err
	}
	consumeNewLine(reader)

	GammaBeta2, err := ReadG2(reader)
	if err != nil {
		return err
	}
	consumeNewLine(reader)

	Z, err := ReadG2(reader)
	if err != nil {
		return err
	}
	consumeNewLine(reader)

	ic := new(SparseVector)
	err = ic.ParseFromReader(reader)
	if err != nil {
		return nil
	}

	vk.A = A
	vk.B = B
	vk.C = C
	vk.Gamma = Gamma
	vk.GammaBeta1 = GammaBeta1
	vk.GammaBeta2 = GammaBeta2
	vk.Z = Z
	vk.IC = ic
	return nil
}
