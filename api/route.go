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

	upload := global.Router.Group("/upload")
	upload.Use(global.CheckAuth)
	upload.POST("/public", UploadPublicFile)
	upload.POST("/avatar", UploadAvatar)

	problem := global.Router.Group("/problem")
	problem.Use(global.CheckAuth)

	choiceProblem := problem.Group("/choice")
	choiceProblem.GET("", GetChoiceProblems)
	choiceProblem.POST("/create", CreateChoiceProblem)
	choiceProblem.PUT("/update", UpdateChoiceProblem)
	choiceProblem.DELETE("/delete/:id", DeleteChoiceProblem)

	note := global.Router.Group("/note")
	note.Use(global.CheckAuth)
	note.GET("", GetNotes)
	note.POST("/create", CreateNote)
	note.PUT("/update", UpdateNote)
	note.DELETE("/delete/:id", DeleteNote)
	/*
		TODO: 以下路由需要添加
		problem := global.Router.Group("/problem")
		problem.Use(global.CheckAuth)
		problem.GET("/all", GetAllProblems)
		problem.GET("/my", GetMyProblems)

		fillProblem := problem.Group("/fill")
		judgeProblem := problem.Group("/judge")

		fillProblem.GET("", GetFillProblems)
		judgeProblem.GET("", GetJudgeProblems)

		fillProblem.POST("/create", CreateFillProblem)
		judgeProblem.POST("/create", CreateJudgeProblem)

		fillProblem.POST("/update", UpdateFillProblem)
		judgeProblem.POST("/update", UpdateJudgeProblem)

		fillProblem.POST("/delete", DeleteFillProblem)
		judgeProblem.POST("/delete", DeleteJudgeProblem)

		problemset := global.Router.Group("/problemset")
		problemset.Use(global.CheckAuth)
		problemset.GET("/all", GetAllProblemsets)
		problemset.GET("/my", GetMyProblemsets)

		problemset.POST("/add-problem", AddProblem)
		problemset.POST("/remove-problem", RemoveProblem)

		problemset.POST("/create", CreateProblemset)
		problemset.POST("/update", UpdateProblemset)
		problemset.POST("/delete", DeleteProblemset)
		problemset.POST("/favorite", FavoriteProblemset)
	*/
}
