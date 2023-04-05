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
	user.GET("/note", GetUserNotes)
	user.GET("/wrong_record", GetUserWrongRecords)
	user.GET("/favorite/problem", GetUserFavoriteProblems)
	user.GET("/favorite/problemset", GetUserFavoriteProblemsets)
	user.GET("/favorite/note", GetUserFavoriteNotes)
	user.GET("/problem/choice", GetUserChoiceProblems)
	user.GET("/problem/blank", GetUserBlankProblems)
	user.GET("/problemset", GetUserProblemsets)

	upload := global.Router.Group("/upload")
	upload.Use(global.CheckAuth)
	upload.POST("/public", UploadPublicFile)
	upload.POST("/avatar", UploadAvatar)

	note := global.Router.Group("/note")
	note.Use(global.CheckAuth)
	global.Router.GET("/note/all", GetNotes)
	note.POST("/create", CreateNote)
	note.PUT("/update", UpdateNote)
	note.DELETE("/delete/:id", DeleteNote)
	note.POST("/like/:id", LikeNote)
	note.POST("/unlike/:id", UnlikeNote)
	note.POST("/favorite/:id", FavoriteNote)
	note.POST("/unfavorite/:id", UnfavoriteNote)

	wrongRecord := global.Router.Group("/wrong_record")
	wrongRecord.Use(global.CheckAuth)
	wrongRecord.POST("/create/:id", CreateWrongRecord)
	wrongRecord.PUT("/increase/:id", IncreaseWrongRecord)
	wrongRecord.PUT("/decrease/:id", DecreaseWrongRecord)
	wrongRecord.DELETE("/delete/:id", DeleteWrongRecord)

	problem := global.Router.Group("/problem")
	problem.Use(global.CheckAuth)
	problem.DELETE("/unfavorite/:id", RemoveProblemFromFavorite)
	problem.POST("/favorite/:id", AddProblemToFavorite)

	choiceProblem := problem.Group("/choice")
	global.Router.GET("/problem/choice/all", GetChoiceProblems)
	global.Router.GET("/problem/choice/:id", GetChoiceProblem)
	choiceProblem.POST("/create", CreateChoiceProblem)
	choiceProblem.PUT("/update", UpdateChoiceProblem)
	choiceProblem.DELETE("/delete/:id", DeleteChoiceProblem)

	blankProblem := problem.Group("/blank")
	global.Router.GET("/problem/blank/all", GetBlankProblems)
	global.Router.GET("/problem/blank/:id", GetBlankProblem)
	blankProblem.POST("/create", CreateBlankProblem)
	blankProblem.PUT("/update", UpdateBlankProblem)
	blankProblem.DELETE("/delete/:id", DeleteBlankProblem)

	problemSet := global.Router.Group("/problemSet")
	problemSet.Use(global.CheckAuth)
	global.Router.GET("/problemSet/all", GetProblemsets)
	problemSet.POST("/create", CreateProblemset)
	problemSet.DELETE("/delete/:id", DeleteProblemset)
	problemSet.GET("/:id/all", GetProblemsInProblemset)
	problemSet.PUT("/:id/add", AddProblemToProblemset)
	problemSet.PUT("/:id/remove", RemoveProblemFromProblemset)
	problemSet.POST("/favorite/:id", AddProblemSetToFavorite)
	problemSet.DELETE("/unfavorite/:id", RemoveProblemSetFromFavorite)

	noteReview := global.Router.Group("/note_review")
	noteReview.Use(global.CheckAuth)
	noteReview.POST("/add", AddNoteReview)
	noteReview.POST("/remove", RemoveNoteReview)
	noteReview.GET("/get", GetNoteReviews)
}
