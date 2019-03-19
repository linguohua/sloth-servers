package pay

const (
	// OwnerModDiamondCreateRoom 房主创建房间扣钱
	ownerModDiamondCreateRoom = 10000010
	// OwnerModDiamondReturn 房主创建房间返还
	ownerModDiamondReturn = 10000015

	// AAModDiamondCreateRoom AA创建房间扣除钻石
	aaModDiamondCreateRoom = 10000011
	// AAModDiamondReturn AA创建房间返还
	aaModDiamondReturn = 10000016

	// ModDiamondDonate 道具消耗钻石
	modDiamondDonate = 10000021

	// ModDiamondCreateRoomForOther 替人开房
	modDiamondCreateRoomForOther = 10000022

	// ModDiamondCreateRoomForOtherReturn 替人开房返还
	modDiamondCreateRoomForOtherReturn = 10000023

	// AddDiamondFromBackend 运营后台发放
	addDiamondFromBackend = 10000008

	// ownerModDiamondCreateRoomForGroup 牌友群房主创建房间
	ownerModDiamondCreateRoomForGroup = 10000032
	// aaModDiamondCreateRoomForGroup 牌友群AA创建房间
	aaModDiamondCreateRoomForGroup = 10000033
	// ownerModDiamondCreateRoomForGroupReturn 牌友群房主创建返还
	ownerModDiamondCreateRoomForGroupReturn = 10000034
	// aaModDiamondCreateRoomForGroupReturn 牌友群AA创建返还
	aaModDiamondCreateRoomForGroupReturn = 10000035

	// masterModDiamondCreateRoomForGroup 群主代开房
	masterModDiamondCreateRoomForGroup = 10000038
	// masterModDiamondCreateRoomForGroupReturn 群主代开房 返还
	masterModDiamondCreateRoomForGroupReturn = 10000039

	diamondNotEnoughMsg = "更新数量不足"

)
