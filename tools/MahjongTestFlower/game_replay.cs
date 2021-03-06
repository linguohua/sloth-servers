//------------------------------------------------------------------------------
// <auto-generated>
//     This code was generated by a tool.
//
//     Changes to this file may cause incorrect behavior and will be lost if
//     the code is regenerated.
// </auto-generated>
//------------------------------------------------------------------------------

// Generated from: game_replay.proto
namespace mahjong
{
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgReplayPlayerInfo")]
  public partial class MsgReplayPlayerInfo : global::ProtoBuf.IExtensible
  {
    public MsgReplayPlayerInfo() {}
    
    private string _userID;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"userID", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public string userID
    {
      get { return _userID; }
      set { _userID = value; }
    }

    private string _nick = "";
    [global::ProtoBuf.ProtoMember(2, IsRequired = false, Name=@"nick", DataFormat = global::ProtoBuf.DataFormat.Default)]
    [global::System.ComponentModel.DefaultValue("")]
    public string nick
    {
      get { return _nick; }
      set { _nick = value; }
    }
    private int _chairID;
    [global::ProtoBuf.ProtoMember(3, IsRequired = true, Name=@"chairID", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int chairID
    {
      get { return _chairID; }
      set { _chairID = value; }
    }

    private int _totalScore = default(int);
    [global::ProtoBuf.ProtoMember(4, IsRequired = false, Name=@"totalScore", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    [global::System.ComponentModel.DefaultValue(default(int))]
    public int totalScore
    {
      get { return _totalScore; }
      set { _totalScore = value; }
    }

    private uint _sex = default(uint);
    [global::ProtoBuf.ProtoMember(5, IsRequired = false, Name=@"sex", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    [global::System.ComponentModel.DefaultValue(default(uint))]
    public uint sex
    {
      get { return _sex; }
      set { _sex = value; }
    }

    private string _headIconURI = "";
    [global::ProtoBuf.ProtoMember(6, IsRequired = false, Name=@"headIconURI", DataFormat = global::ProtoBuf.DataFormat.Default)]
    [global::System.ComponentModel.DefaultValue("")]
    public string headIconURI
    {
      get { return _headIconURI; }
      set { _headIconURI = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgReplayPlayerScoreSummary")]
  public partial class MsgReplayPlayerScoreSummary : global::ProtoBuf.IExtensible
  {
    public MsgReplayPlayerScoreSummary() {}
    
    private int _chairID;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"chairID", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int chairID
    {
      get { return _chairID; }
      set { _chairID = value; }
    }
    private int _score;
    [global::ProtoBuf.ProtoMember(2, IsRequired = true, Name=@"score", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int score
    {
      get { return _score; }
      set { _score = value; }
    }
    private int _winType;
    [global::ProtoBuf.ProtoMember(3, IsRequired = true, Name=@"winType", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int winType
    {
      get { return _winType; }
      set { _winType = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgReplayRecordSummary")]
  public partial class MsgReplayRecordSummary : global::ProtoBuf.IExtensible
  {
    public MsgReplayRecordSummary() {}
    
    private string _recordUUID;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"recordUUID", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public string recordUUID
    {
      get { return _recordUUID; }
      set { _recordUUID = value; }
    }
    private readonly global::System.Collections.Generic.List<mahjong.MsgReplayPlayerScoreSummary> _playerScores = new global::System.Collections.Generic.List<mahjong.MsgReplayPlayerScoreSummary>();
    [global::ProtoBuf.ProtoMember(2, Name=@"playerScores", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public global::System.Collections.Generic.List<mahjong.MsgReplayPlayerScoreSummary> playerScores
    {
      get { return _playerScores; }
    }
  
    private uint _endTime;
    [global::ProtoBuf.ProtoMember(3, IsRequired = true, Name=@"endTime", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public uint endTime
    {
      get { return _endTime; }
      set { _endTime = value; }
    }

    private string _shareAbleID = "";
    [global::ProtoBuf.ProtoMember(4, IsRequired = false, Name=@"shareAbleID", DataFormat = global::ProtoBuf.DataFormat.Default)]
    [global::System.ComponentModel.DefaultValue("")]
    public string shareAbleID
    {
      get { return _shareAbleID; }
      set { _shareAbleID = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgReplayRoom")]
  public partial class MsgReplayRoom : global::ProtoBuf.IExtensible
  {
    public MsgReplayRoom() {}
    
    private int _recordRoomType;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"recordRoomType", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int recordRoomType
    {
      get { return _recordRoomType; }
      set { _recordRoomType = value; }
    }
    private uint _startTime;
    [global::ProtoBuf.ProtoMember(2, IsRequired = true, Name=@"startTime", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public uint startTime
    {
      get { return _startTime; }
      set { _startTime = value; }
    }
    private uint _endTime;
    [global::ProtoBuf.ProtoMember(3, IsRequired = true, Name=@"endTime", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public uint endTime
    {
      get { return _endTime; }
      set { _endTime = value; }
    }
    private string _roomNumber;
    [global::ProtoBuf.ProtoMember(4, IsRequired = true, Name=@"roomNumber", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public string roomNumber
    {
      get { return _roomNumber; }
      set { _roomNumber = value; }
    }
    private readonly global::System.Collections.Generic.List<mahjong.MsgReplayPlayerInfo> _players = new global::System.Collections.Generic.List<mahjong.MsgReplayPlayerInfo>();
    [global::ProtoBuf.ProtoMember(5, Name=@"players", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public global::System.Collections.Generic.List<mahjong.MsgReplayPlayerInfo> players
    {
      get { return _players; }
    }
  
    private readonly global::System.Collections.Generic.List<mahjong.MsgReplayRecordSummary> _records = new global::System.Collections.Generic.List<mahjong.MsgReplayRecordSummary>();
    [global::ProtoBuf.ProtoMember(6, Name=@"records", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public global::System.Collections.Generic.List<mahjong.MsgReplayRecordSummary> records
    {
      get { return _records; }
    }
  
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
}