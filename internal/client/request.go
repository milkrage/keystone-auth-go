package client

type Request struct {
	Auth Auth `json:"auth"`
}

type Auth struct {
	Identity Identity `json:"identity"`
	Scope    *Scope   `json:"scope,omitempty"`
}

type Identity struct {
	Methods  []string  `json:"methods"`
	Password *Password `json:"password,omitempty"`
}

type Scope struct {
	Domain  *Domain  `json:"domain,omitempty"`
	Project *Project `json:"project,omitempty"`
	System  *System  `json:"system,omitempty"`
}

type Password struct {
	User *User `json:"user,omitempty"`
}

type User struct {
	ID       string  `json:"id,omitempty"`
	Name     string  `json:"name,omitempty"`
	Password string  `json:"password,omitempty"`
	Domain   *Domain `json:"domain,omitempty"`
}

type Domain struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Project struct {
	ID     string  `json:"id,omitempty"`
	Name   string  `json:"name,omitempty"`
	Domain *Domain `json:"domain,omitempty"`
}

type System struct {
	All bool `json:"all,omitempty"`
}
