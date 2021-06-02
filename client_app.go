package arpio

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const AppPollPeriod = 5 * time.Second

// NewApp creates an App struct with the Terraform AppType for the account
// the Client is configured to use, but does not create it in the Arpio
// service (use CreateApp).
func (c *Client) NewApp() (app App) {
	app.AccountID = c.AccountID
	app.AppType = TerraformAppType
	return app
}

// CreateApp creates an application in the Arpio service in the account the
// Client is configured to use.
func (c *Client) CreateApp(a App) (created App, err error) {
	u := fmt.Sprintf("/accounts/%s/applications", c.AccountID)

	_, err = c.apiPost(u, a, &created)
	if err != nil {
		return created, err
	}

	return created, nil
}

// ListApps lists all the applications in the account the Client is configured
// to use.
func (c *Client) ListApps() (apps []App, err error) {
	u := fmt.Sprintf("/accounts/%s/applications", c.AccountID)

	_, err = c.apiGet(u, &apps)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

// GetApp gets the application with the specified ID in the account the Client
// is configured to use.
func (c *Client) GetApp(appID string) (app *App, err error) {
	u := c.appPath(appID)

	var a App
	status, err := c.apiGet(u, &a)
	if status == http.StatusNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// UpdateApp updates the mutable properties of the specified application.
func (c *Client) UpdateApp(a App) (app App, err error) {
	u := c.appPath(a.AppID)

	_, err = c.apiPut(u, a, &app)
	if err != nil {
		return app, err
	}

	return app, nil
}

// DeleteApp deletes the application with the specified ID.  If the application
// does not exist, no error is returned.
func (c *Client) DeleteApp(appID string) error {
	u := c.appPath(appID)

	status, err := c.apiDelete(u, nil)
	if status == http.StatusNotFound {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

// GetAppByName gets the one and only app with the specified name.  If there
// are multiple applications with the name, an error is returned.  If no
// applications exist with the name, a nil app is returned.
func (c *Client) GetAppByName(name string) (app *App, err error) {
	apps, err := c.ListApps()
	if err != nil {
		return app, err
	}

	for i, a := range apps {
		if a.Name == name {
			if app != nil {
				err = fmt.Errorf("more than one Arpio app "+
					"exists with the name %q; use the Arpio web interface to "+
					"rename the unrelated apps, then retry",
					name)
				return nil, err
			}
			app = &apps[i]
		}
	}

	return app, nil
}

// MustGetAppByName finds the one and only app with the specified name.
// If timeout is > 0, the function tries to find a matching app
// until the timeout has elapsed.  An error is returned if no matching app
// could be found.
func (c *Client) MustGetAppByName(name string, timeout time.Duration) (app *App, err error) {
	const zeroDuration = time.Duration(0)
	for timeoutAt := time.Now().Add(timeout); timeout == zeroDuration || time.Now().Before(timeoutAt); {
		app, err = c.GetAppByName(name)
		if err != nil {
			return app, err
		}
		if app != nil || timeout == zeroDuration {
			break
		}
		log.Printf("[DEBUG] Waiting for a matching app to exist")
		time.Sleep(AppPollPeriod)
	}

	// If we didn't find an app, prepare an error
	if app == nil {
		err = fmt.Errorf("there is no Arpio application named %q", name)
	}

	return app, err
}

func (c *Client) appPath(appID string) string {
	return fmt.Sprintf(
		"/accounts/%s/applications/%s",
		c.AccountID,
		appID,
	)
}
