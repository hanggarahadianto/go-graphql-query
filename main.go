package main

import (
	"fmt"
	"golang-graphql-crud/db"
	"golang-graphql-crud/models"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}


func main(){
	fmt.Println("server running")
	var err  error
	db.Run()

	laptopType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Laptop",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Int),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if laptop, ok := p.Source.(*models.Laptop); ok {
						return laptop.ID, nil
					}
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if laptop, ok := p.Source.(*models.Laptop); ok {
						return laptop.Name, nil
					}
					return nil, nil
				},
			},
			"model": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if laptop, ok := p.Source.(*models.Laptop); ok {
						return laptop.Model, nil
					}
					return nil, nil
				},
			},
			"created_at": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if laptop, ok := p.Source.(*models.Laptop); ok {
						return laptop.CreatedAt, nil
					}
					return nil, nil
				},
			},
		},
	})


	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"laptop": &graphql.Field{
				Type:        laptopType,
				Description: "Get an Laptop.",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)

					laptop := &models.Laptop{}

					err = db.DB.QueryRow("select id, name, model from laptops where id = $1", id).Scan(&laptop.ID, &laptop.Name, &laptop.Model)
					checkErr(err)

					return laptop, nil
				},
			},
			"laptops": &graphql.Field{
				Type:        graphql.NewList(laptopType),
				Description: "List of laptops.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					rows, err := db.DB.Query("SELECT id, name, model FROM laptops")
					checkErr(err)
					var laptops []*models.Laptop

					for rows.Next() {
						laptop := &models.Laptop{}

						err = rows.Scan(&laptop.ID, &laptop.Name, &laptop.Model)
						checkErr(err)
						laptops = append(laptops, laptop)
					}
					return laptops, nil
				},
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createLaptop": &graphql.Field{
				Type:        laptopType,
				Description: "Create new Laptop",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"model": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					name, _ := params.Args["name"].(string)
					model, _ := params.Args["model"].(string)
					createdAt := time.Now()

					var lastInsertId int
					err = db.DB.QueryRow("INSERT INTO laptops(name, model, created_at) VALUES($1, $2, $3) returning id;", name, model, createdAt).Scan(&lastInsertId)
					checkErr(err)

					newLaptop := &models.Laptop{
						ID:        lastInsertId,
						Name:      name,
						Model:     model,
						CreatedAt: createdAt,
					}

					return newLaptop, nil
				},
			},
			"updateLaptop": &graphql.Field{
				Type:        laptopType,
				Description: "Update an author",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"model": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)
					name, _ := params.Args["name"].(string)
					model, _ := params.Args["model"].(string)

					stmt, err := db.DB.Prepare("UPDATE laptops SET name = $1, model = $2 WHERE id = $3")
					checkErr(err)

					_, err2 := stmt.Exec(name, model, id)
					checkErr(err2)

					newLaptop := &models.Laptop{
						ID:    id,
						Name:  name,
						Model: model,
					}

					return newLaptop, nil
				},
			},
			"deleteLaptop": &graphql.Field{
				Type:        laptopType,
				Description: "Delete an laptop",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(int)

					stmt, err := db.DB.Prepare("DELETE FROM laptops WHERE id = $1")
					checkErr(err)

					_, err2 := stmt.Exec(id)
					checkErr(err2)

					return nil, nil
				},
			},
		},
	})
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	http.Handle("/graphql", h)
	http.ListenAndServe(":9090", nil)

}