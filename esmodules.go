// / ESModules bundles the project.
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

func main() {
	var (
		watchMode  = flag.Bool("watch", false, "Watch files for changes")
		withServer = flag.Bool("server", false, "Server HTTP files from public")
		port       = flag.Int("port", 8080, "HTTP server port")
	)
	flag.Parse()

	// Build ES Modules
	ctx, buildErr := esbuild.Context(esbuild.BuildOptions{
		EntryPoints: []string{"src/app.ts"},
		Bundle:      true,
		Format:      esbuild.FormatESModule,
		Outdir:      "public",
		Write:       true,
		Engines: []esbuild.Engine{
			{Name: esbuild.EngineChrome, Version: "64"},
			{Name: esbuild.EngineEdge, Version: "16"},
			{Name: esbuild.EngineFirefox, Version: "60"},
			{Name: esbuild.EngineSafari, Version: "11"},
		},
	})
	defer ctx.Dispose()

	if buildErr != nil {
		panic(buildErr)
	}

	// Build once
	result := ctx.Rebuild()

	if len(result.Errors) > 0 {
		for _, err := range result.Errors {
			println(err.Text)
		}

		os.Exit(1)
	}

	// Find HTML file and copy it to the public directory
	htmlSource, err := os.Open("index.html")
	if err != nil {
		panic(err)
	}
	defer htmlSource.Close()

	destination, err := os.Create("public/index.html")
	if err != nil {
		panic(err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, htmlSource)
	if err != nil {
		panic(err)
	}

	if *watchMode {
		err := ctx.Watch(esbuild.WatchOptions{})
		if err != nil {
			panic(err)
		}
		println("watching for changes...")

		// Block I/O if the server is not running
		if !*withServer {
			<-make(chan struct{})
		}
	}

	// Open a FileServer
	if *withServer {
		fmt.Printf("starting server on :%d\n", *port)
		err = http.ListenAndServe(fmt.Sprintf(":%d", *port), http.FileServer(http.Dir("public")))
		if err != nil {
			panic(err)
		}
	}
}
