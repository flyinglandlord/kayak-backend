package api

import "kayak-backend/global"

func InitRoute() {
	global.Router.GET("/ping", Ping)
	global.Router.GET("/logout", Logout)
	global.Router.POST("/login", Login)
	global.Router.POST("/register", Register)
	global.Router.POST("/reset-password", ResetPassword)

	user := global.Router.Group("/user")
	user.Use(global.CheckAuth)
	user.GET("/info", GetUserInfo)
	user.PUT("/update", UpdateUserInfo)
	user.GET("/problem/choice", GetUserChoiceProblems)
	user.GET("/problem/blank", GetUserBlankProblems)
	user.GET("problemset", GetUserProblemsets)
	user.GET("/note", GetUserNotes)

	upload := global.Router.Group("/upload")
	upload.Use(global.CheckAuth)
	upload.POST("/public", UploadPublicFile)
	upload.POST("/avatar", UploadAvatar)

	problem := global.Router.Group("/problem")
	problem.Use(global.CheckAuth)

	choiceProblem := problem.Group("/choice")
	global.Router.GET("/problem/choice/all", GetChoiceProblems)
	choiceProblem.POST("/create", CreateChoiceProblem)
	choiceProblem.PUT("/update", UpdateChoiceProblem)
	choiceProblem.DELETE("/delete/:id", DeleteChoiceProblem)

	blankProblem := problem.Group("/blank")
	global.Router.GET("/problem/blank/all", GetBlankProblems)
	blankProblem.POST("/create", CreateBlankProblem)
	blankProblem.PUT("/update", UpdateBlankProblem)
	blankProblem.DELETE("/delete/:id", DeleteBlankProblem)

	problemset := global.Router.Group("/problemset")
	problemset.Use(global.CheckAuth)
	global.Router.GET("/problemset/all", GetProblemsets)
	problemset.POST("/create", CreateProblemset)
	problemset.DELETE("/delete/:id", DeleteProblemset)
	problemset.PUT("/:id/add", AddProblemToProblemset)
	problemset.PUT("/:id/remove", RemoveProblemFromProblemset)

	note := global.Router.Group("/note")
	note.Use(global.CheckAuth)
	global.Router.GET("/note/all", GetNotes)
	note.POST("/create", CreateNote)
	note.PUT("/update", UpdateNote)
	note.DELETE("/delete/:id", DeleteNote)
}
