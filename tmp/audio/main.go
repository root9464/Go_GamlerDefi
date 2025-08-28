package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

func main() {
	// Папка с исходными файлами
	inputDir := "../audiorecorder/"
	// Временный склеенный файл
	mergedOgg := "merged.ogg"
	// Итоговый MP3 файл
	outputMP3 := "output.mp3"

	// 1. Получаем список OGG файлов и сортируем их
	files, err := getSortedOGGFiles(inputDir)
	if err != nil {
		fmt.Printf("Ошибка при получении списка файлов: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("Не найдено OGG файлов для обработки")
		return
	}

	fmt.Println("Найдены файлы:", files)

	// 2. Склеиваем файлы
	if err := mergeOGGFiles(files, mergedOgg); err != nil {
		fmt.Printf("Ошибка при склеивании файлов: %v\n", err)
		return
	}
	fmt.Println("Файлы успешно склеены в", mergedOgg)

	// 3. Конвертируем в MP3
	if err := convertToMP3(mergedOgg, outputMP3); err != nil {
		fmt.Printf("Ошибка при конвертации в MP3: %v\n", err)
		return
	}
	fmt.Println("Файл успешно конвертирован в", outputMP3)

	// 4. Удаляем временный файл
	if err := os.Remove(mergedOgg); err != nil {
		fmt.Printf("Не удалось удалить временный файл: %v\n", err)
	}

	fmt.Println("Готово!")
}

// Получаем отсортированный список OGG файлов
func getSortedOGGFiles(dir string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "1_mixed.ogg"))
	if err != nil {
		return nil, err
	}

	// Сортируем файлы по номеру
	sort.Slice(files, func(i, j int) bool {
		numI := extractNumber(files[i])
		numJ := extractNumber(files[j])
		return numI < numJ
	})

	return files, nil
}

// Извлекаем номер из имени файла
func extractNumber(filename string) int {
	base := filepath.Base(filename)
	var num int
	fmt.Sscanf(base, "segment_%d.ogg", &num)
	return num
}

// Склеиваем OGG файлы
func mergeOGGFiles(inputFiles []string, outputFile string) error {
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, file := range inputFiles {
		in, err := os.Open(file)
		if err != nil {
			return err
		}

		_, err = io.Copy(out, in)
		in.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// Конвертируем OGG в MP3 с помощью ffmpeg
func convertToMP3(inputFile, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-acodec", "libmp3lame", "-q:a", "2", outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// package main
//
// import (
// 	"fmt"
// 	"io"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"sort"
// )
//
// func main() {
// 	// Папка с исходными файлами
// 	inputDir := "../../scripts/audio/1"
// 	// Временный склеенный файл
// 	mergedOgg := "merged.ogg"
// 	// Итоговый MP3 файл
// 	outputMP3 := "output.mp3"
//
// 	// 1. Получаем список OGG файлов и сортируем их
// 	files, err := getSortedOGGFiles(inputDir)
// 	if err != nil {
// 		fmt.Printf("Ошибка при получении списка файлов: %v\n", err)
// 		return
// 	}
//
// 	if len(files) == 0 {
// 		fmt.Println("Не найдено OGG файлов для обработки")
// 		return
// 	}
//
// 	fmt.Println("Найдены файлы:", files)
//
// 	// 2. Склеиваем файлы
// 	if err := mergeOGGFiles(files, mergedOgg); err != nil {
// 		fmt.Printf("Ошибка при склеивании файлов: %v\n", err)
// 		return
// 	}
// 	fmt.Println("Файлы успешно склеены в", mergedOgg)
//
// 	// 3. Конвертируем в MP3 с очисткой звука
// 	if err := convertAndCleanAudio(mergedOgg, outputMP3); err != nil {
// 		fmt.Printf("Ошибка при обработке аудио: %v\n", err)
// 		return
// 	}
// 	fmt.Println("Файл успешно обработан и сохранен как", outputMP3)
//
// 	// 4. Удаляем временный файл
// 	if err := os.Remove(mergedOgg); err != nil {
// 		fmt.Printf("Не удалось удалить временный файл: %v\n", err)
// 	}
//
// 	fmt.Println("Готово!")
// }
//
// // Получаем отсортированный список OGG файлов
// func getSortedOGGFiles(dir string) ([]string, error) {
// 	files, err := filepath.Glob(filepath.Join(dir, "room_1_segment_*.ogg"))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Сортируем файлы по номеру
// 	sort.Slice(files, func(i, j int) bool {
// 		numI := extractNumber(files[i])
// 		numJ := extractNumber(files[j])
// 		return numI < numJ
// 	})
//
// 	return files, nil
// }
//
// // Извлекаем номер из имени файла (теперь для room_1_segment_*.ogg)
// func extractNumber(filename string) int {
// 	base := filepath.Base(filename)
// 	var num int
// 	// Обновленный паттерн для room_1_segment_X.ogg
// 	fmt.Sscanf(base, "room_1_segment_%d.ogg", &num)
// 	return num
// }
//
// // Склеиваем OGG файлы
// func mergeOGGFiles(inputFiles []string, outputFile string) error {
// 	out, err := os.Create(outputFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()
//
// 	for _, file := range inputFiles {
// 		in, err := os.Open(file)
// 		if err != nil {
// 			return err
// 		}
//
// 		_, err = io.Copy(out, in)
// 		in.Close()
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	return nil
// }
//
// // Конвертируем и очищаем аудио
// func convertAndCleanAudio(inputFile, outputFile string) error {
// 	// Команда FFmpeg с фильтрами для очистки звука:
// 	// 1. highpass - убирает низкочастотный шум
// 	// 2. lowpass - убирает высокочастотный шум
// 	// 3. afftdn - подавляет фоновый шум
// 	// 4. loudnorm - нормализует громкость
// 	cmd := exec.Command("ffmpeg",
// 		"-i", inputFile,
// 		"-af", "highpass=f=100,lowpass=f=5000,afftdn=nf=-25,loudnorm=I=-16:LRA=11:TP=-1.5",
// 		"-acodec", "libmp3lame",
// 		"-q:a", "1", // Лучшее качество (0-9, где 0 - наилучшее)
// 		"-ar", "44100", // Стандартная частота дискретизации
// 		"-ac", "2", // Стерео звучание
// 		"-y", // Перезаписывать выходной файл без подтверждения
// 		outputFile,
// 	)
//
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	return cmd.Run()
// }
