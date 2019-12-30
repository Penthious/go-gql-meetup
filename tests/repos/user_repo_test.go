package repos

import (
	_ "github.com/lib/pq"

	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/khaiql/dbcleaner/engine"
	"github.com/penthious/go-gql-meetup/database"
	"github.com/penthious/go-gql-meetup/domain"
	"github.com/penthious/go-gql-meetup/domain/sevices"
	"github.com/penthious/go-gql-meetup/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"

)

type UserRepoSuite struct {
	suite.Suite
	Domain *domain.Domain
}

func (s *UserRepoSuite) SetupSuite() {
	fmt.Println("setup the suite")
	DB := database.New(&pg.Options{
		User: "tleffew",
		Password:"postgres",
		Database:"meetup_test",
	})

	graphqlDB := domain.DB{
		UserRepo:   database.NewUserRepo(DB),
		MeetupRepo: database.NewMeetupRepo(DB),
		DB: DB,
	}
	g := &domain.Domain{DB: graphqlDB}

	s.Domain = g


	DB.AddQueryHook(database.DBLogger{})

	PGDB := engine.NewPostgresEngine("postgres://tleffew:postgres@localhost:5432/meetup_test?sslmode=disable")
	Cleaner.SetEngine(PGDB)
}

func (s *UserRepoSuite) SetupTest() {
	Cleaner.Acquire("users")
}

func (s *UserRepoSuite) TearDownTest() {
	fmt.Println("teardown")
	Cleaner.Clean("users")
}

func (s *UserRepoSuite) TearDownSuite() {
	fmt.Println("teardown full suite")
	s.Domain.DB.DB.Close()
}

func (s *UserRepoSuite) TestUserRepo_Create() {
	user := &models.User{
		Username: "bob",
		Email: "bob@bob.com",
		Password: "password",
	}

	CreateTestUser(s, user)

	assert.NotEmpty(s.T(), user.ID)
	assert.Equal(s.T(), user.Username, "bob")
	assert.Equal(s.T(), user.Email, "bob@bob.com")
	assert.NotEqual(s.T(), user.Password, "password")
	//assert.Nil(s.T(), user.DeletedAt, "deleted_at should be nil")
}

func (s *UserRepoSuite) TestUserRepo_All() {
	CreateTestUser(s, &models.User{Username: "bob1", Email: "bob1@bob.com", Password: "password"})
	CreateTestUser(s, &models.User{Username: "bob2", Email: "bob2@bob.com", Password: "password"})
	CreateTestUser(s, &models.User{Username: "bob3", Email: "bob3@bob.com", Password: "password"})

	dbUsers, _ := s.Domain.DB.UserRepo.All()

	assert.NotEmpty(s.T(), dbUsers[0].ID)
	assert.Equal(s.T(), dbUsers[0].Username, "bob1")
	assert.Equal(s.T(), dbUsers[1].Username, "bob2")
	assert.Equal(s.T(), dbUsers[2].Username, "bob3")
	assert.Equal(s.T(), 3, len(dbUsers))
}

func (s *UserRepoSuite) TestUserRepo_GetByIDs() {
	user1 := CreateTestUser(s, &models.User{Username: "bob1", Email: "bob1@bob.com", Password: "password"})
	user2 := CreateTestUser(s, &models.User{Username: "bob2", Email: "bob2@bob.com", Password: "password"})
	ids := []string{user1.ID, user2.ID}

	dbUsers, _ := s.Domain.DB.UserRepo.GetByIDs(ids)

	assert.NotEqual(s.T(), dbUsers[0].ID, dbUsers[1].ID)
	assert.Equal(s.T(), 2, len(dbUsers))
}

func (s *UserRepoSuite) TestUserRepo_GetByKey_email() {
	CreateTestUser(s, &models.User{Username: "bob", Email: "bob@bob.com", Password: "password"})

	dbUser, _ := s.Domain.DB.UserRepo.GetByKey("email", "bob@bob.com")

	assert.NotEmpty(s.T(), dbUser.ID)
	assert.Equal(s.T(), dbUser.Username, "bob")
	assert.Equal(s.T(), dbUser.Email, "bob@bob.com")
	assert.NotEqual(s.T(), dbUser.Password, "password")
	//assert.Nil(s.T(), user.DeletedAt, "deleted_at should be nil")
}

func (s *UserRepoSuite) TestUserRepo_GetByKey_username() {
	CreateTestUser(s, &models.User{Username: "bob", Email: "bob@bob.com", Password: "password"})

	dbUser, _ := s.Domain.DB.UserRepo.GetByKey("username", "bob")

	assert.NotEmpty(s.T(), dbUser.ID)
	assert.Equal(s.T(), dbUser.Username, "bob")
	assert.Equal(s.T(), dbUser.Email, "bob@bob.com")
	assert.NotEqual(s.T(), dbUser.Password, "password")
	//assert.Nil(s.T(), user.DeletedAt, "deleted_at should be nil")
}

func (s *UserRepoSuite) TestUserRepo_GetByKey_id() {
	user := CreateTestUser(s, &models.User{Username: "bob", Email: "bob@bob.com", Password: "password"})

	dbUser, _ := s.Domain.DB.UserRepo.GetByKey("id", user.ID)

	assert.NotEmpty(s.T(), dbUser.ID)
	assert.Equal(s.T(), dbUser.Username, "bob")
	assert.Equal(s.T(), dbUser.Email, "bob@bob.com")
	assert.NotEqual(s.T(), dbUser.Password, "password")
	//assert.Nil(s.T(), user.DeletedAt, "deleted_at should be nil")
}

func (s *UserRepoSuite) TestUserRepo_GetByKey_err_no_rows() {

	_, err := s.Domain.DB.UserRepo.GetByKey("id", "1")

	assert.Error(s.T(), err, "No results")
}

func (s *UserRepoSuite) TestUserRepo_Update_username() {
	user := CreateTestUser(s, &models.User{})
	user.Username = "updatedBob"

	s.Domain.DB.UserRepo.Update(user)

	assert.NotEmpty(s.T(), user.ID)
	assert.Equal(s.T(), user.Username, "updatedBob")
	assert.Equal(s.T(), user.Email, "bob@bob.com")
	assert.NotEqual(s.T(), user.Password, "password")
	//assert.Nil(s.T(), user.DeletedAt, "deleted_at should be nil")
}
func (s *UserRepoSuite) TestUserRepo_Update_email() {
	user := CreateTestUser(s, &models.User{})
	user.Email = "updatedEmail@email.com"

	s.Domain.DB.UserRepo.Update(user)

	assert.NotEmpty(s.T(), user.ID)
	assert.Equal(s.T(), user.Username, "bob")
	assert.Equal(s.T(), user.Email, "updatedEmail@email.com")
	assert.NotEqual(s.T(), user.Password, "password")
	//assert.Nil(s.T(), user.DeletedAt, "deleted_at should be nil")
}


func TestUserRepoSuite(t *testing.T) {
	suite.Run(t, new(UserRepoSuite))
}

func CreateTestUser(s *UserRepoSuite, user *models.User) *models.User{
	if user.Username == "" {user.Username = "bob"}
	if user.Email == "" {user.Email = "bob@bob.com"}
	if user.Password == "" {user.Password = "password"}
	passwordHash, _ := sevices.SetPassword(user.Password)
	user.Password = *passwordHash
	s.Domain.DB.UserRepo.Create(user)

	return user
}