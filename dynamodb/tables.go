package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDBTables is a map whose keys are DynamoDB table names and whose values are `dynamodb.CreateTableInput` instances.
var DynamoDBTables = map[string]*dynamodb.CreateTableInput{
	"libraryofcongress": &dynamodb.CreateTableInput{
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       aws.String("HASH"), // partition key
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Label"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Source"),
				AttributeType: aws.String("S"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("label"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("Label"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("Id"),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
			},
		},
		BillingMode: aws.String("PAY_PER_REQUEST"),
		// TableName:   set inline below
	},
}
