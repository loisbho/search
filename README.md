# Search

Web server that searches the data in the json files (tickets.json, organizations.json, users.json)
and returns the results in a human-readable format. 

## Approach

1) import the data in the json file into sql database upon starting the server
2) fetch the data from the sql database
3) display the data in json format 

Alternative solutions could be using elasticsearch or a data structure to hold the data.
This can implemented by updating the dao files with the new implementation.  The interface
wouldn't be changed. 

This is still a work in progress. Given more time, I would improve error handling and include more 
tests.

## Getting Started

To run the program, download the repo and run the command in the project directory.
```
go run main.go app
```

For help

```
go run main.go
```

### API Endpoints

Endpoints can be found in web/server.go.  


Search for field in organization.
```
/organizations?field={field}&value={value}
```

Search for field in users. 
```
/users?field={field}&value={value}
```

Search for field in tickets.
```
/tickets?field={field}&value={value}
```


List searchable fields for organization.
```
/org-details
```

List searchable fields for users.
```
/user-details
```

List searchable fields for tickets.
```
/ticket-details
```


### Example

Sample query

```
http://localhost:3000/organizations?field=id&value=1
```


Response
```
[
{
"_id": 1,
"url": "http://initech.zendesk.com/api/v2/users/1.json",
"external_id": "74341f74-9c79-49d5-9611-87ef9b6eb75f",
"name": "Francisca Rasmussen",
"alias": "Miss Coffey",
"created_at": "2020-02-05 12:47:36.286052 -0800 -0800",
"active": true,
"verified": true,
"shared": false,
"locale": "en-AU",
"timezone": "Sri Lanka",
"last_login_at": "2020-02-05 12:47:36.286053 -0800 -0800",
"email": "coffeyrasmussen@flotonic.com",
"phone": "8335-422-718",
"signature": "Don't Worry Be Happy!",
"organization_id": 119,
"tags": [
"Springville",
"Sutton",
"Hartsville/Hartley",
"Diaperville"
],
"suspended": true,
"role": "admin"
}
]
```      

### Source files

Implementation of API endpoints. 
```
web
```

Implementation of organizations, tickets, and users interface 
```
domains
```


## Running the tests

Run this command in the project directory.
```
 go test ./...
```
