package auth_service

import (
	repository "blizzflow/backend/domain/repositories"
	"blizzflow/backend/infrastructure/database"
	"os"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestAuthServiceSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Auth Service Test Suite")
}

const testDBPath = "test.db"

var (
	DB                    *gorm.DB
	authService           *AuthService
	userRepo              *repository.UserRepository
	sessionRepo           *repository.SessionRepository
	securityQuestionsRepo *repository.SecurityQuestionRepository
)

var _ = ginkgo.BeforeSuite(func() {
	os.Remove(testDBPath)
	database.InitDB(testDBPath)
	DB = database.DB

	userRepo = repository.NewUserRepository(DB)
	sessionRepo = repository.NewSessionRepository(DB)
	securityQuestionsRepo = repository.NewSecurityQuestionRepository(DB)

	authService = NewAuthService(
		userRepo,
		sessionRepo,
		securityQuestionsRepo,
	)
})

var _ = ginkgo.AfterSuite(func() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
	os.Remove(testDBPath)
})

var _ = ginkgo.Describe("Auth Service", func() {
	ginkgo.BeforeEach(func() {
		DB.Exec("DELETE FROM users")
		DB.Exec("DELETE FROM sessions")
		DB.Exec("DELETE FROM security_questions")
	})

	ginkgo.It("should register user successfully", func() {
		err := authService.Register("testuser", "password123")
		gomega.Expect(err).To(gomega.BeNil())

		user, err := userRepo.GetUserByUsername("testuser")
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(user.Username).To(gomega.Equal("testuser"))
	})

	ginkgo.It("should fail registration with empty credentials", func() {
		err := authService.Register("", "")
		gomega.Expect(err).To(gomega.Equal(ErrEmptyCredentials))
	})

	ginkgo.It("should login user successfully", func() {
		err := authService.Register("testuser2", "password123")
		gomega.Expect(err).To(gomega.BeNil())

		session, err := authService.Login("testuser2", "password123")
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(session).ToNot(gomega.BeNil())
	})

	ginkgo.It("should fail login with wrong password", func() {
		err := authService.Register("testuser3", "password123")
		gomega.Expect(err).To(gomega.BeNil())

		_, err = authService.Login("testuser3", "wrongpassword")
		gomega.Expect(err).To(gomega.Equal(ErrInvalidCredentials))
	})

	ginkgo.It("should set security questions", func() {
		err := authService.Register("testuser4", "password123")
		gomega.Expect(err).To(gomega.BeNil())

		questions := map[string]string{
			"What is your pet's name?": "Rex",
			"Where were you born?":     "London",
		}

		err = authService.SetSecurityQuestions("testuser4", questions)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("should recover password with correct answers", func() {
		err := authService.Register("testuser5", "password123")
		gomega.Expect(err).To(gomega.BeNil())

		questions := map[string]string{
			"What is your pet's name?": "Rex",
		}

		err = authService.SetSecurityQuestions("testuser5", questions)
		gomega.Expect(err).To(gomega.BeNil())

		err = authService.RecoverPassword("testuser5", questions, "newpassword123")
		gomega.Expect(err).To(gomega.BeNil())

		// Verify new password works
		session, err := authService.Login("testuser5", "newpassword123")
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(session).ToNot(gomega.BeNil())
	})

	ginkgo.It("should fail password recovery with wrong answers", func() {
		err := authService.Register("testuser6", "password123")
		gomega.Expect(err).To(gomega.BeNil())

		questions := map[string]string{
			"What is your pet's name?": "Rex",
		}

		err = authService.SetSecurityQuestions("testuser6", questions)
		gomega.Expect(err).To(gomega.BeNil())

		wrongAnswers := map[string]string{
			"What is your pet's name?": "Wrong",
		}

		err = authService.RecoverPassword("testuser6", wrongAnswers, "newpassword123")
		gomega.Expect(err).To(gomega.Equal(ErrInvalidAnswers))
	})

})
