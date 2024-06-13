package acl

import (
	"context"

	"purchase/domain/service"
	"purchase/infra/request"
)

const AiSpamCheckURL = "https://www.baidu.com"

type AiSpamChecker struct {
	client *request.HttpClient
}

func NewAiSpamChecker(cli *request.HttpClient) service.SpamChecker {
	return &AiSpamChecker{client: cli}
}

func (c *AiSpamChecker) Check(ctx context.Context, content string) error {
	req := map[string]string{
		"content": content,
	}
	resp := ""
	return c.client.NewRequest(AiSpamCheckURL, req, &resp).Get(ctx)
}
