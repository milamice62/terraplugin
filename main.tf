provider "movie" {
  address = "http://localhost"
  port    = "3000"
  token   = "superSecretToken"
}

resource "movie_genres" "comedy" {
  name = "comedy"
}

# resource "example_item" "test" {
#   name        = "this_is_an_item"
#   description = "this is an item"
#   tags = [
#     "hello",
#     "world"
#   ]
# }
