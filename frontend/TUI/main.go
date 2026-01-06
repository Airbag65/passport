package main

import (
	"fmt"
	"pwd-manager-tui/auth"
	pass "pwd-manager-tui/passport"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Println("Exited passport")
		}
	}()

	fmt.Print("\033[H\033[2J")
	for !auth.ValidTokenExists() {
		startScreenModel := auth.NewStartScreenModel(new(int))
		startScreen := tea.NewProgram(startScreenModel, tea.WithAltScreen())
		startScreen.Run()

		switch startScreenModel.GetValue() {
		case 0:
			loginModel := auth.NewLoginModel()
			loginScreen := tea.NewProgram(loginModel, tea.WithAltScreen())
			loginScreen.Run()
			_, err := auth.Login(loginModel.GetValues()[0], loginModel.GetValues()[1])
			if err != nil {
				fmt.Printf("Could not login: %v", err)
			}
		case 1:
			signUpModel := auth.NewSignUpModel()
			signUpScreen := tea.NewProgram(signUpModel, tea.WithAltScreen())
			signUpScreen.Run()
			_, err := auth.SignUp(signUpModel.GetValues())
			if err != nil {
				fmt.Printf("Could not sign up: %v", err)
			}
		case 2:
			panic("exit")
		}
	}
	mainScreenModel := pass.NewMainScreenModel()
	mainscreen := tea.NewProgram(mainScreenModel, tea.WithAltScreen())
	mainscreen.Run()
}
