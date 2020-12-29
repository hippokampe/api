package search

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/blevesearch/bleve"
	"github.com/hippokampe/api/models"
)

type Search struct {
	filename string
}

func New(path, filename string) *Search {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, os.ModePerm)
	}

	current := generateName()
	filename += current + ".bleve"
	fullPath := filepath.Join(path, filename)

	return &Search{
		filename: fullPath,
	}
}

func (s *Search) IndexProjects(projects []models.Project) error {
	scope := "indexing"
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		mapping := bleve.NewIndexMapping()
		index, err := bleve.New(s.filename, mapping)
		if err != nil {
			return err
		}

		batch := index.NewBatch()

		for _, project := range projects {
			if err := batch.Index(project.ID, project); err != nil {
				return errors.Wrap(err, scope)
			}

			if err := index.Batch(batch); err != nil {
				return errors.Wrap(err, scope)
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

func (s *Search) GetProjects(query string, limit int) (models.ProjectsResultSearch, error) {
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		return models.ProjectsResultSearch{}, errors.New("you need to index first in order to create the mapping")
	}

	queryArray := strings.Split(query, "_")
	query = strings.Join(queryArray, " ")

	index, err := bleve.Open(s.filename)
	if err != nil {
		return models.ProjectsResultSearch{}, err
	}

	queryString := bleve.NewQueryStringQuery(query)
	searchRequest := bleve.NewSearchRequest(queryString)
	searchResult, err := index.Search(searchRequest)

	defer index.Close()
	if err != nil {
		return models.ProjectsResultSearch{}, err
	}

	if searchResult.Total == 0 {
		return models.ProjectsResultSearch{}, nil
	}

	var projects []models.ProjectSearch
	for idx, hit := range searchResult.Hits {
		projects = append(projects, models.ProjectSearch{
			ID:    hit.ID,
			Score: hit.Score,
		})

		if idx+1 == limit {
			break
		}
	}

	return models.ProjectsResultSearch{Query: query, Results: projects}, nil
}

func generateName() string {
	return "/" + time.Now().Format("2006_01_02")
}
