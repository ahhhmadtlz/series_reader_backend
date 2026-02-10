package series

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/entity"
	"github.com/ahhhmadtlz/series_reader_backend/internal/domain/series/param"
	"github.com/ahhhmadtlz/series_reader_backend/internal/pkg/richerror"
)

type PostgresRepository struct {
	db *sql.DB
}

func New(db *sql.DB)*PostgresRepository {
	return &PostgresRepository{db:db}
}

func (r *PostgresRepository) Create(ctx context.Context,series entity.Series)(entity.Series,error){
	const op=richerror.Op("repository.postgres.series.Create")
	genresJSON,err:=json.Marshal(series.Genres)
	
	if err!=nil{
		return entity.Series{},richerror.New(op).WithErr(err).WithMessage("failed to marshal alternative titles")
	}
 altTitlesJSON,err:=json.Marshal(series.AlternativeTitles)

 if err !=nil{
	 return entity.Series{},richerror.New(op).WithErr(err).WithMessage("failed to marshal alternative titles")
 }

	query := `
			INSERT INTO series (
				title, slug, slug_id, full_slug, description, author, artist, status, type,
				genres, alternative_titles, cover_image_url, publication_year,
				view_count, rating, is_published
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
			RETURNING id, created_at, updated_at
		`
 var createdSeries entity.Series

err = r.db.QueryRowContext(
		ctx,
		query,
		series.Title,
		series.Slug,
		series.SlugID,
		series.FullSlug,
		series.Description,
		series.Author,
		series.Artist,
		series.Status,
		series.Type,
		genresJSON,
		altTitlesJSON,
		series.CoverImageURL,
		series.PublicationYear,
		series.ViewCount,
		series.Rating,
		series.IsPublished,
	).Scan(&createdSeries.ID, &createdSeries.CreatedAt, &createdSeries.UpdatedAt)
	
	if err!=nil{
		return entity.Series{},richerror.New(op).WithErr(err).WithMessage("failed to create series")
	}

	createdSeries.Title=series.Title
	createdSeries.Slug=series.Slug
	createdSeries.SlugID=series.SlugID
	createdSeries.FullSlug=series.FullSlug
  createdSeries.Description=series.Description
	createdSeries.Author=series.Author
	createdSeries.Artist = series.Artist
	createdSeries.Status = series.Status
	createdSeries.Type = series.Type
	createdSeries.Genres = series.Genres
	createdSeries.AlternativeTitles = series.AlternativeTitles
	createdSeries.CoverImageURL = series.CoverImageURL
	createdSeries.PublicationYear = series.PublicationYear
	createdSeries.ViewCount = series.ViewCount
	createdSeries.Rating = series.Rating
	createdSeries.IsPublished = series.IsPublished

	return createdSeries, nil

}


func (r *PostgresRepository) GetByID(ctx context.Context,id uint)(entity.Series,error){
	const op=richerror.Op("repository.postgres.series.GetByID")

	query:=`
	SELECT id,title,slug,slug_id,full_slug,description,author,artist,status,type,
	genres,alternative_titles,cover_image_url,publication_year,view_count,rating,is_published,created_at,updated_at
	FROM series
	WHERE id=$1
	`

	var series entity.Series
	var genresJSON,altTitlesJSON[]byte

	err:=r.db.QueryRowContext(ctx,query,id).Scan(
		&series.ID,
    &series.Title,
    &series.Slug,
		&series.SlugID,
    &series.FullSlug,
		&series.Description,
		&series.Author,
		&series.Artist,
		&series.Status,
		&series.Type,
    &genresJSON,
		&altTitlesJSON,
		&series.CoverImageURL,
		&series.PublicationYear,
		&series.ViewCount,
		&series.Rating,
		&series.IsPublished,
		&series.CreatedAt,
		&series.UpdatedAt,
	)
	if err==sql.ErrNoRows {
		return entity.Series{},richerror.New(op).WithErr(err).WithMessage("series not found").WithKind(richerror.KindNotFound)
	}

	if err!=nil{
		return entity.Series{},richerror.New(op).WithErr(err).WithMessage("failed to get series")
	}

	if err:=json.Unmarshal(genresJSON,&series.Genres);err!=nil{
		return entity.Series{},richerror.New(op).WithErr(err).WithMessage("failed to unmarshal genres")
	}
	if err:=json.Unmarshal(altTitlesJSON,&series.AlternativeTitles);err!=nil{
		return  entity.Series{},richerror.New(op).WithErr(err).WithMessage("failed to unmarshal alternative titles")
	}

	return  series,nil
}

func (r *PostgresRepository)GetList(ctx context.Context,req param.GetListRequest)([]entity.Series,int,error){
	const op=richerror.Op("repository.postgres.series.GetList")

	whereClauses :=[]string{}
	args:=[]interface{}{}
	argCounter:=1

	if req.Filter.Status !=""{
		whereClauses=append(whereClauses, fmt.Sprintf("Status= $%d",argCounter))
		args = append(args, req.Filter.Status)
		argCounter++
	}

	if req.Filter.Type!=""{
		whereClauses=append(whereClauses, fmt.Sprintf("type= $%d",argCounter))
		args = append(args, req.Filter.Type)
		argCounter++
	}

	if req.Filter.IsPublished !=nil{
		whereClauses=append(whereClauses, fmt.Sprintf("is_published = $%d",argCounter))
		args = append(args, *req.Filter.IsPublished)
		argCounter++
	}

	if len(req.Filter.Genres)>0 {
		genresJSON,_:=json.Marshal(req.Filter.Genres)
		whereClauses=append(whereClauses, fmt.Sprintf("genres @> $%d::jsonb",argCounter))
		args=append(args, genresJSON)
		argCounter++
	}

	if req.Filter.Search != "" {
		searchPattern := "%" + req.Filter.Search + "%"
		whereClauses = append(whereClauses, fmt.Sprintf("(title ILIKE $%d OR author ILIKE $%d OR artist ILIKE $%d)", argCounter, argCounter, argCounter))
		args = append(args, searchPattern)
		argCounter++
	}

	whereClause  := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	countQuery:=fmt.Sprintf("SELECT COUNT(*) FROM series %s",whereClause)
  var totalCount int
	err:=r.db.QueryRowContext(ctx,countQuery,args...).Scan(&totalCount)
	if err !=nil{
		return nil,0,richerror.New(op).WithErr(err).WithMessage("failed to count series")
	}

	orderBy:="created_at DESC"
	if req.Sort.SortBy !=""{
		direction:="DESC"
		if req.Sort.SortOrder=="asc"{
			direction="ASC"
		}
		switch req.Sort.SortBy{
			case "rating","view_count","created_at","title":
         orderBy=fmt.Sprintf("%s %s",req.Sort.SortBy,direction)
		}
	}

	page:=req.Pagination.Page

	if page <1{
		page=1
	}

	pageSize:=req.Pagination.PageSize


	if pageSize <1 || pageSize >100 {
		pageSize=20
	}

	offset:=(page-1)*pageSize

	query:=fmt.Sprintf(`
	  SELECT id ,title,slug,slug_id,full_slug,description,author,artist,status,type,
		genres,alternative_titles,cover_image_url,publication_year,
		view_count,rating,is_published,created_at ,updated_at
		FROM series
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`,whereClause,orderBy,argCounter,argCounter+1)

	args=append(args, pageSize,offset)

	rows,err:=r.db.QueryContext(ctx,query,args...)

	if err !=nil{
		return nil,0,richerror.New(op).WithErr(err).WithMessage("failed to query series list")
	}

	defer rows.Close()

  seriesList:=[]entity.Series{}

	for rows.Next(){
		var series entity.Series
		var genresJSON,altTitlesJSON []byte

		err:=rows.Scan(
			&series.ID,
			&series.Title,
			&series.Slug,
			&series.SlugID,
			&series.FullSlug,
			&series.Description,
			&series.Author,
			&series.Artist,
			&series.Status,
			&series.Type,
			&genresJSON,
			&altTitlesJSON,
			&series.CoverImageURL,
			&series.PublicationYear,
			&series.ViewCount,
			&series.Rating,
			&series.IsPublished,
			&series.CreatedAt,
			&series.UpdatedAt,
		)

		if err !=nil{
			return nil,0,richerror.New(op).WithErr(err).WithMessage("failed to scan series row")
		}

		if err:=json.Unmarshal(genresJSON,&series.Genres);err!=nil{
			return nil,0,richerror.New(op).WithErr(err).WithMessage("failed to unmarshal genres")
		}

		if err:=json.Unmarshal(altTitlesJSON,&series.AlternativeTitles);err!=nil{
			return  nil,0,richerror.New(op).WithErr(err).WithMessage("failed to unamarshal alternative titles")
		}

		seriesList=append(seriesList, series)
	}

	if err=rows.Err();err !=nil{
		return nil,0,richerror.New(op).WithErr(err).WithMessage("error iterating series rows")
	}

	return seriesList,totalCount,nil
}

func (r *PostgresRepository) Update(ctx context.Context, id uint, series entity.Series) (entity.Series, error) {
	const op =richerror.Op("repository.postgres.series.Update")

	// Convert slices to JSON
	genresJSON, err := json.Marshal(series.Genres)
	if err != nil {
		return entity.Series{}, richerror.New(op).WithErr(err).WithMessage("failed to marshal genres")
	}

	altTitlesJSON, err := json.Marshal(series.AlternativeTitles)
	if err != nil {
		return entity.Series{}, richerror.New(op).WithErr(err).WithMessage("failed to marshal alternative titles")
	}

	query := `
		UPDATE series SET
			title = $1,
			slug = $2,
			slug_id = $3,
			full_slug = $4,
			description = $5,
			author = $6,
			artist = $7,
			status = $8,
			type = $9,
			genres = $10,
			alternative_titles = $11,
			cover_image_url = $12,
			publication_year = $13,
			is_published = $14
		WHERE id = $15
		RETURNING view_count, rating, created_at, updated_at
	`

	var updatedSeries entity.Series
	err = r.db.QueryRowContext(
		ctx,
		query,
		series.Title,
		series.Slug,
		series.SlugID,
		series.FullSlug,
		series.Description,
		series.Author,
		series.Artist,
		series.Status,
		series.Type,
		genresJSON,
		altTitlesJSON,
		series.CoverImageURL,
		series.PublicationYear,
		series.IsPublished,
		id,
	).Scan(&updatedSeries.ViewCount, &updatedSeries.Rating, &updatedSeries.CreatedAt, &updatedSeries.UpdatedAt)

	if err == sql.ErrNoRows {
		return entity.Series{}, richerror.New(op).WithErr(err).WithMessage("series not found").WithKind(richerror.KindNotFound)
	}
	if err != nil {
		return entity.Series{}, richerror.New(op).WithErr(err).WithMessage("failed to update series")
	}

	// Copy fields
	updatedSeries.ID = id
	updatedSeries.Title = series.Title
	updatedSeries.Slug = series.Slug
	updatedSeries.SlugID=series.SlugID
	updatedSeries.FullSlug=series.FullSlug
	updatedSeries.Description = series.Description
	updatedSeries.Author = series.Author
	updatedSeries.Artist = series.Artist
	updatedSeries.Status = series.Status
	updatedSeries.Type = series.Type
	updatedSeries.Genres = series.Genres
	updatedSeries.AlternativeTitles = series.AlternativeTitles
	updatedSeries.CoverImageURL = series.CoverImageURL
	updatedSeries.PublicationYear = series.PublicationYear
	updatedSeries.IsPublished = series.IsPublished

	return updatedSeries, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uint) error {
	const op = richerror.Op("repository.postgres.series.Delete")

	query := `DELETE FROM series WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return richerror.New(op).WithErr(err).WithMessage("failed to delete series")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return richerror.New(op).WithErr(err).WithMessage("failed to get rows affected")
	}

	if rowsAffected == 0 {
		return richerror.New(op).WithMessage("series not found").WithKind(richerror.KindNotFound)
	}

	return nil
}

func (r *PostgresRepository) IncrementViewCount(ctx context.Context, id uint) error {
	const op = richerror.Op("repository.postgres.series.IncrementViewCount")

	query := `UPDATE series SET view_count = view_count + 1 WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return richerror.New(op).WithErr(err).WithMessage("failed to increment view count")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return richerror.New(op).WithErr(err).WithMessage("failed to get rows affected")
	}

	if rowsAffected == 0 {
		return richerror.New(op).WithMessage("series not found").WithKind(richerror.KindNotFound)
	}

	return nil
}

func (r *PostgresRepository) IsSlugExists(ctx context.Context, slug string) (bool, error) {
	const op = richerror.Op("repository.postgres.series.IsSlugExists")

	query := `SELECT EXISTS(SELECT 1 FROM series WHERE slug = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, slug).Scan(&exists)
	if err != nil {
		return false, richerror.New(op).WithErr(err).WithMessage("failed to check slug existence")
	}

	return exists, nil
}

func (r *PostgresRepository) GetByFullSlug(ctx context.Context, fullSlug string) (entity.Series, error) {
	const op = richerror.Op("repository.postgres.series.GetBySlug")

	query := `
	SELECT id, title, slug, slug_id, full_slug, description, author, artist, status, type,
	       genres, alternative_titles, cover_image_url, publication_year,
	       view_count, rating, is_published, created_at, updated_at
	FROM series
	WHERE full_slug = $1
	LIMIT 1
	`

	var s entity.Series
	var genresJSON, altTitlesJSON []byte

	err := r.db.QueryRowContext(ctx, query, fullSlug).Scan(
		&s.ID, &s.Title, &s.Slug, &s.SlugID, &s.FullSlug,
		&s.Description, &s.Author, &s.Artist, &s.Status, &s.Type,
		&genresJSON, &altTitlesJSON, &s.CoverImageURL, &s.PublicationYear,
		&s.ViewCount, &s.Rating, &s.IsPublished, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return entity.Series{}, richerror.New(op).WithErr(err).WithMessage("series not found").WithKind(richerror.KindNotFound)
	}
	if err != nil {
		return entity.Series{}, richerror.New(op).WithErr(err).WithMessage("failed to get series by slug").WithKind(richerror.KindUnexpected)
	}

	if err := json.Unmarshal(genresJSON, &s.Genres); err != nil {
		return entity.Series{}, richerror.New(op).WithErr(err).WithMessage("failed to unmarshal genres").WithKind(richerror.KindUnexpected)
	}
	if err := json.Unmarshal(altTitlesJSON, &s.AlternativeTitles); err != nil {
		return entity.Series{}, richerror.New(op).WithErr(err).WithMessage("failed to unmarshal alternative titles").WithKind(richerror.KindUnexpected)
	}

	return s, nil
}
