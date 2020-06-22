package provider

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/milamice62/terraplugin/api/client"
)

func RentalItem() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"dateout": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of checkout",
				ForceNew:    true,
			},
			"customer": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The customer information",
				MaxItems:    1,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the customer",
							ForceNew:    true,
						},
						"isgold": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The status of the customer",
							ForceNew:    true,
						},
						"phone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The phone number of customer",
							ForceNew:    true,
						},
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the customer",
							ForceNew:    true,
						},
					}},
			},
			"movie": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The movie information",
				MaxItems:    1,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dailyrentalrate": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The daily rental rate of the movie",
							ForceNew:    true,
						},
						"title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The title of the movie",
							ForceNew:    true,
						},
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The id of the movie",
							ForceNew:    true,
						},
					}},
			},
		},
		Create: createRental,
		Read:   readRental,
		Delete: deleteRental,
		Exists: existRental,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func expandCustomer(customer []interface{}) (*client.Customer, error) {
	cus := &client.Customer{}

	if len(customer) == 0 || customer[0] == nil {
		return nil, fmt.Errorf("Customer parameters should be specified, but get %v", customer)
	}

	in, ok := customer[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Error fetching customer element: %v", in)
	}

	cus.CustomerID = in["id"].(string)

	return cus, nil
}

func flattenCustomer(customer *client.Customer, d *schema.ResourceData) []interface{} {
	m := make(map[string]interface{})
	m["id"] = customer.CustomerID
	m["isgold"] = customer.IsGold
	m["name"] = customer.Name
	m["phone"] = customer.Phone

	return []interface{}{m}
}

func expandMovie(movie []interface{}) (*client.Movie, error) {
	mov := &client.Movie{}

	if len(movie) == 0 || movie[0] == nil {
		return nil, fmt.Errorf("Movie parameters should be specified, but get %v", movie)
	}

	in, ok := movie[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Error fetching movie element: %v", in)
	}

	mov.MovieID = in["id"].(string)

	return mov, nil
}

func flattenMovie(movie *client.Movie, d *schema.ResourceData) []interface{} {
	m := make(map[string]interface{})
	m["id"] = movie.MovieID
	m["dailyrentalrate"] = movie.Rate
	m["title"] = movie.Title

	return []interface{}{m}
}

func updateRental(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	rentalID := d.Id()
	rental, err := apiClient.GetRental(rentalID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding rental with id %s", rentalID)
		}
	}

	if d.HasChange("rentalid") {
		p := d.Get("rentalid").(string)
		rental.RentalID = p
	}

	err = apiClient.UpdateRental(rental)
	if err != nil {
		return err
	}

	return readRental(d, m)
}

func createRental(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	rental := client.Rental{}

	customer, err := expandCustomer(d.Get("customer").([]interface{}))
	if err != nil {
		fmt.Printf("%v", err)
	}

	movie, err := expandMovie(d.Get("movie").([]interface{}))
	if err != nil {
		fmt.Printf("%v", err)
	}

	rentalID := client.RentalID{
		CustomerID: customer.CustomerID,
		MovieID:    movie.MovieID,
	}

	resBody, err := apiClient.NewRental(&rentalID)

	if err != nil {
		return err
	}

	err = json.NewDecoder(*resBody).Decode(&rental)
	if err != nil {
		return err
	}

	d.SetId(rental.RentalID)
	customerRes := flattenCustomer(rental.Customer, d)
	movieRes := flattenMovie(rental.Movie, d)

	d.SetId(rental.RentalID)
	if d.Set("customer", customerRes); err != nil {
		return err
	}
	if d.Set("movie", movieRes); err != nil {
		return err
	}
	if d.Set("dateout", rental.DateOut); err != nil {
		return err
	}
	return nil
}

func readRental(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	rentalID := d.Id()
	rental, err := apiClient.GetRental(rentalID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding rental with id %s", rentalID)
		}
	}

	customer := flattenCustomer(rental.Customer, d)
	movie := flattenMovie(rental.Movie, d)

	d.SetId(rental.RentalID)
	if d.Set("customer", customer); err != nil {
		return err
	}
	if d.Set("movie", movie); err != nil {
		return err
	}
	if d.Set("dateout", rental.DateOut); err != nil {
		return err
	}
	return nil
}

func existRental(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	rentalID := d.Id()
	_, err := apiClient.GetRental(rentalID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func deleteRental(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	rentalID := d.Id()

	err := apiClient.DeleteRental(rentalID)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
