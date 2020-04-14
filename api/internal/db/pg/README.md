# Postgres DB Driver Implementation
Because I wanted to focus on the logic within the service, solid unit testing, and to illustrate the ability for the service to stand up without needing a specific backend (and also due to time), I have opted to not implement a database driver for this project.

This implementation, however, allows for easily adding drivers for any datastores easily by following the interface described in `db.go`.

## Sample SQL Queries
I do want to at least illustrate that I do, in fact, know SQL, so, in lieu of a true implementation, I will provide the query structure I would use for different driver methods

### GetMessage

`SELECT * FROM messages WHERE id = $1;`

### CreateMessage

- Creating a message should also create a conversation, if one doesn't exist
- Creating a message should also add the message to the conversation using an UPDATE request

`INSERT INTO messages (id, sender, recipient, message, date) VALUES ( $1, $2, $3, $4, $5 );`

### ListMessages

`SELECT * FROM messages WHERE recipient=$1 AND date BETWEEN $2 AND $3 LIMIT $4;`

### GetConversation

- Would expect a single result

`SELECT * FROM conversations WHERE (sender = $1 AND recipient = $2) OR (recipient = $1 AND sender =$2) AND updated BETWEEN $3 AND $4;`

### CreateConversation

- Conversation Messages would be a Postgres Array of strings (uuids)
- Conversation shouldn't create a new conversation if one already exists (but is reversed). This could be handled with two statements (a select then an insert, if appropriate) but could also be done with a trickier single statement. Depending on the use cases, I would likely opt for the simpler two-statement approach as it is easier to read, but to prevent a race condition a lock would need to be put on the table during the select

`INSERT INTO conversations (id, sender, recipient, updated, messages) VALUES ($1, $2, $3, $4, $5);`

### ListConversations

`SELECT * FROM conversations WHERE reciepient = $1 AND updated BETWEEN $2 AND $3 LIMIT $4 ORDER BY sender;`

### GetUser

`SELECT (id, email, username) FROM users WHERE id = $1 AND archived_on IS NULL;`

### DeleteUser

`UPDATE users SET archived_on = $2 WHERE id = $1;`