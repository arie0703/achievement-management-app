package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// MockDynamoDBClient DynamoDBクライアントのモック
type MockDynamoDBClient struct {
	putItemFunc           func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	getItemFunc           func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	updateItemFunc        func(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	scanFunc              func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	deleteItemFunc        func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	transactWriteItemsFunc func(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
}

func (m *MockDynamoDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if m.putItemFunc != nil {
		return m.putItemFunc(ctx, params, optFns...)
	}
	return &dynamodb.PutItemOutput{}, nil
}

func (m *MockDynamoDBClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if m.getItemFunc != nil {
		return m.getItemFunc(ctx, params, optFns...)
	}
	return &dynamodb.GetItemOutput{}, nil
}

func (m *MockDynamoDBClient) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	if m.updateItemFunc != nil {
		return m.updateItemFunc(ctx, params, optFns...)
	}
	return &dynamodb.UpdateItemOutput{}, nil
}

func (m *MockDynamoDBClient) Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	if m.scanFunc != nil {
		return m.scanFunc(ctx, params, optFns...)
	}
	return &dynamodb.ScanOutput{}, nil
}

func (m *MockDynamoDBClient) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	if m.deleteItemFunc != nil {
		return m.deleteItemFunc(ctx, params, optFns...)
	}
	return &dynamodb.DeleteItemOutput{}, nil
}

func (m *MockDynamoDBClient) TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
	if m.transactWriteItemsFunc != nil {
		return m.transactWriteItemsFunc(ctx, params, optFns...)
	}
	return &dynamodb.TransactWriteItemsOutput{}, nil
}

// TestItem テスト用のアイテム構造体
type TestItem struct {
	ID    string `dynamodbav:"id"`
	Name  string `dynamodbav:"name"`
	Value int    `dynamodbav:"value"`
}

func TestDynamoDBRepository_PutItem(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockDynamoDBClient{}
	repo := NewDynamoDBRepositoryWithClient(ctx, mockClient)

	testItem := TestItem{
		ID:    "test-id",
		Name:  "test-name",
		Value: 100,
	}

	err := repo.PutItem("test-table", testItem)
	if err != nil {
		t.Errorf("PutItem failed: %v", err)
	}
}

func TestDynamoDBRepository_GetItem(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockDynamoDBClient{
		getItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{
				Item: map[string]types.AttributeValue{
					"id":    &types.AttributeValueMemberS{Value: "test-id"},
					"name":  &types.AttributeValueMemberS{Value: "test-name"},
					"value": &types.AttributeValueMemberN{Value: "100"},
				},
			}, nil
		},
	}
	repo := NewDynamoDBRepositoryWithClient(ctx, mockClient)

	key := map[string]interface{}{
		"id": "test-id",
	}

	var result TestItem
	err := repo.GetItem("test-table", key, &result)
	if err != nil {
		t.Errorf("GetItem failed: %v", err)
	}

	if result.ID != "test-id" || result.Name != "test-name" || result.Value != 100 {
		t.Errorf("GetItem returned unexpected result: %+v", result)
	}
}

func TestDynamoDBRepository_GetItem_NotFound(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockDynamoDBClient{
		getItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{
				Item: nil, // アイテムが見つからない場合
			}, nil
		},
	}
	repo := NewDynamoDBRepositoryWithClient(ctx, mockClient)

	key := map[string]interface{}{
		"id": "non-existent-id",
	}

	var result TestItem
	err := repo.GetItem("test-table", key, &result)
	if err == nil {
		t.Error("GetItem should have failed for non-existent item")
	}
}

func TestDynamoDBRepository_TransactWrite(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockDynamoDBClient{}
	repo := NewDynamoDBRepositoryWithClient(ctx, mockClient)

	items := []TransactWriteItem{
		{
			TableName: "test-table",
			Item: TestItem{
				ID:    "test-id-1",
				Name:  "test-name-1",
				Value: 100,
			},
			Operation: "PUT",
		},
		{
			TableName: "test-table",
			Item: TestItem{
				ID:    "test-id-2",
				Name:  "test-name-2",
				Value: 200,
			},
			Operation: "PUT",
		},
	}

	err := repo.TransactWrite(items)
	if err != nil {
		t.Errorf("TransactWrite failed: %v", err)
	}
}

func TestDynamoDBRepository_WithRetry(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockDynamoDBClient{}
	repo := NewDynamoDBRepositoryWithClient(ctx, mockClient)

	callCount := 0
	operation := func() error {
		callCount++
		if callCount < 3 {
			return errors.New("temporary error")
		}
		return nil
	}

	err := repo.WithRetry(operation, 3)
	if err != nil {
		t.Errorf("WithRetry failed: %v", err)
	}

	if callCount != 3 {
		t.Errorf("Expected 3 calls, got %d", callCount)
	}
}

func TestDynamoDBRepository_WithRetry_MaxRetriesExceeded(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockDynamoDBClient{}
	repo := NewDynamoDBRepositoryWithClient(ctx, mockClient)

	operation := func() error {
		return errors.New("persistent error")
	}

	err := repo.WithRetry(operation, 2)
	if err == nil {
		t.Error("WithRetry should have failed after max retries")
	}
}