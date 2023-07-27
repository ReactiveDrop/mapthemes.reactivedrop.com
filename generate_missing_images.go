package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	apiHost = flag.String("api-host", "", "AUTOMATIC1111 webui origin name")
	count   = flag.Int("count", 10, "number of seeds per adjective/noun combo")
)

func main() {
	flag.Parse()

	if *apiHost == "" {
		panic("-api-host flag is required")
	}

	adjectives := readLines("www/adjective.txt")
	nouns := readLines("www/noun.txt")

	err := os.MkdirAll("www/images", 0755)
	if err != nil {
		panic(err)
	}

	total := len(adjectives) * len(nouns) * *count

	writeIndexFile(adjectives, nouns)

	fmt.Printf("%d adjectives; %d nouns; %d variants\n%d total images\n", len(adjectives), len(nouns), *count, total)

	generated := 0

	for _, noun := range nouns {
		for _, adjective := range adjectives {
			generated += generateMissingImages(adjective, noun)
		}
	}

	fmt.Printf("%d images generated / %d images already present\n", generated, total-generated)
}

func readLines(filename string) []string {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.ReplaceAll(string(b), "\r", ""), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		return lines[:len(lines)-1]
	}

	return lines
}

func writeIndexFile(adjectives, nouns []string) {
	tmpl := template.Must(template.ParseFiles("www/images/index.tmpl"))

	f, err := os.Create("www/images/index.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	variants := make([]int, *count)
	for i := range variants {
		variants[i] = i
	}

	data := struct {
		Adjectives []string
		Nouns      []string
		Variants   []int
	}{
		Adjectives: adjectives,
		Nouns:      nouns,
		Variants:   variants,
	}

	err = tmpl.Execute(f, &data)
	if err != nil {
		panic(err)
	}
}

func generateMissingImages(adjective, noun string) int {
	generated := 0

	for i := 0; i < *count; i++ {
		filename := fmt.Sprintf("www/images/%s-%s-%04d.avif", adjective, noun, i)
		_, err := os.Stat(filename)
		if err == nil {
			continue
		}

		if !os.IsNotExist(err) {
			panic(err)
		}

		fmt.Printf("Generating image %d/%d for %s %s... ", i+1, *count, adjective, noun)
		start := time.Now()

		png := generateImage(adjective, noun, i)
		err = os.WriteFile(filename+".png", png, 0644)
		if err != nil {
			panic(err)
		}

		err = exec.Command("avifenc", "--speed", "0", filename+".png", filename).Run()
		if err != nil {
			panic(err)
		}

		err = os.Remove(filename + ".png")
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v\n", time.Since(start))

		generated++
	}

	return generated
}

func generateImage(adjective, noun string, number int) []byte {
	query, err := json.Marshal(map[string]interface{}{
		"sampler_name":       "Euler a",
		"steps":              50,
		"cfg_scale":          7,
		"width":              640,
		"height":             360,
		"enable_hr":          true,
		"hr_scale":           2,
		"hr_upscaler":        "Latent",
		"denoising_strength": 0.7,
		"seed":               563560 + number,
		"prompt":             adjective + " " + noun + ", level design render, wide view, dim volumetric lighting, retrofuturism",
		"negative_prompt":    "text, 2d, screenshot, watermark",
	})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, *apiHost+"/sdapi/v1/txt2img", bytes.NewReader(query))
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", "ReactiveDropMapThemeScript/1.0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("unexpected response status " + resp.Status + " for A:" + adjective + " N:" + noun + " V:" + strconv.Itoa(number))
	}

	var data struct {
		Images     [][]byte        `json:"images"`
		Parameters json.RawMessage `json:"parameters"`
		Info       string          `json:"info"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	return data.Images[0]
}
