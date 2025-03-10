package utils

import (
	"context"
	"github.com/sriramr98/todo_task_service/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"strings"
)

func AuthHandler(handler ApiHandler) ApiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		authToken := r.Header.Get("Authorization")
		isValid, userId := validateTokenAndGetUserID(authToken)
		if !isValid {
			return ApiError{
				StatusCode: http.StatusUnauthorized,
				Code:       ErrUnauthorized,
				Message:    "Unauthorized",
			}
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		return handler(w, r.WithContext(ctx))
	}
}

func validateTokenAndGetUserID(token string) (bool, int64) {
	token = strings.TrimSpace(token)

	if token == "" {
		return false, 0
	}

	if !strings.HasPrefix(token, "Bearer ") {
		return false, 0
	}

	token = strings.Split(token, " ")[1]

	conn, err := grpc.NewClient(os.Getenv("GRPC_HOST"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return false, 0
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)

	authClient := services.NewAuthServiceClient(conn)

	userInfo, err := authClient.GetUserInfo(context.Background(), &services.UserRequest{
		Token: token,
	})
	if err != nil {
		log.Println(err)
		return false, 0
	}

	if userInfo == nil {
		log.Println("User Info Not found")
		return false, 0
	}

	return !userInfo.GetBlocked(), userInfo.GetUserId()

}
