
# sso-v2

## Overview
This is a simple REST-ish session management tool for keeping an arbitrary set of values associated with an ID.  It also handles very basic user authentication.  This application uses a single Redis store to store both users and sessions.  Sessions are designed to time out automatically after 60 minutes using Redis' key timeout functionality. 

This should not be used in production as it lacks critical identity validation steps associated with retrieval of session values.

## Routes
#### POST /v1/user/
Creates a new user

Request Structure
```json
{
  "username": string,
  "password": string
}
```

#### POST /v1/user/doAuth
Authorizes a user and creates a new session.  Returns the new session id in a response header `X-Session-Id`

Request Body Structure
```json
{
  "username": string,
  "password": string
}
```

***

#### GET /v1/sessions/:sessionId
Retrieve current session data for a sessionId

#### PUT /v1/sessions/:sessionId
Sets the set of session variables in the session data.

Request Body Structure
```json
{
  "sessionVars": {
    "key"(string): "value"(string)
  }
}
```

#### DELETE /v1/sessions/:sessionId
Destroys a session by removing it from the Redis store explicitly.

## Heroku Configuration
This is set up to be run as a docker container on the Heroku platform.  Please contact me for a live demo link if you desire.  