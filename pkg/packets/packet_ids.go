package packets

// Packet type ID constants
const (
	FAILURE                           = 0
	TELEPORT                          = 1
	CLAIMDAILYREWARD                  = 3
	DELETEPET                         = 4
	REQUESTTRADE                      = 5
	QUESTFETCHRESPONSE                = 6
	JOINGUILD                         = 7
	PING                              = 8
	PLAYERTEXT                        = 9
	NEWTICK                           = 10
	SHOWEFFECT                        = 11
	SERVERPLAYERSHOOT                 = 12
	USEITEM                           = 13
	TRADEACCEPTED                     = 14
	GUILDREMOVE                       = 15
	PETUPGRADEREQUEST                 = 16
	GOTO                              = 18
	INVENTORYDROP                     = 19
	OTHERHIT                          = 20
	NAMERESULT                        = 21
	BUYRESULT                         = 22
	HATCHPET                          = 23
	ACTIVEPETUPDATEREQUEST            = 24
	ENEMYHIT                          = 25
	GUILDRESULT                       = 26
	EDITACCOUNTLIST                   = 27
	TRADECHANGED                      = 28
	PLAYERSHOOT                       = 30
	PONG                              = 31
	CHANGEPETSKIN                     = 33
	TRADEDONE                         = 34
	ENEMYSHOOT                        = 35
	ACCEPTTRADE                       = 36
	CHANGEGUILDRANK                   = 37
	PLAYSOUND                         = 38
	SQUAREHIT                         = 40
	NEWABILITY                        = 41
	UPDATE                            = 42
	TEXT                              = 44
	RECONNECT                         = 45
	DEATH                             = 46
	USEPORTAL                         = 47
	GOTOQUESTROOM                     = 48
	ALLYSHOOT                         = 49
	RESKIN                            = 51
	RESETDAILYQUESTS                  = 52
	INVENTORYSWAP                     = 55
	CHANGETRADE                       = 56
	CREATE                            = 57
	QUESTREDEEM                       = 58
	CREATEGUILD                       = 59
	SETCONDITION                      = 60
	LOAD                              = 61
	MOVE                              = 62
	KEYINFORESPONSE                   = 63
	AOE                               = 64
	GOTOACK                           = 65
	NOTIFICATION                      = 67
	CLIENTSTAT                        = 69
	HELLO                             = 74
	DAMAGE                            = 75
	ACTIVEPET                         = 76
	INVITEDTOGUILD                    = 77
	PETYARDUPDATE                     = 78
	PASSWORDPROMPT                    = 79
	UPDATEACK                         = 81
	QUESTOBJECTID                     = 82
	PIC                               = 83
	HEROLEFT                          = 84
	BUY                               = 85
	TRADESTART                        = 86
	EVOLVEDPET                        = 87
	TRADEREQUESTED                    = 88
	AOEACK                            = 89
	PLAYERHIT                         = 90
	CANCELTRADE                       = 91
	MAPINFO                           = 92
	KEYINFOREQUEST                    = 94
	INVENTORYRESULT                   = 95
	QUESTREDEEMRESPONSE               = 96
	CHOOSENAME                        = 97
	QUESTFETCHASK                     = 98
	ACCOUNTLIST                       = 99
	CREATESUCCESS                     = 101
	CHECKCREDITS                      = 102
	GROUNDDAMAGE                      = 103
	GUILDINVITE                       = 104
	ESCAPE                            = 105
	FILE                              = 106
	UNLOCKCUSTOMIZATION               = 107
	NEWCHARACTERINFORMATION           = 108
	UNLOCKNEWSLOT                     = 109
	QUEUE                             = 112
	QUEUECANCEL                       = 113
	EXALTATIONBONUSCHANGED            = 114
	REDEEMEXALTATIONREWARD            = 115
	EXALTATIONREDEEMINFO              = 116
	VAULTCONTENT                      = 117
	FORGEREQUEST                      = 118
	FORGERESULT                       = 119
	FORGEUNLOCKEDBLUEPRINTS           = 120
	SHOOTACKCOUNTER                   = 121
	CHANGEALLYSHOOT                   = 122
	PLAYERSLIST                       = 123
	MODERATORACTION                   = 124
	GETPLAYERSLIST                    = 125
	CREEPMOVE                         = 126
	CUSTOMMAPDELETE                   = 129
	CUSTOMMAPDELETERESPONSE           = 130
	CUSTOMMAPLIST                     = 131
	CUSTOMMAPLISTRESPONSE             = 132
	CREEPHIT                          = 133
	PLAYERCALLOUT                     = 134
	REFINERESULT                      = 135
	BUYREFINEMENT                     = 136
	STARTUSE                          = 137
	ENDUSE                            = 138
	STACKS                            = 139
	BUYITEM                           = 140
	BUYITEMRESULT                     = 141
	DRAWDEBUGSHAPE                    = 142
	DRAWDEBUGARROW                    = 143
	DASHRESET                         = 144
	FAVORPET                          = 145
	SKINRECYCLE                       = 146
	SKINRECYCLERESPONSE               = 147
	DAMAGEBOOST                       = 148
	CLAIMBPMILESTONE                  = 149
	CLAIMBPMILESTONERESULT            = 150
	BOOSTBPMILESTONE                  = 151
	BOOSTBPMILESTONERESULT            = 152
	ACCELERATORADDED                  = 153
	UNSEASONREQUEST                   = 154
	RETITLE                           = 155
	SETGRAVESTONE                     = 156
	SETABILITY                        = 157
	MISSIONPROGRESSUPDATE             = 158
	EMOTE                             = 159
	BUYEMOTE                          = 160
	SETTRACKEDSEASON                  = 162
	CLAIMMISSION                      = 163
	CLAIMMISSIONRESULT                = 164
	MULTIPLEMISSIONSPROGRESSUPDATE    = 165
	DAMAGEWITHEFFECT                  = 166
	SETDISCOVERABLE                   = 167
	REALMSCOREUPDATE                  = 169
	CLAIMREWARDSINFOPROMPT            = 170
	CLAIMCHESTREWARDSUBMIT            = 171
	CHESTREWARDRESULT                 = 172
	UNLOCKENCHANTMENTSLOT             = 173
	UNLOCKENCHANTMENTSLOTRESULT       = 174
	UNLOCKENCHANTMENT                 = 175
	UNLOCKENCHANTMENTRESULT           = 176
	APPLYENCHANTMENT                  = 177
	APPLYENCHANTMENTRESULT            = 178
	ACCELERATORUPDATEDMESSAGE         = 179
	ENABLECRUCIBLE                    = 180
	CRUCIBLERESULT                    = 181
	GETDEFINITION                     = 182
	RESULTDEFINITION                  = 183
	TUTORIALSTATECHANGED              = 184
	UPGRADEENCHANTER                  = 185
	UPGRADEENCHANTERRESULT            = 186
	UPGRADEENCHANTMENT                = 187
	UPGRADEENCHANTMENTRESULT          = 188
	REROLLENCHANTMENTS                = 189
	REROLLENCHANTMENTSRESULT          = 190
	RESETENCHANTMENTSREROLLCOUNT      = 191
	RESETENCHANTMENTREROLLCOUNTRESULT = 192
	PURCHASEPETSHADER                 = 193
	PETSHADERPURCHASERESULT           = 194
	DISMANTLEITEMMESSAGE              = 195
	UNKNOWN196                        = 196
	UNKNOWN197                        = 197
	UNKNOWN198                        = 198
	UNKNOWN199                        = 199
	CREATEPARTYMESSAGE                = 200
	PARTYACTION                       = 204
	PARTYACTIONRESULT                 = 207
	INCOMINGPARTINVITATION            = 208
	PARTYINVITATIONRESPONSE           = 209
	INCOMINGPARTYMEMBERINFO           = 210
	PARTYMEMBERADDED                  = 212
	PARTYLISTMESSAGE                  = 214
	PARTYJOINREQUEST                  = 215
	PARTYREQUESTRESPONSE              = 217
	FORRECONNECTMESSAGE               = 218
	UNDEFINED                         = 254
	UNKNOWN                           = 255
)
