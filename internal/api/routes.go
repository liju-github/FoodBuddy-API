package api

import (
	"foodbuddy/internal/controllers"
	"foodbuddy/view"

	"github.com/gin-gonic/gin"
)

func ServerHealth(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "server status ok",
		})
	})
}
func AuthenticationRoutes(router *gin.Engine) {
	// Authentication Endpoints
	//admin
	router.POST("/api/v1/auth/admin/login", controllers.AdminLogin) //

	//user
	router.POST("/api/v1/auth/user/email/login", controllers.EmailLogin)   //
	router.POST("/api/v1/auth/user/email/signup", controllers.EmailSignup) //
	router.GET("/api/v1/auth/google/login", controllers.GoogleHandleLogin) //
	router.GET("/api/v1/googlecallback", controllers.GoogleHandleCallback) //

	//additional endpoints for email verification and password reset
	router.GET("/api/v1/auth/verifyemail/:role/:email/:otp", controllers.VerifyEmail) //

	router.POST("/api/v1/auth/passwordreset/step1", controllers.Step1PasswordReset) //
	router.GET("/api/v1/auth/passwordreset", controllers.LoadPasswordReset)         //
	router.POST("/api/v1/auth/passwordreset/step2", controllers.Step2PasswordReset) //

	//restaurant
	router.POST("/api/v1/auth/restaurant/signup", controllers.RestaurantSignup) //
	router.POST("/api/v1/auth/restaurant/login", controllers.RestaurantLogin)   //
}

func UserRoutes(router *gin.Engine) {
	userRoutes := router.Group("/api/v1/user")
	{
		// User Profile Management
		userRoutes.GET("/profile", controllers.GetUserProfile)       //
		userRoutes.POST("/edit", controllers.UpdateUserInformation)  //
		userRoutes.GET("/wallet/all", controllers.GetUserWalletData) //

		// Favorite Products
		userRoutes.GET("/favorites/all", controllers.GetUsersFavouriteProduct)     //
		userRoutes.POST("/favorites/add", controllers.AddFavouriteProduct)         //
		userRoutes.DELETE("/favorites/delete", controllers.RemoveFavouriteProduct) //

		// User Address Management
		userRoutes.GET("/address/all", controllers.GetUserAddress)          //
		userRoutes.POST("/address/add", controllers.AddUserAddress)         //
		userRoutes.PATCH("/address/edit", controllers.EditUserAddress)      //
		userRoutes.DELETE("/address/delete", controllers.DeleteUserAddress) //

		// Cart Management
		userRoutes.POST("/cart/add", controllers.AddToCart)               //
		userRoutes.GET("/cart/all", controllers.GetCartTotal)             //
		userRoutes.DELETE("/cart/delete/", controllers.ClearCart)         //
		userRoutes.DELETE("/cart/remove", controllers.RemoveItemFromCart) //
		userRoutes.PUT("/cart/update/", controllers.UpdateQuantity)       //

		// Order Management
		userRoutes.POST("/order/step1/placeorder", controllers.PlaceOrder)
		userRoutes.POST("/order/step2/initiatepayment", controllers.InitiatePayment)
		userRoutes.POST("/order/step3/razorpaycallback/:orderid", controllers.RazorPayGatewayCallback)
		userRoutes.GET("/order/step3/stripecallback", controllers.StripeCallback)
		userRoutes.POST("/order/cancel", controllers.CancelOrderedProduct)
		userRoutes.POST("/order/history", controllers.UserOrderHistory)
		userRoutes.GET("/order/invoice/:orderid", controllers.GetOrderInfoByOrderIDAndGeneratePDF)
		userRoutes.POST("/order/paymenthistory", controllers.PaymentDetailsByOrderID)
		userRoutes.POST("/order/review", controllers.UserReviewonOrderItem)
		userRoutes.POST("/order/rating", controllers.UserRatingOrderItem)
		userRoutes.GET("/coupon/cart/:couponcode", controllers.ApplyCouponOnCart) //

		// Referral System
		userRoutes.GET("/referral/code", controllers.GetRefferalCode)
		userRoutes.PATCH("/referral/activate", controllers.ActivateReferral)
		userRoutes.GET("/referral/claim", controllers.ClaimReferralRewards)
		userRoutes.GET("/referral/stats", controllers.GetReferralStats)
	}
}

func RestaurantRoutes(router *gin.Engine) {
	restaurantRoutes := router.Group("/api/v1/restaurants")
	{
		// Restaurant Management
		restaurantRoutes.POST("/edit", controllers.EditRestaurant)                 //
		restaurantRoutes.POST("/products/add", controllers.AddProduct)             //
		restaurantRoutes.POST("/products/edit", controllers.EditProduct)           //
		restaurantRoutes.DELETE("/products/:productid", controllers.DeleteProduct) //

		// Order History and Status Updates
		restaurantRoutes.GET("/order/history/:status", controllers.OrderHistoryRestaurants) //
		restaurantRoutes.POST("/order/nextstatus", controllers.UpdateOrderStatusForRestaurant)

		// Product Offers
		restaurantRoutes.POST("/product/offer/add", controllers.AddProductOffer)                 //
		restaurantRoutes.PUT("/product/offer/remove/:productid", controllers.RemoveProductOffer) //

		//orderitem information in excel
		restaurantRoutes.GET("/orderitems/excel/all", controllers.OrderInformationsCSVFileForRestaurant) //

		//restaurant wallet balance and history
		restaurantRoutes.GET("/wallet/all", controllers.GetRestaurantWalletData) //
		// restaurantRoutes.POST("/profile/update",controllers.RestaurantProfileUpdate)
	}
}

func AdminRoutes(router *gin.Engine) {
	adminRoutes := router.Group("/api/v1/admin")
	{
		// User Management
		adminRoutes.GET("/users", controllers.GetUserList)                 //
		adminRoutes.GET("/users/blocked", controllers.GetBlockedUserList)  //
		adminRoutes.PUT("/users/block/:userid", controllers.BlockUser)     //
		adminRoutes.PUT("/users/unblock/:userid", controllers.UnblockUser) //

		// Category Management
		adminRoutes.POST("/categories/add", controllers.AddCategory)                     //
		adminRoutes.PATCH("/categories/edit", controllers.EditCategory)                  //
		adminRoutes.DELETE("/categories/delete/:categoryid", controllers.DeleteCategory) //

		// Restaurant Management
		adminRoutes.GET("/restaurants", controllers.GetRestaurants)
		// adminRoutes.DELETE("/restaurants/:restaurantid", controllers.DeleteRestaurant)
		adminRoutes.PUT("/restaurants/block/:restaurantid", controllers.BlockRestaurant)     //
		adminRoutes.PUT("/restaurants/unblock/:restaurantid", controllers.UnblockRestaurant) //

		// Coupon Management
		adminRoutes.POST("/coupon/create", controllers.CreateCoupon)  //
		adminRoutes.PATCH("/coupon/update", controllers.UpdateCoupon) //
	}
}

func PublicRoutes(router *gin.Engine) {
	// Public API Endpoints
	publicRoute := router.Group("/api/v1/public")
	{
		publicRoute.GET("/coupon/all", controllers.GetAllCoupons)                                     //
		publicRoute.GET("/categories", controllers.GetCategoryList)                                   //
		publicRoute.GET("/categories/products", controllers.GetCategoryProductList)                   //
		publicRoute.GET("/products", controllers.GetProductList)                                      //
		publicRoute.GET("/restaurants", controllers.GetRestaurants)                                   //
		publicRoute.GET("/restaurants/products/:restaurantid", controllers.GetProductsByRestaurantID) //
		publicRoute.GET("/products/onlyveg", controllers.OnlyVegProducts)                             //
		publicRoute.GET("/products/newarrivals", controllers.NewArrivals)                             //
		publicRoute.GET("/products/lowtohigh", controllers.PriceLowToHigh)                            //
		publicRoute.GET("/products/hightolow", controllers.PriceHighToLow)                            //
		publicRoute.GET("/products/offerproducts", controllers.GetProductOffers)                      //
		publicRoute.GET("/report/products/:productid", controllers.ProductReport)                     //
		publicRoute.GET("/report/products/best", controllers.BestSellingProducts)                     //
		publicRoute.GET("/report/totalorders/all", controllers.PlatformOverallSalesReport)            //

	}
}

func AdditionalRoutes(router *gin.Engine) {
	// Additional Endpoints
	router.GET("/api/v1/user/profileimage", view.LoadUpload)                                 //
	router.POST("/api/v1/user/profileimage", controllers.UserProfileImageUpload)             //
	router.GET("/api/v1/restaurant/profileimage", view.LoadUpload)                           //
	router.POST("/api/v1/restaurant/profileimage", controllers.RestaurantProfileImageUpload) //
	router.GET("/api/v1/logout", controllers.Logout)                                         //
}

//public routes related to sales -- total sales
//User wallet and history
//restaurant eallet and history
//change in wallet amounts
//order invoice pdf
//order informations csv file
//*Referral
//coupon minimum amount
