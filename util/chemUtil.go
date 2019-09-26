package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"tcmio/models"
)

func SDFToDatabase(SDFile string) (int64, error) {
	fmt.Println("Loading SDF ")
	file, err := os.Open(SDFile)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer file.Close()

	var strucCnt int64
	strucCnt = 1
	var str string
	var bytes []byte
	sdReader := bufio.NewReader(file)
	for {
		// read one structure
		var mol string
		var molName string
		// read molecular name
		bytes, _, err = sdReader.ReadLine()
		str = string(bytes)
		if err != nil {
			if err == io.EOF {
				return strucCnt, nil
			} else {
				return strucCnt, err
			}
		}
		if molName = str; molName == "" {
			molName = "noname"
		}
		mol += molName + "\n"

		// read mol data
		for {
			bytes, _, err = sdReader.ReadLine()
			str = string(bytes)
			if err != nil {
				if err == io.EOF {
					return strucCnt, nil
				} else {
					return strucCnt, err
				}
			}
			//str = strings.TrimSpace(str)
			mol += str + "\n"
			// make it right for win(\r\n) or linux(\n)
			if str == "M  END" {
				break
			}
		}
		var lig models.Ligand
		lig.Name = molName
		lig.Mol = mol

		// generate inchikey for molecule
		//zinc.Inchi = GenInchiKey(tmpFolder, molName, mol)

		if err != nil {
			// go to "$$$$"
			for {
				bytes, _, err = sdReader.ReadLine()
				str = string(bytes)
				if err != nil {
					if err == io.EOF {
						return strucCnt, nil
					} else {
						return strucCnt, err
					}
				}
				// make it right for win(\r\n) or linux(\n)
				if str == "$$$$" {
					break
				}
			}
			// continue to read next molecule
			continue
		}

		// read property, need to do ...

		// go to "$$$$"
		for {
			bytes, _, err = sdReader.ReadLine()
			str = string(bytes)
			if err != nil {
				if err == io.EOF {
					return strucCnt, nil
				} else {
					return strucCnt, err
				}
			}
			if strings.Contains(str, "ChEMBLID") {
				bytes, _, err = sdReader.ReadLine()
				str = string(bytes)
				lig.ChemblId = str
				fmt.Println(str)
			}
			if strings.Contains(str, "Formula") {
				bytes, _, err = sdReader.ReadLine()
				str = string(bytes)
				lig.Formula = str
				fmt.Println(str)
			}
			if strings.Contains(str, "mol_weight") {
				bytes, _, err = sdReader.ReadLine()
				str = string(bytes)
				lig.MolWeight, _ = strconv.ParseFloat(strings.TrimSpace(str), 64)
				fmt.Println(str)
			}
			if strings.Contains(str, "Smiles") {
				bytes, _, err = sdReader.ReadLine()
				str = string(bytes)
				lig.Smiles = str
				fmt.Println(str)
			}
			if strings.Contains(str, "> <Inchi>") {
				bytes, _, err = sdReader.ReadLine()
				str = string(bytes)
				lig.Inchi = str
				fmt.Println(str)
			}
			if strings.Contains(str, "> <InchiKey>") {
				bytes, _, err = sdReader.ReadLine()
				str = string(bytes)
				lig.Inchikey = str
				fmt.Println(str)
			}
			if strings.Contains(str, "> <HBA>") {
				bytes, _, err = sdReader.ReadLine()
				str = string(bytes)
				lig.Hba, _ = strconv.ParseInt(strings.TrimSpace(str), 10, 64)
				fmt.Println(str)
			}
			if strings.Contains(str, "> <HBD>") {
				bytes, _, err = sdReader.ReadLine()
				str = string(bytes)
				lig.Hbd, _ = strconv.ParseInt(strings.TrimSpace(str), 10, 64)
				fmt.Println(str)
			}
			if strings.Contains(str, "> <RTB>") {
				bytes, _, err = sdReader.ReadLine()
				str = string(bytes)
				lig.Rtb, _ = strconv.ParseInt(strings.TrimSpace(str), 10, 64)
				fmt.Println(str)
			}
			if strings.Contains(str, "> <ALOGP>") {
				bytes, _, err = sdReader.ReadLine()
				str = string(bytes)
				lig.Alogp, _ = strconv.ParseFloat(strings.TrimSpace(str), 64)
				fmt.Println(str)
			}
			// make it right for win(\r\n) or linux(\n)
			if str == "$$$$" {
				break
			}
		}
		lig.Id = strucCnt
		err = lig.Insert()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		strucCnt++
	}
	fmt.Println(strucCnt)
	return strucCnt, nil
}
