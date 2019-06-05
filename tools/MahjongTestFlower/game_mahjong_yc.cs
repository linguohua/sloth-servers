//------------------------------------------------------------------------------
// <auto-generated>
//     This code was generated by a tool.
//
//     Changes to this file may cause incorrect behavior and will be lost if
//     the code is regenerated.
// </auto-generated>
//------------------------------------------------------------------------------

// Generated from: game_mahjong_yc.proto
namespace mahjong
{
    [global::ProtoBuf.ProtoContract(Name=@"YCMiniWinType")]
    public enum YCMiniWinType
    {
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCMiniWinType_None", Value=0)]
      enumYCMiniWinType_None = 0,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCMiniWinType_DoubleSuits", Value=1)]
      enumYCMiniWinType_DoubleSuits = 1,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCMiniWinType_SuitSequence", Value=2)]
      enumYCMiniWinType_SuitSequence = 2,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCMiniWinType_SevenPair", Value=4)]
      enumYCMiniWinType_SevenPair = 4,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCMiniWinType_PureSame", Value=8)]
      enumYCMiniWinType_PureSame = 8,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCMiniWinType_PongKong", Value=16)]
      enumYCMiniWinType_PongKong = 16,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCMiniWinType_FrontClear", Value=32)]
      enumYCMiniWinType_FrontClear = 32,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCMiniWinType_RichiWinChucker", Value=64)]
      enumYCMiniWinType_RichiWinChucker = 64,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCMiniWinType_FlyRichiWinChucker", Value=128)]
      enumYCMiniWinType_FlyRichiWinChucker = 128,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCMiniWinType_FlyRichi", Value=256)]
      enumYCMiniWinType_FlyRichi = 256,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCMiniWinType_Heaven", Value=512)]
      enumYCMiniWinType_Heaven = 512,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCMiniWinType_NoFlowerWhenRichi", Value=1024)]
      enumYCMiniWinType_NoFlowerWhenRichi = 1024,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCMiniWinType_FirstReadyHand", Value=2048)]
      enumYCMiniWinType_FirstReadyHand = 2048
    }
  
    [global::ProtoBuf.ProtoContract(Name=@"YCRichiFlag")]
    public enum YCRichiFlag
    {
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCRichi_Normal", Value=0)]
      enumYCRichi_Normal = 0,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCRichi_Fly", Value=1)]
      enumYCRichi_Fly = 1,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCRichi_FlyNoFlower", Value=2)]
      enumYCRichi_FlyNoFlower = 2,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCRichi_FirstReadyHand", Value=4)]
      enumYCRichi_FirstReadyHand = 4
    }
  
    [global::ProtoBuf.ProtoContract(Name=@"YCFlyRichiExtraScoreType")]
    public enum YCFlyRichiExtraScoreType
    {
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_None", Value=0)]
      enumYCFR_None = 0,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_SuitSpecialX1", Value=1)]
      enumYCFR_SuitSpecialX1 = 1,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_SuitSpecialX2", Value=2)]
      enumYCFR_SuitSpecialX2 = 2,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_SuitSpecialX3", Value=4)]
      enumYCFR_SuitSpecialX3 = 4,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_SuitSpecialX4", Value=8)]
      enumYCFR_SuitSpecialX4 = 8,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_ReadyHandDoublePair", Value=16)]
      enumYCFR_ReadyHandDoublePair = 16,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_ReadyHandDaDuziX1", Value=32)]
      enumYCFR_ReadyHandDaDuziX1 = 32,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_ReadyHandDaDuziX2", Value=64)]
      enumYCFR_ReadyHandDaDuziX2 = 64,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_ReadyHandSinglePair", Value=128)]
      enumYCFR_ReadyHandSinglePair = 128,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_ReadyHandSinglePairWith2", Value=256)]
      enumYCFR_ReadyHandSinglePairWith2 = 256,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_ReadyHandSinglePairWith3", Value=512)]
      enumYCFR_ReadyHandSinglePairWith3 = 512,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_GreatSevenPairX1", Value=1024)]
      enumYCFR_GreatSevenPairX1 = 1024,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_GreatSevenPairX2", Value=2048)]
      enumYCFR_GreatSevenPairX2 = 2048,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_GreatSevenPairX3", Value=4096)]
      enumYCFR_GreatSevenPairX3 = 4096,
            
      [global::ProtoBuf.ProtoEnum(Name=@"enumYCFR_SevenPair", Value=8192)]
      enumYCFR_SevenPair = 8192
    }
  
}