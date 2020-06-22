package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/milamice62/terraplugin/api/client"
)

func testAccCheckGenreDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "store_genres" {
			continue
		}

		_, err := apiClient.GetGenre(rs.Primary.ID)
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
func Test_Genre_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGenreDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGenreBasic(), // equal to 'Terraform Apply'
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleGenreExists("store_genres.kind"),
					resource.TestCheckResourceAttr(
						"store_genres.kind", "name", "comedy"),
				),
			},
		},
	})
}

func testAccCheckGenreBasic() string {
	return fmt.Sprintf(`
resource "store_genres" "kind" {
  name = "comedy"
}
`)
}

func testAccCheckExampleGenreExists(resource string) resource.TestCheckFunc {
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
		_, err := apiClient.GetGenre(id)
		if err != nil {
			return fmt.Errorf("error fetching genre with resource %s. %s", resource, err)
		}
		return nil
	}
}
