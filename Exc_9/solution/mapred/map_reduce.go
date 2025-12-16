package mapred

import (
	"regexp"
	"strings"
	"sync"
)

type MapReduce struct {
}

var nonAlphaNum = regexp.MustCompile(`[^a-zA-Z]+`)

func (mr *MapReduce) wordCountMapper(text string) []KeyValue {
	text = strings.ToLower(text)
	text = nonAlphaNum.ReplaceAllString(text, " ")

	words := strings.Fields(text)
	var result []KeyValue

	for _, word := range words {
		result = append(result, KeyValue{
			Key:   word,
			Value: 1,
		})
	}
	return result
}

func (mr *MapReduce) wordCountReducer(key string, values []int) KeyValue {
	sum := 0
	for _, v := range values {
		sum += v
	}
	return KeyValue{
		Key:   key,
		Value: sum,
	}
}

func (mr *MapReduce) Run(input []string) map[string]int {
	mapOut := make(chan KeyValue)
	shuffle := make(map[string][]int)

	var mapWG sync.WaitGroup

	// Map phase (concurrent)
	for _, line := range input {
		mapWG.Add(1)
		go func(text string) {
			defer mapWG.Done()
			for _, kv := range mr.wordCountMapper(text) {
				mapOut <- kv
			}
		}(line)
	}

	go func() {
		mapWG.Wait()
		close(mapOut)
	}()

	// Shuffle phase
	for kv := range mapOut {
		shuffle[kv.Key] = append(shuffle[kv.Key], kv.Value)
	}

	// Reduce phase (concurrent)
	results := make(map[string]int)
	var reduceWG sync.WaitGroup
	var mu sync.Mutex

	for key, values := range shuffle {
		reduceWG.Add(1)
		go func(k string, vals []int) {
			defer reduceWG.Done()
			kv := mr.wordCountReducer(k, vals)

			mu.Lock()
			results[kv.Key] = kv.Value
			mu.Unlock()
		}(key, values)
	}

	reduceWG.Wait()
	return results
}
