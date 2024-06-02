easyjson:
	easyjson -pkg internal/auth
	easyjson -pkg internal/feed
	easyjson -pkg internal/image

mocks:
	mockgen --source=internal/auth/interfaces.go --destination=internal/auth/mocks/core_mocks.go --package=mock
	mockgen --source=internal/feed/interfaces.go --destination=internal/feed/mocks/core_mocks.go --package=mock
	mockgen --source=internal/image/interfaces.go --destination=internal/image/mocks/core_mocks.go --package=mock
