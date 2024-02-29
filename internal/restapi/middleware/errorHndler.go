package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, ginErr := range c.Errors {
			log.Println("HEREE", ginErr)
			// // log.Println("whoops", ginErr)
			// if restAPIErr, ok := ginErr.Err.(*restapierror.RestAPIError); ok {
			// 	switch restAPIErr.Code {
			// 	case 400:
			// 		c.AbortWithStatusJSON(restAPIErr.Code, restAPIErr)
			// 	}
			// 	log.Println("RestAPIError:", restAPIErr.Code, restAPIErr.Message, restAPIErr.Details)
			// }
		}
		// log.Println(c.Errors)
		// c.JSON(gin)

	}
}
