package search

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/hippokampe/api/app/models"
)

type Search struct {
	filename string
}

func New(filename string) *Search {
	path, _ := filepath.Split(filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, os.ModePerm)
	}

	current := generateName()
	filename += current + ".bleve"

	return &Search{
		filename: filename,
	}
}

func (s *Search) IndexProjects(projects []models.Project) error {
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		mapping := bleve.NewIndexMapping()
		index, err := bleve.New(s.filename, mapping)
		if err != nil {
			return err
		}

		batch := index.NewBatch()

		for _, project := range projects {
			if err := batch.Index(project.ID, project); err != nil {
				return err
			}

			if err := index.Batch(batch); err != nil {
				return err
			}
		}

		defer index.Close()
	}

	return nil
}

func (s *Search) GetProjectID(query string) (string, error) {
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		return "", errors.New("you need to index first in order to create the mapping")
	}

	queryArray := strings.Split(query, "_")
	query = strings.Join(queryArray, " ")

	index, err := bleve.Open(s.filename)
	if err != nil {
		return "", err
	}

	queryString := bleve.NewQueryStringQuery(query)
	searchRequest := bleve.NewSearchRequest(queryString)
	searchResult, err := index.Search(searchRequest)

	defer index.Close()
	if err != nil {
		return "", err
	}

	if searchResult.Total == 0 {
		return "", nil
	}

	return searchResult.Hits[0].ID, nil
}

func generateName() string {
	return "/" + time.Now().Format("2006_01_02")
}
