package request

type CreateBid struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	TenderId    string `json:"tenderId" validate:"required"`
	AuthorType  string `json:"authorType" validate:"required,oneof=Organization User"`
	AuthorId    string `json:"authorId" validate:"required"`
}
type GetBidsByUser struct {
	Limit    int32  `query:"limit" validate:"gte=0"`
	Offset   int32  `query:"offset" validate:"gte=0"`
	Username string `query:"username" validate:"required,max=50"`
}
type BidsForTender struct {
	TenderId string `param:"tenderId" validate:"required,uuid4"`
	Limit    int32  `query:"limit" validate:"gte=0"`
	Offset   int32  `query:"offset" validate:"gte=0"`
	Username string `query:"username" validate:"required,max=50"`
}
type BidStatus struct {
	BidId    string `param:"bidId" validate:"required,uuid4"`
	Username string `query:"username" validate:"required,max=50"`
}
type UpdateBidStatus struct {
	BidId    string `param:"bidId" validate:"required,uuid4"`
	Status   string `query:"status" validate:"required,oneof=Created Published Canceled"`
	Username string `query:"username" validate:"required,max=50"`
}
type EditBid struct {
	BidId       string `param:"bidId" validate:"required,uuid4"`
	Username    string `query:"username" validate:"required,max=50"`
	Name        string `json:"name" validate:"max=100"`
	Description string `json:"description" validate:"max=500"`
}
type SubmitDecision struct {
	BidId    string `param:"bidId" validate:"required,uuid4"`
	Decision string `query:"decision" validate:"required,oneof=Approved Rejected"`
	Username string `query:"username" validate:"required,max=50"`
}
type Feedback struct {
	BidId       string `param:"bidId" validate:"required,uuid4"`
	BidFeedback string `query:"bidFeedback" validate:"required,max=1000"`
	Username    string `query:"username" validate:"required,max=50"`
}
type RollbackBid struct {
	BidId    string `param:"bidId" validate:"required,uuid4"`
	Version  int32  `param:"version" validate:"required,gt=0"`
	Username string `query:"username" validate:"required,max=50"`
}
type Reviews struct {
	TenderId          string `param:"tenderId" validate:"required,uuid4" `
	AuthorUsername    string `query:"authorUsername" validate:"required,max=50"`
	RequesterUsername string `query:"requesterUsername" validate:"required,max=50"`
	Limit             int32  `query:"limit" validate:"gte=0"`
	Offset            int32  `query:"offset" validate:"gte=0"`
}
