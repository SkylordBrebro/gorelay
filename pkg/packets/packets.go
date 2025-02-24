package packets

import (
	"gorelay/pkg/models"
)

// Packet is the interface that all packets must implement
type Packet interface {
	ID() int32
}

// WorldPosData represents a position in the game world
type WorldPosData struct {
	X float32
	Y float32
}

// SquareDistanceTo returns the square of the distance to another position
func (w *WorldPosData) SquareDistanceTo(other *WorldPosData) float32 {
	dx := w.X - other.X
	dy := w.Y - other.Y
	return dx*dx + dy*dy
}

// PlayerData contains player information
type PlayerData struct {
	Name         string
	ObjectID     int32
	Pos          *WorldPosData
	HP           int32
	MaxHP        int32
	MP           int32
	MaxMP        int32
	Level        int32
	Exp          int32
	NextLevelExp int32
	Fame         int32
	CurrentFame  int32
	Stars        int32
	AccountID    string
	GuildName    string
	GuildRank    int32
	Stats        map[string]int32
	Inventory    []int32
	HPPots       int32
	MPPots       int32
	HasBackpack  bool
}

// StatData represents a single stat value
type StatData struct {
	StatType    models.StatType
	StatValue   int32
	StringValue string
}

// UpdateData contains entity update information
type UpdateData struct {
	ObjectID int32
	Pos      *WorldPosData
	Stats    []StatData
}

// GroundTileData represents map tile information
type GroundTileData struct {
	X    int32
	Y    int32
	Type int32
}

// ObjectData contains information about game objects
type ObjectData struct {
	ObjectType int32
	Status     ObjectStatusData
}

// ObjectStatusData contains object status information
type ObjectStatusData struct {
	ObjectID int32
	Pos      *WorldPosData
	Stats    []StatData
}

// SlotObjectData represents an inventory slot
type SlotObjectData struct {
	ObjectID   int32
	SlotID     int32
	ObjectType int32
}

// TradeItem represents an item in a trade
type TradeItem struct {
	Item      int32
	SlotType  int32
	Tradeable bool
	Included  bool
}

// QuestData represents a quest
type QuestData struct {
	ID           string
	Name         string
	Description  string
	Expiration   int32
	Category     int32
	Requirements []int32
	Rewards      []int32
	Completed    bool
}

// MoveRecord represents a movement record
type MoveRecord struct {
	Time int32
	X    float32
	Y    float32
}

// FailureCode represents different types of connection failures
type FailureCode int32

const (
	IncorrectVersion FailureCode = iota
	BadKey
	InvalidTeleportTarget
	EmailVerificationNeeded
	InvalidCharacter
)

// Packet Definitions

// HelloPacket is sent after TCP connection to initiate communication
type HelloPacket struct {
	BuildVersion  string
	GameID        int32
	GUID          string
	Random1       int32
	Password      string
	Random2       int32
	Secret        string
	KeyTime       int32
	Key           []byte
	MapJSON       string
	EntryTag      string
	GameNet       string
	GameNetUserID string
	PlayPlatform  string
	PlatformToken string
	UserToken     string
}

func (p *HelloPacket) ID() int32 { return HELLO }

// FailurePacket indicates a connection or game error
type FailurePacket struct {
	ErrorID          FailureCode
	ErrorDescription string
}

func (p *FailurePacket) ID() int32 { return FAILURE }

// ClaimLoginRewardMsgPacket for claiming daily login rewards
type ClaimLoginRewardMsgPacket struct {
	ClaimKey string
	Type     string
}

func (p *ClaimLoginRewardMsgPacket) ID() int32 { return CLAIMDAILYREWARD }

// DeletePetPacket for deleting a pet
type DeletePetPacket struct {
	PetID int32
}

func (p *DeletePetPacket) ID() int32 { return DELETEPET }

// RequestTradePacket for initiating a trade
type RequestTradePacket struct {
	Name string
}

func (p *RequestTradePacket) ID() int32 { return REQUESTTRADE }

// QuestFetchResponsePacket for quest data
type QuestFetchResponsePacket struct {
	Quests []QuestData
}

func (p *QuestFetchResponsePacket) ID() int32 { return QUESTFETCHRESPONSE }

// JoinGuildPacket for joining a guild
type JoinGuildPacket struct {
	GuildName string
}

func (p *JoinGuildPacket) ID() int32 { return JOINGUILD }

// PingPacket for connection keepalive
type PingPacket struct {
	Serial int32
}

func (p *PingPacket) ID() int32 { return PING }

// NewTickPacket contains game tick information
type NewTickPacket struct {
	TickID                  int32
	TickTime                int32
	ServerRealTimeMS        int32
	ServerLastRealTimeRTTMS int32
	Updates                 []UpdateData
}

func (p *NewTickPacket) ID() int32 { return NEWTICK }

// PlayerTextPacket for chat messages
type PlayerTextPacket struct {
	Text string
}

func (p *PlayerTextPacket) ID() int32 { return PLAYERTEXT }

// UseItemPacket is sent to use an inventory item
type UseItemPacket struct {
	Time       int32
	SlotObject *SlotObjectData
	ItemUsePos *WorldPosData
	UseType    int32
}

func (p *UseItemPacket) ID() int32 { return USEITEM }

// ServerPlayerShootPacket represents another player shooting
type ServerPlayerShootPacket struct {
	BulletID      int32
	OwnerID       int32
	ContainerType int32
	StartingPos   *WorldPosData
	Angle         float32
	Damage        int32
}

func (p *ServerPlayerShootPacket) ID() int32 { return SERVERPLAYERSHOOT }

// ShowEffectPacket for visual effects
type ShowEffectPacket struct {
	EffectType     int32
	TargetObjectID int32
	Pos1           *WorldPosData
	Pos2           *WorldPosData
	Color          int32
	Duration       float32
}

func (p *ShowEffectPacket) ID() int32 { return SHOWEFFECT }

// TradeAcceptedPacket for accepting trades
type TradeAcceptedPacket struct {
	ClientOffer  []bool
	PartnerOffer []bool
}

func (p *TradeAcceptedPacket) ID() int32 { return TRADEACCEPTED }

// GuildRemovePacket for removing guild members
type GuildRemovePacket struct {
	Name string
}

func (p *GuildRemovePacket) ID() int32 { return GUILDREMOVE }

// PetUpgradeRequestPacket for upgrading pets
type PetUpgradeRequestPacket struct {
	PetTransType int32
	PetID        int32
	SlotObject   *SlotObjectData
}

func (p *PetUpgradeRequestPacket) ID() int32 { return PETUPGRADEREQUEST }

// GotoPacket indicates an object moving to a new position
type GotoPacket struct {
	ObjectID int32
	Position *WorldPosData
}

func (p *GotoPacket) ID() int32 { return GOTO }

// InvSwapPacket for inventory management
type InvSwapPacket struct {
	Time        int32
	Position    *WorldPosData
	SlotObject1 *SlotObjectData
	SlotObject2 *SlotObjectData
}

func (p *InvSwapPacket) ID() int32 { return INVENTORYSWAP }

// OtherHitPacket for hit notifications
type OtherHitPacket struct {
	Time     int32
	BulletID int32
	ObjectID int32
	TargetID int32
}

func (p *OtherHitPacket) ID() int32 { return OTHERHIT }

// NameResultPacket for name change results
type NameResultPacket struct {
	Success   bool
	ErrorText string
}

func (p *NameResultPacket) ID() int32 { return NAMERESULT }

// BuyResultPacket for purchase results
type BuyResultPacket struct {
	Result       int32
	ResultString string
}

func (p *BuyResultPacket) ID() int32 { return BUYRESULT }

// HatchPetPacket for hatching pets
type HatchPetPacket struct {
	PetName string
	PetSkin int32
}

func (p *HatchPetPacket) ID() int32 { return HATCHPET }

// ActivePetUpdateRequestPacket for updating active pets
type ActivePetUpdateRequestPacket struct {
	CommandType int32
	InstanceID  int32
}

func (p *ActivePetUpdateRequestPacket) ID() int32 { return ACTIVEPETUPDATEREQUEST }

// EnemyHitPacket for enemy hit notifications
type EnemyHitPacket struct {
	Time     int32
	BulletID int32
	TargetID int32
	Kill     bool
}

func (p *EnemyHitPacket) ID() int32 { return ENEMYHIT }

// GuildResultPacket for guild operation results
type GuildResultPacket struct {
	Success         bool
	LineBuilderJSON string
}

func (p *GuildResultPacket) ID() int32 { return GUILDRESULT }

// EditAccountListPacket for managing account lists
type EditAccountListPacket struct {
	AccountListID int32
	Add           bool
	ObjectID      int32
}

func (p *EditAccountListPacket) ID() int32 { return EDITACCOUNTLIST }

// TradeChangedPacket for trade updates
type TradeChangedPacket struct {
	Offer []bool
}

func (p *TradeChangedPacket) ID() int32 { return TRADECHANGED }

// PlayerShootPacket is sent when the player shoots
type PlayerShootPacket struct {
	Time          int32
	BulletID      int32
	ContainerType int32
	StartingPos   *WorldPosData
	Angle         float32
}

func (p *PlayerShootPacket) ID() int32 { return PLAYERSHOOT }

// PongPacket for responding to ping
type PongPacket struct {
	Serial int32
	Time   int32
}

func (p *PongPacket) ID() int32 { return PONG }

// PetChangeSkinMsgPacket for changing pet skins
type PetChangeSkinMsgPacket struct {
	PetID    int32
	SkinType int32
}

func (p *PetChangeSkinMsgPacket) ID() int32 { return CHANGEPETSKIN }

// TradeDonePacket for completed trades
type TradeDonePacket struct {
	Code        int32
	Description string
}

func (p *TradeDonePacket) ID() int32 { return TRADEDONE }

// EnemyShootPacket represents an enemy shooting projectiles
type EnemyShootPacket struct {
	BulletID    int32
	OwnerID     int32
	BulletType  int32
	StartingPos *WorldPosData
	Angle       float32
	Damage      int32
	NumShots    int32
	AngleInc    float32
}

func (p *EnemyShootPacket) ID() int32 { return ENEMYSHOOT }

// AcceptTradePacket for accepting trades
type AcceptTradePacket struct {
	ClientOffer  []bool
	PartnerOffer []bool
}

func (p *AcceptTradePacket) ID() int32 { return ACCEPTTRADE }

// ChangeGuildRankPacket for changing guild ranks
type ChangeGuildRankPacket struct {
	Name      string
	GuildRank int32
}

func (p *ChangeGuildRankPacket) ID() int32 { return CHANGEGUILDRANK }

// PlaySoundPacket for playing sound effects
type PlaySoundPacket struct {
	OwnerID int32
	SoundID int32
}

func (p *PlaySoundPacket) ID() int32 { return PLAYSOUND }

// SquareHitPacket for square hit notifications
type SquareHitPacket struct {
	Time     int32
	BulletID int32
	ObjectID int32
}

func (p *SquareHitPacket) ID() int32 { return SQUAREHIT }

// NewAbilityPacket for new ability unlocks
type NewAbilityPacket struct {
	Type int32
}

func (p *NewAbilityPacket) ID() int32 { return NEWABILITY }

// MovePacket relays player position to server
type MovePacket struct {
	TickID      int32
	Time        int32
	NewPosition *WorldPosData
	Records     []MoveRecord
}

func (p *MovePacket) ID() int32 { return MOVE }

// TextPacket represents a chat message
type TextPacket struct {
	Name       string
	ObjectID   int32
	NumStars   int32
	BubbleTime int32
	Recipient  string
	Text       string
	CleanText  string
}

func (p *TextPacket) ID() int32 { return TEXT }

// ReconnectPacket for server reconnection
type ReconnectPacket struct {
	Name        string
	Host        string
	Port        int32
	GameID      int32
	KeyTime     int32
	Key         []byte
	IsFromArena bool
}

func (p *ReconnectPacket) ID() int32 { return RECONNECT }

// DeathPacket for player deaths
type DeathPacket struct {
	AccountID  string
	CharID     int32
	KilledBy   string
	ZombieID   int32
	ZombieType int32
	IsZombie   bool
}

func (p *DeathPacket) ID() int32 { return DEATH }

// UsePortalPacket for using portals
type UsePortalPacket struct {
	ObjectID int32
}

func (p *UsePortalPacket) ID() int32 { return USEPORTAL }

// QuestRoomMsgPacket for quest room messages
type QuestRoomMsgPacket struct {
	Message string
}

func (p *QuestRoomMsgPacket) ID() int32 { return GOTOQUESTROOM }

// AllyShootPacket for ally shooting
type AllyShootPacket struct {
	BulletID      int32
	OwnerID       int32
	ContainerType int32
	Angle         float32
}

func (p *AllyShootPacket) ID() int32 { return ALLYSHOOT }

// ReskinPacket for character reskinning
type ReskinPacket struct {
	SkinID int32
}

func (p *ReskinPacket) ID() int32 { return RESKIN }

// ResetDailyQuestsPacket for resetting daily quests
type ResetDailyQuestsPacket struct{}

func (p *ResetDailyQuestsPacket) ID() int32 { return RESETDAILYQUESTS }

// InvDropPacket for dropping inventory items
type InvDropPacket struct {
	SlotObject *SlotObjectData
}

func (p *InvDropPacket) ID() int32 { return INVENTORYDROP }

// LoadPacket for loading into maps
type LoadPacket struct {
	CharID      int32
	IsFromArena bool
}

func (p *LoadPacket) ID() int32 { return LOAD }

// CreateGuildPacket for creating guilds
type CreateGuildPacket struct {
	Name string
}

func (p *CreateGuildPacket) ID() int32 { return CREATEGUILD }

// CreatePacket for creating new characters
type CreatePacket struct {
	ClassType int32
	SkinType  int32
}

func (p *CreatePacket) ID() int32 { return CREATE }

// UpdatePacket contains entity updates
type UpdatePacket struct {
	Tiles      []GroundTileData
	NewObjects []ObjectData
	Drops      []int32
}

func (p *UpdatePacket) ID() int32 { return UPDATE }

// NotificationPacket for player notifications
type NotificationPacket struct {
	ObjectID int32
	Message  string
	Color    int32
}

func (p *NotificationPacket) ID() int32 { return NOTIFICATION }

// AoePacket represents an area of effect attack
type AoePacket struct {
	Pos           *WorldPosData
	Radius        float32
	Damage        int32
	Effect        int32
	Duration      float32
	OrigType      int32
	Color         int32
	ArmorPiercing bool
}

func (p *AoePacket) ID() int32 { return AOE }

// GotoAckPacket acknowledges a GotoPacket
type GotoAckPacket struct {
	Time int32
}

func (p *GotoAckPacket) ID() int32 { return GOTOACK }

// ClientStatPacket for client statistics
type ClientStatPacket struct {
	Name  string
	Value int32
}

func (p *ClientStatPacket) ID() int32 { return CLIENTSTAT }

// TeleportPacket for player teleportation
type TeleportPacket struct {
	ObjectID int32
}

func (p *TeleportPacket) ID() int32 { return TELEPORT }

// DamagePacket for damage notifications
type DamagePacket struct {
	TargetID     int32
	Effects      []int32
	DamageAmount int32
	Kill         bool
	BulletID     int32
	ObjectID     int32
}

func (p *DamagePacket) ID() int32 { return DAMAGE }

// ActivePetUpdatePacket for active pet updates
type ActivePetUpdatePacket struct {
	PetID int32
}

func (p *ActivePetUpdatePacket) ID() int32 { return ACTIVEPET }

// InvitedToGuildPacket for guild invites
type InvitedToGuildPacket struct {
	Name      string
	GuildName string
}

func (p *InvitedToGuildPacket) ID() int32 { return INVITEDTOGUILD }

// PetYardUpdatePacket for pet yard updates
type PetYardUpdatePacket struct {
	Type int32
}

func (p *PetYardUpdatePacket) ID() int32 { return PETYARDUPDATE }

// PasswordPromptPacket for password prompts
type PasswordPromptPacket struct {
	CleanPasswordStatus int32
}

func (p *PasswordPromptPacket) ID() int32 { return PASSWORDPROMPT }

// UpdateAckPacket acknowledges an UpdatePacket
type UpdateAckPacket struct{}

func (p *UpdateAckPacket) ID() int32 { return UPDATEACK }

// QuestObjIDPacket for quest object IDs
type QuestObjIDPacket struct {
	ObjectID int32
}

func (p *QuestObjIDPacket) ID() int32 { return QUESTOBJECTID }

// PICPacket for account security
type PICPacket struct {
	PICData string
}

func (p *PICPacket) ID() int32 { return PIC }

// MapInfoPacket contains information about the current map
type MapInfoPacket struct {
	Width               int32
	Height              int32
	Name                string
	Seed                int32
	Background          int32
	Difficulty          int32
	AllowPlayerTeleport bool
	ShowDisplays        bool
	ClientXML           []byte
	ExtraXML            []byte
}

func (p *MapInfoPacket) ID() int32 { return MAPINFO }

// KeyInfoRequestPacket for key info requests
type KeyInfoRequestPacket struct {
	ItemType int32
}

func (p *KeyInfoRequestPacket) ID() int32 { return KEYINFOREQUEST }

// InvResultPacket for inventory operation results
type InvResultPacket struct {
	Result int32
}

func (p *InvResultPacket) ID() int32 { return INVENTORYRESULT }

// QuestRedeemResponsePacket for quest redemption responses
type QuestRedeemResponsePacket struct {
	OK      bool
	Message string
}

func (p *QuestRedeemResponsePacket) ID() int32 { return QUESTREDEEMRESPONSE }

// ChooseNamePacket for choosing character names
type ChooseNamePacket struct {
	Name string
}

func (p *ChooseNamePacket) ID() int32 { return CHOOSENAME }

// QuestFetchAskPacket for fetching quests
type QuestFetchAskPacket struct{}

func (p *QuestFetchAskPacket) ID() int32 { return QUESTFETCHASK }

// AccountListPacket for account lists
type AccountListPacket struct {
	AccountListID int32
	AccountIDs    []string
	LockAction    int32
}

func (p *AccountListPacket) ID() int32 { return ACCOUNTLIST }

// ShootAckPacket for shoot acknowledgments
type ShootAckPacket struct {
	Time int32
}

func (p *ShootAckPacket) ID() int32 { return SHOOTACKCOUNTER }

// CreateSuccessPacket for successful character creation
type CreateSuccessPacket struct {
	ObjectID int32
	CharID   int32
}

func (p *CreateSuccessPacket) ID() int32 { return CREATESUCCESS }

// CheckCreditsPacket for checking credits
type CheckCreditsPacket struct{}

func (p *CheckCreditsPacket) ID() int32 { return CHECKCREDITS }

// GroundDamagePacket for ground damage
type GroundDamagePacket struct {
	Time     int32
	Position *WorldPosData
}

func (p *GroundDamagePacket) ID() int32 { return GROUNDDAMAGE }

// GuildInvitePacket for guild invites
type GuildInvitePacket struct {
	Name string
}

func (p *GuildInvitePacket) ID() int32 { return GUILDINVITE }

// EscapePacket for escaping to nexus
type EscapePacket struct{}

func (p *EscapePacket) ID() int32 { return ESCAPE }

// FilePacket for file transfers
type FilePacket struct {
	Name  string
	Bytes []byte
}

func (p *FilePacket) ID() int32 { return FILE }

// ReskinUnlockPacket for unlocking skins
type ReskinUnlockPacket struct {
	SkinID int32
}

func (p *ReskinUnlockPacket) ID() int32 { return UNLOCKCUSTOMIZATION }

// NewCharacterInfoPacket for new character information
type NewCharacterInfoPacket struct {
	CharXML string
}

func (p *NewCharacterInfoPacket) ID() int32 { return NEWCHARACTERINFORMATION }

// UnlockInfoPacket for unlock information
type UnlockInfoPacket struct {
	UnlockType int32
	UnlockID   int32
}

func (p *UnlockInfoPacket) ID() int32 { return UNLOCKNEWSLOT }

// QueueInfoPacket for queue information
type QueueInfoPacket struct {
	Position    int32
	Count       int32
	CurrentTime int32
}

func (p *QueueInfoPacket) ID() int32 { return QUEUE }

// ExaltationUpdatePacket for exaltation updates
type ExaltationUpdatePacket struct {
	ObjType int32
	Stars   int32
}

func (p *ExaltationUpdatePacket) ID() int32 { return EXALTATIONBONUSCHANGED }

// VaultInfoPacket for vault information
type VaultInfoPacket struct {
	VaultContents  []int32
	GiftContents   []int32
	PotionContents []int32
}

func (p *VaultInfoPacket) ID() int32 { return VAULTCONTENT }

// ForgeRequestPacket for forge requests
type ForgeRequestPacket struct {
	SlotObject1 *SlotObjectData
	SlotObject2 *SlotObjectData
}

func (p *ForgeRequestPacket) ID() int32 { return FORGEREQUEST }

// ForgeResponsePacket for forge responses
type ForgeResponsePacket struct {
	Success bool
	Message string
}

func (p *ForgeResponsePacket) ID() int32 { return FORGERESULT }

// ShowAllyShootPacket for showing ally shots
type ShowAllyShootPacket struct {
	BulletID int32
	OwnerID  int32
	Angle    float32
}

func (p *ShowAllyShootPacket) ID() int32 { return CHANGEALLYSHOOT }

// ChangeTradePacket for modifying trade offers
type ChangeTradePacket struct {
	OfferedItems []bool
}

func (p *ChangeTradePacket) ID() int32 { return CHANGETRADE }

// QuestRedeemPacket for redeeming quests
type QuestRedeemPacket struct {
	QuestID string
}

func (p *QuestRedeemPacket) ID() int32 { return QUESTREDEEM }

// SetConditionPacket for setting entity conditions
type SetConditionPacket struct {
	ConditionEffect   int32
	ConditionDuration float32
}

func (p *SetConditionPacket) ID() int32 { return SETCONDITION }

// KeyInfoResponsePacket for key information responses
type KeyInfoResponsePacket struct {
	Name        string
	Description string
	Creator     string
}

func (p *KeyInfoResponsePacket) ID() int32 { return KEYINFORESPONSE }

// RealmHeroesResponsePacket for realm heroes information
type RealmHeroesResponsePacket struct {
	NumberOfRealmHeroes int32
	Heroes              []string
}

func (p *RealmHeroesResponsePacket) ID() int32 { return HEROLEFT }

// BuyPacket for purchasing items
type BuyPacket struct {
	ObjectID int32
	Quantity int32
}

func (p *BuyPacket) ID() int32 { return BUY }

// TradeStartPacket for initiating trades
type TradeStartPacket struct {
	ClientItems []TradeItem
	PartnerName string
}

func (p *TradeStartPacket) ID() int32 { return TRADESTART }

// EvolvePetPacket for evolving pets
type EvolvePetPacket struct {
	PetID       int32
	InitialSkin int32
	FinalSkin   int32
}

func (p *EvolvePetPacket) ID() int32 { return EVOLVEDPET }

// TradeRequestedPacket for trade requests
type TradeRequestedPacket struct {
	Name string
}

func (p *TradeRequestedPacket) ID() int32 { return TRADEREQUESTED }

// AoeAckPacket for acknowledging area of effect
type AoeAckPacket struct {
	Time     int32
	Position *WorldPosData
}

func (p *AoeAckPacket) ID() int32 { return AOEACK }

// PlayerHitPacket for player hit notifications
type PlayerHitPacket struct {
	BulletID int32
	ObjectID int32
}

func (p *PlayerHitPacket) ID() int32 { return PLAYERHIT }

// CancelTradePacket for canceling trades
type CancelTradePacket struct{}

func (p *CancelTradePacket) ID() int32 { return CANCELTRADE }

// Encode serializes a packet for network transmission
func Encode(packet Packet) ([]byte, error) {
	// TODO: Implement proper packet encoding based on RotMG protocol
	return nil, nil
}
