package main

import (
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

var (
	BackgroundColor = tcell.NewRGBColor(0x87, 0x87, 0x5f)
	FormBolderColor = tcell.NewRGBColor(0xff, 0xf6, 0xe9)

	loginFormColors = loginForm{
		Background:       BackgroundColor,
		Title:            tcell.ColorWhiteSmoke,
		Border:           tcell.ColorGray,
		Label:            tcell.ColorBlack,
		ButtonBackground: tcell.ColorDarkGray,
		ButtonText:       tcell.ColorWhite,
		FieldBackground:  tcell.ColorBeige,
		FieldText:        tcell.ColorBlack,
	}

	confirmFormColors = confirmForm{
		Background:       BackgroundColor,
		Title:            tcell.ColorWhiteSmoke,
		Border:           tcell.ColorGray,
		Label:            tcell.ColorBlack,
		ButtonBackground: tcell.ColorDarkGray,
		ButtonText:       tcell.ColorWhite,
		FieldBackground:  tcell.ColorBeige,
		FieldText:        tcell.ColorBlack,
	}

	passwordMgmgFormColors = passwordMgmtForm{
		Background:       BackgroundColor,
		TitleAdd:         tcell.ColorWhiteSmoke,
		TitleUpdate:      tcell.ColorMistyRose,
		BorderAdd:        tcell.ColorGray,
		BorderUpdate:     tcell.ColorMistyRose,
		Label:            tcell.ColorBlack,
		ButtonBackground: tcell.ColorDarkGray,
		ButtonText:       tcell.ColorWhite,
		FieldBackground:  tcell.ColorBeige,
		FieldText:        tcell.ColorBlack,
	}

	globalScreenColors = globalScreen{
		Background: BackgroundColor,
		Title:      tcell.ColorNavajoWhite,
		Border:     tcell.ColorGray,
		HelpText:   tcell.ColorLightGray,
	}

	passwordsTableColors = passwordsTable{
		Background:          BackgroundColor,
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
	}

	messagesColors = messages{
		SuccessBackground: BackgroundColor,
		SuccessTitle:      tcell.ColorNavajoWhite,
		SuccessBorder:     tcell.ColorGray,
		SuccessText:       tcell.ColorGreen,
		FailureBackground: BackgroundColor,
		FailureTitle:      tcell.ColorRed,
		FailureBorder:     tcell.ColorRed,
		FailureText:       tcell.ColorRed,
	}
)
