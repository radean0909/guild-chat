# guild-chat

A simple messaging API for the guild code challenge

## requirements

- go version 1.13 or greater
- git
- restful client to process data

## installation

- Clone the repository: `git clone https://github.com/radean0909/guild-chat.git && cd guild-chat`
- Build the project: `go build ./cmd/main.go -o guild-chat`
- Run the project: `./guild-chat` or `guild-chat.exe`

## usage

- using curl, postman, or something similar make requests to the api that is now running on localhost:8000
- refer to the documentation by running `go-swagger serve ./swagger/guild-chat

## notes and improvements

I made several design decisions when implementing this code challenge that should be considered when reading and writing the code.
I used Golang. Part of this is due to familiarity and recent use, though such an app could be developed using many different languages without significant change to my methodology.

Regarding that methodology:

### the database
- I implemented a Driver as a data abstraction layer. This driver is an interface that allows to easily reproduce the same results, regardless of datastore implemented.
- I leveraged a simple in-memory db mock as part of my implementation. This is the same approach I usually take when first starting a project so that I am able to quickly get a local environment set up without relying on other resources.
- The in-memory datastore also allows for the features of the application to be developed rapidly without worrying about strictly defining a schema or datastructure early on. Later, once things have been proven, it is simple to add another driver to connect to whatever database is ultimately used
- Finally, instead of relying on external db mocking tools or standing up external resources, seeding data, etc, unit tests can be accomplished with the in-memory store. 

### the testing story
- The service relies on unit tests at the 'db' level (in-memory store), at the handler level, and for important utility functions
- In a production app, an end-to-end runner would be ideal for regression testing
- Integration testing can occur by standing up the service and excercising the endpoints directly

### improvements
- leverage containerization through docker/k8s or a serverless architecture. I opted to avoid this at present, because, though critical, adding complexities to app at this stage doesn't reveal much, and ratchets up the complexity
- add a websocket implementation. For a truly effective, robust, and modern approach to chat, using websockets as opposed to a RESTful API (or in addition to) is much prefered. 
- minor improvements are noted in comments throughout the code

## routes

### messages

#### GET /message/:id

Gets a single message by id. Errors if message not found or id missing.

On success returns Message JSON

``` JSON
{
    sender: uuid,
    recipient: uuid,
    date: date,
    message: string,
}
```

On failure returns error JSON

Returns: 200, 400, 404, 500

#### POST /message

Creates a single message from JSON body content (application/json)
Errors if message is missing sender, recipient, or message string.

Input body
``` JSON
{
    sender: uuid,
    recipient: uuid,
    message: string,
}
```

On success returns Message JSON
``` JSON
{
    id: uuid,
    sender: uuid,
    recipient: uuid,
    date: date,
    message: string,
}
```

On failure returns error JSON

Returns: 200, 400, 404, 500

### users

#### GET /user/:id

Retrieves a single user by id. Errors if user is not found or id is missing.

On success returns User JSON

``` JSON
{
    id: uuid,
    username: string,
    email: string
}
```

On error, returns an error object

Returns: 200, 400, 404, 500

#### POST /user

Creates a new user from JSON body content (application/json)
Errors if missing username or email, or if username is taken

Input body:
``` JSON
{
    username: string,
    email: string
}
```

On success returns User JSON

``` JSON
{
    id: uuid,
    username: string,
    email: string
}
```

On error, returns an error object

Returns: 200, 400, 404, 500

#### DELETE /user/:id

(Soft) Deletes a single user. Errors if userid cannot be found

On success returns no content

On error, returns an error object

Returns: 200, 400, 404, 500

### conversations

### GET /conversation/:to/:from?start=YYYY-MM-DD&until=YYYY-MM-DD&limit=100

Gets a conversation between two users. If there are no messages sent to the recipient, returns 404.

Params: 
- to - path - uuid
- from - path - uuid
- start - query - date in YYYY-MM-DD format for earliest message (defaults to 30 days ago)
- until - query - date in YYYY-MM-DD format for most recent message (defaults to now)
- limit - query - maximum number of messages to return, defaults to 100

On success returns an array of message JSON object. Only includes messages sent *to* the recipient. If the sending user is deleted, redacts uuid with `deleted`

``` JSON
[
    {
        id: uuid,
        sender: uuid,
        recipient: uuid,
        message: string,
        date: date
    }, 
    ...
]
```

On error returns error message

Returns: 200, 404, 500

### GET /conversation/:to?start=YYYY-MM-DD&until=YYYY-MM-DD&limit=100

Gets all conversations sent to a user.

Params: 
- to - path - uuid
- start - query - date in YYYY-MM-DD format for earliest message (defaults to 30 days ago)
- until - query - date in YYYY-MM-DD format for most recent message (defaults to now)
- limit - query - maximum number of messages to return, defaults to 100

On success returns an array of message JSON objects. Only includes messages sent *to* the recipient. If the sending user is deleted, redacts uuid with `deleted`

``` JSON
[
    {
        id: uuid,
        sender: uuid,
        recipient: uuid,
        message: string,
        date: date
    }, 
    ...
]

```

On error returns error message

Returns: 200, 404, 500