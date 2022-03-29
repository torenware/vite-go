# Test Program for the Vite-Go Module

Here's a very simple program to demonstrate the module. To run it, first add your Vue files. To automatically install the Vue "starter project" on Mac or Linux, you can run the command. The command will put Vue's various pieces in the right places for the test program.

```shell
$ ./install-vue.sh
```

Once you have your Vue files installed, just run the make script to build your Vue project, compile the test program and run it:

```shell
$ make run
```

The makefile will do the following:

1. Build a `dist/` directory from the starter app.
2. Run the Go web server in `main.go` and post a sample page at [http://localhost:4000](http://localhost:4000).

The files demonstrate the basics of using vite-go in your go program:

* How to embed your `dist/` directory so Go can find it.
* How to initialize the module and load your Vue files.
* How to set up the standard router in Go to find and serve your Vue-related files.

You'll need to have the make utility and npm installed (you do, right?) 
