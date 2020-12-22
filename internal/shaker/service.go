package shaker

// need to mock ? implement this interface in unit tests
type Service interface {
	Search(query string) chan SearchResult
	SearchDB(query, table string) chan SearchResult
}

type serviceImpl struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repository: repo}
}

func (s serviceImpl) Search(query string) chan SearchResult {
	return s.repository.Search(query)
}

func (s serviceImpl) SearchDB(query, table string) chan SearchResult {
	return s.repository.SearchDB(query, table)
}
