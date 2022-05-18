# API Documentation

# Authentication

You must authenticate against the REST API using an API token which is passed as `basic` token in the `Authorization` header.

```
Authorization: basic m5Vu9/oyUOBmWu5Yn6ksUkFVKmOxzujN
```

The API token can be obtained from the web interface in the `/settings` route.

![](/.github/media/ss/token.png)

# REST API

## Obtain Access Token

> `GET /api/v1/auth`

Obtain an access token which can be used to authenticate a web socket connection.

### Example Request

```
> GET /api/v1/auth/refresh HTTP/2
> Host: localhost:8080
> Authorization: basic t39GIXJlBYOvMu9XdYWWGGFg7pWWUUvT
> Accept: */*
```

### Response 

```json
{
	"access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJle...",
	"deadline": "2023-05-17T14:17:38.020185058+02:00"
}
```

## Obtain OTA Token

> `GET /api/v1/auth/ota/token`

Obtain an one time authorization token and QR code.

### Response 

```json
{
	"deadline": "2022-05-17T14:50:11.170340454+02:00",
	"token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2...",
	"qrcode_data": "data:image/png;base64,iVBORw0KGgoAAAANSUhEU..."
}
```

## Check Authentication

> `GET /api/v1/auth/check`

Check your authentication state. Returns `401 Unauthorized` when not authenticated and otherwise, returns `200 OK`.

### Response 

```json
{
	"status": 200,
	"message": "Ok"
}
```

## List Sounds

> `GET /api/v1/sounds`

Check your authentication state. Returns `401 Unauthorized` when not authenticated and otherwise, returns `200 OK`.

### Response 

```json
{
	"status": 200,
	"message": "Ok"
}
```