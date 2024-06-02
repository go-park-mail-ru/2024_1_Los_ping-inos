package delivery

//. "main.go/internal/logs"

// func TestImageHandler_GetApi_LogsError(t *testing.T) {
// 	logger := InitLog()
// 	api := GetApi(&usecase.UseCase{}, logger)

// 	grpcConn, _ := grpc.Dial("auth:50051", grpc.WithInsecure())
// 	authManager := gen.NewAuthHandlClient(grpcConn)

// 	req, _ := http.NewRequest(http.MethodGet, "/api/v1/getImage", nil)
// 	req = req.WithContext(context.WithValue(req.Context(), Logg, logger))
// 	req = req.WithContext(context.WithValue(req.Context(), 1, types.UserID(1)))

// 	err := api.mx.ServeHTTP(httptest.NewRecorder(), req)
// 	if err == nil {
// 		t.Errorf("Expected error, but got nil")
// 	}

// 	entries := logger.Entries
// 	if len(entries) == 0 {
// 		t.Errorf("Expected at least one log entry, but got none")
// 	}

// 	entry := entries[len(entries)-1]
// 	if entry.Level != logrus.ErrorLevel {
// 		t.Errorf("Expected log entry to be at error level, but got %s", entry.Level)
// 	}

// 	if !strings.Contains(entry.Message, "error.IsNotExistentEntity") {
// 		t.Errorf("Expected error message to contain 'error.IsNotExistentEntity', but got %s", entry.Message)
// 	}
// }
