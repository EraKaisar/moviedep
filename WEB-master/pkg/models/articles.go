package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoMovie   = errors.New("models: no matching movie found")
	ErrDuplicate = errors.New("models: duplicate movie title")
)

type Movie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Genre       string             `bson:"genre"`
	Rating      int                `bson:"rating"`
	SessionTime time.Time          `bson:"sessionTime"`
}

type MovieModel struct {
	Collection *mongo.Collection
}

func (m *MovieModel) Create(title, genre string, rating int, sessionTime time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем новый документ для фильма
	movie := Movie{
		Title:       title,
		Genre:       genre,
		Rating:      rating,
		SessionTime: sessionTime,
	}

	// Вставляем документ в коллекцию MongoDB
	_, err := m.Collection.InsertOne(ctx, movie)
	if err != nil {
		if isDuplicateError(err) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}

func (m *MovieModel) Update(id primitive.ObjectID, title, genre string, rating int, sessionTime time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Обновление данных фильма
	update := bson.M{
		"$set": bson.M{
			"title":       title,
			"genre":       genre,
			"rating":      rating,
			"sessionTime": sessionTime,
		},
	}
	res, err := m.Collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}

	// Проверка результата обновления
	if res.MatchedCount == 0 {
		return ErrNoMovie
	}

	return nil
}
func (m *MovieModel) Delete(_id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Преобразование id из строки в ObjectID

	// Выполняем запрос удаления по ID
	res, err := m.Collection.DeleteOne(ctx, bson.M{"_id": _id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNoMovie
	}
	return nil
}

func (m *MovieModel) Get(id primitive.ObjectID) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Выполняем запрос поиска по ID
	var movie Movie
	err := m.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoMovie
		}
		return nil, err
	}
	return &movie, nil
}

func (m *MovieModel) Latest(limit int64) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Выполняем запрос на поиск последних фильмов
	opts := options.Find().SetSort(bson.D{{"sessionTime", -1}}).SetLimit(limit)
	cur, err := m.Collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// Обрабатываем результаты запроса
	var movies []*Movie
	for cur.Next(ctx) {
		var movie Movie
		if err := cur.Decode(&movie); err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}

func (m *MovieModel) GetMovieByGenre(genre string) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Выполняем запрос на поиск фильмов по жанру
	opts := options.Find().SetSort(bson.D{{"sessionTime", -1}})
	cur, err := m.Collection.Find(ctx, bson.M{"genre": genre}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// Обрабатываем результаты запроса
	var movies []*Movie
	for cur.Next(ctx) {
		var movie Movie
		if err := cur.Decode(&movie); err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}

func isDuplicateError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "E11000")
}
