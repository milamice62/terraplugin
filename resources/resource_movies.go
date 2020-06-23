package provider

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/milamice62/terraplugin/api/client"
)

func MovieItem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"title": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The movie title",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"genre": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The movie genre",
				MaxItems:    1,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the genre",
							ForceNew:    true,
						},
						"_id": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The id of the genre",
							ForceNew:     true,
							ValidateFunc: validateName,
						},
					}},
			},
			"stock": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The movie stock",
				ForceNew:     true,
				ValidateFunc: validateInt,
			},
			"daily_rate": {
				Type:         schema.TypeFloat,
				Required:     true,
				Description:  "The movie daily rental rate",
				ForceNew:     true,
				ValidateFunc: validateFloat,
			},
		},
		Create: createMovie,
		Read:   readMovie,
		Delete: deleteMovie,
		Exists: existMovie,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func expandGenre(genre []interface{}) (*client.Genre, error) {
	gen := &client.Genre{}

	if len(genre) == 0 || genre[0] == nil {
		return nil, fmt.Errorf("Genre should be specified, but get %v", genre)
	}

	in, ok := genre[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Error fetching genre element: %v", in)
	}

	gen.Name = in["name"].(string)
	gen.ID = in["_id"].(string)
	return gen, nil
}

func flattenGenre(movie *client.Movie, d *schema.ResourceData) []interface{} {
	m := make(map[string]interface{})
	m["name"] = movie.Genre.Name
	m["_id"] = movie.Genre.ID

	return []interface{}{m}
}

func createMovie(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	genre, err := expandGenre(d.Get("genre").([]interface{}))
	if err != nil {
		fmt.Printf("%v", err)
	}

	movie := client.Movie{}
	movie.Title = d.Get("title").(string)
	movie.Stock = d.Get("stock").(int)
	movie.Rate = d.Get("daily_rate").(float64)
	movie.Genre = genre

	body, err := apiClient.NewMovie(&movie)

	if err != nil {
		return err
	}

	err = json.NewDecoder(*body).Decode(&movie)
	if err != nil {
		return err
	}

	d.SetId(movie.MovieID)
	gen := flattenGenre(&movie, d)
	d.Set("genre", gen)

	return nil
}

func readMovie(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	movieID := d.Id()
	movie, err := apiClient.GetMovie(movieID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding movie with id %s", movieID)
		}
	}

	genre := flattenGenre(movie, d)
	d.SetId(movieID)
	d.Set("title", movie.Title)
	d.Set("daily_rate", movie.Rate)
	d.Set("stock", movie.Stock)
	d.Set("genre", genre)

	return nil
}

func existMovie(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	movieID := d.Id()
	_, err := apiClient.GetMovie(movieID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func deleteMovie(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	movieID := d.Id()

	err := apiClient.DeleteMovie(movieID)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
