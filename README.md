# go-crypto-price-tracker

## Description

This is an application to track crypto price written in golang

## Tech stack

* [gofiber](https://docs.gofiber.io/)
* sqlite3 for database
* [go-sqlite3](https://github.com/mattn/go-sqlite3) for sqlite driver
* JWT for tokenization
* [coincap.io](https://docs.coincap.io/) to get price
* [currencyapi](https://currencyapi.com/) to convert currency

## How to use

### Migration

To migrate the database, use [golang-migrate](https://github.com/golang-migrate/migrate/)

After `golang-migrate` has installed, use this command to do the migration

```
migrate -database "sqlite3://./go-crypto-price-tracker.db" -path ./database/migrations up
```

### Run Application

Run the application using this command

```
go run main.go
```

## Endpoints

### POST /user/signup

This endpoint used to register user. Before using the application, user need to be registered first

#### Body

| Field | Type | Required (Y/N) | Description |
| ----- | ---- | -------------- | ----------- |
| email | string | Y | email to register, must be unique
| password | string | Y | password needed to sign in
| confirm_password | string | Y | password confirmation, need to be same as password

Example

```
{
    "email": "test@test.com",
    "password": "test123",
    "confirm_password": "test123"
}
```

#### Response Body

##### 201 CREATED

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |
| user | user interface | contains user information (email) |

##### 400 BAD REQUEST

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

##### 500 INTERNAL SERVER ERROR

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

#### Response Cookies

| Field | Type | Description |
| ----- | ---- | -------------- |
| token | string | token used to use tracker and coin APIs |

<br />

### POST /user/signin

This endpoint used to sign in. Before using the application, user need to be signed in first to obtain the token

#### Body

| Field | Type | Required (Y/N) | Description |
| ----- | ---- | -------------- | ----------- |
| email | string | Y | email registered |
| password | string | Y | password needed to sign in |

Example

```
{
    "email": "test@test.com",
    "password": "test123"
}
```

#### Response Body

##### 200 OK

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |
| user | user interface | contains user information (email) |

##### 400 BAD REQUEST

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

##### 500 INTERNAL SERVER ERROR

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

#### Response Cookies

| Field | Type | Description |
| ----- | ---- | -------------- |
| token | string | token used to use tracker and coin APIs |

<br />

### GET /user/signout

This endpoint used to sign out

#### Response Body

##### 200 OK

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

##### 400 BAD REQUEST

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

##### 500 INTERNAL SERVER ERROR

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

<br />

### GET /coin

This endpoint used to get coins

#### Query Params

| Field | Type | Required (Y/N) | Description |
| ----- | ---- | -------------- | ----------- |
| limit | string | N | limit of the data |
| offset | string | N | offset of the data |

Example

```
{
    "limit": "3",
    "offset": "0"
}
```

#### Cookies

| Field | Type | Required (Y/N) | Description |
| ----- | ---- | -------------- | ----------- |
| token | string | Y | token to authorize user |

#### Response Body

##### 200 OK

| Field | Type | Description |
| ----- | ---- | -------------- |
| data | array of coin | contain array of coin information (id, name, priceUsd, priceIdr) |

##### 400 BAD REQUEST

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

##### 500 INTERNAL SERVER ERROR

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

<br />

### GET /coin/:coinId

This endpoint used to get a coin information

#### Cookies

| Field | Type | Required (Y/N) | Description |
| ----- | ---- | -------------- | ----------- |
| token | string | Y | token to authorize user |

#### Response Body

##### 200 OK

| Field | Type | Description |
| ----- | ---- | -------------- |
| data | coin interface | contain coin information (id, name, priceUsd, priceIdr)

##### 400 BAD REQUEST

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

##### 404 NOT FOUND

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

##### 500 INTERNAL SERVER ERROR

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

<br />

### POST /tracker

This endpoint used to create a new tracker of a coin

#### Body

| Field | Type | Required (Y/N) | Description |
| ----- | ---- | -------------- | ----------- |
| coin_id | string | Y | coin id to be tracked |

Example

```
{
    "coin_id": "bitcoin"
}
```

#### Cookies

| Field | Type | Required (Y/N) | Description |
| ----- | ---- | -------------- | ----------- |
| token | string | Y | token to authorize user |

#### Response Body

##### 200 OK

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |
| tracker | trakcer interface | contains trakcer information (coin_id, user_email) |

##### 400 BAD REQUEST

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

##### 500 INTERNAL SERVER ERROR

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

<br />

### GET /tracker

This endpoint used to get user's coin tracker

#### Cookies

| Field | Type | Required (Y/N) | Description |
| ----- | ---- | -------------- | ----------- |
| token | string | Y | token to authorize user |

#### Response Body

##### 200 OK

| Field | Type | Description |
| ----- | ---- | -------------- |
| data | array of coin | contain array of coin information (id, name, priceUsd, priceIdr) |

##### 400 BAD REQUEST

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

##### 500 INTERNAL SERVER ERROR

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

<br />

### DELETE /tracker

This endpoint used to delete a coin from user's tracker

#### Body

| Field | Type | Required (Y/N) | Description |
| ----- | ---- | -------------- | ----------- |
| coin_id | string | Y | coin id to be deleted from tracker |

Example

```
{
    "coin_id": "bitcoin"
}
```

#### Cookies

| Field | Type | Required (Y/N) | Description |
| ----- | ---- | -------------- | ----------- |
| token | string | Y | token to authorize user |

#### Response Body

##### 200 OK

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |
| tracker | trakcer interface | contains trakcer information (coin_id, user_email) |

##### 400 BAD REQUEST

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |

##### 500 INTERNAL SERVER ERROR

| Field | Type | Description |
| ----- | ---- | -------------- |
| message | string | message of the operation |
