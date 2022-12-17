package controller

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	"time"
)

type Table struct {
	ID       int64 `json:"id"`
	Capacity int64 `json:"capacity" validate:"required,gt=0,lt=4294967295"`
}

type GuestName struct {
	Name string `json:"name" validate:"min=0,max=100"`
}

type GuestList struct {
	Table              int64  `json:"table" validate:"required,gt=0,lt=4294967295"`
	AccompanyingGuests int64  `json:"accompanying_guests" validate:"omitempty,gte=0,lt=4294967295"`
	Name               string `json:"name"`
}

type ArrivedGuests struct {
	TimeArrived        time.Time `json:"time_arrived"`
	AccompanyingGuests int64     `json:"accompanying_guests"`
	Name               string    `json:"name"`
}

type EmptySeats struct {
	SeatsEmpty int64 `json:"seats_empty"`
}

func validateCapacity(capacity int64) (string, error) {
	return validatePositiveInteger(capacity, false)
}

func validateAccompanyingGuests(accompanyingGuests int64) (string, error) {
	return validatePositiveInteger(accompanyingGuests, true)
}

func validatePositiveInteger(integer int64, zeroAllowed bool) (string, error) {
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		return "", fmt.Errorf("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		return "", err
	}

	var err error
	if zeroAllowed {
		err = v.Var(integer, "gte=0,lt=4294967295")
	} else {
		err = v.Var(integer, "gt=0,lt=4294967295")
	}
	errs := translateError(err, trans)
	var strerrs string
	for _, e := range errs {
		strerrs += e.Error()
	}
	return strerrs, nil
}

func validateName(name string) (string, error) {
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		return "", fmt.Errorf("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		return "", err
	}

	err := v.Var(name, "required,min=0,max=100")
	errs := translateError(err, trans)
	strerrs := ""
	for _, e := range errs {
		strerrs += e.Error()
	}
	return strerrs, nil
}

func validateGuestList(guestList GuestList) (string, error) {
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		return "", fmt.Errorf("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		return "", err
	}

	err := v.Struct(guestList)
	errs := translateError(err, trans)
	strerrs := ""
	for _, e := range errs {
		strerrs += e.Error()
	}
	return strerrs, nil
}

func translateError(err error, trans ut.Translator) (errs []error) {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		errs = append(errs, translatedErr)
	}
	return errs
}
