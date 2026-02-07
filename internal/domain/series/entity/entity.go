package entity

import "time"

const (
	StatusOngoing   = "ongoing"
	StatusCompleted = "completed"
	StatusHiatus    = "hiatus"
	StatusCancelled = "cancelled"
)

const (
	TypeManga   = "manga"
	TypeManhwa  = "manhwa"
	TypeManhua  = "manhua"
	TypeComic   = "comic"
	TypeWebtoon = "webtoon"
)

type Series struct {
	ID                uint      `json:"id"`
	Title             string    `json:"title"`
	Slug              string    `json:"slug"`
	Description       string    `json:"description"`
	Author            string    `json:"author"`
	Artist            string    `json:"artist"`
	Status            string    `json:"status"`
	Type              string    `json:"type"`
	Genres            []string  `json:"genres"`
	AlternativeTitles []string  `json:"alternative_titles"`
	CoverImageURL     string    `json:"cover_image_url"`
	PublicationYear   *int      `json:"publication_year"`
	ViewCount         int       `json:"view_count"`
	Rating            float64   `json:"rating"`
	IsPublished       bool      `json:"is_published"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}


func ValidStatuses()[]string {
	return []string{StatusOngoing,StatusCompleted,StatusHiatus,StatusCancelled}
}

func ValidTypes()[]string {
	return []string{TypeManga,TypeManhwa,TypeManhua,TypeComic,TypeWebtoon}
}

func IsValidStatus(status string)bool{
	for _,s:=range ValidStatuses(){
		if s==status{
			return true
		}
	} 

	return false
}

func IsValidType(t string )bool{
	for _,typ:=range ValidTypes(){
		if typ==t {
			return true
		}
	}
	return false
}
