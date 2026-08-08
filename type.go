package keystoneauth

import "github.com/milkrage/keystone-auth-go/internal/client"

type IdentityMethod string

const (
	MethodPassword IdentityMethod = "password"
)

type Credentials struct {
	Method IdentityMethod
	User   User
	Scope  Scope
}

func (c Credentials) toRequest() client.Request {
	return client.Request{
		Auth: client.Auth{
			Identity: client.Identity{
				Methods:  []string{string(c.Method)},
				Password: &client.Password{User: ptrIfNonEmpty(c.User.toRequest())},
			},
			Scope: ptrIfNonEmpty(c.Scope.toRequest()),
		},
	}
}

type User struct {
	ID       string
	Name     string
	Password string
	Domain   Domain
}

func (u User) toRequest() client.User {
	return client.User{
		ID:       u.ID,
		Name:     u.Name,
		Password: u.Password,
		Domain:   ptrIfNonEmpty(u.Domain.toRequest()),
	}
}

type Scope struct {
	Domain  Domain
	Project Project
	System  System
}

func (s Scope) toRequest() client.Scope {
	return client.Scope{
		Domain:  ptrIfNonEmpty(s.Domain.toRequest()),
		Project: ptrIfNonEmpty(s.Project.toRequest()),
		System:  ptrIfNonEmpty(s.System.toRequest()),
	}
}

type Domain struct {
	ID   string
	Name string
}

func (d Domain) toRequest() client.Domain {
	return client.Domain{
		ID:   d.ID,
		Name: d.Name,
	}
}

type Project struct {
	ID     string
	Name   string
	Domain Domain
}

func (p Project) toRequest() client.Project {
	return client.Project{
		ID:     p.ID,
		Name:   p.Name,
		Domain: ptrIfNonEmpty(p.Domain.toRequest()),
	}
}

type System struct {
	All bool
}

func (s System) toRequest() client.System {
	return client.System{
		All: s.All,
	}
}

func ptrIfNonEmpty[T comparable](v T) *T {
	var empty T

	if v == empty {
		return nil
	}

	return &v
}
