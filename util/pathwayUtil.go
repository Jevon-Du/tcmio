package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"tcmio/models"

	"github.com/astaxie/beego"
)

type EnrichmentResult struct {
	Id             int64
	Category       string
	Term           string
	Count          int64
	Percent        float64
	Pvalue         float64
	Genes          string
	ListTotal      int64
	PopHits        int64
	PopTotal       int64
	FoldEnrichment float64
	Bonferroni     float64
	Benjamini      float64
	FDR            float64
}

func DavaidAnalysis(tar []models.Target) []EnrichmentResult {
	var enrich []EnrichmentResult
	if len(tar) < 5 {
		return enrich
	}
	jobPath := beego.AppConfig.String("jobPath") + "/"
	exePath := beego.AppConfig.String("davidEXE")
	err := os.MkdirAll(jobPath, os.ModePerm)
	targetListPath := jobPath + "/" + "list1.txt"
	resultPath := jobPath + "/" + "list1_result.txt"
	fi, err := os.Create(targetListPath)
	target_uniprots := GetTargetUniProtIDs(tar)
	all := strings.Join(target_uniprots, "\n")
	fi.WriteString(all)
	fi.Close()
	fmt.Println(exePath)

	cmd := exec.Command("python", exePath, targetListPath, resultPath)
	buf, err := cmd.Output()
	fmt.Println(string(buf))
	if err != nil {
		fmt.Println("Fail to call davaid")
		return enrich
	}
	enrich = ReadEnrichmentResult(resultPath)
	return enrich
}

func ReadEnrichmentResult(path string) []EnrichmentResult {
	//path = "E:/Program/GoProject/src/list1_result.txt"
	//path = "E:/GO_WorkSpace/src/list1_result.txt"
	var enrich []EnrichmentResult
	fmt.Println(path)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error:", err)
		return enrich
	}
	defer file.Close()

	br := bufio.NewReader(file)
	var bytes []byte
	bytes, _, err = br.ReadLine()

	for {
		bytes, _, err = br.ReadLine()
		str := string(bytes)
		if err != nil {
			fmt.Println(err)
			break
		}
		//fmt.Println(line)
		tmp := strings.Split(str, "\t")
		var e EnrichmentResult
		if len(tmp) != 13 {
			continue
		}
		e.Category = tmp[0]
		e.Term = tmp[1]
		e.Count, _ = strconv.ParseInt(tmp[2], 10, 64)
		e.Percent, _ = strconv.ParseFloat(tmp[3], 64)
		//Trunc
		//e.Percent = strconv.FormatFloat(tmp[3], 'f',4,64)
		e.Pvalue, _ = strconv.ParseFloat(tmp[4], 64)
		e.Genes = tmp[5]
		e.ListTotal, _ = strconv.ParseInt(tmp[6], 10, 64)
		e.PopHits, _ = strconv.ParseInt(tmp[7], 10, 64)
		e.PopTotal, _ = strconv.ParseInt(tmp[8], 10, 64)
		e.FoldEnrichment, _ = strconv.ParseFloat(tmp[9], 64)
		e.Bonferroni, _ = strconv.ParseFloat(tmp[10], 64)
		e.Benjamini, _ = strconv.ParseFloat(tmp[11], 64)
		e.FDR, _ = strconv.ParseFloat(tmp[12], 64)
		enrich = append(enrich, e)
	}
	//fmt.Printf("pathway:\n%v\n", enrich)
	return enrich
}
