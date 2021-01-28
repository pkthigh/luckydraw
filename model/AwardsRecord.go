package model


type AwardsRecord struct {
	RID int
	*User
	*Awards
}