package bruteforce

import (
	"crypto/sha256"
	"fmt"
	"math"
	"sync"
)

func calculateSHA256(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}

func BruteForceSHA256(hash string, threadsAmount int) {
	var wg sync.WaitGroup
	stopChannel := make(chan bool)

	totalCombinations := int(math.Pow(float64(dictionaryLength), 5));
	partSize := totalCombinations / threadsAmount

	worker := func(id int, start, end int) {
		defer wg.Done()
		for i := start; i < end; i++ {
			if calculateSHA256(intToPassword(i)) == hash {
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