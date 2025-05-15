#!/bin/bash

# Generate mocks for authentication service
mockgen -source=internal/authentication/client.go -destination=internal/authentication/mock/client_mock.go -package=mock

# Generate mocks for image analysis service
mockgen -source=internal/image_analysis/client.go -destination=internal/image_analysis/mock/client_mock.go -package=mock

# Generate mocks for middleware
mockgen -source=internal/middleware/authentication.go -destination=internal/middleware/mock/authentication_mock.go -package=mock 