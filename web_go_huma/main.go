package main

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/gofiber/fiber/v2"
)

// Options for the CLI. Pass `--port` or set the `SERVICE_PORT` env var.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func main() {
	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Create a new router & API
		router := fiber.New()
		api := humafiber.New(router, huma.DefaultConfig("My API", "1.0.0"))

		// Add the operation handler to the API.
		huma.Get(api, "/greeting/{name}", func(ctx context.Context, input *struct {
			Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
		},
		) (*GreetingOutput, error) {
			resp := &GreetingOutput{}
			resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
			return resp, nil
		})

		// Tell the CLI how to start your router.
		hooks.OnStart(func() {
			router.Listen(fmt.Sprintf(":%d", options.Port))
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
