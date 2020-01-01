package resolvers

import (
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-pg/pg/v9"
	"github.com/joho/godotenv"
	"github.com/khaiql/dbcleaner/engine"
	_ "github.com/lib/pq"
	"github.com/penthious/go-gql-meetup/database"
	"github.com/penthious/go-gql-meetup/domain"
	"github.com/penthious/go-gql-meetup/graphql/dataloaders"
	"github.com/penthious/go-gql-meetup/graphql/resolvers"
	"github.com/penthious/go-gql-meetup/models"
	"github.com/stretchr/testify/suite"
	"testing"
)


type UserResolverSuite struct {
	suite.Suite
	Domain *domain.Domain
}
func (s *UserResolverSuite) SetupSuite() {
	fmt.Println("setup the suite")
	godotenv.Load("./environments/.env_test")
	//DB := database.New(&pg.Options{
	//	User:     os.Getenv("DATABASE_USER"),
	//	Password: os.Getenv("DATABASE_PASSWORD"),
	//	Database: os.Getenv("DATABASE"),
	//})
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
func (s *UserResolverSuite) SetupTest() {
	cmd := "TRUNCATE TABLE users RESTART IDENTITY CASCADE;"
	s.Domain.DB.DB.Exec(cmd)
	Cleaner.Acquire("users")
	//stmt, _ := s.Domain.DB.DB.Prepare("TRUNCATE 'users' RESTART IDENTITY;")
	//_, err := stmt.Exec()
	//if err != nil {
	//	fmt.Println(err)
	//}
}

func (s *UserResolverSuite) TearDownTest() {
	fmt.Println("teardown")
	Cleaner.Clean("users")
}

func (s *UserResolverSuite) TearDownSuite() {
	fmt.Println("teardown full suite")
	s.Domain.DB.DB.Close()
}

func SetupServer(d *domain.Domain) *client.Client{

	c := resolvers.Config{
		Resolvers: &resolvers.Resolver{Domain: *d},
	}
	srv := handler.NewDefaultServer(resolvers.NewExecutableSchema(c))


	return client.New(dataloaders.DataloaderMiddleware(d, srv))
}


//var Cleaner = dbcleaner.New()
func (s *UserResolverSuite) TestUserResolver_GetMeetups() {
	query := `
		query GetUser {
		  user(id: "1") {
			id
			email
			username
		  }
		}

	`
	user := &models.User{
		Username: "test",
		Email:    "test@test.com",
		Password: "password",
	}
	s.Domain.DB.UserRepo.Create(user)
	//s.Domain.DB.MeetupRepo.Create(&models.Meetup{
	//	Name:        "Test 11",
	//	Description: "Test",
	//	UserID:      user.ID,
	//})

	var resp struct {
		User *models.User
	}
	client := SetupServer(s.Domain)
	err := client.Post(query, &resp)
	respStr, _ := json.Marshal(resp)
	fmt.Println(respStr)
	if err != nil {
		fmt.Println(err)
		//t.Fatalf("error server: %v", err)
	}
	user.Password = ""

	s.Equal(user, resp.User)
	//assert.Equal(s.T(), `{"User":{"id":"579","username":"test","email":"test@test.com"}}`, string(respStr))
}

func TestUserResolverSuite(t *testing.T){
	suite.Run(t, new(UserResolverSuite))
}

