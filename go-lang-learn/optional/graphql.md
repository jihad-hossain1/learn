# 29. GraphQL

This chapter covers building GraphQL APIs with Go, including schema definition, resolvers, mutations, subscriptions, and integration with existing applications.

## Introduction to GraphQL

GraphQL is a query language for APIs and a runtime for executing those queries against your data. Unlike REST, GraphQL allows clients to request exactly the data they need, making it more efficient and flexible.

Key benefits of GraphQL include:
- **Precise data fetching**: Clients specify exactly what data they need
- **Single request**: Multiple resources can be fetched in a single request
- **Strong typing**: The schema defines what queries are possible
- **Introspection**: The API is self-documenting
- **Version-free**: Fields can be added without breaking existing queries

## Setting Up GraphQL in Go

### Using gqlgen

[gqlgen](https://github.com/99designs/gqlgen) is a Go library for building GraphQL servers. It takes a schema-first approach, generating code based on your schema definition.

Let's set up a basic GraphQL server using gqlgen:

```go
// Install gqlgen
// go get github.com/99designs/gqlgen
// go run github.com/99designs/gqlgen init
```

### Project Structure

A typical gqlgen project structure looks like this:

```
├── gqlgen.yml           # gqlgen configuration
├── graph
│   ├── generated        # generated code
│   ├── model            # custom models
│   ├── resolver.go      # resolver implementations
│   └── schema.graphqls  # GraphQL schema
├── server.go            # server entry point
```

### Schema Definition

Let's create a schema for a book library API:

```graphql
# graph/schema.graphqls
type Book {
  id: ID!
  title: String!
  author: Author!
  publishedYear: Int
  genres: [String!]
  description: String
}

type Author {
  id: ID!
  name: String!
  books: [Book!]!
  bio: String
}

type Query {
  books: [Book!]!
  book(id: ID!): Book
  authors: [Author!]!
  author(id: ID!): Author
}

type Mutation {
  createBook(input: NewBook!): Book!
  updateBook(id: ID!, input: UpdateBook!): Book!
  deleteBook(id: ID!): Boolean!
  
  createAuthor(input: NewAuthor!): Author!
  updateAuthor(id: ID!, input: UpdateAuthor!): Author!
  deleteAuthor(id: ID!): Boolean!
}

input NewBook {
  title: String!
  authorID: ID!
  publishedYear: Int
  genres: [String!]
  description: String
}

input UpdateBook {
  title: String
  authorID: ID
  publishedYear: Int
  genres: [String!]
  description: String
}

input NewAuthor {
  name: String!
  bio: String
}

input UpdateAuthor {
  name: String
  bio: String
}
```

### Generating Code

After defining the schema, generate the code:

```bash
go run github.com/99designs/gqlgen generate
```

This will create the necessary files in the `graph/generated` directory and update your resolvers.

### Implementing Resolvers

Resolvers are functions that fetch the data for each field in your GraphQL schema. Let's implement some basic resolvers:

```go
// graph/resolver.go
package graph

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/yourusername/graphql-example/graph/model"
)

// This file will not be regenerated automatically.
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	books   map[string]*model.Book
	authors map[string]*model.Author
	mutex   sync.RWMutex
}

func NewResolver() *Resolver {
	// Initialize with some sample data
	authors := make(map[string]*model.Author)
	books := make(map[string]*model.Book)
	
	// Create authors
	authorID1 := uuid.New().String()
	authorID2 := uuid.New().String()
	
	authors[authorID1] = &model.Author{
		ID:   authorID1,
		Name: "J.K. Rowling",
		Bio:  "British author best known for the Harry Potter series",
	}
	
	authors[authorID2] = &model.Author{
		ID:   authorID2,
		Name: "George R.R. Martin",
		Bio:  "American novelist best known for A Song of Ice and Fire",
	}
	
	// Create books
	bookID1 := uuid.New().String()
	bookID2 := uuid.New().String()
	bookID3 := uuid.New().String()
	
	books[bookID1] = &model.Book{
		ID:            bookID1,
		Title:         "Harry Potter and the Philosopher's Stone",
		AuthorID:      authorID1,
		PublishedYear: 1997,
		Genres:        []string{"Fantasy", "Young Adult"},
		Description:   "The first book in the Harry Potter series",
	}
	
	books[bookID2] = &model.Book{
		ID:            bookID2,
		Title:         "A Game of Thrones",
		AuthorID:      authorID2,
		PublishedYear: 1996,
		Genres:        []string{"Fantasy", "Epic"},
		Description:   "The first book in A Song of Ice and Fire series",
	}
	
	books[bookID3] = &model.Book{
		ID:            bookID3,
		Title:         "A Clash of Kings",
		AuthorID:      authorID2,
		PublishedYear: 1998,
		Genres:        []string{"Fantasy", "Epic"},
		Description:   "The second book in A Song of Ice and Fire series",
	}
	
	return &Resolver{
		books:   books,
		authors: authors,
	}
}
```

### Query Resolvers

Implement the query resolvers to fetch data:

```go
// graph/schema.resolvers.go
package graph

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/yourusername/graphql-example/graph/generated"
	"github.com/yourusername/graphql-example/graph/model"
)

// Books returns all books
func (r *queryResolver) Books(ctx context.Context) ([]*model.Book, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	books := make([]*model.Book, 0, len(r.books))
	for _, book := range r.books {
		books = append(books, book)
	}
	
	return books, nil
}

// Book returns a book by ID
func (r *queryResolver) Book(ctx context.Context, id string) (*model.Book, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	book, exists := r.books[id]
	if !exists {
		return nil, errors.New("book not found")
	}
	
	return book, nil
}

// Authors returns all authors
func (r *queryResolver) Authors(ctx context.Context) ([]*model.Author, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	authors := make([]*model.Author, 0, len(r.authors))
	for _, author := range r.authors {
		authors = append(authors, author)
	}
	
	return authors, nil
}

// Author returns an author by ID
func (r *queryResolver) Author(ctx context.Context, id string) (*model.Author, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	author, exists := r.authors[id]
	if !exists {
		return nil, errors.New("author not found")
	}
	
	return author, nil
}

// Book resolver
func (r *bookResolver) Author(ctx context.Context, obj *model.Book) (*model.Author, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	author, exists := r.authors[obj.AuthorID]
	if !exists {
		return nil, errors.New("author not found")
	}
	
	return author, nil
}

// Author resolver
func (r *authorResolver) Books(ctx context.Context, obj *model.Author) ([]*model.Book, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var books []*model.Book
	for _, book := range r.books {
		if book.AuthorID == obj.ID {
			books = append(books, book)
		}
	}
	
	return books, nil
}

func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }
func (r *Resolver) Book() generated.BookResolver { return &bookResolver{r} }
func (r *Resolver) Author() generated.AuthorResolver { return &authorResolver{r} }

type queryResolver struct{ *Resolver }
type bookResolver struct{ *Resolver }
type authorResolver struct{ *Resolver }
```

### Mutation Resolvers

Implement the mutation resolvers to modify data:

```go
// graph/schema.resolvers.go (continued)

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }

// CreateBook creates a new book
func (r *mutationResolver) CreateBook(ctx context.Context, input model.NewBook) (*model.Book, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if author exists
	_, exists := r.authors[input.AuthorID]
	if !exists {
		return nil, errors.New("author not found")
	}
	
	// Create new book
	id := uuid.New().String()
	book := &model.Book{
		ID:            id,
		Title:         input.Title,
		AuthorID:      input.AuthorID,
		PublishedYear: input.PublishedYear,
		Genres:        input.Genres,
		Description:   input.Description,
	}
	
	r.books[id] = book
	return book, nil
}

// UpdateBook updates an existing book
func (r *mutationResolver) UpdateBook(ctx context.Context, id string, input model.UpdateBook) (*model.Book, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if book exists
	book, exists := r.books[id]
	if !exists {
		return nil, errors.New("book not found")
	}
	
	// Update fields if provided
	if input.Title != nil {
		book.Title = *input.Title
	}
	
	if input.AuthorID != nil {
		// Check if author exists
		_, exists := r.authors[*input.AuthorID]
		if !exists {
			return nil, errors.New("author not found")
		}
		book.AuthorID = *input.AuthorID
	}
	
	if input.PublishedYear != nil {
		book.PublishedYear = *input.PublishedYear
	}
	
	if input.Genres != nil {
		book.Genres = input.Genres
	}
	
	if input.Description != nil {
		book.Description = *input.Description
	}
	
	return book, nil
}

// DeleteBook deletes a book
func (r *mutationResolver) DeleteBook(ctx context.Context, id string) (bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if book exists
	_, exists := r.books[id]
	if !exists {
		return false, errors.New("book not found")
	}
	
	delete(r.books, id)
	return true, nil
}

// CreateAuthor creates a new author
func (r *mutationResolver) CreateAuthor(ctx context.Context, input model.NewAuthor) (*model.Author, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Create new author
	id := uuid.New().String()
	author := &model.Author{
		ID:   id,
		Name: input.Name,
		Bio:  input.Bio,
	}
	
	r.authors[id] = author
	return author, nil
}

// UpdateAuthor updates an existing author
func (r *mutationResolver) UpdateAuthor(ctx context.Context, id string, input model.UpdateAuthor) (*model.Author, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if author exists
	author, exists := r.authors[id]
	if !exists {
		return nil, errors.New("author not found")
	}
	
	// Update fields if provided
	if input.Name != nil {
		author.Name = *input.Name
	}
	
	if input.Bio != nil {
		author.Bio = *input.Bio
	}
	
	return author, nil
}

// DeleteAuthor deletes an author
func (r *mutationResolver) DeleteAuthor(ctx context.Context, id string) (bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if author exists
	_, exists := r.authors[id]
	if !exists {
		return false, errors.New("author not found")
	}
	
	// Check if author has books
	for _, book := range r.books {
		if book.AuthorID == id {
			return false, errors.New("cannot delete author with books")
		}
	}
	
	delete(r.authors, id)
	return true, nil
}
```

### Server Setup

Set up the GraphQL server:

```go
// server.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/yourusername/graphql-example/graph"
	"github.com/yourusername/graphql-example/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	resolvers := graph.NewResolver()
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

## Advanced GraphQL Features

### Subscriptions

Subscriptions allow clients to receive real-time updates when data changes. Let's add a subscription to our schema:

```graphql
# Add to schema.graphqls
type Subscription {
  bookAdded: Book!
}
```

Implement the subscription resolver:

```go
// Add to resolver.go
type Resolver struct {
	books          map[string]*model.Book
	authors        map[string]*model.Author
	mutex          sync.RWMutex
	bookSubscribers map[string]chan *model.Book
}

func NewResolver() *Resolver {
	// ... existing code ...
	return &Resolver{
		books:           books,
		authors:         authors,
		bookSubscribers: make(map[string]chan *model.Book),
	}
}

// Add to schema.resolvers.go
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) BookAdded(ctx context.Context) (<-chan *model.Book, error) {
	id := uuid.New().String()
	bookChan := make(chan *model.Book, 1)
	
	r.mutex.Lock()
	r.bookSubscribers[id] = bookChan
	r.mutex.Unlock()
	
	go func() {
		<-ctx.Done()
		r.mutex.Lock()
		delete(r.bookSubscribers, id)
		close(bookChan)
		r.mutex.Unlock()
	}()
	
	return bookChan, nil
}

// Modify CreateBook to notify subscribers
func (r *mutationResolver) CreateBook(ctx context.Context, input model.NewBook) (*model.Book, error) {
	// ... existing code ...
	
	// Notify subscribers
	for _, ch := range r.bookSubscribers {
		ch <- book
	}
	
	return book, nil
}
```

### Middleware

Add middleware for authentication, logging, and error handling:

```go
// middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		
		if auth != "" {
			parts := strings.Split(auth, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				// In a real app, validate the token
				userID := "user-123" // Get from token
				
				ctx := context.WithValue(r.Context(), UserIDKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}
		
		next.ServeHTTP(w, r)
	})
}

// Get user ID from context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
```

Update the server to use middleware:

```go
// server.go
func main() {
	// ... existing code ...
	
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers}))
	
	// Add middleware
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", middleware.Auth(srv))
	
	// ... existing code ...
}
```

### Directives

Directives allow you to add custom behavior to your schema. Let's add an `@auth` directive to restrict access to certain fields:

```graphql
# Add to schema.graphqls
directive @auth on FIELD_DEFINITION

type Mutation {
  # ... existing mutations ...
  
  # Only authenticated users can create books
  createBook(input: NewBook!): Book! @auth
}
```

Implement the directive:

```go
// directives/auth.go
package directives

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/yourusername/graphql-example/middleware"
)

func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	// Check if user is authenticated
	_, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, errors.New("not authenticated")
	}
	
	return next(ctx)
}
```

Update the server to use the directive:

```go
// server.go
func main() {
	// ... existing code ...
	
	c := generated.Config{Resolvers: resolvers}
	c.Directives.Auth = directives.Auth
	
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))
	
	// ... existing code ...
}
```

## Integrating with Databases

Let's integrate our GraphQL server with a database using GORM:

```go
// database/db.go
package database

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Book struct {
	ID            string `gorm:"primaryKey"`
	Title         string
	AuthorID      string
	PublishedYear int
	Genres        string // Stored as JSON
	Description   string
}

type Author struct {
	ID   string `gorm:"primaryKey"`
	Name string
	Bio  string
}

type Database struct {
	db *gorm.DB
}

func New() (*Database, error) {
	db, err := gorm.Open(sqlite.Open("library.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// Migrate the schema
	if err := db.AutoMigrate(&Book{}, &Author{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	
	return &Database{db: db}, nil
}

// Book methods
func (d *Database) GetBooks() ([]Book, error) {
	var books []Book
	result := d.db.Find(&books)
	return books, result.Error
}

func (d *Database) GetBook(id string) (Book, error) {
	var book Book
	result := d.db.First(&book, "id = ?", id)
	return book, result.Error
}

func (d *Database) CreateBook(book Book) error {
	return d.db.Create(&book).Error
}

func (d *Database) UpdateBook(book Book) error {
	return d.db.Save(&book).Error
}

func (d *Database) DeleteBook(id string) error {
	return d.db.Delete(&Book{}, "id = ?", id).Error
}

// Author methods
func (d *Database) GetAuthors() ([]Author, error) {
	var authors []Author
	result := d.db.Find(&authors)
	return authors, result.Error
}

func (d *Database) GetAuthor(id string) (Author, error) {
	var author Author
	result := d.db.First(&author, "id = ?", id)
	return author, result.Error
}

func (d *Database) CreateAuthor(author Author) error {
	return d.db.Create(&author).Error
}

func (d *Database) UpdateAuthor(author Author) error {
	return d.db.Save(&author).Error
}

func (d *Database) DeleteAuthor(id string) error {
	return d.db.Delete(&Author{}, "id = ?", id).Error
}

func (d *Database) GetBooksByAuthor(authorID string) ([]Book, error) {
	var books []Book
	result := d.db.Where("author_id = ?", authorID).Find(&books)
	return books, result.Error
}
```

Update the resolver to use the database:

```go
// graph/resolver.go
type Resolver struct {
	db *database.Database
	bookSubscribers map[string]chan *model.Book
}

func NewResolver(db *database.Database) *Resolver {
	return &Resolver{
		db:              db,
		bookSubscribers: make(map[string]chan *model.Book),
	}
}

// Update query resolvers
func (r *queryResolver) Books(ctx context.Context) ([]*model.Book, error) {
	dbBooks, err := r.db.GetBooks()
	if err != nil {
		return nil, err
	}
	
	books := make([]*model.Book, len(dbBooks))
	for i, dbBook := range dbBooks {
		books[i] = &model.Book{
			ID:            dbBook.ID,
			Title:         dbBook.Title,
			AuthorID:      dbBook.AuthorID,
			PublishedYear: dbBook.PublishedYear,
			// Parse genres from JSON
			Description:   dbBook.Description,
		}
	}
	
	return books, nil
}

// Update other resolvers similarly
```

## DataLoader for Batch Loading

DataLoader helps solve the N+1 query problem by batching and caching database queries:

```go
// dataloader/dataloader.go
package dataloader

import (
	"context"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/yourusername/graphql-example/database"
)

type ctxKey string

const loadersKey = ctxKey("dataloaders")

type Loaders struct {
	AuthorLoader *dataloader.Loader
	BooksByAuthorLoader *dataloader.Loader
}

func NewLoaders(db *database.Database) *Loaders {
	authorLoader := dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		ids := make([]string, len(keys))
		for i, key := range keys {
			ids[i] = key.String()
		}
		
		// Fetch all authors in one query
		authors := make(map[string]database.Author)
		for _, id := range ids {
			author, err := db.GetAuthor(id)
			if err == nil {
				authors[id] = author
			}
		}
		
		// Build results in the same order as keys
		results := make([]*dataloader.Result, len(keys))
		for i, key := range keys {
			id := key.String()
			if author, ok := authors[id]; ok {
				results[i] = &dataloader.Result{Data: author, Error: nil}
			} else {
				results[i] = &dataloader.Result{Data: nil, Error: nil}
			}
		}
		
		return results
	}, dataloader.WithClearCacheOnBatch())
	
	booksByAuthorLoader := dataloader.NewBatchedLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		authorIDs := make([]string, len(keys))
		for i, key := range keys {
			authorIDs[i] = key.String()
		}
		
		// Map to store books by author ID
		booksByAuthor := make(map[string][]database.Book)
		
		// Fetch books for each author
		for _, authorID := range authorIDs {
			books, err := db.GetBooksByAuthor(authorID)
			if err == nil {
				booksByAuthor[authorID] = books
			}
		}
		
		// Build results
		results := make([]*dataloader.Result, len(keys))
		for i, key := range keys {
			authorID := key.String()
			if books, ok := booksByAuthor[authorID]; ok {
				results[i] = &dataloader.Result{Data: books, Error: nil}
			} else {
				results[i] = &dataloader.Result{Data: []database.Book{}, Error: nil}
			}
		}
		
		return results
	}, dataloader.WithClearCacheOnBatch())
	
	return &Loaders{
		AuthorLoader:        authorLoader,
		BooksByAuthorLoader: booksByAuthorLoader,
	}
}

// Middleware to add dataloaders to context
func Middleware(db *database.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loaders := NewLoaders(db)
			ctx := context.WithValue(r.Context(), loadersKey, loaders)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Get loaders from context
func GetLoaders(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
```

Update the resolvers to use DataLoader:

```go
// Update book resolver
func (r *bookResolver) Author(ctx context.Context, obj *model.Book) (*model.Author, error) {
	loaders := dataloader.GetLoaders(ctx)
	thunk := loaders.AuthorLoader.Load(ctx, dataloader.StringKey(obj.AuthorID))
	
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	
	dbAuthor := result.(database.Author)
	return &model.Author{
		ID:   dbAuthor.ID,
		Name: dbAuthor.Name,
		Bio:  dbAuthor.Bio,
	}, nil
}

// Update author resolver
func (r *authorResolver) Books(ctx context.Context, obj *model.Author) ([]*model.Book, error) {
	loaders := dataloader.GetLoaders(ctx)
	thunk := loaders.BooksByAuthorLoader.Load(ctx, dataloader.StringKey(obj.ID))
	
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	
	dbBooks := result.([]database.Book)
	books := make([]*model.Book, len(dbBooks))
	for i, dbBook := range dbBooks {
		books[i] = &model.Book{
			ID:            dbBook.ID,
			Title:         dbBook.Title,
			AuthorID:      dbBook.AuthorID,
			PublishedYear: dbBook.PublishedYear,
			// Parse genres from JSON
			Description:   dbBook.Description,
		}
	}
	
	return books, nil
}
```

## Integrating with REST APIs

You can wrap existing REST APIs with GraphQL:

```go
// rest/client.go
package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RESTClient struct {
	baseURL string
	client  *http.Client
}

func NewRESTClient(baseURL string) *RESTClient {
	return &RESTClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type Book struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	AuthorID      string   `json:"author_id"`
	PublishedYear int      `json:"published_year"`
	Genres        []string `json:"genres"`
	Description   string   `json:"description"`
}

type Author struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

func (c *RESTClient) GetBooks() ([]Book, error) {
	resp, err := c.client.Get(c.baseURL + "/books")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	var books []Book
	if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
		return nil, err
	}
	
	return books, nil
}

// Implement other methods for the REST client
```

Update the resolver to use the REST client:

```go
// graph/resolver.go
type Resolver struct {
	restClient *rest.RESTClient
	bookSubscribers map[string]chan *model.Book
}

func NewResolver(restClient *rest.RESTClient) *Resolver {
	return &Resolver{
		restClient:      restClient,
		bookSubscribers: make(map[string]chan *model.Book),
	}
}

// Update query resolvers
func (r *queryResolver) Books(ctx context.Context) ([]*model.Book, error) {
	restBooks, err := r.restClient.GetBooks()
	if err != nil {
		return nil, err
	}
	
	books := make([]*model.Book, len(restBooks))
	for i, restBook := range restBooks {
		books[i] = &model.Book{
			ID:            restBook.ID,
			Title:         restBook.Title,
			AuthorID:      restBook.AuthorID,
			PublishedYear: restBook.PublishedYear,
			Genres:        restBook.Genres,
			Description:   restBook.Description,
		}
	}
	
	return books, nil
}

// Update other resolvers similarly
```

## Best Practices

### Schema Design

1. **Use meaningful types**: Design your schema around your domain model, not your database schema.
2. **Consistent naming**: Use consistent naming conventions (e.g., camelCase for fields, PascalCase for types).
3. **Avoid deeply nested queries**: Limit the depth of your schema to prevent performance issues.
4. **Use interfaces and unions**: For polymorphic responses.
5. **Pagination**: Implement cursor-based pagination for large collections.

### Performance

1. **Use DataLoader**: Batch and cache database queries to avoid the N+1 problem.
2. **Limit query complexity**: Implement query complexity analysis to prevent expensive queries.
3. **Caching**: Cache expensive operations and results.
4. **Query timeouts**: Set timeouts for long-running queries.
5. **Optimize resolvers**: Make resolvers as efficient as possible.

### Security

1. **Input validation**: Validate all input data.
2. **Authentication and authorization**: Implement proper auth checks.
3. **Rate limiting**: Prevent abuse with rate limiting.
4. **Query depth limiting**: Prevent deeply nested queries that could cause performance issues.
5. **Error handling**: Don't expose sensitive information in error messages.

## Exercises

### Exercise 1: Build a GraphQL API
Create a GraphQL API for a blog with posts, comments, and users. Implement queries, mutations, and subscriptions.

### Exercise 2: Database Integration
Integrate your GraphQL API with a database of your choice (PostgreSQL, MySQL, MongoDB).

### Exercise 3: Authentication and Authorization
Implement authentication and authorization in your GraphQL API using JWT tokens.

### Exercise 4: Performance Optimization
Optimize your GraphQL API using DataLoader and caching strategies.

## Key Takeaways

- GraphQL provides a flexible and efficient way to query and manipulate data
- Schema-first development helps create a clear contract between client and server
- Resolvers map schema fields to data sources
- DataLoader solves the N+1 query problem
- Subscriptions enable real-time updates
- Proper error handling and security are essential
- Performance optimization is crucial for production applications

## Next Steps

Next, we'll explore [WebAssembly](30-webassembly.md) to learn how to run Go code in the browser using WebAssembly.