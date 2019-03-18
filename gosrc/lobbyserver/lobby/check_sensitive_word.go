package lobby

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"

	"github.com/huayuego/wordfilter/trie"
)

var (
	myTrie *trie.Trie
)

func loadSensitiveWordDictionary(filename string) {
	log.Println("loadSensitiveWord:", filename)

	if filename == "" {
		log.Println("sensitive word file not exist")
		return
	}

	if myTrie == nil {
		myTrie = trie.NewTrie()
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Println("open file error:", err)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var contentString = scanner.Text()
		var ws = strings.Split(contentString, ",")
		for _, w := range ws {
			word := strings.TrimSpace(w)
			if word != "" {
				myTrie.Add(word)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}

func isContainSensitiveWord(word string) bool {
	if myTrie == nil {
		log.Println("myTrie == nil")
		return false
	}

	ok, keyword, _ := myTrie.Query(word)
	if ok {
		log.Println("sensitive word:", keyword)
		return true
	}

	return false
}

func replaceSensitiveWord(word string, replace string) (bool, string) {
	if myTrie == nil {
		log.Println("myTrie == nil")
		return false, word
	}

	ok, keyword, _ := myTrie.Query(word)
	if ok {
		log.Println("find sensitive word:", keyword)
		var newWord = word
		for _, kw := range keyword {
			newWord = strings.Replace(newWord, kw, replace, -1)
		}
		return true, newWord
	}

	return false, word
}
