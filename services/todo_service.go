package services

import (
	"Shawty/models"
	"errors"
	"sync"
	"time"
)

type ToDoService struct {
	lists map[uint]*models.ToDoList
	items map[uint]*models.ToDoItem
	mu    sync.RWMutex
}

func NewToDoService() *ToDoService {
	return &ToDoService{
		lists: make(map[uint]*models.ToDoList),
		items: make(map[uint]*models.ToDoItem),
		mu:    sync.RWMutex{},
	}
}

func (s *ToDoService) CreateList(userID uint, title string) *models.ToDoList {
	s.mu.Lock()
	defer s.mu.Unlock()

	list := &models.ToDoList{
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Title:     title,
	}

	list.ID = uint(len(s.lists) + 1)
	s.lists[list.ID] = list

	return list
}

func (s *ToDoService) GetListByID(id uint, userID uint, usertype int) []*models.ToDoList {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var lists []*models.ToDoList
	for _, list := range s.lists {
		if list.ID == id && list.UserID == userID || list.ID == id && usertype == 2 {
			lists = append(lists, list)
		}
	}
	return lists
}
func (s *ToDoService) GetListsByUserID(userID uint) []*models.ToDoList {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var lists []*models.ToDoList
	for _, list := range s.lists {
		if list.UserID == userID {
			lists = append(lists, list)
		}
	}
	return lists
}

func (s *ToDoService) GetAllLists() []*models.ToDoList {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var lists []*models.ToDoList
	for _, list := range s.lists {
		lists = append(lists, list)
	}
	return lists
}
func (s *ToDoService) DeleteList(userID uint, listID uint, usertype int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Liste var mı kontrol et
	list, ok := s.lists[listID]
	if !ok {
		return errors.New("list not found")
	}

	// Kullanıcı yetkisi kontrol et
	if userID != list.UserID && usertype != 2 {
		return errors.New("unauthorized")
	}

	// Silinmiş işareti koy
	list.Deleted = true
	list.DeletedAt = time.Now()

	return nil
}
func (s *ToDoService) UpdateListTitle(userID uint, listID uint, newTitle string, userType int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Liste var mı kontrol et
	list, ok := s.lists[listID]
	if !ok {
		return errors.New("list not found")
	}

	// Kullanıcı yetkisi kontrol et
	if userID != list.UserID && userType != 2 {
		return errors.New("unauthorized")
	}

	// Başlık güncelle
	list.Title = newTitle
	list.UpdatedAt = time.Now()

	return nil
}

func (s *ToDoService) GetItemsByTaskID(taskID uint, userID uint, userType int) []*models.ToDoItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var items []*models.ToDoItem
	for _, item := range s.items {
		if item.ID == taskID {
			if item.UserID == userID || userType == 2 {
				items = append(items, item)
			}
		}
	}
	return items
}

func (s *ToDoService) GetItemsByListID(userID uint, listID uint, userType int) []*models.ToDoItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var items []*models.ToDoItem
	for _, item := range s.items {
		if item.ListID == listID {
			if userType == 1 && item.UserID == userID {
				items = append(items, item)
			} else if userType == 2 {
				items = append(items, item)
			}
		}
	}
	return items
}

func (s *ToDoService) CreateTask(userID, listID uint, task string, userType int) (*models.ToDoItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Listeyi kontrol et
	list, ok := s.lists[listID]
	if !ok {
		return nil, errors.New("list not found")
	}

	// Kullanıcı yetkisi kontrol et
	if userID != list.UserID && userType != 2 {
		return nil, errors.New("unauthorized")
	}

	item := &models.ToDoItem{
		ID:        uint(len(s.items) + 1),
		UserID:    userID,
		ListID:    listID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Task:      task,
		Completed: false,
		Deleted:   false,
	}

	s.items[item.ID] = item

	return item, nil
}
