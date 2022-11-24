package main

import "fetch_rewards/routes"

// Main function //
func main() {

	r := routes.SetUpRouter()
	r.Run()

}
