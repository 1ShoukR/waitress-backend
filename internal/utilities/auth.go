package utilities

import (
	"errors"
	"fmt"
	"net/http"
	"waitress-backend/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserType string
const (
    Dev         UserType = "dev"
    AdminSuper  UserType = "admin_super"
    Admin       UserType = "admin"
    StaffSuper  UserType = "staff_super"
    Staff       UserType = "staff"
    Customer    UserType = "customer"
)

// UserGroups holds sets of UserType for different groups.
type UserGroups struct {
    Dev         map[UserType]struct{}
    Admin       map[UserType]struct{}
    Staff       map[UserType]struct{}
    All         map[UserType]struct{}
    AllOrdered  []UserType
}

// Transform the map keys into a slice of strings for easier access.
func (ug *UserGroups) GetAdminTypes() []string {
    adminTypes := make([]string, 0, len(ug.Admin))
    for k := range ug.Admin {
        adminTypes = append(adminTypes, string(k))
    }
    return adminTypes
}

func (ug *UserGroups) GetStaffTypes() []string {
    staffTypes := make([]string, 0, len(ug.Staff))
    for k := range ug.Staff {
        staffTypes = append(staffTypes, string(k))
    }
    return staffTypes
}

// AuthGroups defines the permission hierarchy for each user type.
type AuthGroups struct {
    Dev     map[string]map[UserType]struct{}
    Admin   map[string]map[UserType]struct{}
    Staff   map[string]map[UserType]struct{}
    Customer map[string]map[UserType]struct{}
    All     map[UserType]struct{}
}

func printSessionValues(c *gin.Context) {
    session := sessions.Default(c)
    userID := session.Get("userID") // Replace "userID" with the actual key you used to store the user ID
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
        Dev:         map[UserType]struct{}{Dev: {}},
        Admin:       map[UserType]struct{}{AdminSuper: {}, Admin: {}},
        Staff:       map[UserType]struct{}{StaffSuper: {}, Staff: {}},
        All:         allUsers,
        AllOrdered:  []UserType{Dev, AdminSuper, Admin, StaffSuper, Staff, Customer},
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


func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

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

func UserRequired(authGroups AuthGroups, group, subgroup string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authType, err := getAuthTypeFromSession(c)
        if err != nil {
            fmt.Println("Error getting auth type from session:", err)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
            c.Abort()
            return
        }

        fmt.Println("Auth type fetched:", authType)

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

        fmt.Println("User is permitted, proceeding")
        c.Next()
    }
}
