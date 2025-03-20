package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type User struct {
	ID        int       `json:"id"`
	Nome      string    `json:"nome"`
	Email     string    `json:"email"`
	CPF       string    `json:"cpf"`
	Nascimento string   `json:"nascimento"`
	Senha     string    `json:"senha,omitempty"`
}

var db *sql.DB
var jwtKey = []byte("chave_secreta")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar .env")
	}
	var err error
	connStr := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME")
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := gin.Default()
	r.POST("/register", registerUser)
	r.POST("/login", loginUser)
	r.GET("/users", getUsers)

	r.Run(":8080")
}

func registerUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	query := "INSERT INTO users (nome, email, cpf, nascimento, senha) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, user.Nome, user.Email, user.CPF, user.Nascimento, user.Senha)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao registrar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário registrado com sucesso"})
}

func loginUser(c *gin.Context) {
	var user User
	var dbUser User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	query := "SELECT id, email, senha FROM users WHERE email = ?"
	err := db.QueryRow(query, user.Email).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Senha)
	if err != nil || user.Senha != dbUser.Senha {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func getUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, nome, email, cpf, nascimento FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuários"})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Nome, &user.Email, &user.CPF, &user.Nascimento); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar dados"})
			return
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}
