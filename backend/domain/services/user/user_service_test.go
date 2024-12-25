package user_service

import (
	"blizzflow/backend/infrastructure/database"
	"os"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestUserServiceSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Database Test Suite")
}

const testDBPath = "test.db"

var (
	DB          *gorm.DB
	userService *UserService
)

var _ = ginkgo.BeforeSuite(func() {
	os.Remove(testDBPath)
	database.InitDB(testDBPath)
	DB = database.DB
	userService = NewUserService(DB)
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

var _ = ginkgo.Describe("User Service", func() {
	ginkgo.BeforeEach(func() {
		// Clean test data instead of recreating connection
		DB.Exec("DELETE FROM users")
	})

	ginkgo.It("should create a user successfully", func() {
		user, err := userService.CreateUser("testuser", "password123")
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(user.Username).To(gomega.Equal("testuser"))
	})

	ginkgo.It("should fail creating duplicate username", func() {
		_, err := userService.CreateUser("testuser2", "password123")
		gomega.Expect(err).To(gomega.BeNil())

		_, err = userService.CreateUser("testuser2", "password123")
		gomega.Expect(err).To(gomega.MatchError("username already exists"))
	})

	ginkgo.It("should get user by ID", func() {
		created, err := userService.CreateUser("testuser3", "password123")
		gomega.Expect(err).To(gomega.BeNil())

		found, err := userService.GetUserByID(created.ID)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(found.Username).To(gomega.Equal("testuser3"))
	})

	ginkgo.It("should get user by username", func() {
		_, err := userService.CreateUser("testuser4", "password123")
		gomega.Expect(err).To(gomega.BeNil())

		found, err := userService.GetUserByUsername("testuser4")
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(found.Username).To(gomega.Equal("testuser4"))
	})

	ginkgo.It("should update user", func() {
		user, err := userService.CreateUser("testuser5", "password123")
		gomega.Expect(err).To(gomega.BeNil())

		user.Username = "updated_user"
		err = userService.UpdateUser(user)
		gomega.Expect(err).To(gomega.BeNil())

		found, err := userService.GetUserByID(user.ID)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(found.Username).To(gomega.Equal("updated_user"))
	})

	ginkgo.It("should delete user", func() {
		user, err := userService.CreateUser("testuser6", "password123")
		gomega.Expect(err).To(gomega.BeNil())

		err = userService.DeleteUser(user.ID)
		gomega.Expect(err).To(gomega.BeNil())

		_, err = userService.GetUserByID(user.ID)
		gomega.Expect(err).ToNot(gomega.BeNil())
	})

})
