package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Admin represents an admin user who can access the system
type Admin struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"` // e.g., "admin", "manager", "viewer"
	Active       bool      `json:"active"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// JWT Claims structure
type Claims struct {
	AdminID  int    `json:"admin_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Login request structure
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register request structure
type RegisterAdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"` // Optional, defaults to "admin"
}

// Login response structure
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	Admin     AdminInfo `json:"admin"`
}

type AdminInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// Get JWT secret from environment or use a default (change in production!)
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-this-in-production"
	}
	return []byte(secret)
}

// Hash password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Check if password matches hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate a random API key
func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Generate JWT token for admin
func generateToken(adminID int, username, role string) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		AdminID:  adminID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "middleware-finance-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// Verify JWT token and return claims
func verifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// Authentication middleware - requires JWT token
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sendErrorResponse(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if it's a Bearer token (JWT)
		if !strings.HasPrefix(authHeader, "Bearer ") {
			sendErrorResponse(w, "Invalid authorization format. Use 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := verifyToken(tokenString)
		if err != nil {
			sendErrorResponse(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Verify admin is still active
		var active bool
		err = db.QueryRow("SELECT active FROM admins WHERE id = $1", claims.AdminID).Scan(&active)
		if err != nil || !active {
			sendErrorResponse(w, "Admin account is not active", http.StatusUnauthorized)
			return
		}

		// Add admin info to request context
		ctx := context.WithValue(r.Context(), "adminID", claims.AdminID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "role", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Login handler for admins
func adminLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate input
	if req.Username == "" || req.Password == "" {
		sendErrorResponse(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Get admin from database
	var admin Admin
	err = db.QueryRow(`
		SELECT id, username, password_hash, role, active 
		FROM admins 
		WHERE username = $1
	`, req.Username).Scan(&admin.ID, &admin.Username, &admin.PasswordHash, &admin.Role, &admin.Active)

	if err != nil {
		sendErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check if admin is active
	if !admin.Active {
		sendErrorResponse(w, "Admin account is not active", http.StatusUnauthorized)
		return
	}

	// Verify password
	if !checkPasswordHash(req.Password, admin.PasswordHash) {
		sendErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, expiresAt, err := generateToken(admin.ID, admin.Username, admin.Role)
	if err != nil {
		sendErrorResponse(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Send response
	response := DataResp[LoginResponse]{
		D: LoginResponse{
			Token:     token,
			ExpiresAt: expiresAt,
			Admin: AdminInfo{
				ID:       admin.ID,
				Username: admin.Username,
				Role:     admin.Role,
			},
		},
		Msg: "Login successful",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Register new admin (should be protected - only existing admins can create new admins)
func registerAdmin(w http.ResponseWriter, r *http.Request) {
	// Get admin info from context (to verify they have permission)
	role, ok := r.Context().Value("role").(string)
	if !ok || role != "admin" {
		sendErrorResponse(w, "Only admins can register new admin users", http.StatusForbidden)
		return
	}

	var req RegisterAdminRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate input
	if req.Username == "" || req.Password == "" {
		sendErrorResponse(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		sendErrorResponse(w, "Password must be at least 6 characters long", http.StatusBadRequest)
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "admin"
	}

	// Validate role
	validRoles := map[string]bool{"admin": true, "manager": true, "viewer": true}
	if !validRoles[req.Role] {
		sendErrorResponse(w, "Invalid role. Must be 'admin', 'manager', or 'viewer'", http.StatusBadRequest)
		return
	}

	// Hash the password
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		sendErrorResponse(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// Insert admin
	var adminID int
	err = db.QueryRow(`
		INSERT INTO admins (username, password_hash, role, active, created_at, updated_at)
		VALUES ($1, $2, $3, true, NOW(), NOW())
		RETURNING id
	`, req.Username, passwordHash, req.Role).Scan(&adminID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			sendErrorResponse(w, "Username already exists", http.StatusConflict)
		} else {
			sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Send response
	response := DataResp[AdminInfo]{
		D: AdminInfo{
			ID:       adminID,
			Username: req.Username,
			Role:     req.Role,
		},
		Msg: "Admin registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Get current admin info
func getCurrentAdmin(w http.ResponseWriter, r *http.Request) {
	adminID, ok := r.Context().Value("adminID").(int)
	if !ok {
		sendErrorResponse(w, "Admin not authenticated", http.StatusUnauthorized)
		return
	}

	var admin Admin
	err := db.QueryRow(`
		SELECT id, username, role, active, created_at, updated_at 
		FROM admins 
		WHERE id = $1
	`, adminID).Scan(&admin.ID, &admin.Username, &admin.Role, &admin.Active, &admin.CreatedAt, &admin.UpdatedAt)

	if err != nil {
		sendErrorResponse(w, "Admin not found", http.StatusNotFound)
		return
	}

	response := DataResp[Admin]{
		D:   admin,
		Msg: "Admin retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Get all admins (admin only)
func getAllAdmins(w http.ResponseWriter, r *http.Request) {
	// Check if requester is admin
	role, ok := r.Context().Value("role").(string)
	if !ok || role != "admin" {
		sendErrorResponse(w, "Only admins can view all users", http.StatusForbidden)
		return
	}

	rows, err := db.Query(`
		SELECT id, username, role, active, created_at, updated_at 
		FROM admins 
		ORDER BY id
	`)
	if err != nil {
		sendErrorResponse(w, "Failed to fetch admins", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var admins []Admin
	for rows.Next() {
		var admin Admin
		err := rows.Scan(&admin.ID, &admin.Username, &admin.Role, &admin.Active, &admin.CreatedAt, &admin.UpdatedAt)
		if err != nil {
			sendErrorResponse(w, "Failed to scan admin", http.StatusInternalServerError)
			return
		}
		admins = append(admins, admin)
	}

	response := DataResp[[]Admin]{
		D:   admins,
		Msg: "Admins retrieved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Update admin
func updateAdmin(w http.ResponseWriter, r *http.Request) {
	// Check if requester is admin
	role, ok := r.Context().Value("role").(string)
	if !ok || role != "admin" {
		sendErrorResponse(w, "Only admins can update users", http.StatusForbidden)
		return
	}

	// Get admin ID from URL
	vars := mux.Vars(r)
	adminID := vars["id"]

	// Check if target is the super admin
	var targetUsername string
	err := db.QueryRow("SELECT username FROM admins WHERE id = $1", adminID).Scan(&targetUsername)
	if err != nil {
		sendErrorResponse(w, "Admin not found", http.StatusNotFound)
		return
	}

	var req struct {
		Username *string `json:"username"`
		Password *string `json:"password"`
		Role     *string `json:"role"`
		Active   *bool   `json:"active"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Prevent modifying super admin's critical fields
	if targetUsername == "admin" {
		if req.Role != nil || req.Active != nil {
			sendErrorResponse(w, "Cannot modify role or active status of super admin", http.StatusForbidden)
			return
		}
	}

	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}
	paramCount := 1

	if req.Username != nil {
		updates = append(updates, "username = $"+fmt.Sprint(paramCount))
		args = append(args, *req.Username)
		paramCount++
	}

	if req.Password != nil {
		hashedPassword, err := hashPassword(*req.Password)
		if err != nil {
			sendErrorResponse(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		updates = append(updates, "password_hash = $"+fmt.Sprint(paramCount))
		args = append(args, hashedPassword)
		paramCount++
	}

	if req.Role != nil {
		// Validate role
		validRoles := map[string]bool{"admin": true, "manager": true, "viewer": true}
		if !validRoles[*req.Role] {
			sendErrorResponse(w, "Invalid role. Must be: admin, manager, or viewer", http.StatusBadRequest)
			return
		}
		updates = append(updates, "role = $"+fmt.Sprint(paramCount))
		args = append(args, *req.Role)
		paramCount++
	}

	if req.Active != nil {
		updates = append(updates, "active = $"+fmt.Sprint(paramCount))
		args = append(args, *req.Active)
		paramCount++
	}

	if len(updates) == 0 {
		sendErrorResponse(w, "No fields to update", http.StatusBadRequest)
		return
	}

	// Add updated_at
	updates = append(updates, "updated_at = NOW()")

	// Add adminID as last parameter
	args = append(args, adminID)

	query := "UPDATE admins SET " + strings.Join(updates, ", ") + " WHERE id = $" + fmt.Sprint(paramCount)

	result, err := db.Exec(query, args...)
	if err != nil {
		sendErrorResponse(w, "Failed to update admin", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		sendErrorResponse(w, "Admin not found", http.StatusNotFound)
		return
	}

	// Fetch updated admin
	var admin Admin
	err = db.QueryRow(`
		SELECT id, username, role, active, created_at, updated_at 
		FROM admins 
		WHERE id = $1
	`, adminID).Scan(&admin.ID, &admin.Username, &admin.Role, &admin.Active, &admin.CreatedAt, &admin.UpdatedAt)

	if err != nil {
		sendErrorResponse(w, "Failed to fetch updated admin", http.StatusInternalServerError)
		return
	}

	response := DataResp[Admin]{
		D:   admin,
		Msg: "User updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Delete admin
func deleteAdmin(w http.ResponseWriter, r *http.Request) {
	// Check if requester is admin
	role, ok := r.Context().Value("role").(string)
	if !ok || role != "admin" {
		sendErrorResponse(w, "Only admins can delete users", http.StatusForbidden)
		return
	}

	// Get admin ID from URL
	vars := mux.Vars(r)
	adminID := vars["id"]

	// Check if target is the super admin
	var targetUsername string
	err := db.QueryRow("SELECT username FROM admins WHERE id = $1", adminID).Scan(&targetUsername)
	if err != nil {
		sendErrorResponse(w, "Admin not found", http.StatusNotFound)
		return
	}

	// Prevent deleting super admin
	if targetUsername == "admin" {
		sendErrorResponse(w, "Cannot delete the super admin account", http.StatusForbidden)
		return
	}

	// Prevent self-deletion
	currentAdminID, ok := r.Context().Value("adminID").(int)
	if ok && fmt.Sprint(currentAdminID) == adminID {
		sendErrorResponse(w, "Cannot delete your own account", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM admins WHERE id = $1", adminID)
	if err != nil {
		sendErrorResponse(w, "Failed to delete admin", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		sendErrorResponse(w, "Admin not found", http.StatusNotFound)
		return
	}

	response := MsgResp{
		Msg: "User deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
