# AI Agent Command: Go JWT Authentication & Authorization Implementation

## 🎯 Context & Objective
We are building a logistics/courier application named **Kuroute** using Go (Golang). 
The database is PostgreSQL, and we are using GORM as the ORM.
Primary Keys for all models (including Users) use **UUID** (via `github.com/google/uuid` and PostgreSQL `gen_random_uuid()`).

Your task is to implement the complete **Authentication** and **Authorization** feature, following a strict layered architecture. 
Implement the components in this exact order: **Service Layer -> Handler -> Middleware -> Router**.

---

## 🛠️ Technical Specifications & Rules

1. **User Struct Reference:**
   The user model embeds a `BaseWithUpdate` struct which handles the UUID `ID`, `CreatedAt`, and `UpdatedAt` automatically via GORM defaults.
2. **Token Handover (Mobile App Friendly):**
   Since the frontend is a Mobile App (Expo), **DO NOT use HTTP-Only Cookies**. 
   Both `access_token` and `refresh_token` MUST be returned inside the **JSON Response Body**.
3. **Token Lifespan:**
   - Access Token: 15 Minutes (Short-lived, contains payload: `user_id`, `role`, `hub_id`).
   - Refresh Token: 7 Days (Long-lived, tracked in the database for revocation).
   - List role: staff_sortir and kurir
4. **No Self-Register for Couriers:**
   Couriers are registered by admins. When a courier logs in for the first time, the backend should detect it (e.g., via `must_change_password: true`) and force a password setup.

---

## 🚀 Step-by-Step Implementation Instructions

Please generate the code step-by-step for the following layers:

### 1. Service Layer (`/internal/service/auth_service.go`)
Create an `AuthService` that handles the business logic:
- `Login(email, password)`: Verifies credentials. If it's a first-time login, returns a short-lived `setup_token`. If verified, returns `access_token` and `refresh_token`.
- `SetupFirstPassword(userId, newPassword)`: Updates the temporary password to a permanent one, marks the user as verified, and returns the real Access & Refresh tokens.
- `RefreshToken(oldRefreshToken)`: Validates the refresh token against the database, implements *Refresh Token Rotation* (invalidates the old one, issues a new pair).
- `Logout(refreshToken)`: Deletes or invalidates the refresh token in the database.

### 2. Handler / Controller Layer (`/internal/handler/auth_handler.go`)
Create the HTTP handlers that parse requests and return standard JSON payloads.
- **Login Response Structure (Success):**
  ```json
  {
    "status": "success",
    "data": {
      "user": { "id": "uuid-string", "name": "...", "role": "courier" },
      "tokens": { "access_token": "...", "refresh_token": "...", "expires_in": 900 }
    }
  }




{
  "success": true,
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJhYTJiMjJhYi03Y2M4LTRiMTMtOWIxYi05ZTYyMWM1MTI0NjkiLCJyb2xlIjoic3RhZmZfc29ydGlyIiwiaHViSWQiOiIzYTk2M2ZiZS1iZGQwLTQ5YzAtYWE3OC1hZDJjZTgzNDQyMWEiLCJleHAiOjE3ODMwMjIxNjYsImlhdCI6MTc4MzAyMTI2Nn0.J4nJaDYtsXr7H_AUazxeDAmGQQ58R2d9YSy2MnLC9uc",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJhYTJiMjJhYi03Y2M4LTRiMTMtOWIxYi05ZTYyMWM1MTI0NjkiLCJyb2xlIjoic3RhZmZfc29ydGlyIiwiaHViSWQiOiIzYTk2M2ZiZS1iZGQwLTQ5YzAtYWE3OC1hZDJjZTgzNDQyMWEiLCJleHAiOjE3ODU2MTMyNjYsImlhdCI6MTc4MzAyMTI2Nn0.eRGmN1SEczdnyw9dgF1TgJ7JeEUUS0dwjtr4tJNaMBs",
    "expiresIn": 900,
    "user": {
      "id": "aa2b22ab-7cc8-4b13-9b1b-9e621c512469",
      "hubId": "3a963fbe-bdd0-49c0-aa78-ad2ce834421a",
      "name": "Admin",
      "email": "kradmin@example.com",
      "phone": "08123456789",
      "role": "staff_sortir"
    }
  }
}