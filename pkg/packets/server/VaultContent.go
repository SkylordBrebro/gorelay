package server

import (
	"gorelay/pkg/packets/interfaces"
)

// VaultContent represents the server packet for vault contents
type VaultContent struct {
	LastVaultUpdate            bool
	VaultChestObjectId         int32
	MaterialChestObjectID      int32
	GiftChestObjectId          int32
	PotionStorageObjectId      int32
	SeasonalSpoilChestObjectId int32
	VaultContents              []int32
	MaterialContents           []int32
	GiftContents               []int32
	PotionContents             []int32
	SeasonalSpoilContent       []int32
	VaultUpgradeCost           int16
	MaterialUpgradeCost        int16
	PotionUpgradeCost          int16
	CurrentPotionMax           int16
	NextPotionMax              int16
	VaultChestEnchants         string
	GiftChestEnchants          string
	SpoilsChestEnchants        string
}

// Type returns the packet type for VaultContent
func (p *VaultContent) Type() interfaces.PacketType {
	return interfaces.VaultContent
}

// Read reads the packet data from the provided reader
func (p *VaultContent) Read(r interfaces.Reader) error {
	var err error

	// Read boolean and object IDs
	p.LastVaultUpdate, err = r.ReadBool()
	if err != nil {
		return err
	}

	vaultChestId, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}
	p.VaultChestObjectId = int32(vaultChestId)

	materialChestId, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}
	p.MaterialChestObjectID = int32(materialChestId)

	giftChestId, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}
	p.GiftChestObjectId = int32(giftChestId)

	potionStorageId, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}
	p.PotionStorageObjectId = int32(potionStorageId)

	seasonalSpoilId, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}
	p.SeasonalSpoilChestObjectId = int32(seasonalSpoilId)

	// Read vault contents
	vaultCount, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}
	p.VaultContents = make([]int32, vaultCount)
	for i := 0; i < vaultCount; i++ {
		itemId, err := r.ReadCompressedInt()
		if err != nil {
			return err
		}
		p.VaultContents[i] = int32(itemId)
	}

	// Read material contents
	materialCount, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}
	p.MaterialContents = make([]int32, materialCount)
	for i := 0; i < materialCount; i++ {
		itemId, err := r.ReadCompressedInt()
		if err != nil {
			return err
		}
		p.MaterialContents[i] = int32(itemId)
	}

	// Read gift contents
	giftCount, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}
	p.GiftContents = make([]int32, giftCount)
	for i := 0; i < giftCount; i++ {
		itemId, err := r.ReadCompressedInt()
		if err != nil {
			return err
		}
		p.GiftContents[i] = int32(itemId)
	}

	// Read potion contents
	potionCount, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}
	p.PotionContents = make([]int32, potionCount)
	for i := 0; i < potionCount; i++ {
		itemId, err := r.ReadCompressedInt()
		if err != nil {
			return err
		}
		p.PotionContents[i] = int32(itemId)
	}

	// Read seasonal spoil content
	spoilCount, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}
	p.SeasonalSpoilContent = make([]int32, spoilCount)
	for i := 0; i < spoilCount; i++ {
		itemId, err := r.ReadCompressedInt()
		if err != nil {
			return err
		}
		p.SeasonalSpoilContent[i] = int32(itemId)
	}

	// Read upgrade costs and potion max values
	p.VaultUpgradeCost, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.MaterialUpgradeCost, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.PotionUpgradeCost, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.CurrentPotionMax, err = r.ReadInt16()
	if err != nil {
		return err
	}

	p.NextPotionMax, err = r.ReadInt16()
	if err != nil {
		return err
	}

	// Read enchant strings
	p.VaultChestEnchants, err = r.ReadString()
	if err != nil {
		return err
	}

	p.GiftChestEnchants, err = r.ReadString()
	if err != nil {
		return err
	}

	p.SpoilsChestEnchants, err = r.ReadString()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *VaultContent) Write(w interfaces.Writer) error {
	// Write boolean and object IDs
	if err := w.WriteBool(p.LastVaultUpdate); err != nil {
		return err
	}

	if err := w.WriteCompressedInt(int(p.VaultChestObjectId)); err != nil {
		return err
	}

	if err := w.WriteCompressedInt(int(p.MaterialChestObjectID)); err != nil {
		return err
	}

	if err := w.WriteCompressedInt(int(p.GiftChestObjectId)); err != nil {
		return err
	}

	if err := w.WriteCompressedInt(int(p.PotionStorageObjectId)); err != nil {
		return err
	}

	if err := w.WriteCompressedInt(int(p.SeasonalSpoilChestObjectId)); err != nil {
		return err
	}

	// Write vault contents
	if err := w.WriteCompressedInt(len(p.VaultContents)); err != nil {
		return err
	}
	for _, item := range p.VaultContents {
		if err := w.WriteCompressedInt(int(item)); err != nil {
			return err
		}
	}

	// Write material contents
	if err := w.WriteCompressedInt(len(p.MaterialContents)); err != nil {
		return err
	}
	for _, item := range p.MaterialContents {
		if err := w.WriteCompressedInt(int(item)); err != nil {
			return err
		}
	}

	// Write gift contents
	if err := w.WriteCompressedInt(len(p.GiftContents)); err != nil {
		return err
	}
	for _, item := range p.GiftContents {
		if err := w.WriteCompressedInt(int(item)); err != nil {
			return err
		}
	}

	// Write potion contents
	if err := w.WriteCompressedInt(len(p.PotionContents)); err != nil {
		return err
	}
	for _, item := range p.PotionContents {
		if err := w.WriteCompressedInt(int(item)); err != nil {
			return err
		}
	}

	// Write seasonal spoil content
	if err := w.WriteCompressedInt(len(p.SeasonalSpoilContent)); err != nil {
		return err
	}
	for _, item := range p.SeasonalSpoilContent {
		if err := w.WriteCompressedInt(int(item)); err != nil {
			return err
		}
	}

	// Write upgrade costs and potion max values
	if err := w.WriteInt16(p.VaultUpgradeCost); err != nil {
		return err
	}

	if err := w.WriteInt16(p.MaterialUpgradeCost); err != nil {
		return err
	}

	if err := w.WriteInt16(p.PotionUpgradeCost); err != nil {
		return err
	}

	if err := w.WriteInt16(p.CurrentPotionMax); err != nil {
		return err
	}

	if err := w.WriteInt16(p.NextPotionMax); err != nil {
		return err
	}

	// Write enchant strings
	if err := w.WriteString(p.VaultChestEnchants); err != nil {
		return err
	}

	if err := w.WriteString(p.GiftChestEnchants); err != nil {
		return err
	}

	if err := w.WriteString(p.SpoilsChestEnchants); err != nil {
		return err
	}

	return nil
}

func (p *VaultContent) ID() int32 {
	return int32(interfaces.VaultContent)
}