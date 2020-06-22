package provider

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/milamice62/terraplugin/api/client"
)

func CustomerItem() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the customer",
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"isgold": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "The status of the customer",
				ForceNew:    true,
			},
			"phone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The phone number of customer",
			},
		},
		Create: createCustomer,
		Read:   readCustomer,
		Delete: deleteCustomer,
		Exists: existCustomer,
		Update: updateCustomer,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func updateCustomer(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	customerID := d.Id()
	customer, err := apiClient.GetCustomer(customerID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding customer with id %s", customerID)
		}
	}

	if d.HasChange("phone") {
		p := d.Get("phone").(string)
		customer.Phone = p
	}

	err = apiClient.UpdateCustomer(customer)
	if err != nil {
		return err
	}

	return readCustomer(d, m)
}

func createCustomer(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	customer := client.Customer{
		Name:  d.Get("name").(string),
		Phone: d.Get("phone").(string),
	}

	resBody, err := apiClient.NewCustomer(&customer)

	if err != nil {
		return err
	}

	err = json.NewDecoder(*resBody).Decode(&customer)
	if err != nil {
		return err
	}

	d.SetId(customer.CustomerID)
	return nil
}

func readCustomer(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	customerID := d.Id()
	customer, err := apiClient.GetCustomer(customerID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding customer with id %s", customerID)
		}
	}

	d.SetId(customer.CustomerID)
	if d.Set("name", customer.Name); err != nil {
		return err
	}
	if d.Set("phone", customer.Phone); err != nil {
		return err
	}
	if d.Set("isGold", customer.IsGold); err != nil {
		return err
	}
	return nil
}

func existCustomer(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	customerID := d.Id()
	_, err := apiClient.GetCustomer(customerID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func deleteCustomer(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	customerID := d.Id()

	err := apiClient.DeleteCustomer(customerID)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
