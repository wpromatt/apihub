package account_test

import (
	"github.com/apihub/apihub/account"
	"github.com/apihub/apihub/errors"
	. "gopkg.in/check.v1"
)

func (s *S) TestCreateApp(c *C) {
	err := app.Create(owner, team)
	c.Assert(err, IsNil)
	defer app.Delete(owner)
}

func (s *S) TestCreateAppWithDuplicateCliendId(c *C) {
	err := app.Create(owner, team)
	c.Check(err, IsNil)

	err = app.Create(owner, team)
	_, ok := err.(errors.ValidationError)
	c.Assert(ok, Equals, true)
	defer app.Delete(owner)
}

func (s *S) TestCreateAppWithoutRequiredFields(c *C) {
	app = account.App{}
	err := app.Create(owner, team)
	_, ok := err.(errors.ValidationError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestUpdateApp(c *C) {
	app.Name = "ApiHub App"
	err := app.Create(owner, team)
	c.Assert(err, IsNil)
	c.Assert(app.Name, Equals, "ApiHub App")

	app.Name = "Another name"
	err = app.Update()
	c.Assert(err, IsNil)
	c.Assert(app.Name, Equals, "Another name")
	defer app.Delete(owner)
}

func (s *S) TestUpdateAppWithoutRequiredFields(c *C) {
	app.Name = "ApiHub App"
	err := app.Create(owner, team)
	c.Assert(err, IsNil)
	c.Assert(app.Name, Equals, "ApiHub App")
	defer app.Delete(owner)

	app.Name = ""
	err = app.Update()
	_, ok := err.(errors.ValidationError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestAppExists(c *C) {
	app.Create(owner, team)
	c.Assert(app.Exists(), Equals, true)
	defer app.Delete(owner)
}

func (s *S) TestAppExistsNotFound(c *C) {
	app = account.App{ClientId: "not_found"}
	c.Assert(app.Exists(), Equals, false)
}

func (s *S) TestDeleteApp(c *C) {
	app.Create(owner, team)
	c.Assert(app.Exists(), Equals, true)
	app.Delete(owner)
	c.Assert(app.Exists(), Equals, false)
}

func (s *S) TestDeleteAppNotOwner(c *C) {
	app.Create(alice, team)
	c.Assert(app.Exists(), Equals, true)
	defer app.Delete(alice)

	err := app.Delete(owner)
	_, ok := err.(errors.ForbiddenError)
	c.Assert(ok, Equals, true)
}

func (s *S) TestDeleteAppsByTeam(c *C) {
	app.Create(alice, team)
	c.Assert(app.Exists(), Equals, true)

	newApp := account.App{ClientId: "android", ClientSecret: "secret", Name: "Andoird App", Team: team.Alias, Owner: owner.Email, RedirectUris: []string{"http://www.example.org/auth"}}
	newApp.Create(alice, team)

	account.DeleteAppsByTeam(team, alice)

	_, err := account.FindAppByClientId("ios")
	c.Check(err, Not(IsNil))
	_, err = account.FindAppByClientId("android")
	c.Check(err, Not(IsNil))
}

func (s *S) TestFindAppByClientId(c *C) {
	err := app.Create(owner, team)

	a, err := account.FindAppByClientId(app.ClientId)
	c.Check(a, Not(IsNil))
	c.Check(err, IsNil)
	defer app.Delete(owner)
}

func (s *S) TestFindAppByClientIdNotFound(c *C) {
	a, err := account.FindAppByClientId("not-found")
	c.Check(a, IsNil)
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}
