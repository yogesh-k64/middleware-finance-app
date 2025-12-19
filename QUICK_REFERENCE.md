# 🚀 Quick Reference: Database Insertion Method

## Create Your First Admin (Production Safe)

### Step 1: Generate Password Hash
```bash
go run scripts/create_admin_hash.go YourSecurePassword123
```

**Output will be:**
```
✅ Password hash generated successfully!

Password: YourSecurePassword123
Hash: $2a$14$HoQOsqOAHPJ/ME5WvR2AXe4ca6RbEiZAo2TCUVlEvPbnUt7plIZge

📋 SQL Insert Statement:
INSERT INTO admins (username, email, password_hash, role, active, created_at, updated_at)
VALUES (
    'admin',
    'admin@example.com',
    '$2a$14$HoQOsqOAHPJ/ME5WvR2AXe4ca6RbEiZAo2TCUVlEvPbnUt7plIZge',
    'admin',
    true,
    NOW(),
    NOW()
);
```

### Step 2: Insert into Database
```bash
# Connect to your database
psql $DATABASE_URL

# Paste the INSERT statement (update username/email if needed)
INSERT INTO admins (username, email, password_hash, role, active, created_at, updated_at)
VALUES ('admin', 'admin@example.com', '$2a$14$...your-hash...', 'admin', true, NOW(), NOW());

# Verify
SELECT id, username, email, role, active FROM admins;
```

### Step 3: Test Login
```bash
curl -X POST http://localhost:9000/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourSecurePassword123"}'
```

**Expected Response:**
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

---

## Frontend Implementation Summary

### 1. Save Token on Login
```javascript
// After successful login
localStorage.setItem('auth_token', data.data.token);
localStorage.setItem('admin_info', JSON.stringify(data.data.admin));
localStorage.setItem('token_expiry', data.data.expiresAt);
```

### 2. Include Token in Every Request
```javascript
const token = localStorage.getItem('auth_token');

fetch('http://localhost:9000/customers', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
})
```

### 3. Handle Token Expiry
```javascript
// If you get 401 response
if (response.status === 401) {
  localStorage.clear();
  window.location.href = '/login';
}
```

### 4. Logout
```javascript
localStorage.removeItem('auth_token');
localStorage.removeItem('admin_info');
localStorage.removeItem('token_expiry');
window.location.href = '/login';
```

---

## Complete Frontend Flow

```
User Opens App
      ↓
Check localStorage for token
      ↓
   ┌──┴──┐
   No   Yes
   ↓     ↓
 Login  Dashboard
   ↓
Enter username/password
   ↓
POST /admin/login
   ↓
Save token + admin info
   ↓
Redirect to Dashboard
   ↓
Every API call includes:
Authorization: Bearer <token>
   ↓
If 401 → Clear storage → Back to Login
```

---

## API Endpoints Reference

### Public
- `POST /admin/login` - Get JWT token

### Protected (Require `Authorization: Bearer <token>`)
**Admin:**
- `POST /admin/register` - Create new admin (admin role only)
- `GET /admin/me` - Get current admin info

**Customers:**
- `GET /customers` - List all
- `POST /customers` - Create
- `GET /customers/{id}` - Get one
- `PUT /customers/{id}` - Update
- `DELETE /customers/{id}` - Delete
- `GET /customers/{id}/handouts` - Customer's handouts
- `GET /customers/{id}/referred-by` - Who referred
- `POST /customers/{id}/referral` - Link referral

**Handouts:**
- `GET /handouts` - List all
- `POST /handouts` - Create
- `GET /handouts/{id}` - Get one
- `PUT /handouts/{id}` - Update
- `DELETE /handouts/{id}` - Delete
- `GET /handouts/{id}/collections` - Handout collections

**Collections:**
- `GET /collections` - List all
- `POST /collections` - Create
- `PUT /collections/{id}` - Update
- `DELETE /collections/{id}` - Delete

---

## Common Issues & Solutions

### "Authorization header required"
```javascript
// Make sure you include the header:
headers: {
  'Authorization': `Bearer ${token}`
}
```

### "Invalid or expired token"
```javascript
// Token expired (24h) - user must login again
localStorage.clear();
window.location.href = '/login';
```

### CORS Error
```javascript
// Backend already configured for:
// - http://localhost:5173
// - https://yogesh-k64.github.io

// If using different origin, update main.go:
allowedOrigins := handlers.AllowedOrigins([]string{
  "http://localhost:5173",
  "https://yogesh-k64.github.io",
  "http://localhost:3000",  // Add your origin
})
```

---

## Security Checklist

✅ **Never store passwords in code**
✅ **Always use HTTPS in production**
✅ **Set strong JWT_SECRET** (min 32 chars)
✅ **Tokens expire after 24 hours**
✅ **Use httpOnly cookies for production** (better than localStorage)
✅ **Validate on backend** (never trust frontend)

---

## Next Steps

1. ✅ Create admin with hash script
2. ✅ Insert into database
3. ✅ Test login with curl
4. 📝 Build login page in frontend
5. 📝 Save token in localStorage
6. 📝 Create protected routes
7. 📝 Build customer management UI
8. 📝 Build handout management UI

---

## Complete Example: React Login Component

```jsx
import { useState } from 'react';

function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleLogin = async (e) => {
    e.preventDefault();
    
    try {
      const response = await fetch('http://localhost:9000/admin/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Login failed');
      }

      // Save token
      localStorage.setItem('auth_token', data.data.token);
      localStorage.setItem('admin_info', JSON.stringify(data.data.admin));
      
      // Redirect
      window.location.href = '/dashboard';
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <form onSubmit={handleLogin}>
      <input 
        value={username} 
        onChange={(e) => setUsername(e.target.value)}
        placeholder="Username"
        required
      />
      <input 
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
        required
      />
      {error && <div>{error}</div>}
      <button type="submit">Login</button>
    </form>
  );
}
```

---

**You're all set! 🎉**

For detailed frontend implementation, see: [FRONTEND_GUIDE.md](FRONTEND_GUIDE.md)
