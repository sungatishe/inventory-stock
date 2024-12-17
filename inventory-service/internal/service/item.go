package service

import (
	"context"
	"inventory-service/api/proto"
	"inventory-service/internal/models"
	"log"
)

type ItemRepository interface {
	AddItem(item *models.Item) error
	UpdateItem(item *models.Item) error
	GetItemByID(itemID, userID string) (*models.Item, error)
	DeleteItem(itemID, userID string) error
	GetAllItems(userID string) ([]*models.Item, error)
}

type Producer interface {
	Close()
	PublishEvent(ctx context.Context, key string, event map[string]interface{}) error
}

type ItemService struct {
	repo     ItemRepository
	producer Producer
	proto.UnimplementedInventoryServiceServer
}

func NewItemService(repo ItemRepository, producer Producer) *ItemService {
	return &ItemService{repo: repo, producer: producer}
}

func (s *ItemService) AddItem(ctx context.Context, req *proto.AddItemRequest) (*proto.AddItemResponse, error) {
	item := &models.Item{
		Name:     req.GetItem().GetName(),
		Quantity: req.GetItem().GetQuantity(),
		Price:    req.GetItem().GetPrice(),
		UserID:   req.GetItem().GetUserId(),
	}

	err := s.repo.AddItem(item)
	if err != nil {
		log.Printf("Failed to add item: %v", err)
		return &proto.AddItemResponse{
			Message: "Error in adding item",
			Success: false,
		}, err
	}

	event := map[string]interface{}{
		"action": "create",
		"item":   item,
	}
	err = s.producer.PublishEvent(ctx, "item-create", event)
	if err != nil {
		log.Printf("Failed to publish create event: %v", err)
	}

	return &proto.AddItemResponse{
		Message: "Item Created Successfully",
		Success: true,
	}, nil
}

func (s *ItemService) UpdateItem(ctx context.Context, req *proto.UpdateItemRequest) (*proto.UpdateItemResponse, error) {
	item := &models.Item{
		ID:       req.GetItem().GetId(),
		Name:     req.GetItem().GetName(),
		Quantity: req.GetItem().GetQuantity(),
		Price:    req.GetItem().GetPrice(),
		UserID:   req.GetItem().GetUserId(),
	}

	err := s.repo.UpdateItem(item)
	if err != nil {
		log.Printf("Failed to update item: %v", err)
		return &proto.UpdateItemResponse{
			Message: "Error in updating item",
			Success: false,
		}, err
	}

	event := map[string]interface{}{
		"action": "update",
		"item":   item,
	}

	if err := s.producer.PublishEvent(ctx, "item-update", event); err != nil {
		log.Printf("Failed to publish update event: %v", err)
	}

	return &proto.UpdateItemResponse{
		Message: "Item Updated successfully",
		Success: true,
	}, nil
}

func (s *ItemService) GetItem(ctx context.Context, req *proto.GetItemRequest) (*proto.GetItemResponse, error) {
	item, err := s.repo.GetItemByID(req.GetId(), req.GetUserId())
	if err != nil {
		log.Printf("Failed to get item: %v", err)
		return nil, err
	}

	return &proto.GetItemResponse{Item: &proto.Item{
		Id:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		Price:    item.Price,
		UserId:   item.UserID,
	}}, nil
}

func (s *ItemService) DeleteItem(ctx context.Context, req *proto.DeleteItemRequest) (*proto.DeleteItemResponse, error) {
	err := s.repo.DeleteItem(req.GetId(), req.GetUserId())
	if err != nil {
		log.Printf("Failed to deleting item: %v", err)
		return &proto.DeleteItemResponse{
			Message: "Error in deleting item",
			Success: false,
		}, err
	}

	return &proto.DeleteItemResponse{
		Message: "Item deleted successfully",
		Success: true,
	}, nil
}

func (s *ItemService) GetInventory(ctx context.Context, req *proto.GetInventoryRequest) (*proto.GetInventoryResponse, error) {
	items, err := s.repo.GetAllItems(req.GetUserId())
	if err != nil {
		log.Printf("Failed to get all items: %v", err)
		return nil, err
	}

	var protoItems []*proto.Item
	for _, item := range items {
		protoItems = append(protoItems, &proto.Item{
			Id:       item.ID,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
			UserId:   item.UserID,
		})
	}

	return &proto.GetInventoryResponse{Items: protoItems}, nil
}
