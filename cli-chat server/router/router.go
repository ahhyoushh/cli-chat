package router

import (
	"github.com/TanishkBansode/cli-chat/controller"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/getall", controller.GetAllMsgs).Methods("POST")

	// JSON DATA TO POST for /api/getall:

	// {
	// 	"userid": "6628d540cbd359c2a35858e6",
	// 	"password": "ayush"
	//   }

	router.HandleFunc("/api/getunread", controller.GetUnreadMsgs).Methods("POST")

	// JSON DATA TO POST for /api/getunread:

	// {
	// 	"userid": "6628d540cbd359c2a35858e6",
	// 	"password": "ayush"
	//   }

	router.HandleFunc("/api/send", controller.CreateMsg).Methods("POST")

	// JSON DATA TO POST for /api/send:

	// {
	// 	"sender": "bandu",
	// 	"receiver": "ayush",
	// 	"message": "ntg"
	//   }

	// Currently there are only 2 users established:
	// bandu, ID: 6628d59ccbd359c2a35858e7, password: bandu
	// ayush, ID: 6628d540cbd359c2a35858e6, password: ayush

	router.HandleFunc("/api/signup", controller.Signup).Methods("POST")

	// JSON DATA TO POST for /api/signup:

	// {
	// 	"username": "xyz",
	// 	"password": "xyz",
	//   }

	router.HandleFunc("/api/login", controller.Login).Methods("POST")

	// JSON DATA TO POST for /api/login:

	// {
	// 	"username": "xyz",
	// 	"password": "xyz",
	//   }

	return router
}
