package logger

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Init() {
	// Configuración básica del logger.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Info(msg string) {
	log.Println("INFO: " + msg)
}

func Error(msg string) {
	log.Println("ERROR: " + msg)
}

// GinLogger es un middleware para logging de peticiones HTTP
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Inicio del tiempo
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Procesar la petición
		c.Next()

		// Calcular tiempo de respuesta
		latency := time.Since(start)

		// Obtener información de la respuesta
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		// Loggear la información
		log.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s %s",
			start.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			errorMessage,
		)
	}
}
