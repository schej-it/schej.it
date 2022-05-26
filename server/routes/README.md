# routes
This package contains all the routes for the schej.it API

To view the docs, visit http://localhost:3002/swagger/index.html

## How to document routes
Visit https://github.com/swaggo/swag for a comprehensive overview of the swagger comment structure

To generate swagger docs, make sure you have swag installed:
```
go install github.com/swaggo/swag/cmd/swag@v1.6.7
```

Then, in the root directory run `swag init` every time you make a change for it to appear in the docs