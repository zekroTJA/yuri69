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
> Accept: application/json
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

List all available sounds.

### Query Parameters

| Name      | Type       | Example     | Description                                                       |
| --------- | ---------- | ----------- | ----------------------------------------------------------------- |
| `order`   | `string`   | `created`   | Sort sounds either alphabetically by `name` or by `created` date. |
| `include` | `string[]` | `meme,loud` | Filter sounds by tags which must be included.                     |
| `exclude` | `string[]` | `nsfw,nsft` | Filter sounds by tags which must not be included.                 |

### Response

```json
[
  {
    "uid": "mathemann",
    "display_name": "",
    "created_date": "2022-05-16T21:11:05.019837Z",
    "creator_id": "221905671296253953",
    "tags": null
  },
  {
    "uid": "cringe",
    "display_name": "",
    "created_date": "2022-05-16T21:11:04.304456Z",
    "creator_id": "221905671296253953",
    "tags": null
  },
  {
    "uid": "sheesh",
    "display_name": "",
    "created_date": "2022-05-16T21:11:03.615009Z",
    "creator_id": "221905671296253953",
    "tags": null
  },
  {
    "uid": "sus",
    "display_name": "",
    "created_date": "2022-05-16T21:11:02.973601Z",
    "creator_id": "221905671296253953",
    "tags": null
  }
]
```

## Upload Sound

> `PUT /api/v1/sounds/upload`

Upload sound file as multipart/form-data.

### Query Parameters

| Name   | Type     | Example      | Description                                                                                                            |
| ------ | -------- | ------------ | ---------------------------------------------------------------------------------------------------------------------- |
| `type` | `string` | `audio/mpeg` | Optional soecifier for the content type when it can not be inferred from the file itself or the `Content-Type` header. |

### Example Request

```
> PUT /api/v1/sounds/upload HTTP/2
> Host: localhost:8080
> authorization: basic t39GIXJlBYOvMu9XdYWWGGFg7pWWUUvT
> content-type: multipart/form-data; boundary=-----------------------------3928448259432856853608407134
> Accept: application/json
> content-length: 22652

-----------------------------3928448259432856853608407134
Content-Disposition: form-data; name="file"; filename="heyo.mp3"
Content-Type: audio/mpeg

[data]
```

### Response

```json
{
  "upload_id": "ca2f3an9f40lstb13gpg",
  "deadline": "2022-05-18T13:21:26.515045192Z"
}
```

## Create Sound

> `POST /api/v1/sounds/create`

Create a sound by passing meta data with the upload ID.

### Example Payload

```json
{ 
  "upload_id": "ca2f3an9f40lstb13gpg", 
  "uid": "heyo",
  "display_name": "heyo",
  "tags": [
    "meme",
    "borderlands"
  ],
  "normalize": true,
}
```

### Response

```json
{
  "uid":"heyo",
  "display_name":"heyo",
  "created_date":"2022-05-18T13:26:56.000702818Z",
  "creator_id":"221905671296253953",
  "tags": [
    "meme",
    "borderlands"
  ]
}
```

## Get Sound

> `GET /api/v1/sounds/<uid>`

Get a sound by UID.

### Response

```json
{
  "uid":"heyo",
  "display_name":"heyo",
  "created_date":"2022-05-18T13:26:56.000702818Z",
  "creator_id":"221905671296253953",
  "tags": [
    "meme",
    "borderlands"
  ]
}
```

## Update Sound

> `POST /api/v1/sounds/<id>`

Update a sound by its UID.

### Example Payload

```json
{ 
  "upload_id": "ca2f3an9f40lstb13gpg", 
  "uid": "heyo",
  "display_name": "heyo",
  "tags": [
    "meme",
    "borderlands"
  ],
  "normalize": true,
}
```

### Response

```json
{
  "uid":"heyo",
  "display_name":"heyo",
  "created_date":"2022-05-18T13:26:56.000702818Z",
  "creator_id":"221905671296253953",
  "tags": [
    "meme",
    "borderlands"
  ]
}
```