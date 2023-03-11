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

	/*
		TODO: 以下路由需要添加
		problem := global.Router.Group("/problem")
		problem.Use(global.CheckAuth)
		problem.GET("/all", GetAllProblems)
		problem.GET("/my", GetMyProblems)

		choiceProblem := problem.Group("/choice")
		fillProblem := problem.Group("/fill")
		judgeProblem := problem.Group("/judge")

		choiceProblem.GET("", GetChoiceProblems)
		fillProblem.GET("", GetFillProblems)
		judgeProblem.GET("", GetJudgeProblems)

		choiceProblem.POST("/create", CreateChoiceProblem)
		fillProblem.POST("/create", CreateFillProblem)
		judgeProblem.POST("/create", CreateJudgeProblem)

		choiceProblem.POST("/update", UpdateChoiceProblem)
		fillProblem.POST("/update", UpdateFillProblem)
		judgeProblem.POST("/update", UpdateJudgeProblem)

		choiceProblem.POST("/delete", DeleteChoiceProblem)
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
