package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
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

var verboseMode bool
var help bool
var directory string

var appName = "videos_duration"

var helpMessage = fmt.Sprintf(`
Flags:
	-v | Modo verboso.
	-d | Diretório a ser escaniado.
	-h | Modo de ajuda.

Exemplos de uso:
	%s -d="/caminho/do/diretorio" << Escaneia o diretório passado.

	%s -d="."  << Escaneia o diretório atual.

	%s -v -d="."  << Escaneia em modo verboso.

`, color.CyanString(appName), color.CyanString(appName), color.CyanString(appName))

func main() {

	flag.BoolVar(&verboseMode, "v", false, "Modo verboso")
	flag.BoolVar(&help, "h", false, "Modo de ajuda")
	flag.StringVar(&directory, "d", "", "Diretório a ser escaniado")
	flag.Parse()

	if help || directory == "" {
		fmt.Print(helpMessage)
		return
	}

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

			if verboseMode {
				fmt.Printf("Nome: %s | Duração: %s\n", color.CyanString(info.Name()), color.GreenString(duration.String()))
			}

			totalDuration += duration
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nDuração total dos vídeos: %s\n", color.GreenString(totalDuration.String()))
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
