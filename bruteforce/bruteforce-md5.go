package bruteforce

import (
	"crypto/md5"
	"fmt"
	"math"
	"sync"
)

// Функция для вычисления MD5-хэша строки
func calculateMD5(s string) string {
	hash := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func BruteForceMD5(hash string, threadsAmount int) {
	var wg sync.WaitGroup
	stopChannel := make(chan bool)

	totalCombinations := int(math.Pow(float64(dictionaryLength), 5));
	partSize := totalCombinations / threadsAmount

	worker := func(id int, start, end int) {
		defer wg.Done()
		for i := start; i < end; i++ {
			fmt.Printf("Поток %d: Найден пароль - %s\n", id, intToPassword(i))
			if calculateMD5(intToPassword(i)) == hash {
				fmt.Printf("Поток %d: Найден пароль - %s\n", id, intToPassword(i))
				close(stopChannel)
				return
			}
			select {
			case <-stopChannel:
				return
			default:
			}
		}
	}

	wg.Add(threadsAmount)
	for i := 0; i < threadsAmount; i++ {
		start := i * partSize
		end := start + partSize
		if i == threadsAmount-1 {
			end = totalCombinations
		}
		go worker(i, start, end)
	}

	wg.Wait()
}