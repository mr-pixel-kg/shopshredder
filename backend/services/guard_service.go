package services

import (
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/config"
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/database/models"
	"github.com/mr-pixel-kg/shopware-sandbox-plattform/database/repository"
	"log"
	"log/slog"
)

type GuardService struct {
	sessionRepository *repository.SessionRepository
	guardConfig       config.GuardConfig
}

func NewGuardService(sessionRepository *repository.SessionRepository, guardConfig config.GuardConfig) *GuardService {
	guardService := &GuardService{
		sessionRepository: sessionRepository,
		guardConfig:       guardConfig,
	}
	guardService.startupCheck()
	return guardService
}

// Gets all sandbox sessions for a given IP address
func (s *GuardService) GetSessions(ipAddress string) []models.Session {
	sessions, err := s.sessionRepository.GetSessionsForIp(ipAddress)
	if err != nil {
		log.Printf("Error getting sessions for IP: %v", err)
		return make([]models.Session, 0)
	}
	return sessions
}

// Records a new session after new sandbox is created
func (s *GuardService) RegisterSession(ipAddress string, userAgent string, username *string, sandboxId string) error {
	err := s.sessionRepository.Create(ipAddress, userAgent, username, sandboxId)
	if err != nil {
		return err
	}
	return nil
}

// Removes a sandbox session when the sandbox is deleted
func (s *GuardService) UnregisterSession(sandboxId string) error {
	err := s.sessionRepository.Remove(sandboxId)
	if err != nil {
		return err
	}
	return nil
}

// Checks if the current session (based on IP address) has exceeded the limit of concurrent sandboxes
// Returns true if the limit is exceeded
// Returns false if the limit is not exceeded
func (s *GuardService) IsNewSessionAllowed(ipAddress string) bool {
	sessions := s.GetSessions(ipAddress)
	if len(sessions) < s.guardConfig.MaxSandboxesPerIP {
		return true
	}
	return false
}

// Checks if the current IP address has exceeded the limit of concurrent sandboxes
// Returns true and records a new session if the limit is not exceeded
// Returns false if the limit is exceeded for that IP address
func (s *GuardService) CheckAndRegisterSession(ipAddress string, userAgent string, username *string, sandboxId string) (bool, error) {
	if s.IsNewSessionAllowed(ipAddress) {
		err := s.RegisterSession(ipAddress, userAgent, username, sandboxId)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (s *GuardService) startupCheck() {
	log.Println("*** Executing guard service startup check ***")

	sessions, err := s.sessionRepository.GetAll()
	if err != nil {
		log.Panicf("Failed to list sandbox sessions: %v", err)
	}
	for _, session := range sessions {
		slog.Info("Remove old session from database", "sessionId", session.ID, "sandboxId", session.SandboxID)

		// Quick fix to remove all sessions on startup
		s.sessionRepository.Remove(session.SandboxID)

		/*_, contErr := s.sandboxService.GetSandbox(context.Background(), session.SandboxID)
		if contErr != nil {
			slog.Warn("Found dangling session database record", "sessionId", session.ID, "sandboxId", session.SandboxID, "err", contErr )
			s.sessionRepository.Remove(session.SandboxID)
		}*/
	}
}
