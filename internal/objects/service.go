package objects

import "context"

// Service — бизнес-логика объектов и типов событий.
type Service struct {
	repo Repository
}

// NewService создаёт сервис объектов.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetObjectData возвращает данные одного объекта по ID.
func (s *Service) GetObjectData(ctx context.Context, id int) (*SingleObjectInfo, error) {
	return s.repo.GetByID(ctx, id)
}

// GetObjectsList возвращает список объектов с фильтрами.
func (s *Service) GetObjectsList(ctx context.Context, eventTypeIDs []int, dateFrom, dateTo string) ([]ObjectInfo, error) {
	return s.repo.List(ctx, eventTypeIDs, dateFrom, dateTo)
}

