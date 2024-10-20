# Todo CLI app

This todo cli app is for learning purpose of go and cobra cli

## Authors

- [@sachinrathore]https://github.com/12sachin95/

## Run Locally

Clone the project

```bash
  git clone https://github.com/12sachin95/go-todo-cli-app
```

Go to the project directory

```bash
  cd go-todo-cli-app
```

Install dependencies

```bash
  go mod tidy
```

Start the server

```bash
  go run main.go serve
```

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`MONGODB_URI= Your mongo db uri here`

`DATABASE_NAME=go-todo-db`

`SECRET_KEY=your secret key here for jwt`

`PORT=8080`

`TODO_SERVER_PATH=http://localhost:8080/todo-app/api/v1`

## Demo

#### Start server

`go mod tiddy`

`go run main.go serve`

#### Available comands-

Register user

`go run main.go user register  username password`

Login User

`go run main.go user login  username password`

Logout User

`go run main.go user logout --user_id userId`

Get User details

`go run main.go user details --user_id userId`

Create Todo

`go run main.go todo create --user_id userId --title title1 --completed=true`

Get all todos

`go run main.go todo get todoId --user_id userId`

Get one Todo

`go run main.go todo getOne todoId todoId --user_id userId`

Update Todo

`go run main.go todo update todoId --user_id userId --title title1 --completed=true`

Delete Todo

`go run main.go todo delete todoId --user_id userId`

## Build and Run

#### Build:-

This name might be anything- todo-cli
`go build -o todo-cli`

#### Run the CLI Application:-

`./myapp`

Example:- create todo

`./todo-cli todo create --user_id 6713aace3b5297a43130c713 --title "Hello world"`

The ./ is required because macOS does not automatically look for executables in the current directory unless you specify the path.
The todo-cli is the name of the built binary.
todo create is the CLI command you're running, and --user_id "someUserID" is the flag you're passing.

#### Making the CLI Accessible Globally

If you want to run the CLI without specifying ./ each time, you can move the built binary to a directory in your $PATH, such as /usr/local/bin.

`sudo mv todo-cli /usr/local/bin/`

Now you can run the CLI from anywhere in your terminal without needing ./:

`todo-cli todo get --user_id "someUserID"`

## Tech Stack

**Client:** React

**Server:** Go, cobra, gin, mongodb
