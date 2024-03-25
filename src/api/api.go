package api

import (
	"github.com/VolodymyrShabat/Test_ATN/src/api/handlers"
	"github.com/VolodymyrShabat/Test_ATN/src/app"
	"github.com/VolodymyrShabat/Test_ATN/src/storage/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type APIServer struct {
	config   *models.Config
	logger   *logrus.Logger
	router   *mux.Router
	Handlers *Handlers
}

type Handlers struct {
	user *handlers.UserAccount
	book *handlers.BookHandler
}

func New(config *models.Config, services *app.Services) *APIServer {
	var validate = validator.New()
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
		Handlers: &Handlers{
			user: handlers.NewUserAccountHandler(services.User, config, validate),
			book: handlers.NewBookHandler(services.Book, services.User, validate),
		},
	}
}

func (s *APIServer) Run() error {
	if err := s.ConfigureLogger(); err != nil {
		return err
	}

	s.ConfigureRouter()

	s.logger.Info("Starting api server")
	//"0.0.0.0"
	return http.ListenAndServe(s.config.Domain+s.config.Port, s.router)
}

func (s *APIServer) ConfigureLogger() error {
	level, err := logrus.ParseLevel(s.config.Env)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) ConfigureRouter() {
	s.router.HandleFunc("/user/sign-up", s.Handlers.user.SignUp()).Methods("POST")
	s.router.HandleFunc("/user/sign-in", s.Handlers.user.SignIn()).Methods("POST")
	s.router.HandleFunc("/user/forgot-password", s.Handlers.user.ForgotPassword()).Methods("POST")
	s.router.HandleFunc("/user/get/{id}", s.Handlers.user.GetUserById()).Methods("GET")
	s.router.HandleFunc("/user/reset-password/{reset_token}", s.Handlers.user.ResetPassword()).Methods("POST")

	s.router.HandleFunc("/book/get/{id}", s.Handlers.book.GetByBookId()).Methods("GET")
	s.router.HandleFunc("/book/delete/{id}", s.Handlers.book.DeleteBook()).Methods("DELETE")
	s.router.HandleFunc("/book/update/{id}", s.Handlers.book.UpdateBook()).Methods("POST")
	s.router.HandleFunc("/book/create", s.Handlers.book.CreateBook()).Methods("POST")

}
