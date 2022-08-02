/*
	write a function that takes slice of urls as input and request each url concurrently
	to get the response use 1 second timeout foreach request after aggregate all responses
	print the responses to console if all the calls were successful
*/

package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	mu                          = sync.Mutex{}
	wg                          = sync.WaitGroup{}
	numberOfSuccessfulResponses = 0
)

func main() {
	// You can change the number of urls manually
	numberOfUrls := 100
	urls := UrlMaker(numberOfUrls)
	// please uncomment this line to see the urls
	//fmt.Println("urls: ", urls)

	urlChan := make(chan string, numberOfUrls)
	respChan := make(chan string, numberOfUrls)

	for w := 1; w <= numberOfUrls; w++ {
		wg.Add(1)
		go TranslatorWorker(w, urlChan, respChan)
	}

	for i := 0; i < numberOfUrls; i++ {
		urlChan <- urls[i]
	}
	close(urlChan)

	responses := make([]string, 0)
	for j := 1; j < numberOfUrls; j++ {
		responses = append(responses, <-respChan)
	}

	wg.Wait()

	fmt.Println("numberOfSuccessfulResponses: ", numberOfSuccessfulResponses)
	fmt.Println("numberOfUrls: ", numberOfUrls)
	if numberOfSuccessfulResponses == numberOfUrls {
		for _, resp := range responses {
			fmt.Println(resp)
		}
	}
}

func TranslatorWorker(id int, urlChan <-chan string, respChan chan<- string) {
	defer wg.Done()

	for u := range urlChan {
		responseBody, respErr := call(u)
		if errors.Is(respErr, context.DeadlineExceeded) {
			fmt.Println("ContextDeadlineExceeded in getting the response from the url", respErr)
			continue
		}
		if os.IsTimeout(respErr) {
			fmt.Println(fmt.Errorf(fmt.Sprintf("IsTimeoutError in getting the response from the url %s: ", u), respErr))
			continue
		}
		if respErr != nil {
			fmt.Println(fmt.Errorf(fmt.Sprintf("unable to get the response from the url %s: ", u), respErr))
			continue
		}

		//body, err := ioutil.ReadAll(resp.Body)
		//if err != nil {
		//	fmt.Println(fmt.Errorf(fmt.Sprintf("unable to read the body of response from the url %s: ", u), err))
		//	continue
		//}
		mu.Lock()
		numberOfSuccessfulResponses++
		mu.Unlock()
		respChan <- string(responseBody)
		// Please uncomment the next line to see each result
		//fmt.Println("worker ", id , "response: ", string(body))
	}
}

func call(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	resp, respErr := client.Do(req)
	if respErr != nil {
		return nil, respErr
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(fmt.Errorf(fmt.Sprintf("unable to read the body of response from the url %s: ", url), err))
	}

	return body, nil
}

func UrlMaker(numberOfUrls int) []string {
	urls := make([]string, 0)
	for i := 0; i < numberOfUrls; i++ {
		urls = append(urls, fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", rand.Intn(100)))
	}

	return urls
}
