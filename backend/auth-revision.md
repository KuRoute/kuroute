# AI Agent Command: Revision for Auth Feature & CORS Implementation

## 🛑 Important Corrections & Adjustments
Please revise the previously generated Authentication & Authorization code based on the following structural changes:

1. **NO Separate Refresh Token Table/Repository:**
   - We do NOT have a separate `refresh_tokens` table. 
   - Instead, the `refresh_token` string must be stored **directly inside a column in the existing `users` table** (e.g., `refresh_token` column in the User model).
   - Update the `UserRepository` to handle saving and clearing this token during login, token rotation, and logout. Remove any separate `RefreshTokenRepository`.

2. **Add CORS Middleware:**
   - Implement a secure CORS (Cross-Origin Resource Sharing) middleware in your router configuration (`/cmd/router`).
   - It must allow requests from the Mobile App environment (Expo/Localhost development) and specify allowed methods (`GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`) and allowed headers (`Content-Type`, `Authorization`).

---

## 🛠️ Revised Layer Requirements

### 1. Service Layer (`/internal/service/auth_service.go`)
- Update `Login`, `RefreshToken`, and `Logout` logic to interact only with `UserService` / `UserRepository` to update the `refresh_token` field on the specific User record.

### 2. Middleware & Router Layer (`/cmd/router` or `/internal/middleware`)
- Implement the CORS configuration before registering any application routes.
- Example config for CORS:
  - `AllowedOrigins`: `["*"]` (or specific development origins like `http://localhost:*`)
  - `AllowedMethods`: `["GET", "POST", "PUT", "DELETE", "OPTIONS"]`
  - `AllowedHeaders`: `["Accept", "Authorization", "Content-Type", "X-CSRF-Token"]`

Please refactor the code cleanly based on these rules.