package pgsqldb

import (
	"cache/internal/repository/models"
)

func (r PgSqlRepository) GetArticle(id string) (models.Article, error) {
	query := `SELECT id, url, title
			FROM articles 
			WHERE id = $1 `

	var article models.Article
	if err := r.db.QueryRow(query, id).
		Scan(
			&article.ID,
			&article.URL,
			&article.Title,
		); err != nil {
		r.logger.Errorf("error occurred while getting article by id, err: %s", err)
		return article, err
	}

	return article, nil
}

func (r PgSqlRepository) GetArticles() ([]models.Article, error) {
	rows, err := r.db.Query("SELECT * FROM articles")
	if err != nil {
		r.logger.Errorf("Can't find articles. Err msg: %v", err)
		return nil, err
	}
	defer rows.Close()

	articles := []models.Article{}

	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.Title, &article.URL)
		if err != nil {
			r.logger.Errorf("Error scanning article: %v", err)
			continue
		}
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		r.logger.Errorf("Row scan error: %v", err)
		return nil, err
	}

	return articles, nil
}

func (r PgSqlRepository) CreateArticle(a models.Article) (models.Article, error) {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Errorf("error occurred while creating article, err: %s", err)
		return a, err
	}

	query := `INSERT INTO articles(id, url, title ) values ($1, $2, $3) returning id, url, title`

	if err = tx.QueryRow(query,
		a.ID, a.URL, a.Title).Scan(
		&a.ID,
		&a.URL,
		&a.Title); err != nil {
		r.logger.Errorf("error occurred while creating article, err: %s", err)
		tx.Rollback()
		return a, err
	}
	tx.Commit()

	return a, nil
}
