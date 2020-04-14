# guild-chat

A simple messaging API for the guild code challenge

## requirements

- go version 1.13 or greater
- swagger server (such as go-swagger)

## installation

- Clone the repository: `git clone https://github.com/radean0909/guild-chat.git && cd guild-chat`
- Run the build script: `./scripts/build.sh`
- Run the start script: `./scriptw/start.sh`

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

```{
    sender: uuid,
    recipient: uuid,
    date: date,
    message: string,
}```

On failure returns error JSON

Returns: 200, 400, 404, 500

#### POST /message

Creates a single message. Errors if message is missing sender, recipient, or message string.
On success returns Message JSON
{
    id: uuid,
    sender: uuid,
    recipient: uuid,
    date: date,
    message: string,
}

On failure returns error JSON

Returns: 200, 400, 404, 500