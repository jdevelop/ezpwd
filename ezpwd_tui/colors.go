package main

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell"
)

type loginForm struct {
	Background       tcell.Color
	Title            tcell.Color
	Border           tcell.Color
	Label            tcell.Color
	ButtonBackground tcell.Color
	ButtonText       tcell.Color
	FieldBackground  tcell.Color
	FieldText        tcell.Color
}

type passwordMgmtForm struct {
	Background       tcell.Color
	TitleAdd         tcell.Color
	TitleUpdate      tcell.Color
	BorderAdd        tcell.Color
	BorderUpdate     tcell.Color
	Label            tcell.Color
	ButtonBackground tcell.Color
	ButtonText       tcell.Color
	FieldBackground  tcell.Color
	FieldText        tcell.Color
}

type confirmForm struct {
	Background       tcell.Color
	Title            tcell.Color
	Border           tcell.Color
	Label            tcell.Color
	ButtonBackground tcell.Color
	ButtonText       tcell.Color
	FieldBackground  tcell.Color
	FieldText        tcell.Color
}

type globalScreen struct {
	Background tcell.Color
	Title      tcell.Color
	Border     tcell.Color
	HelpText   tcell.Color
}

type passwordsTable struct {
	Background          tcell.Color
	Title               tcell.Color
	Border              tcell.Color
	Label               tcell.Color
	ButtonBackground    tcell.Color
	ButtonText          tcell.Color
	ButtonAccent        tcell.Color
	FieldBackground     tcell.Color
	FieldText           tcell.Color
	Selection           tcell.Color
	SelectionBackground tcell.Color
	Header              tcell.Color
	CopiedBackground    tcell.Color
	CopiedText          tcell.Color
	CopiedTitle         tcell.Color
	CopiedBorder        tcell.Color
}

type messages struct {
	SuccessBackground tcell.Color
	SuccessBorder     tcell.Color
	SuccessTitle      tcell.Color
	SuccessText       tcell.Color
	FailureBackground tcell.Color
	FailureBorder     tcell.Color
	FailureTitle      tcell.Color
	FailureText       tcell.Color
}

type colorSchema struct {
	LoginFormColors      loginForm
	PasswordMgmtColors   passwordMgmtForm
	ConfirmFormColors    confirmForm
	GlobalScreenColors   globalScreen
	PasswordsTableColors passwordsTable
	MessagesColors       messages
}

var (
	BackgroundColorLight = tcell.NewRGBColor(0x87, 0x87, 0x5f)
	BackgroundColorDark  = tcell.NewRGBColor(0x07, 0x36, 0x42)

	DefaultColorSchema colorSchema

	LightColorSchema = colorSchema{
		LoginFormColors: loginForm{
			Background:       BackgroundColorLight,
			Title:            tcell.ColorWhiteSmoke,
			Border:           tcell.ColorGray,
			Label:            tcell.ColorBlack,
			ButtonBackground: tcell.ColorDarkGray,
			ButtonText:       tcell.ColorWhite,
			FieldBackground:  tcell.ColorBeige,
			FieldText:        tcell.ColorBlack,
		},
		ConfirmFormColors: confirmForm{
			Background:       BackgroundColorLight,
			Title:            tcell.ColorWhiteSmoke,
			Border:           tcell.ColorGray,
			Label:            tcell.ColorBlack,
			ButtonBackground: tcell.ColorDarkGray,
			ButtonText:       tcell.ColorWhite,
			FieldBackground:  tcell.ColorBeige,
			FieldText:        tcell.ColorBlack,
		},
		PasswordMgmtColors: passwordMgmtForm{
			Background:       BackgroundColorLight,
			TitleAdd:         tcell.ColorWhiteSmoke,
			TitleUpdate:      tcell.ColorMistyRose,
			BorderAdd:        tcell.ColorGray,
			BorderUpdate:     tcell.ColorMistyRose,
			Label:            tcell.ColorBlack,
			ButtonBackground: tcell.ColorDarkGray,
			ButtonText:       tcell.ColorWhite,
			FieldBackground:  tcell.ColorBeige,
			FieldText:        tcell.ColorBlack,
		},
		GlobalScreenColors: globalScreen{
			Background: BackgroundColorLight,
			Title:      tcell.ColorNavajoWhite,
			Border:     tcell.ColorGray,
			HelpText:   tcell.ColorLightGray,
		},
		PasswordsTableColors: passwordsTable{
			Background:          BackgroundColorLight,
			Title:               tcell.ColorWhiteSmoke,
			Border:              tcell.ColorGray,
			Label:               tcell.ColorBlack,
			ButtonBackground:    tcell.ColorDarkGray,
			ButtonText:          tcell.ColorWhite,
			FieldBackground:     tcell.ColorBeige,
			FieldText:           tcell.ColorBlack,
			Selection:           tcell.ColorGreen,
			SelectionBackground: tcell.ColorWheat,
			Header:              tcell.ColorBisque,
			ButtonAccent:        tcell.ColorBlueViolet,
			CopiedBackground:    tcell.ColorGold,
			CopiedText:          tcell.ColorBlack,
			CopiedTitle:         tcell.ColorGreen,
			CopiedBorder:        tcell.ColorWhiteSmoke,
		},
		MessagesColors: messages{
			SuccessBackground: BackgroundColorLight,
			SuccessTitle:      tcell.ColorNavajoWhite,
			SuccessBorder:     tcell.ColorGray,
			SuccessText:       tcell.ColorGreen,
			FailureBackground: tcell.ColorBlack,
			FailureTitle:      tcell.ColorRed,
			FailureBorder:     tcell.ColorRed,
			FailureText:       tcell.ColorOrangeRed,
		},
	}
	DarkColorSchema = colorSchema{
		LoginFormColors: loginForm{
			Background:       BackgroundColorDark,
			Title:            tcell.ColorWhiteSmoke,
			Border:           tcell.ColorGray,
			Label:            tcell.ColorWheat,
			ButtonBackground: tcell.ColorDarkGray,
			ButtonText:       tcell.ColorWhite,
			FieldBackground:  tcell.ColorBeige,
			FieldText:        tcell.ColorBlack,
		},
		ConfirmFormColors: confirmForm{
			Background:       BackgroundColorDark,
			Title:            tcell.ColorWhiteSmoke,
			Border:           tcell.ColorGray,
			Label:            tcell.ColorWhiteSmoke,
			ButtonBackground: tcell.ColorDarkGray,
			ButtonText:       tcell.ColorWhite,
			FieldBackground:  tcell.ColorBeige,
			FieldText:        tcell.ColorBlack,
		},
		PasswordMgmtColors: passwordMgmtForm{
			Background:       BackgroundColorDark,
			TitleAdd:         tcell.ColorWhiteSmoke,
			TitleUpdate:      tcell.ColorMistyRose,
			BorderAdd:        tcell.ColorGray,
			BorderUpdate:     tcell.ColorMistyRose,
			Label:            tcell.ColorWhiteSmoke,
			ButtonBackground: tcell.ColorDarkGray,
			ButtonText:       tcell.ColorWhite,
			FieldBackground:  tcell.ColorBeige,
			FieldText:        tcell.ColorBlack,
		},
		GlobalScreenColors: globalScreen{
			Background: BackgroundColorDark,
			Title:      tcell.ColorNavajoWhite,
			Border:     tcell.ColorGray,
			HelpText:   tcell.ColorLightGray,
		},
		PasswordsTableColors: passwordsTable{
			Background:          BackgroundColorDark,
			Title:               tcell.ColorWhiteSmoke,
			Border:              tcell.ColorGray,
			Label:               tcell.ColorWhiteSmoke,
			ButtonBackground:    tcell.ColorLightGray,
			ButtonText:          tcell.ColorBlack,
			FieldBackground:     tcell.ColorBeige,
			FieldText:           tcell.ColorBlack,
			Selection:           tcell.ColorGreen,
			SelectionBackground: tcell.ColorWheat,
			Header:              tcell.ColorBisque,
			ButtonAccent:        tcell.ColorBlueViolet,
			CopiedBackground:    tcell.ColorGold,
			CopiedText:          tcell.ColorBlack,
			CopiedTitle:         tcell.ColorGreen,
			CopiedBorder:        tcell.ColorWhiteSmoke,
		},
		MessagesColors: messages{
			SuccessBackground: BackgroundColorDark,
			SuccessTitle:      tcell.ColorNavajoWhite,
			SuccessBorder:     tcell.ColorGray,
			SuccessText:       tcell.ColorGreen,
			FailureBackground: tcell.ColorBlack,
			FailureTitle:      tcell.ColorRed,
			FailureBorder:     tcell.ColorRed,
			FailureText:       tcell.ColorOrangeRed,
		},
	}
)

func color2Hex(c tcell.Color) string {
	r, g, b := c.RGB()
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func hex2Color(src string) tcell.Color {
	x, err := strconv.ParseInt(src[1:], 16, 32)
	if err != nil {
		panic(err)
	}
	return tcell.NewHexColor(int32(x))
}
