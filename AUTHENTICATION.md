# Finance App - Authentication & Setup Guide

## 📚 Best Practices Learned

### 1. **Clear Naming Conventions**
- **Customers vs Admins**: Renamed "users" to "customers" to avoid confusion
  - `customers` table = Your finance customers/clients
  - `admins` table = System administrators who can log in
- **Consistent Naming**: All related files, functions, and endpoints use consistent naming
  - File: `customers.go`, `customerHelper.go`
  - Functions: `getCustomer()`, `createCustomer()`
  - Endpoints: `/customers`, `/customers/{id}`

### 2. **Separation of Concerns**
- **Authentication** (auth.go) = Handles login, JWT tokens, middleware
- **Business Logic** (customers.go, handouts.go) = Handles your finance operations
- **Database Queries** (querys.go) = Centralized SQL queries
- **Constants** (constants.go) = Reusable messages and values

### 3. **Security Best Practices**
- **JWT Authentication**: Tokens expire after 24 hours
- **Password Hashing**: bcrypt with cost factor 14
- **Protected Routes**: All sensitive endpoints require authentication
- **Role-Based Access**: Admin roles (admin, manager, viewer) for future permissions

---

## 🚀 Setup Instructions

### 1. Run Database Migrations

```bash
# Run migrations in order
psql $DATABASE_URL -f sql/migration-1.sql
psql $DATABASE_URL -f sql/migration-2.sql
psql $DATABASE_URL -f sql/migration-3.sql
psql $DATABASE_URL -f sql/migration-4.sql
psql $DATABASE_URL -f sql/migration-5.sql  # Creates admins table
psql $DATABASE_URL -f sql/migration-6.sql  # Renames users to customers
```

### 2. Set Environment Variables

```bash
export DATABASE_URL="your-postgres-connection-string"
export JWT_SECRET="your-super-secret-jwt-key-change-this"  # Important!
export PORT="9000"  # Optional, defaults to 9000
```

### 3. Create Your First Admin

**Option A: Temporarily expose the register endpoint**
```go
// In main.go, temporarily move this line outside protected routes:
r.HandleFunc("/admin/register", registerAdmin).Methods("POST")
```

Then:
```bash
curl -X POST http://localhost:9000/admin/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "your-secure-password",
    "role": "admin"
  }'
```

**Remember to move the register endpoint back inside protected routes after creating your first admin!**

### 4. Run the Application

```bash
go run .
```

---

## 🔐 API Authentication Guide

### Login

```bash
curl -X POST http://localhost:9000/admin/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "your-password"
  }'
```

**Response:**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresAt": "2025-12-15T10:30:00Z",
    "admin": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin"
    }
  },
  "message": "Login successful"
}
```

### Using the Token

Include the token in all subsequent requests:

```bash
curl -X GET http://localhost:9000/customers \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## 📝 API Endpoints

### Public Endpoints (No Authentication)
- `POST /admin/login` - Admin login

### Protected Endpoints (Requires Authentication)

#### Admin Management
- `POST /admin/register` - Register new admin (admin role only)
- `GET /admin/me` - Get current admin info

#### Customer Management
- `GET /customers` - List all customers
- `POST /customers` - Create new customer
- `GET /customers/{id}` - Get specific customer
- `PUT /customers/{id}` - Update customer
- `DELETE /customers/{id}` - Delete customer
- `GET /customers/{id}/handouts` - Get customer's handouts
- `GET /customers/{id}/referred-by` - Get who referred this customer
- `POST /customers/{id}/referral` - Link customer referral

#### Handout Management
- `GET /handouts` - List all handouts
- `POST /handouts` - Create new handout
- `GET /handouts/{id}` - Get specific handout
- `PUT /handouts/{id}` - Update handout
- `DELETE /handouts/{id}` - Delete handout
- `GET /handouts/{id}/collections` - Get handout collections

#### Collection Management
- `GET /collections` - List all collections
- `POST /collections` - Create new collection
- `PUT /collections/{id}` - Update collection
- `DELETE /collections/{id}` - Delete collection

---

## 💡 Common Patterns

### Creating a Customer
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

### Creating a Handout
```bash
curl -X POST http://localhost:9000/handouts \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customerId": 1,
    "amount": 10000,
    "date": "2025-12-14T00:00:00Z",
    "status": "ACTIVE",
    "bond": true
  }'
```

---

## 🔍 Understanding the Code Structure

```
middleware-finance-app/
├── auth.go              # Authentication logic (JWT, login, middleware)
├── customers.go         # Customer CRUD operations
├── customerHelper.go    # Helper functions for customers
├── handouts.go          # Handout operations
├── handoutHelper.go     # Helper functions for handouts
├── collections.go       # Collection operations
├── querys.go            # SQL queries
├── constants.go         # Constants and messages
├── Interfaces.go        # Data structures (Customer, Handout, etc.)
├── main.go              # Application entry point & routes
└── sql/
    ├── migration-5.sql  # Creates admins table
    └── migration-6.sql  # Renames users to customers
```

---

## 🎓 Learning Points

### 1. Why JWT?
- **Stateless**: Server doesn't need to store sessions
- **Scalable**: Works across multiple servers
- **Secure**: Signed tokens prevent tampering
- **Industry Standard**: Widely used and well-tested

### 2. Why bcrypt?
- **Slow by design**: Makes brute-force attacks impractical
- **Salted**: Each password hash is unique
- **Adaptive**: Can increase cost factor as computers get faster

### 3. Why Middleware?
- **DRY Principle**: Authentication logic in one place
- **Clean Code**: Business logic doesn't handle auth
- **Easy to Modify**: Change auth strategy in one file

### 4. Why Separate Tables?
- **Clear Data Model**: Admins and customers are different entities
- **Different Fields**: Admins have passwords, customers don't
- **Security**: Customer data never mixes with admin credentials
- **Flexibility**: Different access patterns for each type

---

## 🚨 Important Security Notes

1. **Change JWT_SECRET**: Never use default secret in production
2. **HTTPS Only**: Always use HTTPS in production
3. **Token Expiry**: Tokens expire after 24 hours
4. **Password Strength**: Enforce strong passwords (minimum 6 chars, but should be longer)
5. **Rate Limiting**: Consider adding rate limiting for login attempts
6. **Audit Logs**: Consider logging admin actions

---

## 🐛 Troubleshooting

### "Invalid or expired token"
- Token expired (24 hours) - login again
- Token malformed - check Authorization header format: `Bearer TOKEN`
- JWT_SECRET changed - previous tokens are now invalid

### "Admin account is not active"
- Admin was deactivated in database
- Check `active` column in `admins` table

### "Only admins can register new admin users"
- Only users with `role = 'admin'` can create new admins
- Check your current user's role with `GET /admin/me`

---

## 📊 Database Schema

### admins table
```sql
id            SERIAL PRIMARY KEY
username      VARCHAR(50) UNIQUE NOT NULL
email         VARCHAR(255) UNIQUE NOT NULL
password_hash VARCHAR(255) NOT NULL
role          VARCHAR(20) NOT NULL (admin/manager/viewer)
active        BOOLEAN NOT NULL DEFAULT true
created_at    TIMESTAMP
updated_at    TIMESTAMP
```

### customers table (formerly users)
```sql
id            SERIAL PRIMARY KEY
name          VARCHAR(255) NOT NULL
mobile        BIGINT UNIQUE NOT NULL
address       TEXT
info          TEXT
referred_by   INTEGER (references customers.id)
created_at    TIMESTAMP
updated_at    TIMESTAMP
```

---

## 🎯 Next Steps

1. **Add Role-Based Permissions**: Use the `role` field to restrict certain actions
2. **Add Refresh Tokens**: Implement longer-lived refresh tokens
3. **Add Password Reset**: Email-based password reset flow
4. **Add Audit Logging**: Track who did what and when
5. **Add Rate Limiting**: Prevent brute-force attacks
6. **Add API Documentation**: Consider Swagger/OpenAPI
7. **Add Tests**: Write integration tests for auth flows

---

## 📞 Support

For questions or issues:
1. Check this README first
2. Review the code comments
3. Test with curl commands provided above
4. Check server logs for error messages
