package query

type GetEmailsByUserIDs struct {
	IDs []uint `form:"ids"`
}
