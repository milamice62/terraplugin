package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/milamice62/terraplugin/api/client"
)

func Test_Movie_Init(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMovieDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMovieInit(), // equal to 'Terraform Apply'
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleMovieExists("store_movies.movie_example"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "title", "example"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "stock", "100"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "daily_rate", "10"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "genre.#", "1"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "genre.0._id", "5ee19f2a1363f7c0493761e9"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "genre.0.name", "hhhhh"),
				),
			},
		},
	})
}

func Test_Movie_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMovieDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMovieInit(), // equal to 'Terraform Apply'
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleMovieExists("store_movies.movie_example"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "title", "example"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "stock", "100"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "daily_rate", "10"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "genre.#", "1"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "genre.0._id", "5ee19f2a1363f7c0493761e9"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "genre.0.name", "hhhhh"),
				),
			},
			{
				Config: testAccCheckMovieUpdate(), // equal to 'Terraform Apply'
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleMovieExists("store_movies.movie_example"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "title", "example"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "stock", "10"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "daily_rate", "11.1"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "genre.#", "1"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "genre.0._id", "5ee05b02340e2cae12c1bea5"),
					resource.TestCheckResourceAttr(
						"store_movies.movie_example", "genre.0.name", "sci-fic"),
				),
			},
		},
	})
}

func testAccCheckMovieDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "store_movies" {
			continue
		}

		_, err := apiClient.GetMovie(rs.Primary.ID)
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

func testAccCheckExampleMovieExists(resource string) resource.TestCheckFunc {
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
		_, err := apiClient.GetMovie(id)
		if err != nil {
			return fmt.Errorf("error fetching movie with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckMovieInit() string {
	return fmt.Sprintf(`
	resource "store_movies" "movie_example" {
		title = "example"
		genre {
		  _id  = "5ee19f2a1363f7c0493761e9"
		}
		stock      = 100
		daily_rate = 10.00
	  }
`)
}

func testAccCheckMovieUpdate() string {
	return fmt.Sprintf(`
	resource "store_movies" "movie_example" {
		title = "example"
		genre {
		  _id  = "5ee05b02340e2cae12c1bea5"
		}
		stock      = 10
		daily_rate = 11.10
	  }
`)
}
