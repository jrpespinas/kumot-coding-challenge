package domain

type User struct {
	Name        *string `json:"name"`
	Login       *string `json:"login"`
	Company     *string `json:"company"`
	Followers   *int    `json:"followers"`
	PublicRepos *int    `json:"public_repos"`
}
