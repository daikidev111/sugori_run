package entity

// domain では技術的ロジックは何も入れない
type User struct {
	ID        string
	AuthToken string
	Name      string
	HighScore int32
	Coin      int32
}

type Users *[]User
