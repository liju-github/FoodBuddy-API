package controllers

import (
	"fmt"
	"foodbuddy/internal/database"
	"foodbuddy/internal/model"
	"foodbuddy/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// create coupons -admin side
func CreateCoupon(c *gin.Context) { //admin
	// check admin api authentication
	_, role, err := utils.GetJWTClaim(c)
	if role != model.AdminRole || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "unauthorized request",
		})
		return
	}
	var Request model.CouponInventoryRequest
	if err := c.BindJSON(&Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to bind the json",
		})
		return
	}

	if err := utils.Validate(&Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	if CheckCouponExists(Request.CouponCode) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "coupon code already exists",
		})
		return
	}

	if Request.Percentage > model.CouponDiscountPercentageLimit {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "coupon discount percentage should not exceed more than " + strconv.Itoa(model.CouponDiscountPercentageLimit),
		})
		return
	}

	if time.Now().Unix()+12*3600 > int64(Request.Expiry) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "please change the expiry time that is more than a day",
		})
		return
	}

	Coupon := model.CouponInventory{
		CouponCode:    Request.CouponCode,
		Expiry:        Request.Expiry,
		Percentage:    Request.Percentage,
		MaximumUsage:  Request.MaximumUsage,
		MinimumAmount: float64(Request.MinimumAmount),
	}

	if err := database.DB.Create(&Coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to create coupon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully created coupon",
	})
}

func GetAllCoupons(c *gin.Context) { //public
	var Coupons []model.CouponInventory

	if err := database.DB.Find(&Coupons).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to fetch coupon details",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   Coupons,
	})
}

// update coupon
func UpdateCoupon(c *gin.Context) { //admin
	// check admin api authentication
	_, role, err := utils.GetJWTClaim(c)
	if role != model.AdminRole || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "unauthorized request",
		})
		return
	}
	var request model.CouponInventoryRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "failed to bind the json",
		})
		return
	}

	if err := utils.Validate(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	var existingCoupon model.CouponInventory
	err = database.DB.Where("coupon_code = ?", request.CouponCode).First(&existingCoupon).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": "coupon not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to find coupon",
		})
		return
	}

	existingCoupon.Expiry = request.Expiry
	existingCoupon.Percentage = request.Percentage
	existingCoupon.MaximumUsage = request.MaximumUsage

	if err := database.DB.Save(&existingCoupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "failed to update coupon",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "successfully updated coupon",
	})
}

func ApplyCouponOnCart(c *gin.Context) { //user
	// check restaurant api authentication
	email, role, err := utils.GetJWTClaim(c)
	if role != model.UserRole || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "unauthorized request",
		})
		return
	}
	UserID, _ := UserIDfromEmail(email)
	CouponCode := c.Query("couponcode")
	RestaurantID := c.Query("restaurant_id")

	if CouponCode == "" || RestaurantID == "" || RestaurantID == "0" {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "provide couponcode and restaurant_id in the query params"})
		return
	}

	var CartItems []model.CartItems

	if err := database.DB.Where("user_id = ? AND restaurant_id = ?", UserID, RestaurantID).Find(&CartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to fetch cart items. Please try again later.",
		})
		return
	}

	if len(CartItems) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Your cart is empty.",
		})
		return
	}

	// Total price of the cart
	var sum, ProductOfferAmount float64
	for _, item := range CartItems {
		var Product model.Product
		if err := database.DB.Where("id = ?", item.ProductID).First(&Product).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": "Failed to fetch product information. Please try again later.",
			})
			return
		}

		ProductOfferAmount += float64(Product.OfferAmount) * float64((item.Quantity))
		sum += ((Product.Price) * float64(item.Quantity))
	}

	// Apply coupon if provided
	var CouponDiscount float64
	var FinalAmount float64
	if CouponCode != "" {
		var coupon model.CouponInventory
		if err := database.DB.Where("coupon_code = ?", CouponCode).First(&coupon).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": "Invalid coupon code. Please check and try again.",
			})
			return
		}

		// Check coupon expiration
		if time.Now().Unix() > int64(coupon.Expiry) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "The coupon has expired.",
			})
			return
		}

		//check minimum amount
		if sum < coupon.MinimumAmount {
			errmsg := fmt.Sprintf("minimum of %v is needed for using this coupon", coupon.MinimumAmount)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": errmsg,
			})
			return
		}

		// Check coupon usage
		var usage model.CouponUsage
		if err := database.DB.Where("user_id = ? AND coupon_code = ?", UserID, CouponCode).First(&usage).Error; err == nil {
			if usage.UsageCount >= coupon.MaximumUsage {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "The coupon usage limit has been reached.",
				})
				return
			}
		}

		// Calculate discount
		CouponDiscount = float64(sum) * (float64(coupon.Percentage) / 100.0)
		FinalAmount = sum - (CouponDiscount + ProductOfferAmount)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data": gin.H{
			"restaurant_id":        RestaurantID,
			"cart_items":           CartItems,
			"total_amount":         sum,
			"coupon_discount":      CouponDiscount,
			"product_offer_amount": ProductOfferAmount,
			"final_amount":         FinalAmount,
		},
		"message": "Cart items retrieved successfully",
	})
}

func ApplyCouponToOrder(order model.Order, UserID uint, CouponCode string) (bool, string, model.Order) {

	if order.CouponCode != "" {
		errMsg := fmt.Sprintf("%v coupon already exists, remove this coupon to add a new coupon", order.CouponCode)
		return false, errMsg, order
	}

	var coupon model.CouponInventory
	if err := database.DB.Where("coupon_code = ?", CouponCode).First(&coupon).Error; err != nil {
		return false, "coupon not found", order
	}

	if time.Now().Unix() > int64(coupon.Expiry) {
		return false, "coupon has expired", order
	}

	var couponUsage model.CouponUsage
	err := database.DB.Where("coupon_code = ? AND user_id = ?", CouponCode, UserID).First(&couponUsage).Error

	if err == nil {
		if couponUsage.UsageCount >= coupon.MaximumUsage {
			return false, "coupon usage limit reached", order
		}
	} else if err != gorm.ErrRecordNotFound {
		return false, "database error", order
	}

	//check minimum amount
	if order.TotalAmount < coupon.MinimumAmount {
		errmsg := fmt.Sprintf("minimum of %v is needed for using this coupon", coupon.MinimumAmount)
		return false, errmsg, order
	}

	discountAmount := order.TotalAmount * float64(coupon.Percentage) / 100
	finalAmount := order.TotalAmount - (discountAmount + order.ProductOfferAmount)

	order.CouponCode = CouponCode
	order.CouponDiscountAmount = discountAmount
	order.FinalAmount = finalAmount

	if err := database.DB.Where("order_id = ?", order.OrderID).Updates(&order).Error; err != nil {
		return false, "failed to apply coupon to order", order
	}

	if err == gorm.ErrRecordNotFound {
		couponUsage = model.CouponUsage{
			UserID:     UserID,
			CouponCode: CouponCode,
			UsageCount: 1,
		}
		if err := database.DB.Create(&couponUsage).Error; err != nil {
			return false, "failed to create coupon usage record", order
		}
	} else {
		couponUsage.UsageCount++
		if err := database.DB.Where("user_id = ? AND coupon_code = ?", order.UserID, order.CouponCode).Save(&couponUsage).Error; err != nil {
			return false, "failed to update coupon usage record", order
		}
	}

	return true, "coupon applied successfully", order
}

func CheckCouponExists(code string) bool {
	var Coupons []model.CouponInventory
	if err := database.DB.Find(&Coupons).Error; err != nil {
		return false
	}
	fmt.Println(&Coupons)
	for _, c := range Coupons {
		if c.CouponCode == code {
			return true
		}
	}
	return false
}
