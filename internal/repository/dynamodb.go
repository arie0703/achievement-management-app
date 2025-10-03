package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	
	appconfig "achievement-management/internal/config"
)

// DynamoDBAPI DynamoDB操作のインターフェース
type DynamoDBAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
}

// DynamoDBRepository DynamoDB操作の実装
type DynamoDBRepository struct {
	client DynamoDBAPI
	ctx    context.Context
}

// NewDynamoDBRepository DynamoDBリポジトリの作成
func NewDynamoDBRepository(ctx context.Context, appConfig *appconfig.Config) (*DynamoDBRepository, error) {
	// AWS設定を読み込み
	var awsConfig aws.Config
	var err error
	
	if appConfig.AWS.DynamoDBEndpoint != "" {
		// ローカルDynamoDBを使用
		awsConfig, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(appConfig.AWS.Region),
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					if service == dynamodb.ServiceID {
						return aws.Endpoint{
							URL:           appConfig.AWS.DynamoDBEndpoint,
							SigningRegion: appConfig.AWS.Region,
						}, nil
					}
					return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
				})),
		)
	} else {
		// AWS DynamoDBを使用
		if appConfig.AWS.Profile != "" {
			awsConfig, err = config.LoadDefaultConfig(ctx,
				config.WithRegion(appConfig.AWS.Region),
				config.WithSharedConfigProfile(appConfig.AWS.Profile),
			)
		} else if appConfig.AWS.AccessKeyID != "" && appConfig.AWS.SecretAccessKey != "" {
			awsConfig, err = config.LoadDefaultConfig(ctx,
				config.WithRegion(appConfig.AWS.Region),
				config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
					return aws.Credentials{
						AccessKeyID:     appConfig.AWS.AccessKeyID,
						SecretAccessKey: appConfig.AWS.SecretAccessKey,
					}, nil
				})),
			)
		} else {
			awsConfig, err = config.LoadDefaultConfig(ctx,
				config.WithRegion(appConfig.AWS.Region),
			)
		}
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := dynamodb.NewFromConfig(awsConfig)
	
	return &DynamoDBRepository{
		client: client,
		ctx:    ctx,
	}, nil
}

// NewDynamoDBRepositoryWithClient カスタムクライアントでDynamoDBリポジトリを作成
func NewDynamoDBRepositoryWithClient(ctx context.Context, client DynamoDBAPI) *DynamoDBRepository {
	return &DynamoDBRepository{
		client: client,
		ctx:    ctx,
	}
}

// PutItem アイテムを追加
func (r *DynamoDBRepository) PutItem(tableName string, item interface{}) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	_, err = r.client.PutItem(r.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item to table %s: %w", tableName, err)
	}

	return nil
}

// GetItem アイテムを取得
func (r *DynamoDBRepository) GetItem(tableName string, key map[string]interface{}, result interface{}) error {
	keyAv, err := attributevalue.MarshalMap(key)
	if err != nil {
		return fmt.Errorf("failed to marshal key: %w", err)
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       keyAv,
	}

	resp, err := r.client.GetItem(r.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to get item from table %s: %w", tableName, err)
	}

	if resp.Item == nil {
		return fmt.Errorf("item not found in table %s", tableName)
	}

	err = attributevalue.UnmarshalMap(resp.Item, result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return nil
}

// UpdateItem アイテムを更新
func (r *DynamoDBRepository) UpdateItem(tableName string, key map[string]interface{}, updateExpression string, expressionAttributeValues map[string]interface{}) error {
	keyAv, err := attributevalue.MarshalMap(key)
	if err != nil {
		return fmt.Errorf("failed to marshal key: %w", err)
	}

	var eavAv map[string]types.AttributeValue
	if expressionAttributeValues != nil {
		eavAv, err = attributevalue.MarshalMap(expressionAttributeValues)
		if err != nil {
			return fmt.Errorf("failed to marshal expression attribute values: %w", err)
		}
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       keyAv,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: eavAv,
	}

	_, err = r.client.UpdateItem(r.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in table %s: %w", tableName, err)
	}

	return nil
}

// Scan テーブル全体をスキャン
func (r *DynamoDBRepository) Scan(tableName string, result interface{}) error {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	resp, err := r.client.Scan(r.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to scan table %s: %w", tableName, err)
	}

	err = attributevalue.UnmarshalListOfMaps(resp.Items, result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal scan result: %w", err)
	}

	return nil
}

// DeleteItem アイテムを削除
func (r *DynamoDBRepository) DeleteItem(tableName string, key map[string]interface{}) error {
	keyAv, err := attributevalue.MarshalMap(key)
	if err != nil {
		return fmt.Errorf("failed to marshal key: %w", err)
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       keyAv,
	}

	_, err = r.client.DeleteItem(r.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item from table %s: %w", tableName, err)
	}

	return nil
}

// TransactWrite トランザクション書き込み
func (r *DynamoDBRepository) TransactWrite(items []TransactWriteItem) error {
	if len(items) == 0 {
		return fmt.Errorf("no items provided for transaction")
	}

	transactItems := make([]types.TransactWriteItem, 0, len(items))

	for _, item := range items {
		av, err := attributevalue.MarshalMap(item.Item)
		if err != nil {
			return fmt.Errorf("failed to marshal transaction item: %w", err)
		}

		switch item.Operation {
		case "PUT":
			transactItems = append(transactItems, types.TransactWriteItem{
				Put: &types.Put{
					TableName: aws.String(item.TableName),
					Item:      av,
				},
			})
		case "UPDATE":
			// For UPDATE operations, the item should contain key and update expression
			// This is a simplified implementation - in practice, you'd need more complex handling
			return fmt.Errorf("UPDATE operation not implemented in this simplified version")
		case "DELETE":
			transactItems = append(transactItems, types.TransactWriteItem{
				Delete: &types.Delete{
					TableName: aws.String(item.TableName),
					Key:       av,
				},
			})
		default:
			return fmt.Errorf("unsupported transaction operation: %s", item.Operation)
		}
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	}

	_, err := r.client.TransactWriteItems(r.ctx, input)
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	return nil
}

// WithRetry リトライロジック付きで操作を実行
func (r *DynamoDBRepository) WithRetry(operation func() error, maxRetries int) error {
	var lastErr error
	
	for i := 0; i <= maxRetries; i++ {
		err := operation()
		if err == nil {
			return nil
		}
		
		lastErr = err
		
		// 最後の試行でない場合は待機
		if i < maxRetries {
			backoffDuration := time.Duration(i+1) * 100 * time.Millisecond
			time.Sleep(backoffDuration)
		}
	}
	
	return fmt.Errorf("operation failed after %d retries: %w", maxRetries, lastErr)
}