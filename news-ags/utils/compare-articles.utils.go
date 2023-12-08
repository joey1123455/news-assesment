package utils

import "github.com/blevesearch/bleve"

func CalculateSimilarity(string1, string2 string) float64 {
	// Create a new index mapping
	mapping := bleve.NewIndexMapping()

	// Create a new index
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		panic(err)
	}

	// Index the strings
	index.Index("string1", map[string]interface{}{
		"text": string1,
	})
	index.Index("string2", map[string]interface{}{
		"text": string2,
	})

	// Create a new query
	query := bleve.NewMatchQuery(string1)

	// Create a new search request
	request := bleve.NewSearchRequest(query)

	// Search the index
	result, err := index.Search(request)
	if err != nil {
		panic(err)
	}

	// Get the score of the second string
	score := result.Hits[1].Score

	return score
}
