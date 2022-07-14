package main

func main() {
	app := NewWebServer().
		GetWebServer()

	app.Listen(":3000")
}
