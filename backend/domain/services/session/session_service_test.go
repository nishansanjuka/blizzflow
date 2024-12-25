package session_service

import (
	"blizzflow/backend/infrastructure/database"
	"os"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestSessionServiceSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Session Service Test Suite")
}

const testDBPath = "test.db"

var (
	DB             *gorm.DB
	sessionService *SessionService
)

var _ = ginkgo.BeforeSuite(func() {
	os.Remove(testDBPath)
	database.InitDB(testDBPath)
	DB = database.DB
	sessionService = NewSessionService(DB)
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

var _ = ginkgo.Describe("Session Service", func() {
	ginkgo.BeforeEach(func() {
		DB.Exec("DELETE FROM sessions")
	})

	ginkgo.It("should create a session successfully", func() {
		session, err := sessionService.CreateSession(1)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(session.UserID).To(gomega.Equal(uint(1)))
	})

	ginkgo.It("should get session by ID", func() {
		created, err := sessionService.CreateSession(1)
		gomega.Expect(err).To(gomega.BeNil())

		found, err := sessionService.GetSession(created.ID)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(found.UserID).To(gomega.Equal(uint(1)))
	})

	ginkgo.It("should validate existing session", func() {
		created, err := sessionService.CreateSession(1)
		gomega.Expect(err).To(gomega.BeNil())

		valid, err := sessionService.ValidateSession(created.ID)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(valid).To(gomega.BeTrue())
	})

	ginkgo.It("should return false for non-existent session validation", func() {
		valid, err := sessionService.ValidateSession(999)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(valid).To(gomega.BeFalse())
	})

	ginkgo.It("should delete session", func() {
		session, err := sessionService.CreateSession(1)
		gomega.Expect(err).To(gomega.BeNil())

		err = sessionService.DeleteSession(session.ID)
		gomega.Expect(err).To(gomega.BeNil())

		valid, err := sessionService.ValidateSession(session.ID)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(valid).To(gomega.BeFalse())
	})
})
