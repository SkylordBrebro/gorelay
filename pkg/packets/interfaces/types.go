package interfaces

// PacketType represents different types of network packets
type PacketType int

const (
	Unknown                              PacketType = -1
	Failure                              PacketType = 0
	Teleport                             PacketType = 1
	ClaimDailyReward                     PacketType = 3
	DeletePet                            PacketType = 4
	RequestTrade                         PacketType = 5
	QuestFetchResponse                   PacketType = 6
	JoinGuild                            PacketType = 7
	Ping                                 PacketType = 8
	PlayerText                           PacketType = 9
	NewTick                              PacketType = 10
	ShowEffect                           PacketType = 11
	ServerPlayerShoot                    PacketType = 12
	UseItem                              PacketType = 13
	TradeAccepted                        PacketType = 14
	GuildRemove                          PacketType = 15
	PetUpgradeRequest                    PacketType = 16
	NameResult                           PacketType = 21
	BuyResult                            PacketType = 22
	Goto                                 PacketType = 18
	InventoryDrop                        PacketType = 19
	OtherHit                             PacketType = 20
	HatchPet                             PacketType = 23
	ActivePetUpdateRequest               PacketType = 24
	EnemyHit                             PacketType = 25
	GuildResult                          PacketType = 26
	EditAccountList                      PacketType = 27
	TradeChanged                         PacketType = 28
	PlayerShoot                          PacketType = 30
	Pong                                 PacketType = 31
	ChangePetSkin                        PacketType = 33
	TradeDone                            PacketType = 34
	EnemyShoot                           PacketType = 35
	AcceptTrade                          PacketType = 36
	ChangeGuildRank                      PacketType = 37
	PlaySound                            PacketType = 38
	SquareHit                            PacketType = 40
	NewAbility                           PacketType = 41
	Update                               PacketType = 42
	Text                                 PacketType = 44
	Reconnect                            PacketType = 45
	Death                                PacketType = 46
	UsePortal                            PacketType = 47
	GoToQuestRoom                        PacketType = 48
	AllyShoot                            PacketType = 49
	Reskin                               PacketType = 51
	ResetDailyQuests                     PacketType = 52
	InventorySwap                        PacketType = 55
	ChangeTrade                          PacketType = 56
	Create                               PacketType = 57
	QuestRedeem                          PacketType = 58
	CreateGuild                          PacketType = 59
	SetCondition                         PacketType = 60
	Load                                 PacketType = 61
	Move                                 PacketType = 62
	KeyInfoResponse                      PacketType = 63
	AOE                                  PacketType = 64
	GotoAck                              PacketType = 65
	Notification                         PacketType = 67
	ClientStat                           PacketType = 69
	Hello                                PacketType = 74
	Damage                               PacketType = 75
	ActivePet                            PacketType = 76
	InvitedToGuild                       PacketType = 77
	PetYardUpdate                        PacketType = 78
	PasswordPrompt                       PacketType = 79
	UpdateAck                            PacketType = 81
	QuestObjectId                        PacketType = 82
	Pic                                  PacketType = 83
	HeroLeft                             PacketType = 84
	Buy                                  PacketType = 85
	TradeStart                           PacketType = 86
	EvolvedPet                           PacketType = 87
	TradeRequested                       PacketType = 88
	AOEAck                               PacketType = 89
	PlayerHit                            PacketType = 90
	CancelTrade                          PacketType = 91
	MapInfo                              PacketType = 92
	KeyInfoRequest                       PacketType = 94
	InventoryResult                      PacketType = 95
	QuestRedeemResponse                  PacketType = 96
	ChooseName                           PacketType = 97
	QuestFetchAsk                        PacketType = 98
	AccountList                          PacketType = 99
	CreateSuccess                        PacketType = 101
	CheckCredits                         PacketType = 102
	GroundDamage                         PacketType = 103
	GuildInvite                          PacketType = 104
	Escape                               PacketType = 105
	File                                 PacketType = 106
	UnlockCustomization                  PacketType = 107
	NewCharacterInformation              PacketType = 108
	UnlockNewSlot                        PacketType = 109
	Queue                                PacketType = 112
	QueueCancel                          PacketType = 113
	ExaltationBonusChanged               PacketType = 114
	RedeemExaltationReward               PacketType = 115
	ExaltationRedeemInfo                 PacketType = 116
	VaultContent                         PacketType = 117
	ForgeRequest                         PacketType = 118
	ForgeResult                          PacketType = 119
	ForgeUnlockedBlueprints              PacketType = 120
	ShootAckCounter                      PacketType = 121
	ChangeAllyShoot                      PacketType = 122
	PlayersList                          PacketType = 123
	ModeratorAction                      PacketType = 124
	GetPlayersList                       PacketType = 125
	CreepMove                            PacketType = 126
	CustomMapDelete                      PacketType = 129
	CustomMapDeleteResponse              PacketType = 130
	CustomMapList                        PacketType = 131
	CustomMapListResponse                PacketType = 132
	CreepHit                             PacketType = 133
	PlayerCallout                        PacketType = 134
	RefineResult                         PacketType = 135
	BuyRefinement                        PacketType = 136
	StartUse                             PacketType = 137
	EndUse                               PacketType = 138
	Stacks                               PacketType = 139
	BuyItem                              PacketType = 140
	BuyItemResult                        PacketType = 141
	DrawDebugShape                       PacketType = 142
	DrawDebugArrow                       PacketType = 143
	DashReset                            PacketType = 144
	FavorPet                             PacketType = 145
	SkinRecycle                          PacketType = 146
	SkinRecycleResponse                  PacketType = 147
	DamageBoost                          PacketType = 148
	ClaimBPMilestone                     PacketType = 149
	ClaimBPMilestoneResult               PacketType = 150
	BoostBPMilestone                     PacketType = 151
	BoostBPMilestoneResult               PacketType = 152
	AcceleratorAdded                     PacketType = 153
	UnseasonRequest                      PacketType = 154
	Retitle                              PacketType = 155
	SetGravestone                        PacketType = 156
	SetAbility                           PacketType = 157
	MissionProgressUpdate                PacketType = 158
	Emote                                PacketType = 159
	BuyEmote                             PacketType = 160
	SetTrackedSeason                     PacketType = 162
	ClaimMission                         PacketType = 163
	ClaimMissionResult                   PacketType = 164
	MultipleMissionsProgressUpdate       PacketType = 165
	DamageWithEffect                     PacketType = 166
	SetDiscoverable                      PacketType = 167
	RealmScoreUpdate                     PacketType = 169
	ClaimChestReward                     PacketType = 171
	UnlockEnchantment                    PacketType = 175
	ApplyEnchantment                     PacketType = 177
	BaseEnchantmentResult                PacketType = 178
	EnableCrucible                       PacketType = 180
	CrucibleResult                       PacketType = 181
	BuyCustomization                     PacketType = 182
	CrucibleInfo                         PacketType = 183
	TutorialStateChanged                 PacketType = 184
	EnchantReroll                        PacketType = 189
	ResetEnchantmentsRerollCountMessage  PacketType = 191
	ResetEnchantmentsRerollCountResponse PacketType = 192
	DismantleRequest                     PacketType = 195
	DismantleResponse                    PacketType = 196
	PartyCreate                          PacketType = 200
	PartyList                            PacketType = 214
	PartyJoinResponse                    PacketType = 218
	BuyItems                             PacketType = 223

	PartyActionResult        PacketType = -2
	PartyInviteResponse      PacketType = -3
	PartyJoinRequest         PacketType = -4
	PartyAction              PacketType = -5
	PartyJoinRequestResponse PacketType = -6
	PartyMemberAdded         PacketType = -7
	IncomingPartyInvite      PacketType = -8
	IncomingPartyMemberInfo  PacketType = -9
)

// Reader defines the interface for reading packet data
type Reader interface {
	ReadInt16() (int16, error)
	ReadUInt16() (uint16, error)
	ReadInt32() (int32, error)
	ReadUInt32() (uint32, error)
	ReadFloat32() (float32, error)
	ReadString() (string, error)
	ReadUTF32String() (string, error)
	ReadCompressedInt() (int, error)
	ReadByte() (byte, error)
	ReadBytes(n int) ([]byte, error)
	ReadBool() (bool, error)
	RemainingBytes() byte
}

// Writer defines the interface for writing packet data
type Writer interface {
	WriteInt16(value int16) error
	WriteUInt16(value uint16) error
	WriteInt32(value int32) error
	WriteUInt32(value uint32) error
	WriteFloat32(value float32) error
	WriteString(value string) error
	WriteUTF32String(value string) error
	WriteCompressedInt(value int) error
	WriteByte(value byte) error
	WriteBytes(data []byte) error
	WriteBool(value bool) error
	Bytes() []byte
}
