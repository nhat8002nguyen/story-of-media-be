package services

import (
	"database/sql"
	"fmt"
	"mime/multipart"

	"github.com/nhat8002nguyen/story-of-media-be/story-service/src/models"
)

func SaveStory(story models.Story) (models.Story, error) {
	sqlStatement := `
        INSERT INTO stories (id, content)
        VALUES ($1, $2)
        RETURNING id`

	id := ""
	err := models.Db.QueryRow(sqlStatement, story.ID, story.Content).Scan(&id)
	if err != nil {
		return story, err
	}
	story.ID = id
	return story, nil
}

func GetStoryByID(id string) (models.Story, error) {
	sqlStatement := `SELECT id, content FROM stories WHERE id=$1;`
	var story models.Story
	row := models.Db.QueryRow(sqlStatement, id)
	switch err := row.Scan(&story.ID, &story.Content); err {
	case sql.ErrNoRows:
		return story, fmt.Errorf("no records found for story id %s", id)
	case nil:
		return story, nil
	default:
		return story, err
	}
}

func GenerateStory(file *multipart.FileHeader) {
	fmt.Println("Process file")
}
