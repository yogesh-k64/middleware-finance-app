# Quick Start Guide

## ✅ What Was Done

### 1. **Renamed Everything for Clarity**
- **Before**: `users` (confusing - login users or customers?)
- **After**: `customers` (your finance clients) and `admins` (system login users)

### 2. **Added Secure Authentication**
- JWT-based authentication with 24-hour expiry
- bcrypt password hashing (industry standard)
- Middleware that protects all your endpoints
- Role-based access (admin, manager, viewer)

### 3. **Updated All Files**
```
✓ users.go → customers.go
✓ userHelper.go → customerHelper.go
✓ Updated querys.go (all SQL queries)
✓ Updated constants.go (all messages)
✓ Updated Interfaces.go (Customer struct)
✓ Updated handouts.go (customer references)
✓ Updated main.go (new /customers endpoints)
✓ Created auth.go (authentication logic)
✓ Created migration-5.sql (admins table)
✓ Created migration-6.sql (rename users to customers)
```

---

## 🚀 Start Using It (4 Steps)

### Step 1: Install JWT Dependency (Already Done ✓)
```bash
go get github.com/golang-jwt/jwt/v5
```

### Step 2: Run Migrations
```bash
# If using PostgreSQL
psql $DATABASE_URL -f sql/migration-5.sql
psql $DATABASE_URL -f sql/migration-6.sql
```

### Step 3: Create Your First Admin
**Temporarily** expose the register endpoint in [main.go](main.go#L78):
```go
// Move this line BEFORE the protected routes section:
r.HandleFunc("/admin/register", registerAdmin).Methods("POST")
```

Then create admin:
```bash
curl -X POST http://localhost:9000/admin/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "SecurePassword123",
    "role": "admin"
  }'
```

**Important**: Move the endpoint back to protected routes after!

### Step 4: Start Using Your API
```bash
# 1. Login
curl -X POST http://localhost:9000/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "SecurePassword123"}'

# You'll get a token like: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# 2. Use the token for all requests
curl -X GET http://localhost:9000/customers \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## 📍 API Endpoint Changes

### Before → After
```
/users → /customers
/users/{id} → /customers/{id}
/users/{id}/handouts → /customers/{id}/handouts
/users/{id}/referred-by → /customers/{id}/referred-by
/users/{id}/referral → /customers/{id}/referral
```

### New Endpoints
```
POST /admin/login - Login to get token (PUBLIC)
POST /admin/register - Create new admin (PROTECTED)
GET  /admin/me - Get current admin info (PROTECTED)
```

---

## 🔐 How Authentication Works

1. **Login** → Get JWT token
2. **Include token** in `Authorization: Bearer TOKEN` header
3. **Middleware** validates token on every protected request
4. **Token expires** after 24 hours → login again

---

## 🎓 Key Learning Points

### 1. **Naming Matters**
- Clear names prevent confusion
- `customers` = business data
- `admins` = system access
- Good naming makes code self-documenting

### 2. **Security Layers**
```
Request → CORS → Auth Middleware → Business Logic
                      ↓
              (Validates JWT token)
                      ↓
              (Adds admin info to context)
```

### 3. **Separation of Concerns**
```
auth.go        → Authentication only
customers.go   → Customer business logic
handouts.go    → Handout business logic
main.go        → Route definitions
```

### 4. **Database Design**
- Separate tables for different entities
- `admins` have passwords, `customers` don't
- Foreign keys maintain data integrity
- Migrations allow versioned changes

---

## 📊 Your Architecture Now

```
┌─────────────────────────────────────┐
│         Frontend App                │
│   (React, Next.js, etc.)           │
└──────────────┬──────────────────────┘
               │
               │ HTTP + JWT Token
               ↓
┌─────────────────────────────────────┐
│         Go API Server               │
│  ┌─────────────────────────────┐   │
│  │   Auth Middleware            │   │ ← Validates all requests
│  └─────────────────────────────┘   │
│  ┌─────────────────────────────┐   │
│  │   Business Logic             │   │
│  │  • customers.go              │   │
│  │  • handouts.go               │   │
│  │  • collections.go            │   │
│  └─────────────────────────────┘   │
└──────────────┬──────────────────────┘
               │
               ↓
┌─────────────────────────────────────┐
│      PostgreSQL Database            │
│  • admins table (for login)         │
│  • customers table (your clients)   │
│  • handouts table                   │
│  • collections table                │
└─────────────────────────────────────┘
```

---

## 🎯 Best Practices Applied

✅ **Authentication**: JWT tokens (stateless, scalable)
✅ **Password Security**: bcrypt hashing (industry standard)
✅ **Authorization**: Middleware pattern (DRY principle)
✅ **Clear Naming**: No confusion between user types
✅ **Separation**: Auth logic separate from business logic
✅ **Type Safety**: Strong typing with Go structs
✅ **Database**: Separate tables, proper foreign keys
✅ **Documentation**: README with examples

---

## 🔧 Environment Variables Needed

```bash
DATABASE_URL="postgres://user:pass@host:5432/dbname?sslmode=require"
JWT_SECRET="your-super-secret-key-min-32-chars-long"
PORT="9000"  # Optional
```

---

## 📝 Common Tasks

### Create New Admin (after first admin exists)
```bash
curl -X POST http://localhost:9000/admin/register \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "manager",
    "email": "manager@example.com", 
    "password": "SecurePass123",
    "role": "manager"
  }'
```

### Create Customer
```bash
curl -X POST http://localhost:9000/customers \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "mobile": 9876543210,
    "address": "123 Main St",
    "info": "Regular customer"
  }'
```

### Get All Customers
```bash
curl -X GET http://localhost:9000/customers \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 💡 Pro Tips

1. **Token in Environment**: Save your token to avoid retyping
   ```bash
   export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
   curl -H "Authorization: Bearer $TOKEN" http://localhost:9000/customers
   ```

2. **Use Postman/Insomnia**: Save token in environment variables

3. **Frontend Integration**: Store token in localStorage or httpOnly cookie

4. **Token Refresh**: Implement refresh tokens for production

5. **HTTPS Only**: Never send JWT over HTTP in production

---

## ❓ Troubleshooting

**"Authorization header required"**
→ Did you include the token in the header?

**"Invalid or expired token"**
→ Token expired (24h) or malformed. Login again.

**"Only admins can register new admin users"**
→ Your role is not "admin". Check with `/admin/me`

**"Cannot link same customer to each other"**
→ You're trying to set a customer as their own referrer

---

## 📚 Learn More

See [AUTHENTICATION.md](AUTHENTICATION.md) for:
- Detailed API documentation
- Security best practices
- Advanced topics
- Complete examples
- Next steps for production

---

**You're now ready to build a secure finance application! 🎉**

Key changes:
- ✅ Clear separation: customers vs admins
- ✅ Secure authentication with JWT
- ✅ All endpoints protected
- ✅ Industry-standard security practices
- ✅ Clean, maintainable code structure
