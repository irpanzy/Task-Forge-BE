# TaskForge API Documentation

RESTful API documentation for the **TaskForge** backend.

- **Base URL**: `http://localhost:3000/api`
- **Content-Type**: `application/json`
- **Response Format**: JSend-compliant JSON

---

## 1. Security Mechanisms

### A. Double Submit CSRF Protection
For all state-modifying requests (`POST`, `PUT`, `DELETE`), the API is protected by CSRF middleware:
1. Call `GET /api/auth/csrf` before submitting forms, logging in, or registering.
2. The server responds with `csrf_token` in the response body and sets a `csrf_` cookie.
3. Include the following header on all `POST`, `PUT`, `DELETE` requests:
   ```http
   X-CSRF-Token: <csrf_token_value>
   ```
4. Ensure the `csrf_` cookie is sent alongside the request (handled automatically by browsers with `withCredentials: true` or Postman).

### B. JWT Authentication & HTTPOnly Cookies
After successful login, the server sets two cookies with the `HttpOnly` flag:
- **`access_token`**: Used to authenticate protected endpoints (default expiry: 15 minutes).
- **`refresh_token`**: Used to obtain a new `access_token` without logging in again (default expiry: 7 days).

> **Fallback:** Non-browser clients (such as mobile apps) can also send the token via the Authorization header:
> ```http
> Authorization: Bearer <access_token>
> ```

---

## 2. Standard Response Format

### Success Response
```json
{
  "status": "success",
  "code": 200,
  "message": "Descriptive success message",
  "data": { ... }
}
```

### Error Response
```json
{
  "status": "error",
  "code": 400,
  "message": "Descriptive error message",
  "errors": "Detailed error string or null"
}
```

---

## 3. Endpoints

### Auth Module (`/api/auth`)

Overview table for authentication endpoints:

| Method | Endpoint | Access | Required Header | Request Body | Success Status | Description |
| :---: | :--- | :---: | :--- | :---: | :---: | :--- |
| `GET` | `/api/auth/csrf` | Public | - | - | `200` | Retrieve CSRF token & `csrf_` cookie |
| `POST` | `/api/auth/register` | Public | `X-CSRF-Token` | `{ name, email, password }` | `201` | Register new user (default role: `user`) |
| `POST` | `/api/auth/login` | Public | `X-CSRF-Token` | `{ email, password }` | `200` | Login & set `access_token` & `refresh_token` cookies |
| `POST` | `/api/auth/refresh` | Public | `X-CSRF-Token` *(Cookie: `refresh_token`)* | - | `200` | Issue a new `access_token` |
| `POST` | `/api/auth/logout` | Public | `X-CSRF-Token` | - | `200` | Clear session cookies |

---

#### 1. Get CSRF Token
Retrieves a fresh CSRF token to use in subsequent requests.

- **Method**: `GET`
- **URL**: `/api/auth/csrf`
- **Access**: Public
- **Headers**: None

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "CSRF token retrieved successfully",
  "data": {
    "csrf_token": "d2a1941c-f74b-4601-87dc-711cf81825cc"
  }
}
```
*Note: Response header includes `Set-Cookie: csrf_=...; Path=/; SameSite=Lax`.*

---

#### 2. Register
Registers a new user account with default role `user`.

- **Method**: `POST`
- **URL**: `/api/auth/register`
- **Access**: Public
- **Headers**:
  - `Content-Type: application/json`
  - `X-CSRF-Token: <token>`

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Validation:**
- `name`: Required.
- `email`: Required, valid email format, and not already registered.
- `password`: Minimum 6 characters.

**Success Response (201 Created):**
```json
{
  "status": "success",
  "code": 201,
  "message": "Registration successful",
  "data": {
    "public_id": "8f3b2075-e8d9-4b21-9562-b13c7bb61c6b",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2026-09-05T22:30:00Z",
    "updated_at": "2026-09-05T22:30:00Z"
  }
}
```

**Example Error Response (400 Bad Request):**
```json
{
  "status": "error",
  "code": 400,
  "message": "email is already registered",
  "errors": null
}
```

---

#### 3. Login
Authenticates user using email and password credentials.

- **Method**: `POST`
- **URL**: `/api/auth/login`
- **Access**: Public
- **Headers**:
  - `Content-Type: application/json`
  - `X-CSRF-Token: <token>`

**Request Body:**
```json
{
  "email": "admin@taskforge.com",
  "password": "password123"
}
```

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Login successful",
  "data": {
    "public_id": "cc3b4487-901b-45b6-8c32-580de73a8aeb",
    "name": "Administrator",
    "email": "admin@taskforge.com",
    "role": "admin",
    "created_at": "2026-09-04T20:21:41Z",
    "updated_at": "2026-09-04T20:21:41Z"
  }
}
```
*Note: Server sets `access_token` and `refresh_token` cookies (`HttpOnly`, `SameSite=Lax`).*

**Example Error Response (401 Unauthorized):**
```json
{
  "status": "error",
  "code": 401,
  "message": "invalid email or password",
  "errors": null
}
```

---

#### 4. Refresh Token
Issues a new `access_token` when the existing access token has expired.

- **Method**: `POST`
- **URL**: `/api/auth/refresh`
- **Access**: Public (requires `refresh_token` cookie)
- **Headers**:
  - `X-CSRF-Token: <token>`

**Request Body:** None.

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Example Error Response (401 Unauthorized):**
```json
{
  "status": "error",
  "code": 401,
  "message": "Invalid or expired refresh token",
  "errors": null
}
```

---

#### 5. Logout
Clears authentication session cookies on the client.

- **Method**: `POST`
- **URL**: `/api/auth/logout`
- **Access**: Public
- **Headers**:
  - `X-CSRF-Token: <token>`

**Request Body:** None.

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Logout successful",
  "data": null
}
```
*Note: `access_token` and `refresh_token` cookies are cleared.*

---

### User Module (`/api/users`)

All endpoints under the User module require authentication (`Authenticate`).

Overview table for user endpoints:

| Method | Endpoint | Access / Guard | Parameters | Request Body | Success Status | Description |
| :---: | :--- | :---: | :--- | :---: | :---: | :--- |
| `GET` | `/api/users/me` | Authenticated | - | - | `200` | Profile of currently authenticated user |
| `GET` | `/api/users` | `admin` only | Query: `?page=1&limit=10&search=budi` | - | `200` | List of all users with pagination and search |
| `GET` | `/api/users/:id` | Authenticated | Path: `:id` (UUID) | - | `200` | Detail of user by public UUID |
| `PUT` | `/api/users/:id` | Self / Admin | Path: `:id` (UUID) | `{ name?, email? }` | `200` | Update user name or email |
| `DELETE` | `/api/users/:id` | `admin` only | Path: `:id` (UUID) | - | `200` | Delete user (prevents deleting active self) |

---

#### 1. Get Current User Profile (Me)
Retrieves profile data for the currently authenticated user.

- **Method**: `GET`
- **URL**: `/api/users/me`
- **Access**: Authenticated User (all roles)
- **Cookie**: `access_token` (or Header `Authorization: Bearer <token>`)

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Profile retrieved successfully",
  "data": {
    "public_id": "cc3b4487-901b-45b6-8c32-580de73a8aeb",
    "name": "Administrator",
    "email": "admin@taskforge.com",
    "role": "admin",
    "created_at": "2026-09-04T20:21:41Z",
    "updated_at": "2026-09-04T20:21:41Z"
  }
}
```

---

#### 2. Get All Users (Pagination & Search)
Retrieves a paginated and searchable list of all registered users in the system.

- **Method**: `GET`
- **URL**: `/api/users`
- **Access**: **`admin`** only
- **Query Parameters**:
  - `search` (optional): Search keyword across user `name` or `email` (case-insensitive).
  - `page` (optional, default: `1`): Page number.
  - `limit` (optional, default: `10`, max: `100`): Items per page.
- **Example URL**: `/api/users?page=1&limit=10&search=budi`

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Users retrieved successfully",
  "data": {
    "users": [
      {
        "public_id": "cc3b4487-901b-45b6-8c32-580de73a8aeb",
        "name": "Administrator",
        "email": "admin@taskforge.com",
        "role": "admin",
        "created_at": "2026-09-04T20:21:41Z",
        "updated_at": "2026-09-04T20:21:41Z"
      },
      {
        "public_id": "8f3b2075-e8d9-4b21-9562-b13c7bb61c6b",
        "name": "John Doe",
        "email": "john@example.com",
        "role": "user",
        "created_at": "2026-09-05T22:30:00Z",
        "updated_at": "2026-09-05T22:30:00Z"
      }
    ],
    "total_data": 2,
    "current_page": 1,
    "total_pages": 1,
    "limit": 10
  }
}
```

**Example Error Response (403 Forbidden):**
```json
{
  "status": "error",
  "code": 403,
  "message": "Access denied: you do not have permission for this action",
  "errors": null
}
```

---

#### 3. Get User By ID
Retrieves user details by public UUID.

- **Method**: `GET`
- **URL**: `/api/users/:id`
- **Access**: Authenticated User
- **Path Parameters**:
  - `id`: User UUID (e.g. `8f3b2075-e8d9-4b21-9562-b13c7bb61c6b`)

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "User retrieved successfully",
  "data": {
    "public_id": "8f3b2075-e8d9-4b21-9562-b13c7bb61c6b",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2026-09-05T22:30:00Z",
    "updated_at": "2026-09-05T22:30:00Z"
  }
}
```

---

#### 4. Update User Profile
Updates user name or email address.

- **Method**: `PUT`
- **URL**: `/api/users/:id`
- **Access**: Account Owner or **`admin`**
- **Headers**:
  - `Content-Type: application/json`
  - `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `id`: Target user UUID

**Request Body:**
```json
{
  "name": "John Doe Updated",
  "email": "john.updated@example.com"
}
```

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "User updated successfully",
  "data": {
    "public_id": "8f3b2075-e8d9-4b21-9562-b13c7bb61c6b",
    "name": "John Doe Updated",
    "email": "john.updated@example.com",
    "role": "user",
    "created_at": "2026-09-05T22:30:00Z",
    "updated_at": "2026-09-05T23:15:00Z"
  }
}
```

**Example Error Response (403 Forbidden when modifying another user):**
```json
{
  "status": "error",
  "code": 403,
  "message": "Access denied: You do not have permission to modify another user's data",
  "errors": null
}
```

---

#### 5. Delete User
Deletes a user account from the database.

- **Method**: `DELETE`
- **URL**: `/api/users/:id`
- **Access**: **`admin`** only
- **Headers**:
  - `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `id`: Target user UUID

> **Safety Guard:** An administrator cannot delete their own currently active account.

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "User deleted successfully",
  "data": null
}
```

**Example Error Response (400 Bad Request when admin attempts self-deletion):**
```json
{
  "status": "error",
  "code": 400,
  "message": "Cannot delete your own active account",
  "errors": null
}
```

---

### Board Module (`/api/boards`)

All endpoints under the Board module require authentication (`Authenticate`).

Overview table for board endpoints:

| Method | Endpoint | Access / Guard | Parameters | Request Body | Success Status | Description |
| :---: | :--- | :---: | :--- | :---: | :---: | :--- |
| `POST` | `/api/boards` | Authenticated | - | `{ title, description?, due_date? }` | `201` | Create a new board |
| `GET` | `/api/boards` | Authenticated | Query: `?page=1&limit=10&search=...` | - | `200` | List user's boards (owned & member) |
| `GET` | `/api/boards/:id` | Owner / Member / Admin | Path: `:id` (UUID) | - | `200` | Retrieve board detail by UUID (with members) |
| `PUT` | `/api/boards/:id` | Owner / Admin | Path: `:id` (UUID) | `{ title?, description?, due_date? }` | `200` | Update board details |
| `DELETE` | `/api/boards/:id` | Owner / Admin | Path: `:id` (UUID) | - | `200` | Delete board |
| `POST` | `/api/boards/:id/members` | Owner / Admin | Path: `:id` (UUID) | `["user_uuid1", "user_uuid2"]` | `200` | Add members to board |
| `GET` | `/api/boards/:id/members` | Owner / Member / Admin | Path: `:id` (UUID) | - | `200` | List all members of a board |
| `DELETE` | `/api/boards/:id/members/:userId` | Owner / Self / Admin | Path: `:id`, `:userId` (UUID) | - | `200` | Remove member from board |

---

#### 1. Create Board
Creates a new Kanban board for the currently authenticated user.

- **Method**: `POST`
- **URL**: `/api/boards`
- **Access**: Authenticated User
- **Headers**:
  - `Content-Type: application/json`
  - `X-CSRF-Token: <token>`

**Request Body:**
```json
{
  "title": "Project TaskForge Kanban",
  "description": "Main project board for backend and frontend tasks",
  "due_date": "2026-12-31T23:59:59Z"
}
```

**Success Response (201 Created):**
```json
{
  "status": "success",
  "code": 201,
  "message": "Board created successfully",
  "data": {
    "public_id": "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
    "owner_public_id": "cc3b4487-901b-45b6-8c32-580de73a8aeb",
    "title": "Project TaskForge Kanban",
    "description": "Main project board for backend and frontend tasks",
    "due_date": "2026-12-31T23:59:59Z",
    "created_at": "2026-09-06T22:15:00Z",
    "updated_at": "2026-09-06T22:15:00Z"
  }
}
```

---

#### 2. Get User Boards (Pagination & Search)
Retrieves boards owned by the authenticated user with pagination and optional search filter.

- **Method**: `GET`
- **URL**: `/api/boards`
- **Access**: Authenticated User
- **Query Parameters**:
  - `search` (optional): Search keyword on `title` or `description`.
  - `page` (optional, default: `1`): Page number.
  - `limit` (optional, default: `10`, max: `100`): Items per page.
- **Example URL**: `/api/boards?page=1&limit=10&search=Kanban`

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Boards retrieved successfully",
  "data": {
    "boards": [
      {
        "public_id": "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
        "owner_public_id": "cc3b4487-901b-45b6-8c32-580de73a8aeb",
        "title": "Project TaskForge Kanban",
        "description": "Main project board for backend and frontend tasks",
        "due_date": "2026-12-31T23:59:59Z",
        "created_at": "2026-09-06T22:15:00Z",
        "updated_at": "2026-09-06T22:15:00Z"
      }
    ],
    "total_data": 1,
    "current_page": 1,
    "total_pages": 1,
    "limit": 10
  }
}
```

---

#### 3. Get Board By ID
Retrieves details of a specific board.

- **Method**: `GET`
- **URL**: `/api/boards/:id`
- **Access**: Board Owner or **`admin`**
- **Path Parameters**:
  - `id`: Board UUID (e.g. `9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d`)

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Board retrieved successfully",
  "data": {
    "public_id": "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
    "owner_public_id": "cc3b4487-901b-45b6-8c32-580de73a8aeb",
    "title": "Project TaskForge Kanban",
    "description": "Main project board for backend and frontend tasks",
    "due_date": "2026-12-31T23:59:59Z",
    "created_at": "2026-09-06T22:15:00Z",
    "updated_at": "2026-09-06T22:15:00Z"
  }
}
```

**Example Error Response (403 Forbidden):**
```json
{
  "status": "error",
  "code": 403,
  "message": "access denied: you do not have permission to view this board",
  "errors": null
}
```

---

#### 4. Update Board
Updates title, description, or due date of an existing board.

- **Method**: `PUT`
- **URL**: `/api/boards/:id`
- **Access**: Board Owner or **`admin`**
- **Headers**:
  - `Content-Type: application/json`
  - `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `id`: Target board UUID

**Request Body:**
```json
{
  "title": "Updated TaskForge Kanban",
  "description": "Updated project description"
}
```

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Board updated successfully",
  "data": {
    "public_id": "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
    "owner_public_id": "cc3b4487-901b-45b6-8c32-580de73a8aeb",
    "title": "Updated TaskForge Kanban",
    "description": "Updated project description",
    "due_date": "2026-12-31T23:59:59Z",
    "created_at": "2026-09-06T22:15:00Z",
    "updated_at": "2026-09-06T22:20:00Z"
  }
}
```

---

#### 5. Delete Board
Deletes a board from the system (Soft Delete).

- **Method**: `DELETE`
- **URL**: `/api/boards/:id`
- **Access**: Board Owner or **`admin`**
- **Headers**:
  - `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `id`: Target board UUID

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Board deleted successfully",
  "data": null
}
```

---

#### 6. Add Board Members
Adds one or multiple users to the board as members.

- **Method**: `POST`
- **URL**: `/api/boards/:id/members`
- **Access**: Board Owner or **`admin`**
- **Headers**:
  - `Content-Type: application/json`
  - `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `id`: Target board UUID

**Request Body:**
Array of user public UUIDs to add:
```json
[
  "2255affe-8ed3-446a-9efc-978de63d085e",
  "75a25a74-493c-4d51-93ca-9f2c6a14fed8"
]
```

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Members added successfully",
  "data": null
}
```

---

#### 7. Get Board Members
Retrieves all members belonging to a board.

- **Method**: `GET`
- **URL**: `/api/boards/:id/members`
- **Access**: Board Owner, Member, or **`admin`**
- **Path Parameters**:
  - `id`: Target board UUID

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Members retrieved successfully",
  "data": [
    {
      "public_id": "2255affe-8ed3-446a-9efc-978de63d085e",
      "name": "Budi Santoso",
      "email": "budi@example.com",
      "role": "user",
      "joined_at": "2026-09-07T22:30:00Z"
    }
  ]
}
```

---

#### 8. Delete Board Member
Removes a member from a board (or member leaving a board).

- **Method**: `DELETE`
- **URL**: `/api/boards/:id/members/:userId`
- **Access**: Board Owner, the member themselves, or **`admin`**
- **Headers**:
  - `X-CSRF-Token: <token>`
- **Path Parameters**:
  - `id`: Target board UUID
  - `userId`: Target member user UUID

**Success Response (200 OK):**
```json
{
  "status": "success",
  "code": 200,
  "message": "Member removed successfully",
  "data": null
}
```

---

## 4. Role-Based Access Control (RBAC) Matrix

| Endpoint | Public | Role `user` | Role `admin` |
| :--- | :---: | :---: | :---: |
| `GET /api/auth/csrf` | ✅ | ✅ | ✅ |
| `POST /api/auth/register` | ✅ | ✅ | ✅ |
| `POST /api/auth/login` | ✅ | ✅ | ✅ |
| `POST /api/auth/refresh` | ✅ | ✅ | ✅ |
| `POST /api/auth/logout` | ✅ | ✅ | ✅ |
| `GET /api/users/me` | ❌ | ✅ | ✅ |
| `GET /api/users` | ❌ | ❌ | ✅ |
| `GET /api/users/:id` | ❌ | ✅ | ✅ |
| `PUT /api/users/:id` | ❌ | ✅ *(own account only)* | ✅ *(all accounts)* |
| `DELETE /api/users/:id` | ❌ | ❌ | ✅ *(except self)* |
| `POST /api/boards` | ❌ | ✅ | ✅ |
| `GET /api/boards` | ❌ | ✅ *(my boards: owned & member)* | ✅ *(my boards: owned & member)* |
| `GET /api/boards/:id` | ❌ | ✅ *(owner & members)* | ✅ *(all boards)* |
| `PUT /api/boards/:id` | ❌ | ✅ *(owner only)* | ✅ *(all boards)* |
| `DELETE /api/boards/:id` | ❌ | ✅ *(owner only)* | ✅ *(all boards)* |
| `POST /api/boards/:id/members` | ❌ | ✅ *(owner only)* | ✅ *(all boards)* |
| `GET /api/boards/:id/members` | ❌ | ✅ *(owner & members)* | ✅ *(all boards)* |
| `DELETE /api/boards/:id/members/:userId` | ❌ | ✅ *(owner & self)* | ✅ *(all boards)* |


