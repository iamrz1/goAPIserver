# APIServer for Bookstore using using CLI [Cobra](https://github.com/spf13/cobra)
##API endpints for CRUD functions
## Example commands to run

- The following command will run the book server at port 8080 and it requires the authentication from user,
  `go run main.go`
- Login Details (Postman)
  `username : tom95 , pass : pass1`
- Login can be bypassed with -b flag
  `go run main.go  --bypassLogin=true`
  or
  `go run main.go  -b`


- Change the default port (8080) using --port flag

  `go run main.go --port=8081`

  or
  
  `go run main.go --port=8081 --bypassLogin=true`