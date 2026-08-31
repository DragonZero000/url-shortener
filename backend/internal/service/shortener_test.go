package service

import (
	"errors"
	"testing"

	"github.com/DragonZero000/url-shortener/internal/model"
)

type fakeRepo struct {
	byCode        map[string]*model.URL
	byOriginalURL map[string]*model.URL
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		byCode:        make(map[string]*model.URL),
		byOriginalURL: make(map[string]*model.URL),
	}
}

func (f *fakeRepo) Save(code, oriURL string) error {
	entry := &model.URL{
		Code:        code,
		OriginalURL: oriURL,
	}
	f.byCode[code] = entry
	f.byOriginalURL[oriURL] = entry
	return nil
}

func (f *fakeRepo) GetByCode(code string) (*model.URL, error) {
	entry, exi := f.byCode[code]
	if !exi {
		return nil, nil
	}
	return entry, nil
}

func (f *fakeRepo) GetByOriginalURL(oriURL string) (*model.URL, error) {
	entry, exi := f.byOriginalURL[oriURL]
	if !exi {
		return nil, nil
	}
	return entry, nil
}

func setupService() *Service {
	return NewService(newFakeRepo())
}

func TestShorten_NewURL(t *testing.T) {
	svc := setupService()

	code, err := svc.Shorten("https://example.com")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if code == "" {
		t.Errorf("expected non-empty code, got empty string")
	}
}

func TestGetByCode(t *testing.T) {
	svc := setupService()
	testURL := "https://ya.ru"
	code, err := svc.Shorten(testURL)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if code == "" {
		t.Fatalf("expected non-empty code, got empty string")
	}
	oriUrl, err := svc.GetByCode(code)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if oriUrl != testURL {
		t.Fatalf("expected %s, but found %s", testURL, oriUrl)
	}
}

func TestGetByCode_NotFound(t *testing.T) {
	svc := setupService()
	_, err := svc.GetByCode("nonexistent")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestShorten_SameURL(t *testing.T) {
	svc := setupService()
	code1, err := svc.Shorten("https://example.com")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if code1 == "" {
		t.Fatalf("expected non-empty code, got empty string")
	}
	code2, err := svc.Shorten("https://example.com")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if code2 == "" {
		t.Fatalf("expected non-empty code, got empty string")
	}
	if code1 != code2 {
		t.Errorf("code must be same but code 1: %s != code 2: %s", code1, code2)
	}
}

func TestShorten_InvalidURL(t *testing.T) {
	svc := setupService()
	_, err := svc.Shorten("")
	if !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("expected error, but got: %v", err)
	}
	_, err = svc.Shorten("example.com")
	if !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("exprcted error, but got: %v", err)
	}
	_, err = svc.Shorten("http:///")
	if !errors.Is(err, ErrInvalidURL) {
		t.Fatalf("exprcted error, but got: %v", err)
	}
	_, err = svc.Shorten("http://ya.ru/")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
