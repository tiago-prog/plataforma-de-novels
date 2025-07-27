package repository

import (
	"database/sql"
	"fmt"

	"github.com/tiago-prog/novels-api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	connection *sql.DB
}

// GetUserWorks returns the list of work IDs associated with a user
func (repo *UserRepository) GetUserWorks(userID int) ([]int, error) {
	var works []int
	rows, err := repo.connection.Query("SELECT work_id FROM user_works WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var workID int
		if err := rows.Scan(&workID); err != nil {
			return nil, err
		}
		works = append(works, workID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return works, nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{
		connection: db,
	}
}

func (repo *UserRepository) RegisterUser(user model.User) (int, error) {

	var id int

	query, err := repo.connection.Prepare("INSERT INTO users (username, email, password_hash, role, avatar, status, last_login_at, verified, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id")
	if err != nil {
		fmt.Println("Error preparing query:", err)
		return 0, err
	}
	defer query.Close()
	err = query.QueryRow(
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.Avatar,
		user.Status,
		user.LastLoginAt,
		user.Verified,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&id)

	return id, err
}

func (repo *UserRepository) GetAllUsers() ([]model.User, error) {
	users := make([]model.User, 0)
	query := "SELECT id, username, email, role, avatar, verified, status, last_login_at, created_at, updated_at FROM users"
	rows, err := repo.connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Role,
			&user.Avatar,
			&user.Verified,
			&user.Status,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepository) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	var lastLoginAt sql.NullTime

	query := "SELECT id, username, email, role, avatar, verified, status, last_login_at, created_at, updated_at FROM users WHERE email = $1"
	err := repo.connection.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.Avatar,
		&user.Verified,
		&user.Status,
		&lastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return model.User{}, err
	}

	// Set LastLoginAt if it exists
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

func (repo *UserRepository) Login(email, password string) (model.User, error) {
	var user model.User

	query := `
		SELECT id, username, email, password_hash, role, avatar, verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := repo.connection.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Avatar,
		&user.Verified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return model.User{}, fmt.Errorf("email not found: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return model.User{}, fmt.Errorf("invalid password")
	}
	// Limpar a senha antes de retornar o usuário
	user.Password = ""

	return user, nil
}

func (repo *UserRepository) SuspendUser(executorID int, targetID int) error {
	// Check if the executor is a moderator/admin
	query := "SELECT role FROM users WHERE id = $1 AND role IN ('moderator', 'admin')"
	var role string
	err := repo.connection.QueryRow(query, executorID).Scan(&role)
	if err != nil {
		return fmt.Errorf("error getting user role: %w", err)
	}

	// Check if the executor can suspend the target
	if role != "moderator" && role != "admin" {
		return fmt.Errorf("user does not have permission to suspend users")
	}
	// Suspend the target user
	query = "UPDATE users SET status = 'suspended' WHERE id = $1"
	_, err = repo.connection.Exec(query, targetID)
	if err != nil {
		return fmt.Errorf("error suspending user: %w", err)
	}
	return nil
}

func (repo *UserRepository) CreateAdminIfNotExists() error {
	adminEmail := "admin@example.com"

	// Verifica se o admin já existe
	_, err := repo.GetUserByEmail(adminEmail)
	if err == nil {
		return nil // Admin já existe
	}

	// Cria o admin se não existir
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := model.User{
		Username: "admin",
		Email:    adminEmail,
		Password: string(hashedPassword),
		Role:     model.RoleAdmin,
		Verified: true,
		Status:   model.UserStatusActive,
	}

	_, err = repo.RegisterUser(admin)
	return err
}
