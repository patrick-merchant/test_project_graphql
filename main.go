package main

import (
	// "errors"
	"errors"
	"fmt"
	// "strconv"
	"io/ioutil"
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

var postCount = len(posts)

var (
	opts  = []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	posts = []Post{
		{
			ID:     "0-post",
			UserID: "0-user",
			Title:  "First Post",
			Body:   "Here is my first post - please work.",
		},
	}
)

type Post struct {
	ID     graphql.ID
	UserID graphql.ID
	Title  string
	Body   string
}

// INDEX
func (r *RootResolver) Feed() ([]Post, error) {
	return posts, nil
}

// TEST STRING
func (r *RootResolver) Info() (string, error) {
	return "this is a thing", nil // OR newError etc...
}

// GET BY ID
func (r *RootResolver) Search(i struct{ ID graphql.ID }) (Post, error) {
		
	for _, post := range posts {
		if post.ID == i.ID {
			return post, nil
		}
	}

	// Error handling
	
	return posts[0], errors.New("No post with ID " + string(i.ID))

}

// CREATE 
func (r *RootResolver) Post(args struct {
	UserID graphql.ID
	Title  string
	Body   string
}) (Post, error) {
	newPost := Post{
		ID:     graphql.ID(fmt.Sprint(postCount) + "-post"),
		UserID: args.UserID + "-user",
		Title:  args.Title,
		Body:   args.Body,
	}

	posts = append(posts, newPost)
	postCount++
	return newPost, nil
}

// DELETE
func (r *RootResolver) Delete(i struct{ ID graphql.ID }) ([]Post, error) {
	deleted := false
	
	for index, post := range posts {
		if post.ID == i.ID {
			posts = append(posts[:index], posts[index+1:]...)
			deleted = true
		}
	}

	// Error handling
	if deleted {
		return posts, nil
	} else {
		return posts, errors.New("No post with ID " + string(i.ID))
	}
}

// UPDATE 
func (r *RootResolver) Update(args struct {
	ID graphql.ID
	UserID graphql.ID
	Title  string
	Body   string
}, 
	) (Post, error) {

	updated := false

	// Store updated post info in variable
	updatedPost := Post{
		ID:     args.ID,
		UserID: args.UserID + "-user",
		Title:  args.Title,
		Body:   args.Body,
	}

	// loop through posts and replace old with updated post
	for index, post := range posts {
		if post.ID == args.ID {
			posts[index] = updatedPost
			updated = true
		}
	}

	// Error handling
	if updated {
		return updatedPost, nil
	} else {
		return updatedPost, errors.New("No post with ID " + string(args.ID))
	}
}

type RootResolver struct{}

// Reads and parses the schema from file - Associates root resolver. Panics if can't read.
func parseSchema(path string, resolver interface{}) *graphql.Schema {
	bstr, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	schemaString := string(bstr)
	parsedSchema, err := graphql.ParseSchema(
		schemaString,
		resolver,
		opts...,
	)
	if err != nil {
		panic(err)
	}
	return parsedSchema
}

func main() {
	http.Handle("/graphql", &relay.Handler{
		Schema: parseSchema("./schema.graphql", &RootResolver{}),
	})

	fmt.Println("serving on 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
