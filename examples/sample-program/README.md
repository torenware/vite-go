# Test Program for the Vite-Go Module

Here's a very simple program to demonstrate the module. To run it, just do:

```shell
$ make run
```

The makefile will do the following:

1. Copy down the latest starter files for Vue 3 and its Vite 2 build tools.
2. Build a `dist/` directory from the starter app.
3. Run the Go web server in `main.go` and post a sample page at [http://localhost:4000](http://localhost:4000).

The files demonstrate the basics of using vite-go in your go program:

* How to embed your `dist/` directory so Go can find it.
* How to initialize the module and load your Vue files.
* How to set up the standard router in Go to find and serve your Vue-related files.

You'll need to have the make utility and npm installed (you do, right?) 
