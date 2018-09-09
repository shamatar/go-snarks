package verifier

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// 0 20933170699147567579033322641963108696395791188513690084963135910641532928274 2442431204756310761602153127686318554574960794578976757443433976540189788252 0 4887099059762129640302984526748584436465700498472647690070911027557383648471 1537761002038265777208458665898474687934111788158087930192782735780488615370
// 0 16573055444656122843537507139515729097968881044881963088535952613432732109576 15823969213404520820670755358760897637430302424438918422654019939237973220852 19043350362690960802787632939519229156536345342250133376920132587983219143570 5672529627952885536591336430691382520056316077527491874420269967257820619040 0 14461792803667918890803887317933444090830840712748147371492463942958574209725 13071829618513332891552585680079909704651276111106970019151699080596107269576
// 0 4553363573730177428827227824438616079346298937519779886664801330607001971471 21456699004685635406151940504311915468194491877421700836100090035559356021754 0 5750931461328290137207375961192256416814227690413125433502291640571740409967 8930892068663441607094145146857314111761878883131063200134413359379416014906
// 0 5627193368192982618917091473196986665093145528450803546070128784390583134960 7833536160441415171125072713335045341006722851990626699896991335087744993720
// 0 13096706332157259156060110671651801913587695020536048784149956198170577046492 14999138676268127536956531532613144042237861207626366537509655144817612266843

func ParseProofFromFile(filename string) (*Proof, error) {
	fullFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fullContent := string(fullFile)
	return ParseProofFromString(fullContent)
}

func ParseProofFromString(content string) (*Proof, error) {
	fullContent := strings.Replace(content, "\n", " ", -1)
	fullContent = strings.Replace(fullContent, "\t", " ", -1)

	elements := strings.Split(fullContent, " 0 ")
	firstElem := elements[0]
	firstElemSplit := strings.Split(firstElem, " ")
	elements[0] = firstElemSplit[1] + " " + firstElemSplit[2]
	// if len(elements) != 8 {
	// 	return nil, errors.New("Invalid number of elements")
	// }
	i := 0
	proof := &Proof{}
	for _, elem := range elements {
		cleaned := strings.Trim(elem, " \n\t")
		components := strings.Split(cleaned, " ")
		switch i {
		case 0:
			A, err := NewG1FromStrings(components[0], components[1], 10)
			if err != nil {
				return nil, err
			}
			proof.A = A
		case 1:
			Ap, err := NewG1FromStrings(components[0], components[1], 10)
			if err != nil {
				return nil, err
			}
			proof.Ap = Ap
		case 2:
			fmt.Println(elem)
			fmt.Println("Cleaned")
			fmt.Println(cleaned)
			fmt.Println("Components")
			fmt.Println(components)
			B, err := NewG2FromStrings([2]string{components[0], components[1]}, [2]string{components[2], components[3]}, 10)
			if err != nil {
				i++
				continue
				return nil, err
			}
			proof.B = B
		case 3:
			Bp, err := NewG1FromStrings(components[0], components[1], 10)
			if err != nil {
				return nil, err
			}
			proof.Bp = Bp
		case 4:
			C, err := NewG1FromStrings(components[0], components[1], 10)
			if err != nil {
				return nil, err
			}
			proof.C = C
		case 5:
			Cp, err := NewG1FromStrings(components[0], components[1], 10)
			if err != nil {
				return nil, err
			}
			proof.Cp = Cp
		case 6:
			H, err := NewG1FromStrings(components[0], components[1], 10)
			if err != nil {
				return nil, err
			}
			proof.H = H
		case 7:
			K, err := NewG1FromStrings(components[0], components[1], 10)
			if err != nil {
				return nil, err
			}
			proof.K = K
		default:
			return nil, errors.New("Invalid number of elements")
		}
		i++
	}
	if i != 8 {
		return nil, errors.New("Invalid number of elements")
	}
	return proof, nil
}

func ParseProofInLibsnarkFormat(filename string) (*Proof, error) {
	// file, err := os.Open(filename) // For read access.
	// if err != nil {
	// 	return nil, err
	// }
	// defer file.Close()
	fullFile, err := ioutil.ReadFile(filename)
	fullContent := string(fullFile)
	fullContent = strings.Replace(fullContent, "\n", " ", -1)
	fullContent = strings.Replace(fullContent, "\t", " ", -1)
	if err != nil {
		return nil, err
	}
	elements := strings.Split(fullContent, " 0 ")
	firstElem := elements[0]
	firstElemSplit := strings.Split(firstElem, " ")
	elements[0] = firstElemSplit[1] + " " + firstElemSplit[2]
	// if len(elements) != 8 {
	// 	return nil, errors.New("Invalid number of elements")
	// }
	i := 0
	proof := &Proof{}
	for _, elem := range elements {
		cleaned := strings.Trim(elem, " \n\t")
		components := strings.Split(cleaned, " ")
		switch i {
		case 0:
			A, err := NewG1FromStrings(components[0], components[1], 10)
			if err != nil {
				return nil, err
			}
			proof.A = A
		case 1:
			Ap, err := NewG1FromStrings(components[0], components[1], 10)
			if err != nil {
				return nil, err
			}
			proof.Ap = Ap
		case 2:
			fmt.Println(elem)
			fmt.Println("Cleaned")
			fmt.Println(cleaned)
			fmt.Println("Components")
			fmt.Println(components)
			B, err := NewG2FromStrings([2]string{components[0], components[1]}, [2]string{components[2], components[3]}, 10)
			if err != nil {
				i++
				continue
				return nil, err
			}
			proof.B = B
		case 3:
			Bp, err := NewG1FromStrings(components[0], components[1], 10)
			if err != nil {
				return nil, err
			}
			proof.Bp = Bp
		case 4:
			C, err := NewG1FromStrings(components[0], components[1], 10)
			if err != nil {
				return nil, err
			}
			proof.C = C
		case 5:
			Cp, err := NewG1FromStrings(components[0], components[1], 10)
			if err != nil {
				return nil, err
			}
			proof.Cp = Cp
		case 6:
			H, err := NewG1FromStrings(components[0], components[1], 10)
			if err != nil {
				return nil, err
			}
			proof.H = H
		case 7:
			K, err := NewG1FromStrings(components[0], components[1], 10)
			if err != nil {
				return nil, err
			}
			proof.K = K
		default:
			return nil, errors.New("Invalid number of elements")
		}
		i++
	}
	if i != 8 {
		return nil, errors.New("Invalid number of elements")
	}
	return proof, nil
}
