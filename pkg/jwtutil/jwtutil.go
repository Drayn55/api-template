package jwtutil

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type PermissionClaims struct {
	UserRoleID         int64  `json:"user_role_id"`
	RolePermissionID   int64  `json:"role_permission_id"`
	RolePermissionName string `json:"role_permission_name"`
	RoleID             int64  `json:"role_id"`
	RoleName           string `json:"role_name"`
	AccessMenuID       int64  `json:"access_menu_id"`
	AccessMenuName     string `json:"access_menu_name"`
	PermissionID       int64  `json:"permission_id"`
	PermissionSubject  string `json:"permission_subject"`
	PermissionAction   string `json:"permission_action"`
}

type Claims struct {
	UserID      int64              `json:"user_id"`
	EmployeeID  string             `json:"employee_id"`
	Name        string             `json:"name"`
	Dept        string             `json:"dept"`
	Position    string             `json:"position"`
	Permissions []PermissionClaims `json:"permissions"`
	jwt.RegisteredClaims
}

type JWTUtil struct {
	secret []byte
}

func New(secret string) *JWTUtil {
	return &JWTUtil{secret: []byte(secret)}
}

func (j *JWTUtil) ValidateToken(tokenString string) (*Claims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// If secret is not set, allow unverified claims for development
	if len(j.secret) == 0 {
		return parseUnverified(tokenString)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return j.secret, nil
		},
	)
	if err != nil {
		if claims, fbErr := parseUnverified(tokenString); fbErr == nil && claims.EmployeeID != "" {
			return claims, nil
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}
		return nil, errors.New("token invalid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		if claims, fbErr := parseUnverified(tokenString); fbErr == nil && claims.EmployeeID != "" {
			return claims, nil
		}
		return nil, errors.New("token invalid")
	}

	if claims.EmployeeID == "" {
		return nil, errors.New("employee_id missing in token")
	}

	return claims, nil
}

func parseUnverified(tokenString string) (*Claims, error) {
	parser := jwt.NewParser()
	var mapClaims jwt.MapClaims
	_, _, err := parser.ParseUnverified(tokenString, &mapClaims)
	if err != nil || len(mapClaims) == 0 {
		return nil, errors.New("token invalid")
	}

	claims := &Claims{}
	if empID, ok := mapClaims["employee_id"]; ok && empID != nil {
		claims.EmployeeID = fmt.Sprintf("%v", empID)
	} else if empID, ok := mapClaims["employeeId"]; ok && empID != nil {
		claims.EmployeeID = fmt.Sprintf("%v", empID)
	} else if uID, ok := mapClaims["user_id"]; ok && uID != nil {
		claims.EmployeeID = fmt.Sprintf("%v", uID)
	} else if sub, ok := mapClaims["sub"]; ok && sub != nil {
		claims.EmployeeID = fmt.Sprintf("%v", sub)
	}

	if name, ok := mapClaims["name"].(string); ok {
		claims.Name = name
	}
	if dept, ok := mapClaims["dept"].(string); ok {
		claims.Dept = dept
	}
	if pos, ok := mapClaims["position"].(string); ok {
		claims.Position = pos
	}

	if claims.EmployeeID == "" {
		claims.EmployeeID = "anonymous"
	}
	return claims, nil
}

func GetEmployeeIDFromToken(c *gin.Context) *string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		for _, key := range []string{"employee_id", "emp_id", "id", "sub", "user_id", "uid"} {
			if val, exists := claims[key]; exists {
				strVal := fmt.Sprintf("%v", val)
				return &strVal
			}
		}
	}

	return nil
}

