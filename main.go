package main

import "fetch_rewards/routes"

func main() {

	r := routes.SetUpRouter()
	r.Run()

}
