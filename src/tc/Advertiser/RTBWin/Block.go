package main

type BlockType int

const (
	BlockTypeTeaser     BlockType = 10 //Тизерная реклама
	BlockTypeBanner     BlockType = 20 //Баннерная реклама
	BlockTypeBKMessage  BlockType = 30 //ВК-message реклама
	BlockTypeClickUnder BlockType = 40 //Кликандер
	BlockTypeContext    BlockType = 50 //Контекстная реклама
)

type Block struct {
	Id     *int       `json:"id"`
	Type   *BlockType `json:"type"`
	Limit  *int       `json:"limit"`
	Width  *int       `json:"width"`
	Height *int       `json:"height"`
}
