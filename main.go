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

		if !info.IsDir() && (strings.HasSuffix(strings.ToLower(info.Name()), ".mp4")) {
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

func getVideoDuration(filePath string) (time.Duration, error) {
	data, err := ffprobe.GetProbeData(filePath, 120000*time.Millisecond)
	if err != nil {
		log.Printf("Erro ao obter dados: %v", err)
		return 0, err
	}
	return data.Format.Duration(), nil
}
