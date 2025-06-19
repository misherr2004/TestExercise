package httpServer

import (
	"awesomeProject/internal/apiServer/controllers"
	"awesomeProject/internal/config"
	"awesomeProject/internal/userservice"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	router   *chi.Mux
	handlers *controllers.Handler
	log      *log.Logger
}

func NewServer(u userservice.IUserService, log *zap.Logger) *Server {
	server := &Server{
		router:   chi.NewRouter(),
		handlers: controllers.NewHandler(u, log),
	}
	server.setupRoutes()
	return server
}

func (s *Server) Run(cfg *config.Config) {
	const op = "Run"
	log.Println("Cервер запущен")
	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      s.router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Ожидание запросов")
		err := http.ListenAndServe(cfg.HTTPServer.Address, s.router)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error with starting server in %v: %v", op, err)
		}
	}()

	<-stop
	log.Println("Сервер выключается...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("error with stopping server in %s: %v", op, err)

	} else {
		log.Println("Сервер завершен КОРРЕКТНО")
	}

}

func (s *Server) Middleware(next http.Handler) http.Handler {
	const op = "Middleware"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("error with Middleware in %v: %v", op, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) setupRoutes() {
	const op = "setupRoutes"
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Сервер работает")) //проверка работоспособности сервера
		if err != nil {
			log.Printf("error recording the response in %v: %v", op, err)
		}
	})

	s.router.Post("/users", s.handlers.SaveUserHandler)
	s.router.Get("/users/search", s.handlers.GetUserHandler)
	s.router.Get("/users/list", s.handlers.ListUsersHandler)
	//s.router.Delete("/users/delete/soft", s.handlers.SoftDeleteUserHandler)
	//s.router.Delete("/users/delete/hard", s.handlers.HardDeleteUserHandler)

}
