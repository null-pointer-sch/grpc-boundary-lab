package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	pb "github.com/AndySchubert/grpc-boundary-lab/internal/proto"
)

// GrpcBackendClient makes gRPC calls to the backend.
type GrpcBackendClient struct {
	Client pb.PingServiceClient
}

// Ping performs a gRPC call.
func (c *GrpcBackendClient) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return c.Client.Ping(ctx, req)
}

// RestBackendClient makes HTTP calls to the backend.
type RestBackendClient struct {
	TargetURL  string
	HTTPClient *http.Client
}

type pingMessage struct {
	Message string `json:"message"`
}

// Ping performs an HTTP JSON call.
func (c *RestBackendClient) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	reqURL := fmt.Sprintf("%s/api/ping?message=%s", c.TargetURL, url.QueryEscape(req.GetMessage()))
	httpReq, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed pingMessage
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}

	return &pb.PingResponse{Message: parsed.Message}, nil
}
