package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type inMemoryStoreEntry struct {
	sess    SessionData
	expires time.Time
}

type InMemorySessionHandler struct {
	CookieName      string
	SessionLifetime time.Duration
	WhitelistPaths  []string

	store     map[string]*inMemoryStoreEntry
	storeLock sync.Mutex
}

func (s *InMemorySessionHandler) CreateMiddleware() echo.MiddlewareFunc {
	if s.CookieName == "" {
		s.CookieName = "auth"
	}

	s.store = make(map[string]*inMemoryStoreEntry)
	go s.janitor()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			entry, err := s.getEntry(c)
			if err != nil || entry == nil {
				if s.shouldSkip(c) {
					return next(c)
				}

				return echo.NewHTTPError(http.StatusUnauthorized, "Login required")
			}

			c.Set(sessionContextKey, entry.sess)

			return next(c)
		}
	}
}

func (s *InMemorySessionHandler) Start(c echo.Context, sess SessionData) error {
	s.storeLock.Lock()
	defer s.storeLock.Unlock()

	newCookie, err := s.genRandom()
	if err != nil {
		return err
	}

	exp := time.Now().Add(s.SessionLifetime)

	s.store[newCookie] = &inMemoryStoreEntry{
		sess:    sess,
		expires: exp,
	}

	s.setCookie(c, newCookie, exp)
	return err
}

func (s *InMemorySessionHandler) Destroy(c echo.Context) error {
	s.storeLock.Lock()
	defer s.storeLock.Unlock()

	cookie := s.getCookie(c)
	delete(s.store, cookie)
	return nil
}

func (s *InMemorySessionHandler) Refresh(c echo.Context) error {
	s.storeLock.Lock()
	defer s.storeLock.Unlock()

	cookie := s.getCookie(c)
	entry, err := s.getEntry(c)
	if err != nil {
		return err
	}

	entry.expires = time.Now().Add(s.SessionLifetime)
	s.setCookie(c, cookie, entry.expires)

	return nil
}

func (s *InMemorySessionHandler) Rotate(c echo.Context) error {
	s.storeLock.Lock()
	defer s.storeLock.Unlock()

	cookie := s.getCookie(c)

	newCookie, err := s.genRandom()
	if err != nil {
		return err
	}

	entry, err := s.getEntry(c)
	if err != nil {
		return err
	}

	s.store[newCookie] = entry
	delete(s.store, cookie)

	s.setCookie(c, cookie, entry.expires)

	return nil
}

func (s *InMemorySessionHandler) UpdateSessionData(c echo.Context, updateFunc func(*SessionData) error) error {
	s.storeLock.Lock()
	defer s.storeLock.Unlock()

	entry, err := s.getEntry(c)
	if err != nil {
		return err
	}

	return updateFunc(&entry.sess)
}

func (s *InMemorySessionHandler) LogoutAllSessionsForUser(userId string) error {
	s.storeLock.Lock()
	defer s.storeLock.Unlock()

	for sessId, entry := range s.store {
		if entry.sess.UserID == userId {
			delete(s.store, sessId)
		}
	}

	return nil
}

func (s *InMemorySessionHandler) genRandom() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *InMemorySessionHandler) getCookie(c echo.Context) string {
	cookie, err := c.Cookie(s.CookieName)
	if err != nil || cookie == nil {
		return ""
	}
	return cookie.Value
}

// Unsafe: Caller must hold mutex
func (s *InMemorySessionHandler) getEntry(c echo.Context) (*inMemoryStoreEntry, error) {
	cookie := s.getCookie(c)
	if len(cookie) != 64 {
		return nil, fmt.Errorf("%v is an invalid session length (expected 64 characters)", cookie)
	}

	entry, ok := s.store[cookie]
	if !ok {
		return nil, fmt.Errorf("%v is not a valid session", cookie)
	}

	if entry.expires.Before(time.Now()) {
		return nil, fmt.Errorf("%v session expired", cookie)
	}

	return entry, nil
}

func (s *InMemorySessionHandler) setCookie(c echo.Context, val string, expires time.Time) {
	c.SetCookie(&http.Cookie{
		Name:     s.CookieName,
		Value:    val,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Scheme() == "https",
		SameSite: http.SameSiteStrictMode,
	})
}

func (s *InMemorySessionHandler) cleanup() {
	s.storeLock.Lock()
	defer s.storeLock.Unlock()

	for sessId, entry := range s.store {
		if entry.expires.Before(time.Now()) {
			delete(s.store, sessId)
		}
	}
}

func (s *InMemorySessionHandler) shouldSkip(c echo.Context) bool {
	path := c.Request().URL.Path
	for _, bypassPath := range s.WhitelistPaths {
		if strings.HasPrefix(path, bypassPath) {
			return true
		}
	}
	return false
}

func (s *InMemorySessionHandler) janitor() {
	for {
		s.cleanup()
		time.Sleep(time.Minute)
	}
}
