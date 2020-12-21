package shaker

import (
	"bytes"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	fragmentLength = 250
)

type SearchResult struct {
	Fragment string `json:"f"`
	Position int    `json:"p"`
}

// need to mock ? implement this interface in unit tests
type Repository interface {
	Search(query string) chan SearchResult
}

type repositoryImpl struct {
	completeWorks    string
	completeWorksLen int
	suffixArray      *suffixarray.Index
}

func NewRepository(fileName string) (Repository, error) {
	result := &repositoryImpl{}
	err := result.load(fileName)
	return result, err
}

func (r *repositoryImpl) load(fileName string) error {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("error loading filename %q: %w", fileName, err)
	}
	r.completeWorks = string(dat)
	dat = bytes.ToLower(dat)
	r.suffixArray = suffixarray.New(dat)
	r.completeWorksLen = len(dat)
	return nil
}

func (r *repositoryImpl) Search(query string) chan SearchResult {
	result := make(chan SearchResult) // JavaScript "future"
	go func() {                       // JavaScript "async"
		idxs := r.suffixArray.Lookup([]byte(query), -1)
		for _, idx := range idxs {
			fragment := ""
			if idx < fragmentLength {
				fragment = r.completeWorks[:idx+fragmentLength]

			} else if idx+fragmentLength > r.completeWorksLen {
				fragment = r.completeWorks[idx-fragmentLength:]
			} else {
				fragment = r.completeWorks[idx-fragmentLength : idx+fragmentLength]
			}
			logrus.Debugf("delivering result @ %d", idx)
			result <- SearchResult{Fragment: strings.ReplaceAll(fragment, query, "<mark>"+query+"</mark>"), Position: idx}
		}
		close(result)
	}()
	return result // JavaScript "await"
}
