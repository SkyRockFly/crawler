package main

import "os"

type authCookies struct {
	cookieA string
	cookieB string
}

func (a *authCookies) cookieFromEnv() {
	a.cookieA = os.Getenv("COOKIE_A")
	a.cookieB = os.Getenv("COOKIE_B")
}
