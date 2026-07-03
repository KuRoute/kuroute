# AI Agent Command: Strict Stateless Authentication (No DB Token Storage)

## 🛑 CRITICAL ARCHITECTURAL CORRECTION
You added `SaveRefreshToken` and `ClearRefreshToken` to the `UserRepository`. This is **WRONG**. 

As explicitly stated before, the client is a **Mobile App (Expo)**. 
- We do **NOT** have a `refresh_token` column in the `users` table.
- We do **NOT** store, track, or save refresh tokens in the PostgreSQL database at all.
- The authentication architecture must be **100% Stateless**.

---

## 🛠️ Required Refactoring Instructions

### 1. Repository Layer (`/internal/repository/user_repository.go`)
- **REMOVE** the `SaveRefreshToken` and `ClearRefreshToken` functions completely. 
- Do NOT touch or alter the `users` table schema.

### 2. Service Layer (`/internal/service/auth_service.go`)
- **`Login`**: Just verify the user's password. If valid, generate both `access_token` and `refresh_token` using the JWT library, and return them. **Do not save anything to the database.**
- **`Logout`**: Since we are stateless, the backend does not need to clear anything in the DB. The handler can just return a success message (the Mobile App will delete the token from its `SecureStore`). You can remove DB dependency from the Logout service logic.
- **`RefreshToken`**: Validate the incoming `refresh_token`'s cryptography and expiration using the JWT library. If it is valid, parse the claims (e.g., `user_id`, `role`), and immediately generate a new `access_token` (and a new `refresh_token` if rotation is applied). **Zero database checks/updates are required for the token itself.**

Please re-write the Service and Repository layers to strictly adhere to this stateless JWT pattern.