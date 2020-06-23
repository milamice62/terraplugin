package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/milamice62/terraplugin/api/client"
)

func Test_Customer_Init(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCustomerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCustomerInit(), // equal to 'Terraform Apply'
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleCustomerExists("store_customers.customer1"),
					resource.TestCheckResourceAttr(
						"store_customers.customer1", "name", "foobar"),
					resource.TestCheckResourceAttr(
						"store_customers.customer1", "phone", "123456789"),
				),
			},
		},
	})
}

func Test_Customer_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCustomerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCustomerInit(), // equal to 'Terraform Apply'
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleCustomerExists("store_customers.customer1"),
					resource.TestCheckResourceAttr(
						"store_customers.customer1", "name", "foobar"),
					resource.TestCheckResourceAttr(
						"store_customers.customer1", "phone", "123456789"),
				),
			},
			{
				Config: testAccCheckCustomerUpdate(), // equal to 'Terraform Apply'
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleCustomerExists("store_customers.customer1"),
					resource.TestCheckResourceAttr(
						"store_customers.customer1", "name", "barfoo"),
					resource.TestCheckResourceAttr(
						"store_customers.customer1", "phone", "987654321"),
				),
			},
		},
	})
}

func testAccCheckCustomerDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "store_customers" {
			continue
		}

		_, err := apiClient.GetCustomer(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Alert! genre still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}

func testAccCheckExampleCustomerExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		id := rs.Primary.ID
		apiClient := testAccProvider.Meta().(*client.Client)
		_, err := apiClient.GetCustomer(id)
		if err != nil {
			return fmt.Errorf("error fetching customer with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckCustomerInit() string {
	return fmt.Sprintf(`
resource "store_customers" "customer1" {
  name = "foobar"
  phone = "123456789"
}
`)
}

func testAccCheckCustomerUpdate() string {
	return fmt.Sprintf(`
resource "store_customers" "customer1" {
  name = "barfoo"
  phone = "987654321"
}
`)
}
