package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	lFlag := flag.String("locality", "", "locality")
	cFlag := flag.String("country", "", "country")

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		os.Stderr.WriteString("usage: add-loc --country Chile --locality Santiago <image>\n")
		os.Exit(1)
		return
	}

	if len(*cFlag) == 0 {
		os.Stderr.WriteString("country must be provided\n")
		os.Exit(1)
		return
	}
	locality := *lFlag
	country := *cFlag

	file, err := os.Open(args[0])
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
		return
	}

	h := sha256.New()
	_, err = file.WriteTo(h)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
		return
	}

	id := fmt.Sprintf("%x", h.Sum(nil))[:12]

	type loc struct {
		Locality string `json:"locality"`
		Country  string `json:"country"`
	}

	b, err := json.Marshal(loc{
		Locality: locality,
		Country:  country,
	})
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
		return
	}

	body := bytes.NewBuffer(b)

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("https://saws.world/images/%s", id), body)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
		return
	}

	req.Header.Add("Authorization", os.Getenv("SAWS_AUTH"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
		return
	}

	fmt.Fprintf(os.Stdout, "%d %s %s %s\n", resp.StatusCode, country, locality, id)
}
