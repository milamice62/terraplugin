package provider

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/milamice62/terraplugin/api/client"
)

func GenreItem() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the genre",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
		},
		Create: createGenre,
		Read:   readGenre,
		Delete: deleteGenre,
		Exists: existGenre,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func createGenre(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	genre := client.Genre{
		Name: d.Get("name").(string),
	}

	resBody, err := apiClient.NewGenre(&genre)

	if err != nil {
		return err
	}

	err = json.NewDecoder(resBody).Decode(&genre)
	if err != nil {
		return err
	}

	d.SetId(genre.ID)
	return nil
}

func readGenre(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	genreID := d.Id()
	genre, err := apiClient.GetGenre(genreID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Genre with name %s", genreID)
		}
	}

	d.SetId(genre.ID)
	if d.Set("name", genre.Name); err != nil {
		return err
	}
	return nil
}

func existGenre(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	genreID := d.Id()
	_, err := apiClient.GetGenre(genreID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func deleteGenre(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	genreID := d.Id()

	err := apiClient.DeleteGenre(genreID)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
