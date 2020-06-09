package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/milamice62/terraplugin/api/client"
)

func genreItem() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the genre, also acts as it's unique ID",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
		},
		Create: createGenre,
		Read:   readGenre,
		Delete: deleteGenre,
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

	err := apiClient.NewGenre(&genre)

	if err != nil {
		return err
	}
	d.SetId(genre.Name)
	return nil
}

func readGenre(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	genreName := d.Id()
	genre, err := apiClient.GetGenre(genreName)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Genre with name %s", genreName)
		}
	}

	d.SetId(genre.Name)
	d.Set("name", genre.Name)
	return nil
}

func deleteGenre(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	genreName := d.Id()

	err := apiClient.DeleteGenre(genreName)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
