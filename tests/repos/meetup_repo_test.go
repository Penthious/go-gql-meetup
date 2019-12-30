package repos

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/khaiql/dbcleaner/engine"
	"github.com/penthious/go-gql-meetup/database"
	"github.com/penthious/go-gql-meetup/domain"
	"github.com/penthious/go-gql-meetup/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)


type MeetupRepoSuite struct {
	suite.Suite
	Domain *domain.Domain
}

func (s *MeetupRepoSuite) SetupSuite() {
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

func (s *MeetupRepoSuite) SetupTest() {
	Cleaner.Acquire("users")
}

func (s *MeetupRepoSuite) TearDownTest() {
	fmt.Println("teardown")
	Cleaner.Clean("users")
}

func (s *MeetupRepoSuite) TearDownSuite() {
	fmt.Println("teardown full suite")
	s.Domain.DB.DB.Close()
}

func (s *MeetupRepoSuite) TestMeetupRepo_All() {
	meetup1 := CreateTestMeetup(
		s,
		&models.Meetup{Name: "Test 1", Description: "Description 1"},
		&models.User{Username: "user1", Email: "email1@test.com"},
		)
	meetup2 := CreateTestMeetup(
		s,
		&models.Meetup{Name: "Test 2", Description: "Description 2"},
		&models.User{Username: "user2", Email: "email2@test.com"},
	)
	meetup3 := CreateTestMeetup(s,
		&models.Meetup{Name: "Test 3", Description: "Description 3"},
		&models.User{Username: "user3", Email: "email3@test.com"},
	)

	dbMeetups, _ := s.Domain.DB.MeetupRepo.All()
	fmt.Println(meetup1, meetup2, meetup3)

	assert.NotEmpty(s.T(), dbMeetups[0].ID)
	assert.Equal(s.T(), "Test 1", dbMeetups[0].Name)
	assert.Equal(s.T(), "Test 2", dbMeetups[1].Name)
	assert.Equal(s.T(), "Test 3", dbMeetups[2].Name)
	assert.Equal(s.T(), 3, len(dbMeetups))
}

func (s *MeetupRepoSuite) TestMeetupRepo_Create() {
	meetup := &models.Meetup{
		Name: "Test",
		Description: "Something",
	}
	user := &models.User{Username: "user", Email: "email@test.com"}


	CreateTestMeetup(s, meetup, user)

	assert.NotEmpty(s.T(), meetup.ID)
	assert.Equal(s.T(), "Test", meetup.Name)
	assert.Equal(s.T(), "Something", meetup.Description)
	//assert.Nil(s.T(), user.DeletedAt, "deleted_at should be nil")
}

func (s *MeetupRepoSuite) TestMeetupRepo_Delete() {
	meetup := &models.Meetup{
		Name: "Test",
		Description: "Something",
	}
	user := &models.User{Username: "user", Email: "email@test.com"}
	CreateTestMeetup(s, meetup, user)

	s.Domain.DB.MeetupRepo.Delete(meetup)
	_, err := s.Domain.DB.MeetupRepo.GetByKey("id", meetup.ID)

	assert.Error(s.T(), err, "No results")

}

func (s *MeetupRepoSuite) TestMeetupRepo_GetByIDs() {
	meetup1 := CreateTestMeetup(
		s,
		&models.Meetup{Name: "Test 1", Description: "Description 1"},
		&models.User{Username: "user1", Email: "email1@test.com"},
	)
	meetup2 := CreateTestMeetup(
		s,
		&models.Meetup{Name: "Test 2", Description: "Description 2"},
		&models.User{Username: "user2", Email: "email2@test.com"},
	)
	ids := []string{meetup1.ID, meetup2.ID}

	dbUsers, _ := s.Domain.DB.MeetupRepo.GetByIDs(ids)

	assert.NotEqual(s.T(), dbUsers[0].ID, dbUsers[1].ID)
	assert.Equal(s.T(), 2, len(dbUsers))
}

func (s *MeetupRepoSuite) TestMeetupRepo_GetByKey_email() {
	CreateTestMeetup(
		s,
		&models.Meetup{Name: "Test", Description: "Description"},
		&models.User{Username: "user", Email: "email@test.com"},
	)

	dbMeetup, _ := s.Domain.DB.MeetupRepo.GetByKey("name", "Test")

	assert.NotEmpty(s.T(), dbMeetup.ID)
	assert.Equal(s.T(), "Test", dbMeetup.Name)
	assert.Equal(s.T(), "Description", dbMeetup.Description)
	//assert.Nil(s.T(), user.DeletedAt, "deleted_at should be nil")
}
func (s *MeetupRepoSuite) TestMeetupRepo_GetMeetupsByFilter_name() {
	user := &models.User{
		Username: "Test",
		Email:    "test@test.com",
		Password: "password",
	}
	for i := 0; i <= 3; i++ {
		name := fmt.Sprintf("Test %v", i)
		CreateTestMeetup(
			s,
			&models.Meetup{Name: name, Description: "Description"},
			user,
		)
	}
	CreateTestMeetup(
		s,
		&models.Meetup{Name: "Should not match", Description: "Description"},
		user,
	)
	name := "test"
	filters := &models.MeetupFilterPayload{
		Name: &name,
	}

	meetups, _ := s.Domain.DB.MeetupRepo.GetMeetupsByFilter(filters)

	assert.NotEqual(s.T(), meetups[0].ID, meetups[1].ID)
	for i := 0; i < len(meetups); i++ {
		assert.Equal(s.T(), user.ID, meetups[i].UserID)
		assert.True(s.T(), strings.Contains(meetups[i].Name, "Test"))
	}
	assert.Equal(s.T(), 4, len(meetups))

}

func (s *MeetupRepoSuite) TestMeetupRepo_GetMeetupsByFilter_description() {
	user := &models.User{
		Username: "Test",
		Email:    "test@test.com",
		Password: "password",
	}
	for i := 0; i <= 3; i++ {
		name := fmt.Sprintf("Test %v", i)
		description  := fmt.Sprintf("description %v", i)
		CreateTestMeetup(
			s,
			&models.Meetup{Name: name, Description: description},
			user,
		)
	}
	CreateTestMeetup(
		s,
		&models.Meetup{Name: "Next test", Description: "Should not match"},
		user,
	)
	description := "cripti"
	filters := &models.MeetupFilterPayload{
		Description: &description,
	}

	meetups, _ := s.Domain.DB.MeetupRepo.GetMeetupsByFilter(filters)

	assert.NotEqual(s.T(), meetups[0].ID, meetups[1].ID)
	for i := 0; i < len(meetups); i++ {
		assert.Equal(s.T(), user.ID, meetups[i].UserID)
		assert.True(s.T(), strings.Contains(meetups[i].Description, "cripti"))
	}
	assert.Equal(s.T(), 4, len(meetups))

}
func (s *MeetupRepoSuite) TestMeetupRepo_GetMeetupsForUser() {
	user := &models.User{
		Username: "Test",
		Email:    "test@test.com",
		Password: "password",
	}
	s.Domain.DB.UserRepo.Create(user)
	for i := 0; i <= 3; i++ {
		name := fmt.Sprintf("Test %v", i)
		CreateTestMeetup(
			s,
			&models.Meetup{Name: name, Description: "Description"},
			user,
		)
	}

	meetups, _ := s.Domain.DB.MeetupRepo.GetMeetupsForUser(user.ID)

	assert.NotEqual(s.T(), meetups[0].ID, meetups[1].ID)
	for i := 0; i < len(meetups); i++ {
		assert.Equal(s.T(), user.ID, meetups[i].UserID)
	}
	assert.Equal(s.T(), 4, len(meetups))

}

func (s *MeetupRepoSuite) TestMeetupRepo_GetByKey_id() {
	meetup := CreateTestMeetup(
		s,
		&models.Meetup{Name: "Test", Description: "Description"},
		&models.User{Username: "user", Email: "email@test.com"},
	)

	dbMeetup, _ := s.Domain.DB.MeetupRepo.GetByKey("id", meetup.ID)

	assert.NotEmpty(s.T(), dbMeetup.ID)
	assert.Equal(s.T(), "Test", dbMeetup.Name)
	assert.Equal(s.T(), "Description", dbMeetup.Description)
	//assert.Nil(s.T(), user.DeletedAt, "deleted_at should be nil")
}

func (s *MeetupRepoSuite) TestMeetupRepo_Update_name() {
	meetup := CreateTestMeetup(
		s,
		&models.Meetup{Name: "Test", Description: "Description"},
		&models.User{Username: "user", Email: "email@test.com"},
	)
	meetup.Name = "Updated name"

	s.Domain.DB.MeetupRepo.Update(meetup)

	assert.NotEmpty(s.T(), meetup.ID)
	assert.Equal(s.T(), "Updated name", meetup.Name)
	assert.Equal(s.T(), "Description", meetup.Description)
	//assert.Nil(s.T(), user.DeletedAt, "deleted_at should be nil")
}
func (s *MeetupRepoSuite) TestMeetupRepo_Update_description() {
	meetup := CreateTestMeetup(
		s,
		&models.Meetup{Name: "Test", Description: "Description"},
		&models.User{Username: "user", Email: "email@test.com"},
	)
	meetup.Description = "Updated description"

	s.Domain.DB.MeetupRepo.Update(meetup)

	assert.NotEmpty(s.T(), meetup.ID)
	assert.Equal(s.T(), "Test", meetup.Name)
	assert.Equal(s.T(), "Updated description", meetup.Description)
	//assert.Nil(s.T(), user.DeletedAt, "deleted_at should be nil")
}


func TestMeetupRepoSuite(t *testing.T) {
	suite.Run(t, new(MeetupRepoSuite))
}

func CreateTestMeetup(s *MeetupRepoSuite, meetup *models.Meetup, user *models.User) *models.Meetup{
	if meetup.Name == "" {meetup.Name = "Test"}
	if meetup.Description == "" {meetup.Description = "New description"}
	if user.ID == "" {
		if user.Username == "" {user.Username = "bob"}
		if user.Email == "" {user.Email = "bob@bob.com"}
		user.Password = "password"
		s.Domain.DB.UserRepo.Create(user)
		meetup.UserID = user.ID
	}
	meetup.UserID = user.ID

	s.Domain.DB.MeetupRepo.Create(meetup)

	return meetup
}
