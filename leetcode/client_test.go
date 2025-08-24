package leetcode

import (
	"context"
	"net/http"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/assert"
)

func TestGetDailyProblem(t *testing.T) {
	ctx := context.Background()
	client := graphql.NewClient(LeetcodeURL, &http.Client{Transport: &UserAgentTransport{http.DefaultTransport}})

	resp, err := GetActiveDailyCodingChallenge(ctx, client)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestGetProblem(t *testing.T) {
	ctx := context.Background()
	client := graphql.NewClient(LeetcodeURL, &http.Client{Transport: &UserAgentTransport{http.DefaultTransport}})

	resp, err := GetProblem(ctx, client, "two-sum")
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}
