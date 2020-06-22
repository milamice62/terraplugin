provider "store" {
  address = "http://localhost"
  port    = "3000"
  token   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJfaWQiOiI1ZWUwNWY3M2QyZWZhMmFlODU4MGY2Y2QiLCJpc0FkbWluIjp0cnVlLCJpYXQiOjE1OTE3NjI4MDN9.0tgg4WTMp67orDP_v5haxPYmN5NVjXIiyMMpYLNbspw"
}

# resource "store_genres" "kind" {
#   name = "comedy"
# }

# resource "store_movies" "saw" {
#   title = "sawIII"
#   genre {
#     _id  = "5ee19f2a1363f7c0493761e9"
#     name = "hhhhh"
#   }
#   stock      = 10
#   daily_rate = 12.10
# }

# resource "store_customers" "cici" {
#   name  = "selina"
#   phone = "123456789"
# }

resource "store_rentals" "myrental" {
  customer {
    id = "5ee998a7073cfb0d8696fec1"
  }
  movie {
    id = "5ee6fe17de7e8d5eb0ae60ea"
  }
}

# resource "example_item" "test" {
#   name        = "this_is_an_item"
#   description = "this is an item"
#   tags = [
#     "hello",
#     "world"
#   ]
# }
