package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type AWSClient struct {
	sqsClient sqs.Client
	config    *aws.Config
}

func NewAWSClient(config aws.Config) AWSClient {
	sqsClient := sqs.NewFromConfig(config)

	return AWSClient{
		sqsClient: *sqsClient,
		config:    &config,
	}
}

// ListAllQueues retrieves a list of all SQS queue URLs in the configured region.
func (client AWSClient) ListAllQueues(ctx context.Context) ([]string, error) {
	var queueUrls []string
	input := &sqs.ListQueuesInput{}

	// Paginate through results, as ListQueues returns a maximum of 1,000 results at a time.
	for {
		resp, err := client.sqsClient.ListQueues(ctx, input)
		if err != nil {
			return nil, err
		}

		if resp.QueueUrls != nil {
			queueUrls = append(queueUrls, resp.QueueUrls...)
		}

		if resp.NextToken == nil {
			break
		}
		input.NextToken = resp.NextToken
	}

	return queueUrls, nil
}

// GetQueueAttributesBatch retrieves attributes for multiple queues concurrently.
func (client AWSClient) GetQueueAttributesBatch(ctx context.Context, queueUrls []string, attributeNames []types.QueueAttributeName) (map[string]map[string]string, error) {
	var wg sync.WaitGroup
	results := make(chan map[string]map[string]string, len(queueUrls))
	errors := make(chan error, len(queueUrls))

	for _, queueUrl := range queueUrls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			// Fetch attributes for a single queue
			input := &sqs.GetQueueAttributesInput{
				QueueUrl:       aws.String(url),
				AttributeNames: attributeNames,
			}
			resp, err := client.sqsClient.GetQueueAttributes(ctx, input)
			if err != nil {
				errors <- fmt.Errorf("failed to get attributes for queue %s: %w", url, err)
				return
			}

			// Format result: map[queueURL]map[attributeName]attributeValue
			queueAttributes := map[string]map[string]string{
				url: resp.Attributes,
			}
			results <- queueAttributes
		}(queueUrl)
	}

	// Wait for all Go routines to finish
	wg.Wait()
	close(results)
	close(errors)

	// Collect results and check for errors
	combinedResults := make(map[string]map[string]string)
	for res := range results {
		for url, attrs := range res {
			combinedResults[url] = attrs
		}
	}

	if len(errors) > 0 {
		// In a real application, you might want more sophisticated error handling,
		// but for a simple example, return the first error encountered.
		return combinedResults, <-errors
	}

	return combinedResults, nil
}

func queueNameFromURL(url string) string {
	parts := strings.Split(url, "/")

	return parts[len(parts)-1]
}

type queueAttributesMsg struct {
	batch []batchItem
	err   error
}

func (client AWSClient) GetQueueAttributesCmd(ctx context.Context, allItems []list.Item, visibleItems []list.Item) tea.Cmd {
	return func() tea.Msg {
		var batch []batchItem
		indexMap := make(map[string]int)

		for i, cur := range allItems {
			cur := cur.(item)
			indexMap[cur.url] = i
		}

		urls := Map(visibleItems, func(listItem list.Item) string {
			curItem := listItem.(item)
			return curItem.url
		})

		attributes, err := client.GetQueueAttributesBatch(ctx, urls, []types.QueueAttributeName{types.QueueAttributeNameApproximateNumberOfMessages, types.QueueAttributeNameApproximateNumberOfMessagesNotVisible})
		if err != nil {
			return queueAttributesMsg{
				err: err,
			}
		}

		for _, cur := range visibleItems {
			cur := cur.(item)
			cur.available = attributes[cur.url][string(types.QueueAttributeNameApproximateNumberOfMessages)]
			cur.inFlight = attributes[cur.url][string(types.QueueAttributeNameApproximateNumberOfMessagesNotVisible)]

			batch = append(batch, batchItem{
				index:   indexMap[cur.url],
				setItem: cur,
			})
		}

		return queueAttributesMsg{batch: batch}
	}
}
