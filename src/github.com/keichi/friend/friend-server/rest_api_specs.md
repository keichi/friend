# friend-server REST API specs

## Create user

```
POST /users
```

Request body must be `application/json` and include `name`, `password` and `publickey`.

```
{
    "name": "someusername",
    "password": "somepassword",
    "publicKey": "somepublickey"
}
```

## Login

```
POST /login
```

Request body must be `application/json` and include `name` and `password`.

```
{
	"name": "someusername",
	"password": "somepassword"
}
```

On success, session information is returned.

```
{
    "CreatedAt": "2014-07-13T21:02:33.118737091+09:00",
    "Expires": "2014-08-12T21:02:33.118737091+09:00",
    "Id": 0,
    "Token": "0bdae82d43ad56e301eb4705337cb207ace7177ed08ee7e84b7fc45d81918da7",
    "UpdatedAt": "2014-07-13T21:02:33.118737091+09:00",
    "UserId": 0
}
```

The value of field `Token` is a session token, which should be used for other API queries that require authentication. Set session token to the  header `X-Friend-Session-Token`.

## Logout

```
POST /logout
```

`X-Friend-Session-Token` header must be a valid session token.

## Get user

```
GET /users/:name
```

If `X-Friend-Session-Token` is a valid session token, detailed session information is included.

## Delete user

```
DELETE /users/:name
```

`X-Friend-Session-Token` header must be a valid session token.
