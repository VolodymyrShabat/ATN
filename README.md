# Book service


## Available endpoints

## User

### Registration http://localhost:8089/user/sign-up

Request body:
```
{
    "login":"kyrkela",
    "email":"volodymyrshabat@gmail.com",
    "password":"12345",
    "name":"volodymyrshabat",
    "city":"lviv",
    "age":23
}
```
Response:
```
"successfully registered"
```

### Login http://localhost:8089/user/sign-in

Request body:
```
{
    "login":"kyrkela",
    "password":"123"
}
```
Response:
```
logged in successfully
 access token - eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjMwNjUzMjk0MzIsImV4cCI6MTcxMTMwNzkyMywiaWF0IjoxNzExMzA2MTIzfQ.sqRSl7MgGCLkEtjCN5SevxvWFTKjG04XS4RGUR76TtA
```

### Forgot password(sending mail) http://localhost:8089/user/forgot-password
Request body:
```
{
    "email":"volodymyrshabat@gmail.com"
}
```
Response:
```
"Recovery email has been sent to your email"
```

### Password renewal POST http://localhost:8089/user/reset-password/{token}
Request body:
```
{
    "password":"123"
}
```
Response:
```
"Password successfully recovered"
```

## Book

### Book creation POST http://localhost:8089/book/create
Request body:
```
{
    "name":"Eneyida",
    "about":"poem"
}
```
Response:
```
"book successfully created with id: 2236385689"
```

### Book update POST http://localhost:8089/book/update/{id}
Request body:
```
{
    "name":"Harry Potter",
    "about":"book about magic"
}
```
Response:
```
"book successfully updated
```

### Book representation GET http://localhost:8089/book/get/{id}
Request body:empty
Response:
```
Book id: 1668426724
 name: Harry Potter
 about: book about magic
 creator_id: 3065329432
```

### Book deletion DELETE http://localhost:8089/book/delete/{id}
Request body:empty
Response:
```
book successfully deleted
```