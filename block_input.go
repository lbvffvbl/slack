package slack

type InputBlock struct {
	Type    MessageBlockType `json:"type"`
	BlockID string           `json:"block_id,omitempty"`
	Label   *TextBlockObject `json:"label"`
	Element *Accessory       `json:"element"`
}

func (s InputBlock) BlockType() MessageBlockType {
	return s.Type
}

func NewInputBlock(element *Accessory, label *TextBlockObject, blockId string) *InputBlock {
	return &InputBlock{
		Type:    MBTInput,
		Element: element,
		Label:   label,
		BlockID: blockId,
	}
}
