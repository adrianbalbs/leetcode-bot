package leetcode_client

import (
	"context"
	"net/http"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/assert"
)

func TestGetDailyProblem(t *testing.T) {
	ctx := context.Background()
	client := graphql.NewClient(LeetcodeURL, &http.Client{Transport: &userAgentTransport{http.DefaultTransport}})

	resp, err := getActiveDailyCodingChallenge(ctx, client)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestGetProblem(t *testing.T) {
	ctx := context.Background()
	client := graphql.NewClient(LeetcodeURL, &http.Client{Transport: &userAgentTransport{http.DefaultTransport}})

	resp, err := getProblem(ctx, client, "two-sum")
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}
