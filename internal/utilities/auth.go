// This file contains the utilities for authentication and authorization
//
// The utilities here are as follows:
// - UserType
// - UserGroups
// - AuthGroups
// - CheckPasswordHash
// - HashPassword
// - getClientFromRequest
// - getAuthTypeFromSession
// - ClientRequired
// - UserRequired
// - mergeSets
// - NewUserGroups
// - NewAuthGroups
// - printSessionValues

package utilities

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"waitress-backend/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWT secret key for token validation
var jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

// JWTClaims represents the claims stored in JWT tokens
type JWTClaims struct {
	UserID   uint   `json:"userID"`
	Email    string `json:"email"`
	AuthType string `json:"authType"`
	jwt.RegisteredClaims
}

// validateJWTToken verifies and parses a JWT token, returning the claims
func validateJWTToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Check if token is expired
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			return nil, errors.New("token is expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// extractJWTFromHeader extracts JWT token from Authorization header
func extractJWTFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("invalid authorization header format")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

// getAuthTypeFromJWT extracts auth type from JWT token
func getAuthTypeFromJWT(c *gin.Context) (string, error) {
	tokenString, err := extractJWTFromHeader(c)
	if err != nil {
		return "", err
	}

	claims, err := validateJWTToken(tokenString)
	if err != nil {
		return "", err
	}

	return claims.AuthType, nil
}

// UserType defines the different types of users in the system.
type UserType string

// List of user types in the system. This can be expanded as needed.
const (
	Dev        UserType = "dev"
	AdminSuper UserType = "admin_super"
	Admin      UserType = "admin"
	StaffSuper UserType = "staff_super"
	Staff      UserType = "staff"
	Customer   UserType = "customer"
)

// UserGroups holds sets of UserType for different groups.
type UserGroups struct {
	Dev        map[UserType]struct{}
	Admin      map[UserType]struct{}
	Staff      map[UserType]struct{}
	All        map[UserType]struct{}
	AllOrdered []UserType
}

// Transform the map keys into a slice of strings for easier access.
func (ug *UserGroups) GetAdminTypes() []string {
	adminTypes := make([]string, 0, len(ug.Admin))
	for k := range ug.Admin {
		adminTypes = append(adminTypes, string(k))
	}
	return adminTypes
}

// Transform the map keys into a slice of strings for easier access.
func (ug *UserGroups) GetStaffTypes() []string {
	staffTypes := make([]string, 0, len(ug.Staff))
	for k := range ug.Staff {
		staffTypes = append(staffTypes, string(k))
	}
	return staffTypes
}

// AuthGroups defines the permission hierarchy for each user type.
type AuthGroups struct {
	Dev      map[string]map[UserType]struct{}
	Admin    map[string]map[UserType]struct{}
	Staff    map[string]map[UserType]struct{}
	Customer map[string]map[UserType]struct{}
	All      map[UserType]struct{}
}

// printSessionValues prints the values stored in the session for debugging purposes.
func printSessionValues(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")         // Replace "userID" with the actual key you used to store the user ID
	userAuthType := session.Get("authType") // Similarly, for user's auth type

	fmt.Println("UserID in session:", userID)
	fmt.Println("User AuthType in session:", userAuthType)
}

// NewUserGroups initializes and returns a new UserGroups with predefined user types.
func NewUserGroups() UserGroups {
	allUsers := map[UserType]struct{}{
		Dev: {}, AdminSuper: {}, Admin: {}, StaffSuper: {}, Staff: {}, Customer: {},
	}
	return UserGroups{
		Dev:        map[UserType]struct{}{Dev: {}},
		Admin:      map[UserType]struct{}{AdminSuper: {}, Admin: {}},
		Staff:      map[UserType]struct{}{StaffSuper: {}, Staff: {}},
		All:        allUsers,
		AllOrdered: []UserType{Dev, AdminSuper, Admin, StaffSuper, Staff, Customer},
	}
}

// NewAuthGroups initializes and returns a new AuthGroups with permission hierarchies.
func NewAuthGroups(ug UserGroups) AuthGroups {
	return AuthGroups{
		Dev: map[string]map[UserType]struct{}{
			"all": ug.Dev,
		},
		Admin: map[string]map[UserType]struct{}{
			"super": mergeSets(ug.Dev, map[UserType]struct{}{AdminSuper: {}}),
			"all":   mergeSets(ug.Dev, ug.Admin),
		},
		Staff: map[string]map[UserType]struct{}{
			"super": mergeSets(mergeSets(ug.Dev, ug.Admin), map[UserType]struct{}{StaffSuper: {}}),
			"all":   mergeSets(mergeSets(ug.Dev, ug.Admin), ug.Staff),
		},
		Customer: map[string]map[UserType]struct{}{
			"all": mergeSets(mergeSets(mergeSets(ug.Dev, ug.Admin), ug.Staff), ug.All),
		},
		All: ug.All,
	}
}

// mergeSets combines multiple sets of UserType into a single set.
func mergeSets(sets ...map[UserType]struct{}) map[UserType]struct{} {
	result := make(map[UserType]struct{})
	for _, set := range sets {
		for key := range set {
			result[key] = struct{}{}
		}
	}
	return result
}

// CheckPasswordHash compares a password with its hash and returns true if they match.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashPassword generates a hashed password from a plaintext password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// getClientFromRequest retrieves the client from the Gin session.
func getClientFromRequest(c *gin.Context) (models.APIClient, error) {
	session := sessions.Default(c)
	client, ok := session.Get("client").(models.APIClient)
	if !ok {
		return models.APIClient{}, errors.New("client not found in session")
	}
	return client, nil
}

// getAuthTypeFromSession retrieves the auth type from the Gin session.
func getAuthTypeFromSession(c *gin.Context) (string, error) {
	session := sessions.Default(c)
	authType, ok := session.Get("authType").(string)
	if !ok || authType == "" {
		return "", errors.New("auth type not found in session")
	}
	return authType, nil
}

// DEPRECATED: Use UserRequired instead. This will be updated to handle both client and user authentication.
// Gin middleware that ensures the request is made by an authorized client.
func ClientRequired(clientTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		client, err := getClientFromRequest(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Client authentication required"})
			c.Abort() // Abort the request chain
			return
		}

		// Check if client type is within the allowed types, if specified
		if len(clientTypes) > 0 {
			permitted := false
			for _, t := range clientTypes {
				if client.ClientType == t {
					permitted = true
					break
				}
			}
			if !permitted {
				c.JSON(http.StatusForbidden, gin.H{"error": "Client type not permitted"})
				c.Abort() // Abort the request chain
				return
			}
		}

		c.Next() // Continue down the chain if all checks pass
	}
}

// Gin middleware that ensures the request is made by an authorized user.
// Supports both JWT token authentication (mobile) and session authentication (web)
func UserRequired(authGroups AuthGroups, group, subgroup string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var authType string
		var err error

		// Try JWT authentication first (for mobile apps)
		authType, err = getAuthTypeFromJWT(c)
		if err != nil {
			// JWT authentication failed, try session authentication (for web)
			authType, err = getAuthTypeFromSession(c)
			if err != nil {
				fmt.Println("Both JWT and session authentication failed:", err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
				c.Abort()
				return
			}
			fmt.Println("Authenticated via session with auth type:", authType)
		} else {
			fmt.Println("Authenticated via JWT with auth type:", authType)
		}

		var allowedUsers map[UserType]struct{}
		var ok bool

		switch group {
		case "Dev":
			allowedUsers, ok = authGroups.Dev[subgroup]
		case "Admin":
			allowedUsers, ok = authGroups.Admin[subgroup]
		case "Staff":
			allowedUsers, ok = authGroups.Staff[subgroup]
		case "Customer":
			allowedUsers, ok = authGroups.Customer[subgroup]
		default:
			fmt.Println("Invalid user group:", group)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user group"})
			c.Abort()
			return
		}

		if !ok {
			fmt.Println("Invalid subgroup configuration:", subgroup)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid subgroup configuration"})
			c.Abort()
			return
		}

		if _, permitted := allowedUsers[UserType(authType)]; !permitted {
			fmt.Println("User type not permitted:", authType)
			c.JSON(http.StatusForbidden, gin.H{"error": "User type not permitted"})
			c.Abort()
			return
		}

		fmt.Println("User is permitted, proceeding with auth type:", authType)
		c.Next()
	}
}
