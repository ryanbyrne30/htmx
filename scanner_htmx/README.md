# Tailwind, Go and HTMX

## Setting up Tailwind with Go

Accorinding to the [Tailwind](https://tailwindcss.com/docs/installation) docs, we can use the Tailwind cli to build our output.css file used in the Go application.

To install the Tailwind cli, you can either install following the [docs](https://tailwindcss.com/docs/installation) or download the [Tailwind binary](https://tailwindcss.com/blog/standalone-cli).

Once installed you will run
`tailwindcss -i input.css -o output.css --watch`.

For example, if our css files are stored in `/static/styles` you would run

```
tailwindcss -i static/styles/input.css -o static/styles/output.css --watch
```

You can then setup a net/http FileServer to expose the files in `/static` allowing the Go application to import the styles.

```
<head>
  <link href="/static/styles/output.css" rel="stylesheet">
</head>
```