package utils

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/KennethTrinh/go-vite-app/config"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func Error(ctx fiber.Ctx, statusCode int, goError interface{}, clientError interface{}) error {
	var returnObject fiber.Map
	var errorMessage string

	if config.Env.Production {
		// In production, send the client-friendly error message
		if reflect.TypeOf(clientError).Kind() == reflect.String {
			errorMessage = clientError.(string)
		} else {
			errorMessage = fmt.Sprintf("%v", clientError)
		}
	} else {
		// In non-production, send the detailed Go error message
		if reflect.TypeOf(goError).Kind() == reflect.String {
			errorMessage = goError.(string)
		} else {
			errorMessage = fmt.Sprintf("%v", goError)
		}
	}

	if statusCode >= 500 {
		var err error
		if reflect.TypeOf(goError).Kind() == reflect.String {
			err = errors.New(goError.(string))
		} else {
			err = fmt.Errorf("%v", goError)
		}
		log.Error().Err(err).Msg(errorMessage)
	}

	returnObject = fiber.Map{
		"success": false,
		"errors":  []string{errorMessage},
	}

	return ctx.Status(statusCode).JSON(returnObject)
}
