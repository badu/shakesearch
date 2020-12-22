package shaker

import (
	"bytes"
	"database/sql"
	"fmt"
	"html"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type ResultType int

const (
	fragmentLength = 250

	ErrorSearch     ResultType = -1
	StringSearch    ResultType = 1
	ChapterSearch   ResultType = 2
	CharacterSearch ResultType = 3
	ParagraphSearch ResultType = 4
	WorkSearch      ResultType = 5
)

func (r ResultType) MarshalJSON() ([]byte, error) {
	switch r {
	case ErrorSearch:
		return []byte(`"error"`), nil
	case StringSearch:
		return []byte(`"string"`), nil
	case ChapterSearch:
		return []byte(`"chapter"`), nil
	case CharacterSearch:
		return []byte(`"character"`), nil
	case ParagraphSearch:
		return []byte(`"paragraph"`), nil
	case WorkSearch:
		return []byte(`"work"`), nil
	default:
		return []byte(`"unknown type"`), nil
	}
	return nil, nil
}

type SearchResult struct {
	Fragment string     `json:"f"`
	Position int        `json:"p"`
	Type     ResultType `json:"t"`
}

// need to mock ? implement this interface in unit tests
type Repository interface {
	Search(query string) chan SearchResult
	SearchDB(query, table string) chan SearchResult
}

type repositoryImpl struct {
	completeWorks    string
	completeWorksLen int
	suffixArray      *suffixarray.Index
	db               *sql.DB
}

func NewRepository(fileName, sqlFileName string) (Repository, error) {
	var err error
	result := &repositoryImpl{}
	result.db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	if err = result.load(fileName); err != nil {
		return nil, err
	}
	if err := result.loadSqlLite(sqlFileName); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repositoryImpl) loadSqlLite(fileName string) error {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	sqls := strings.Split(string(file), ";")
	for _, sql := range sqls {
		sql = html.UnescapeString(sql)
		sql = strings.ReplaceAll(sql, ";\n", ".\n")
		if _, err := r.db.Exec(sql); err != nil {
			logrus.Errorf("error execting sql %q : %v", sql, err)
			return err
		}
	}
	logrus.Debugf("%d sqls executed against sqlite database", len(sqls))
	return nil
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
			result <- SearchResult{Fragment: strings.ReplaceAll(fragment, query, "<mark>"+query+"</mark>"), Position: idx, Type: StringSearch}
		}
		close(result)
	}()
	return result // JavaScript "await"
}

func (r *repositoryImpl) SearchDB(query, table string) chan SearchResult {
	result := make(chan SearchResult)
	go func() {
		switch table {
		case "chapter":
			row, err := r.db.Query("SELECT * FROM Chapters WHERE Description LIKE '%' || $1 || '%'", query)
			if err != nil {
				log.Fatal(err)
			}
			defer row.Close()
			for row.Next() { // Iterate and fetch the records from result cursor
				var WorkId string
				var ChapterId int
				var Act int
				var Scene int
				var Description string
				if err := row.Scan(&WorkId, &ChapterId, &Act, &Scene, &Description); err != nil {
					logrus.Errorf("error scanning row : %v", err)
					continue
				}
				result <- SearchResult{Position: ChapterId, Fragment: Description, Type: ChapterSearch}
			}
		case "character":
			row, err := r.db.Query("SELECT * FROM Characters WHERE CharName LIKE '%' || $1 || '%'", query)
			if err != nil {
				log.Fatal(err)
			}
			defer row.Close()
			for row.Next() { // Iterate and fetch the records from result cursor
				var CharID string
				var CharName string
				var Abbrev string
				var Works string
				var Description string
				if err := row.Scan(&CharID, &CharName, &Abbrev, &Works, &Description); err != nil {
					logrus.Errorf("error scanning row : %v", err)
					continue
				}
				result <- SearchResult{Fragment: Description, Type: CharacterSearch}
			}
		case "paragraph":
			row, err := r.db.Query("SELECT * FROM Paragraphs WHERE PlainText LIKE '%' || $1 || '%'", query)
			if err != nil {
				log.Fatal(err)
			}
			defer row.Close()
			for row.Next() { // Iterate and fetch the records from result cursor
				var WorkID string
				var ParagraphID int
				var ParagraphNum int
				var CharID string
				var PlainText string
				var Act int
				var Scene int
				if err := row.Scan(&WorkID, &ParagraphID, &ParagraphNum, &CharID, &PlainText, &Act, &Scene); err != nil {
					logrus.Errorf("error scanning row : %v", err)
					continue
				}
				result <- SearchResult{Position: ParagraphID, Fragment: PlainText, Type: ParagraphSearch}
			}
		case "work":
			row, err := r.db.Query("SELECT * FROM Works WHERE LongTitle LIKE '%' || $1 || '%'", query)
			if err != nil {
				log.Fatal(err)
			}
			defer row.Close()
			for row.Next() { // Iterate and fetch the records from result cursor
				var WorkID string
				var Title string
				var LongTitle string
				var Date int
				var GenreType string
				if err := row.Scan(&WorkID, &Title, &LongTitle, &Date, &GenreType); err != nil {
					logrus.Errorf("error scanning row : %v", err)
					continue
				}
				result <- SearchResult{Position: Date, Fragment: LongTitle, Type: WorkSearch}
			}
		default:
			// convention : type indicates error (should be dealt in frontend)
			result <- SearchResult{Type: ErrorSearch, Fragment: fmt.Sprintf("bad request : unrecognized table named %q", table)}
		}
		close(result)
	}()
	return result
}
