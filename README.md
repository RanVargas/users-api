# Mandatory Environment Variables to Declare for API Use

This document details the importance of the environment variables required to use our REST API. Each variable plays a crucial role in configuring the API and its associated database.

## Database Configuration

### POSTGRES_DB
- **Purpose**: Specifies the name of the PostgreSQL database.
- **Importance**: Ensures the API connects to the correct database.

### POSTGRES_PASSWORD
- **Purpose**: Sets the password for the PostgreSQL database user.
- **Importance**: Critical for securing database access. Should be a strong, unique password in production.

### POSTGRES_PORT
- **Purpose**: Defines the port number on which the PostgreSQL server is listening.
- **Importance**: Ensures the API can establish a connection to the database server.

### POSTGRES_HOST
- **Purpose**: Specifies the hostname or IP address of the PostgreSQL server.
- **Importance**: Allows the API to locate and connect to the database server.

### POSTGRES_USERNAME
- **Purpose**: Sets the username for connecting to the PostgreSQL database.
- **Importance**: Provides necessary credentials for database access and determines user privileges.

### POSTGRES_SSL_MODE
- **Purpose**: Configures the SSL mode for the database connection.
- **Importance**: Affects the security of the database connection. "disable" means no SSL, which might be okay for local development but not recommended for production.

### POSTGRES_TIMEZONE
- **Purpose**: Sets the timezone for the database connection.
- **Importance**: Ensures consistent timestamp handling between the API and the database.

## API Configuration

### JWT_SECRET
- **Purpose**: Provides a secret key for JSON Web Token (JWT) generation and validation.

### X-API-KEY
- **Purpose**: The SignUp endpoint was not outlined at the user's request, hence a simple api key was requested when posting tot his endpoint. Only the SignUp endpoint requires this variable.

### API_PORT
- **Purpose**: Specifies the port on which the API will listen for incoming requests.

### API_SERVING_ADDRESS
- **Purpose**: Sets the address on which the API will be served.

Certainly! I'll create a Markdown file detailing the endpoints of your API and the data needed to interact with them.

# API Endpoints Documentation

This document details the endpoints of our REST API and the requirements for interacting with them.

## Authentication

Most endpoints require authentication, which is handled by the `RequireAuth` middleware.

## User Management

### Get User and Role
- **Endpoint**: `GET /users/:uid`
- **Authentication**: Required
- **Description**: Retrieves user information and their role.

### Update User
- **Endpoint**: `PUT /users/:uid`
- **Authentication**: Required
- **Description**: Updates user information.

### Delete User
- **Endpoint**: `DELETE /users/:uid`
- **Authentication**: Required
- **Description**: Deletes a user.

### Get Groups of User
- **Endpoint**: `GET /users/:uid/groups`
- **Authentication**: Required
- **Description**: Retrieves groups associated with a user.

### Create User
- **Endpoint**: `POST /users`
- **Authentication**: Required
- **Description**: Creates a new user.
- **Required Body**:
  ```json
  {
    "name": "string",
    "email": "string",
    "password": "string",
    "status": int16,
    "role_uid": "string",
    "groups_uid": ["string"]
  }
  ```

### Get All Users / Search Users
- **Endpoint**: `GET /users`
- **Authentication**: Required
- **Query Parameters**:
    - `searchTerm` (optional): If provided, searches for users
    - `limit` (optional, default: 10): Number of results to return
    - `orderBy` (optional, default: "id"): Field to order results by
- **Description**: Retrieves all users or searches for users based on query parameters.

## Group Management

### Create Group
- **Endpoint**: `POST /groups`
- **Authentication**: Required
- **Description**: Creates a new group.

### Get All Groups
- **Endpoint**: `GET /groups`
- **Authentication**: Required
- **Description**: Retrieves all groups.

### Get Group
- **Endpoint**: `GET /groups/:uid`
- **Authentication**: Required
- **Description**: Retrieves information about a specific group.

### Update Group
- **Endpoint**: `PUT /groups/:uid`
- **Authentication**: Required
- **Description**: Updates group information.

### Delete Group
- **Endpoint**: `DELETE /groups/:uid`
- **Authentication**: Required
- **Description**: Deletes a group.

## Role Management

### Create Role
- **Endpoint**: `POST /roles`
- **Authentication**: Required
- **Description**: Creates a new role.

### Get All Roles
- **Endpoint**: `GET /roles`
- **Authentication**: Required
- **Description**: Retrieves all roles.

### Get Role
- **Endpoint**: `GET /roles/:uid`
- **Authentication**: Required
- **Description**: Retrieves information about a specific role.

### Update Role
- **Endpoint**: `PUT /roles/:uid`
- **Authentication**: Required
- **Description**: Updates role information.

### Delete Role
- **Endpoint**: `DELETE /roles/:uid`
- **Authentication**: Required
- **Description**: Deletes a role.

### Get Users by Role
- **Endpoint**: `GET /roles/:uid/users`
- **Authentication**: Required
- **Description**: Retrieves users associated with a specific role.

## Authentication Endpoints

### Login
- **Endpoint**: `POST /login`
- **Authentication**: Not required
- **Description**: Authenticates a user and provides a token.

### Signup
- **Endpoint**: `POST /signup`
- **Authentication**: Not required
- **Description**: Registers a new user.
- **Required Body**:
  ```json
  {
    "name": "string",
    "email": "string",
    "password": "string",
    "status": int16,
    "role_uid": "string",
    "groups_uid": ["string"]
  }
  ```

### Logout
- **Endpoint**: `POST /logout`
- **Authentication**: Required
- **Description**: Logs out the current user.

## Notes
- All endpoints with `:uid` in the path require the unique identifier of the resource in question.
- The `RequireAuth` middleware is used for most endpoints to ensure proper authentication.
- Some endpoints like user creation and signup require specific JSON payloads as shown above.
- The search functionality for users is implemented in the GET /users endpoint and is triggered when a `searchTerm` is provided.


This is a test change
