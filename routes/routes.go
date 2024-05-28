package routes

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/jeremyadamsfisher/fake-user-cookie/cli"
	"github.com/labstack/echo/v4"
)

type Route struct {
	HandlerFunc http.HandlerFunc
	Path        string
}

func AddToEcho(e *echo.Echo, cfg *cli.Config) {
	e.GET("/", loginProxyHandler(cfg))
	e.GET("/manage", manageHandler)
	e.GET("/about", aboutHandler)
	e.GET("/oauth/sign_out", logOut)
	e.GET("/oauth/userinfo", userInfo)
	e.POST("/components/login", loginComponent)
}

func loginProxyHandler(cfg *cli.Config) func(*echo.Context) error {
	targetURL, err := url.Parse(cfg.DownstreamURL)
	if err != nil {
		panic("Mal-formatted downstream")
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	return func(c *echo.Context) error {
		if _, err := c.Cookie("email"); err != nil {
			// Intercept and require log-in
			return render(c, http.StatusOK, Home())
		}
		// Forward downstream
		req := c.Request()
		req.Host = targetURL.Host
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		proxy.ServeHTTP(c.Response(), req)
		return nil
	}
}

func aboutHandler(c echo.Context) error {
	return render(c, http.StatusOK, About())
}

func manageHandler(c echo.Context) error {
	cookie, err := c.Cookie("email")
	if err != nil {
		return render(c, http.StatusOK, UserSessionManager("", false))
	}
	return render(c, http.StatusOK, UserSessionManager(cookie.Value, cookie.Value != ""))
}

func logOut(c echo.Context) error {
	cookie, err := c.Cookie("email")
	if err != nil {
		return c.String(http.StatusOK, "Could not access cookie!")
	}
	cookie.Value = ""
	cookie.Path = "/"
	cookie.Secure = false // so that we can use on a non-https domain
	cookie.Expires = time.Now()
	c.SetCookie(cookie)
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func userInfo(c echo.Context) error {
	cookie, err := c.Cookie("email")
	if err != nil {
		return c.String(http.StatusOK, "Could not access cookie!")
	}
	return c.JSON(200, &struct {
		Email string `json:"email"`
	}{
		Email: cookie.Value,
	})
}

func loginComponent(c echo.Context) error {
	email := c.FormValue("email")
	if isEmailValid(email) {
		cookie := new(http.Cookie)
		cookie.Name = "email"
		cookie.Value = email
		cookie.Path = "/"
		cookie.Secure = false // so that we can use on a non-https domain
		cookie.Expires = time.Now().Add(24 * time.Hour)
		c.SetCookie(cookie)
	}
	return render(c, http.StatusOK, Login(email, isEmailValid(email)))
}
