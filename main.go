package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	ffprobe "github.com/vansante/go-ffprobe"
)

var validSuffixes = []string{
	".mp4",
	".ts",
	".mov",
	".wmv",
	".avi",
	".flv",
	".mkv",
	".gif",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Defina o diretório onde estão os vídeos Ex: '/caminho/do/diretorio'")
		os.Exit(1)
	}

	directory := os.Args[1]
	totalDuration := time.Duration(0)

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}

		if !info.IsDir() && hasValidSuffix(info.Name(), validSuffixes) {
			duration, err := getVideoDuration(path)
			if err != nil {
				log.Printf("Erro ao obter a duração do vídeo %s: %v", path, err)
				return nil
			}

			totalDuration += duration
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Duração total dos vídeos: %s\n", totalDuration.String())
}

func hasValidSuffix(name string, suffixes []string) bool {
	name = strings.ToLower(name)
	for _, suffix := range suffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

func getVideoDuration(filePath string) (time.Duration, error) {
	data, err := ffprobe.GetProbeData(filePath, 120000*time.Millisecond)
	if err != nil {
		log.Printf("Erro ao obter dados: %v", err)
		return 0, err
	}
	return data.Format.Duration(), nil
}
