package main

import (
	"crypto/md5"
	"crypto/sha256"
	"flag"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

var dictionary = []rune("abcdefghijklmnopqrstuvwxyz")
var dictionaryLength = len(dictionary)

func intToPassword(n int) string {
	password := make([]rune, 5)
	for i := 4; i >= 0; i-- {
		password[i] = dictionary[n%dictionaryLength]
		n /= dictionaryLength
	}
	return string(password)
}

func calculateMD5(s string) string {
	hash := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func calculateSHA256(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}

func bruteForce(hash string, threadsAmount int, calculateHashFunc func(string) string) {
	var wg sync.WaitGroup
	stopChannel := make(chan struct{})
	var once sync.Once

	totalCombinations := int(math.Pow(float64(dictionaryLength), 5))

	worker := func(id, step int) {
		defer wg.Done()
		for i := id; i < totalCombinations; i += step {
			select {
			case <-stopChannel:
				return
			default:
				if calculateHashFunc(intToPassword(i)) == hash {
					fmt.Printf("Поток %d: Найден пароль - %s\n", id, intToPassword(i))
					once.Do(func() { close(stopChannel) })
					return
				}
			}
		}
	}

	wg.Add(threadsAmount)
	for i := 0; i < threadsAmount; i++ {
		go worker(i, threadsAmount)
	}

	wg.Wait()
}

func main() {
	hashType := flag.String("type", "", "Тип хэш-функции: md5 или sha256")
	hash := flag.String("hash", "", "Хэш для брутфорса")
	threadsAmount := flag.Int("threads", 1, "Количество потоков для брутфорса")

	flag.Parse()

	switch *hashType {
	case "md5":
		if len(*hash) != 32 {
			fmt.Println("Длина хэша, сформированного алгоритмом md5 должна быть равна 32")
			return
		}
	case "sha256":
		if len(*hash) != 64 {
			fmt.Println("Длина хэша, сформированного алгоритмом sha256 должна быть равна 64")
			return
		}
	default:
		fmt.Println("Программа работает только со следующими типами хэширования: md5 или sha256")
		return
	}

	// Проверка количества потоков
	if *threadsAmount <= 0 {
		fmt.Println("Количество потоков должно быть натуральным числом")
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println(*hashType, *hash, *threadsAmount)
	start := time.Now()
	if *hashType == "sha256" {
		bruteForce(*hash, *threadsAmount, calculateSHA256)
	} else {
		bruteForce(*hash, *threadsAmount, calculateMD5)
	}
	elapsed := time.Since(start)

	threadPlural := "потоком"
	if *threadsAmount >= 2 {
		threadPlural = "потоками"
	}

	fmt.Printf("Время перебора %d %s: %s\n", *threadsAmount, threadPlural, elapsed)
}