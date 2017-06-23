package main

type Profile map[string]interface{}

type (
	User struct {
		Id      int     `bson:"_id"`
		Balance float64 `bson:"Balance"`
		// Profile Profile `bson:"Profile"`
		/* Default for bool is false, so it's harder to determine that it's actually false
		Profile struct {
			TnProfit   float64 `bson:"tnprofit"`
			TcEnabled  bool    `bson:"tcenabled"`
			TnEnabled  bool    `bson:"tnenabled"`
			TbEnabled  bool    `bson:"tbenabled"`
			VwEnabled  bool    `bson:"vwenabled"`
			MgEnabled  bool    `bson:"mgenabled"`
			EroEnabled bool    `bson:"eroenabled"`
			WtEnabled  bool    `bson:"wtenabled"`
		} `bson:"Profile"`
		*/

		campaigns *CampaignList
	}
)

func NewDefaultUser(id int) *User {
	return &User{
		Id:      id,
		Balance: 0,
	}
}
