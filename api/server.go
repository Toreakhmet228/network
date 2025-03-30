package main

import (
	pb "chat/api/pb"
	"chat/internal/auth"
	"chat/internal/database"
	"chat/internal/models"
	"chat/pkg/conf"
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	pb.UnimplementedAuthServer
}

var logger *zap.Logger

func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	var existsUser models.User

	if err := database.DB.Where("email = ?", in.GetEmail()).First(&existsUser).Error; err != nil {
		logger.Error("Пользователь не найден", zap.Error(err))
		return nil, fmt.Errorf("неверный email или пароль")
	}

	logger.Info("Найден пользователь", zap.String("email", existsUser.Email)) // ✅ Логируем email

	if !auth.CheckPassword(existsUser.Password, in.GetPassword()) {
		logger.Error("Неверный пароль")
		return nil, fmt.Errorf("неверный email или пароль")
	}

	token, err := auth.GenerateJwt(existsUser.Email)
	if err != nil {
		logger.Error("Ошибка генерации токена", zap.Error(err))
		return nil, err
	}

	logger.Info("Токен успешно создан", zap.String("token", token)) // ✅ Логируем токен

	return &pb.LoginResponse{Token: token}, nil
}

func (s *Server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.ReqisterResponse, error) {
	var existingUser models.User

	if err := database.DB.Where("email = ?", in.GetEmail()).First(&existingUser).Error; err == nil {
		logger.Error("Пользователь с таким email уже существует", zap.String("email", in.GetEmail()))
		return &pb.ReqisterResponse{Succes: false}, fmt.Errorf("пользователь уже зарегистрирован")
	}

	hashedPassword, err := auth.HashPassword(in.GetPassword())
	if err != nil {
		logger.Error("Ошибка хеширования пароля", zap.Error(err))
		return &pb.ReqisterResponse{Succes: false}, err
	}

	newUser := models.User{
		Email:    in.GetEmail(),
		Password: hashedPassword,
	}
	if err := database.DB.Create(&newUser).Error; err != nil {
		logger.Error("Ошибка при создании пользователя", zap.Error(err))
		return &pb.ReqisterResponse{Succes: false}, err
	}

	return &pb.ReqisterResponse{Succes: true}, nil
}

func main() {
	logger = conf.InitLogger()
	database.InitDB()
	database.DB.AutoMigrate(&models.User{}, &models.Comment{}, &models.Like{}, &models.Group{}, &models.Post{}, &models.Like{})
	go func() {
		listener, err := net.Listen("tcp", ":50010")
		if err != nil {
			logger.Fatal("Ошибка запуска gRPC: " + err.Error())
		}
		defer listener.Close()

		server := grpc.NewServer()
		pb.RegisterAuthServer(server, &Server{})

		logger.Info("gRPC сервер запущен на :50010")
		if err := server.Serve(listener); err != nil {
			logger.Fatal("Ошибка работы gRPC: " + err.Error())
		}
	}()

	mux := runtime.NewServeMux()
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterAuthHandlerFromEndpoint(ctx, mux, "localhost:50010", opts); err != nil {
		logger.Fatal("Ошибка запуска gRPC-Gateway", zap.Error(err))
	}

	logger.Info("gRPC-Gateway запущен на :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Fatal("Ошибка запуска HTTP сервера", zap.Error(err))
	}
}
