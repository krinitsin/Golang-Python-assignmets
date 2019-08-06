package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {
	concurrency := flag.Int("c", 5, "How many scrapers to run in parallel.") //Максимальное число горутин.
	sourcefile := flag.String("f", "urls.txt", "File with urls")
	flag.Parse()
	var totalcount int64 //Всего вхождений "Go"

	//Чтение начальных данных
	f, err := os.Open(*sourcefile)
	defer f.Close()
	if err != nil {
		log.Fatalf("could not open file %v: %v", *sourcefile, err)
	}
	scanner := bufio.NewScanner(f)
	//Вывод информации результатов из каждой рутины
	logger := log.New(os.Stdout, "", 0)

	//Канал для тасков.
	tasks := make(chan string)

	go func() {
		for scanner.Scan() {
			url := scanner.Text()
			tasks <- url
		}
		close(tasks)
	}()

	// create workers
	var wg sync.WaitGroup //Группа ожидания выполнения
	wg.Add(*concurrency)
	for i := 0; i < *concurrency; i++ {
		go func() {
			defer wg.Done()
			for url := range tasks {
				count, err := parse(url)
				if err != nil {
					logger.Printf("could not parse %v: %v", url, err)
					continue
				}
				atomic.AddInt64(&totalcount, int64(count))     //Общее количество вхождений подстроки Go threadsafe
				logger.Println("Count for ", url, ": ", count) //Вывод на экран информации по текущему url threadsafe
			}
		}()
	}
	//Ожидание выполнения всех goroutine
	wg.Wait()
	//Итоговый вывод
	fmt.Println("Total: ", totalcount)
}

func parse(url string) (int, error) {
	res, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("could not get %s: %v", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusTooManyRequests {
			return 0, fmt.Errorf("you are being rate limited")
		}

		return 0, fmt.Errorf("bad response from server: %s", res.Status)
	}

	// parse body with Ioutil.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("could not parse page: %v", err)
	}

	// extract info we want.
	count := strings.Count(string(body), "Go") //Кол-во вхождений подстроки Go в теле ответа

	return count, nil
}
