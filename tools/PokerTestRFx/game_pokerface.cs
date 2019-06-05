//------------------------------------------------------------------------------
// <auto-generated>
//     This code was generated by a tool.
//
//     Changes to this file may cause incorrect behavior and will be lost if
//     the code is regenerated.
// </auto-generated>
//------------------------------------------------------------------------------

// Generated from: game_pokerface.proto
namespace pokerface
{
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"GameMessage")]
  public partial class GameMessage : global::ProtoBuf.IExtensible
  {
    public GameMessage() {}
    
    private int _Ops;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"Ops", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int Ops
    {
      get { return _Ops; }
      set { _Ops = value; }
    }

    private byte[] _Data = null;
    [global::ProtoBuf.ProtoMember(2, IsRequired = false, Name=@"Data", DataFormat = global::ProtoBuf.DataFormat.Default)]
    [global::System.ComponentModel.DefaultValue(null)]
    public byte[] Data
    {
      get { return _Data; }
      set { _Data = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgCardHand")]
  public partial class MsgCardHand : global::ProtoBuf.IExtensible
  {
    public MsgCardHand() {}
    
    private int _cardHandType;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"cardHandType", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int cardHandType
    {
      get { return _cardHandType; }
      set { _cardHandType = value; }
    }
    private readonly global::System.Collections.Generic.List<int> _cards = new global::System.Collections.Generic.List<int>();
    [global::ProtoBuf.ProtoMember(2, Name=@"cards", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public global::System.Collections.Generic.List<int> cards
    {
      get { return _cards; }
    }
  
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgPlayerCardList")]
  public partial class MsgPlayerCardList : global::ProtoBuf.IExtensible
  {
    public MsgPlayerCardList() {}
    
    private int _chairID;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"chairID", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int chairID
    {
      get { return _chairID; }
      set { _chairID = value; }
    }
    private int _cardCountOnHand;
    [global::ProtoBuf.ProtoMember(2, IsRequired = true, Name=@"cardCountOnHand", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int cardCountOnHand
    {
      get { return _cardCountOnHand; }
      set { _cardCountOnHand = value; }
    }
    private readonly global::System.Collections.Generic.List<int> _cardsOnHand = new global::System.Collections.Generic.List<int>();
    [global::ProtoBuf.ProtoMember(3, Name=@"cardsOnHand", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public global::System.Collections.Generic.List<int> cardsOnHand
    {
      get { return _cardsOnHand; }
    }
  
    private readonly global::System.Collections.Generic.List<pokerface.MsgCardHand> _discardedHands = new global::System.Collections.Generic.List<pokerface.MsgCardHand>();
    [global::ProtoBuf.ProtoMember(4, Name=@"discardedHands", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public global::System.Collections.Generic.List<pokerface.MsgCardHand> discardedHands
    {
      get { return _discardedHands; }
    }
  
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgDeal")]
  public partial class MsgDeal : global::ProtoBuf.IExtensible
  {
    public MsgDeal() {}
    
    private int _bankerChairID;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"bankerChairID", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int bankerChairID
    {
      get { return _bankerChairID; }
      set { _bankerChairID = value; }
    }
    private int _windFlowerID;
    [global::ProtoBuf.ProtoMember(2, IsRequired = true, Name=@"windFlowerID", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int windFlowerID
    {
      get { return _windFlowerID; }
      set { _windFlowerID = value; }
    }
    private readonly global::System.Collections.Generic.List<pokerface.MsgPlayerCardList> _playerCardLists = new global::System.Collections.Generic.List<pokerface.MsgPlayerCardList>();
    [global::ProtoBuf.ProtoMember(3, Name=@"playerCardLists", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public global::System.Collections.Generic.List<pokerface.MsgPlayerCardList> playerCardLists
    {
      get { return _playerCardLists; }
    }
  
    private int _cardsInWall;
    [global::ProtoBuf.ProtoMember(4, IsRequired = true, Name=@"cardsInWall", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int cardsInWall
    {
      get { return _cardsInWall; }
      set { _cardsInWall = value; }
    }

    private int _dice1 = default(int);
    [global::ProtoBuf.ProtoMember(5, IsRequired = false, Name=@"dice1", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    [global::System.ComponentModel.DefaultValue(default(int))]
    public int dice1
    {
      get { return _dice1; }
      set { _dice1 = value; }
    }

    private int _dice2 = default(int);
    [global::ProtoBuf.ProtoMember(6, IsRequired = false, Name=@"dice2", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    [global::System.ComponentModel.DefaultValue(default(int))]
    public int dice2
    {
      get { return _dice2; }
      set { _dice2 = value; }
    }

    private bool _isContinuousBanker = default(bool);
    [global::ProtoBuf.ProtoMember(7, IsRequired = false, Name=@"isContinuousBanker", DataFormat = global::ProtoBuf.DataFormat.Default)]
    [global::System.ComponentModel.DefaultValue(default(bool))]
    public bool isContinuousBanker
    {
      get { return _isContinuousBanker; }
      set { _isContinuousBanker = value; }
    }

    private int _markup = default(int);
    [global::ProtoBuf.ProtoMember(8, IsRequired = false, Name=@"markup", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    [global::System.ComponentModel.DefaultValue(default(int))]
    public int markup
    {
      get { return _markup; }
      set { _markup = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgAllowPlayerAction")]
  public partial class MsgAllowPlayerAction : global::ProtoBuf.IExtensible
  {
    public MsgAllowPlayerAction() {}
    
    private int _qaIndex;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"qaIndex", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int qaIndex
    {
      get { return _qaIndex; }
      set { _qaIndex = value; }
    }
    private int _actionChairID;
    [global::ProtoBuf.ProtoMember(2, IsRequired = true, Name=@"actionChairID", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int actionChairID
    {
      get { return _actionChairID; }
      set { _actionChairID = value; }
    }
    private int _allowedActions;
    [global::ProtoBuf.ProtoMember(3, IsRequired = true, Name=@"allowedActions", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int allowedActions
    {
      get { return _allowedActions; }
      set { _allowedActions = value; }
    }

    private int _timeoutInSeconds = default(int);
    [global::ProtoBuf.ProtoMember(4, IsRequired = false, Name=@"timeoutInSeconds", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    [global::System.ComponentModel.DefaultValue(default(int))]
    public int timeoutInSeconds
    {
      get { return _timeoutInSeconds; }
      set { _timeoutInSeconds = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgAllowPlayerReAction")]
  public partial class MsgAllowPlayerReAction : global::ProtoBuf.IExtensible
  {
    public MsgAllowPlayerReAction() {}
    
    private int _qaIndex;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"qaIndex", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int qaIndex
    {
      get { return _qaIndex; }
      set { _qaIndex = value; }
    }
    private int _actionChairID;
    [global::ProtoBuf.ProtoMember(2, IsRequired = true, Name=@"actionChairID", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int actionChairID
    {
      get { return _actionChairID; }
      set { _actionChairID = value; }
    }
    private int _allowedActions;
    [global::ProtoBuf.ProtoMember(3, IsRequired = true, Name=@"allowedActions", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int allowedActions
    {
      get { return _allowedActions; }
      set { _allowedActions = value; }
    }

    private int _timeoutInSeconds = default(int);
    [global::ProtoBuf.ProtoMember(4, IsRequired = false, Name=@"timeoutInSeconds", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    [global::System.ComponentModel.DefaultValue(default(int))]
    public int timeoutInSeconds
    {
      get { return _timeoutInSeconds; }
      set { _timeoutInSeconds = value; }
    }

    private int _prevActionChairID = default(int);
    [global::ProtoBuf.ProtoMember(5, IsRequired = false, Name=@"prevActionChairID", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    [global::System.ComponentModel.DefaultValue(default(int))]
    public int prevActionChairID
    {
      get { return _prevActionChairID; }
      set { _prevActionChairID = value; }
    }

    private pokerface.MsgCardHand _prevActionHand = null;
    [global::ProtoBuf.ProtoMember(6, IsRequired = false, Name=@"prevActionHand", DataFormat = global::ProtoBuf.DataFormat.Default)]
    [global::System.ComponentModel.DefaultValue(null)]
    public pokerface.MsgCardHand prevActionHand
    {
      get { return _prevActionHand; }
      set { _prevActionHand = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgPlayerAction")]
  public partial class MsgPlayerAction : global::ProtoBuf.IExtensible
  {
    public MsgPlayerAction() {}
    
    private int _qaIndex;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"qaIndex", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int qaIndex
    {
      get { return _qaIndex; }
      set { _qaIndex = value; }
    }
    private int _action;
    [global::ProtoBuf.ProtoMember(2, IsRequired = true, Name=@"action", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int action
    {
      get { return _action; }
      set { _action = value; }
    }

    private int _flags = default(int);
    [global::ProtoBuf.ProtoMember(3, IsRequired = false, Name=@"flags", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    [global::System.ComponentModel.DefaultValue(default(int))]
    public int flags
    {
      get { return _flags; }
      set { _flags = value; }
    }
    private readonly global::System.Collections.Generic.List<int> _cards = new global::System.Collections.Generic.List<int>();
    [global::ProtoBuf.ProtoMember(4, Name=@"cards", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public global::System.Collections.Generic.List<int> cards
    {
      get { return _cards; }
    }
  
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgActionResultNotify")]
  public partial class MsgActionResultNotify : global::ProtoBuf.IExtensible
  {
    public MsgActionResultNotify() {}
    
    private int _targetChairID;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"targetChairID", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int targetChairID
    {
      get { return _targetChairID; }
      set { _targetChairID = value; }
    }
    private int _action;
    [global::ProtoBuf.ProtoMember(2, IsRequired = true, Name=@"action", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int action
    {
      get { return _action; }
      set { _action = value; }
    }

    private pokerface.MsgCardHand _actionHand = null;
    [global::ProtoBuf.ProtoMember(3, IsRequired = false, Name=@"actionHand", DataFormat = global::ProtoBuf.DataFormat.Default)]
    [global::System.ComponentModel.DefaultValue(null)]
    public pokerface.MsgCardHand actionHand
    {
      get { return _actionHand; }
      set { _actionHand = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgRestore")]
  public partial class MsgRestore : global::ProtoBuf.IExtensible
  {
    public MsgRestore() {}
    
    private pokerface.MsgDeal _msgDeal;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"msgDeal", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public pokerface.MsgDeal msgDeal
    {
      get { return _msgDeal; }
      set { _msgDeal = value; }
    }

    private int _prevActionChairID = default(int);
    [global::ProtoBuf.ProtoMember(2, IsRequired = false, Name=@"prevActionChairID", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    [global::System.ComponentModel.DefaultValue(default(int))]
    public int prevActionChairID
    {
      get { return _prevActionChairID; }
      set { _prevActionChairID = value; }
    }

    private pokerface.MsgCardHand _prevActionHand = null;
    [global::ProtoBuf.ProtoMember(3, IsRequired = false, Name=@"prevActionHand", DataFormat = global::ProtoBuf.DataFormat.Default)]
    [global::System.ComponentModel.DefaultValue(null)]
    public pokerface.MsgCardHand prevActionHand
    {
      get { return _prevActionHand; }
      set { _prevActionHand = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgPlayerScoreGreatWin")]
  public partial class MsgPlayerScoreGreatWin : global::ProtoBuf.IExtensible
  {
    public MsgPlayerScoreGreatWin() {}
    
    private int _baseWinScore;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"baseWinScore", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int baseWinScore
    {
      get { return _baseWinScore; }
      set { _baseWinScore = value; }
    }
    private int _greatWinType;
    [global::ProtoBuf.ProtoMember(2, IsRequired = true, Name=@"greatWinType", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int greatWinType
    {
      get { return _greatWinType; }
      set { _greatWinType = value; }
    }
    private int _greatWinPoints;
    [global::ProtoBuf.ProtoMember(3, IsRequired = true, Name=@"greatWinPoints", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int greatWinPoints
    {
      get { return _greatWinPoints; }
      set { _greatWinPoints = value; }
    }
    private int _trimGreatWinPoints;
    [global::ProtoBuf.ProtoMember(4, IsRequired = true, Name=@"trimGreatWinPoints", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int trimGreatWinPoints
    {
      get { return _trimGreatWinPoints; }
      set { _trimGreatWinPoints = value; }
    }

    private int _continuousBankerExtra = default(int);
    [global::ProtoBuf.ProtoMember(5, IsRequired = false, Name=@"continuousBankerExtra", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    [global::System.ComponentModel.DefaultValue(default(int))]
    public int continuousBankerExtra
    {
      get { return _continuousBankerExtra; }
      set { _continuousBankerExtra = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgPlayerScoreMiniWin")]
  public partial class MsgPlayerScoreMiniWin : global::ProtoBuf.IExtensible
  {
    public MsgPlayerScoreMiniWin() {}
    
    private int _miniWinType;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"miniWinType", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int miniWinType
    {
      get { return _miniWinType; }
      set { _miniWinType = value; }
    }
    private int _miniWinBasicScore;
    [global::ProtoBuf.ProtoMember(2, IsRequired = true, Name=@"miniWinBasicScore", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int miniWinBasicScore
    {
      get { return _miniWinBasicScore; }
      set { _miniWinBasicScore = value; }
    }
    private int _miniWinFlowerScore;
    [global::ProtoBuf.ProtoMember(3, IsRequired = true, Name=@"miniWinFlowerScore", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int miniWinFlowerScore
    {
      get { return _miniWinFlowerScore; }
      set { _miniWinFlowerScore = value; }
    }
    private int _miniMultiple;
    [global::ProtoBuf.ProtoMember(4, IsRequired = true, Name=@"miniMultiple", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int miniMultiple
    {
      get { return _miniMultiple; }
      set { _miniMultiple = value; }
    }
    private int _miniWinTrimScore;
    [global::ProtoBuf.ProtoMember(5, IsRequired = true, Name=@"miniWinTrimScore", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int miniWinTrimScore
    {
      get { return _miniWinTrimScore; }
      set { _miniWinTrimScore = value; }
    }

    private int _continuousBankerExtra = default(int);
    [global::ProtoBuf.ProtoMember(6, IsRequired = false, Name=@"continuousBankerExtra", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    [global::System.ComponentModel.DefaultValue(default(int))]
    public int continuousBankerExtra
    {
      get { return _continuousBankerExtra; }
      set { _continuousBankerExtra = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgPlayerScore")]
  public partial class MsgPlayerScore : global::ProtoBuf.IExtensible
  {
    public MsgPlayerScore() {}
    
    private int _targetChairID;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"targetChairID", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int targetChairID
    {
      get { return _targetChairID; }
      set { _targetChairID = value; }
    }
    private int _winType;
    [global::ProtoBuf.ProtoMember(2, IsRequired = true, Name=@"winType", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int winType
    {
      get { return _winType; }
      set { _winType = value; }
    }
    private int _score;
    [global::ProtoBuf.ProtoMember(3, IsRequired = true, Name=@"score", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int score
    {
      get { return _score; }
      set { _score = value; }
    }
    private int _specialScore;
    [global::ProtoBuf.ProtoMember(4, IsRequired = true, Name=@"specialScore", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int specialScore
    {
      get { return _specialScore; }
      set { _specialScore = value; }
    }

    private pokerface.MsgPlayerScoreGreatWin _greatWin = null;
    [global::ProtoBuf.ProtoMember(5, IsRequired = false, Name=@"greatWin", DataFormat = global::ProtoBuf.DataFormat.Default)]
    [global::System.ComponentModel.DefaultValue(null)]
    public pokerface.MsgPlayerScoreGreatWin greatWin
    {
      get { return _greatWin; }
      set { _greatWin = value; }
    }

    private pokerface.MsgPlayerScoreMiniWin _miniWin = null;
    [global::ProtoBuf.ProtoMember(6, IsRequired = false, Name=@"miniWin", DataFormat = global::ProtoBuf.DataFormat.Default)]
    [global::System.ComponentModel.DefaultValue(null)]
    public pokerface.MsgPlayerScoreMiniWin miniWin
    {
      get { return _miniWin; }
      set { _miniWin = value; }
    }

    private int _fakeWinScore = default(int);
    [global::ProtoBuf.ProtoMember(7, IsRequired = false, Name=@"fakeWinScore", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    [global::System.ComponentModel.DefaultValue(default(int))]
    public int fakeWinScore
    {
      get { return _fakeWinScore; }
      set { _fakeWinScore = value; }
    }
    private readonly global::System.Collections.Generic.List<int> _fakeList = new global::System.Collections.Generic.List<int>();
    [global::ProtoBuf.ProtoMember(8, Name=@"fakeList", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public global::System.Collections.Generic.List<int> fakeList
    {
      get { return _fakeList; }
    }
  

    private bool _isContinuousBanker = default(bool);
    [global::ProtoBuf.ProtoMember(9, IsRequired = false, Name=@"isContinuousBanker", DataFormat = global::ProtoBuf.DataFormat.Default)]
    [global::System.ComponentModel.DefaultValue(default(bool))]
    public bool isContinuousBanker
    {
      get { return _isContinuousBanker; }
      set { _isContinuousBanker = value; }
    }

    private int _continuousBankerMultiple = default(int);
    [global::ProtoBuf.ProtoMember(10, IsRequired = false, Name=@"continuousBankerMultiple", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    [global::System.ComponentModel.DefaultValue(default(int))]
    public int continuousBankerMultiple
    {
      get { return _continuousBankerMultiple; }
      set { _continuousBankerMultiple = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgHandScore")]
  public partial class MsgHandScore : global::ProtoBuf.IExtensible
  {
    public MsgHandScore() {}
    
    private readonly global::System.Collections.Generic.List<pokerface.MsgPlayerScore> _playerScores = new global::System.Collections.Generic.List<pokerface.MsgPlayerScore>();
    [global::ProtoBuf.ProtoMember(1, Name=@"playerScores", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public global::System.Collections.Generic.List<pokerface.MsgPlayerScore> playerScores
    {
      get { return _playerScores; }
    }
  
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
  [global::System.Serializable, global::ProtoBuf.ProtoContract(Name=@"MsgHandOver")]
  public partial class MsgHandOver : global::ProtoBuf.IExtensible
  {
    public MsgHandOver() {}
    
    private int _endType;
    [global::ProtoBuf.ProtoMember(1, IsRequired = true, Name=@"endType", DataFormat = global::ProtoBuf.DataFormat.TwosComplement)]
    public int endType
    {
      get { return _endType; }
      set { _endType = value; }
    }
    private readonly global::System.Collections.Generic.List<pokerface.MsgPlayerCardList> _playerCardLists = new global::System.Collections.Generic.List<pokerface.MsgPlayerCardList>();
    [global::ProtoBuf.ProtoMember(2, Name=@"playerCardLists", DataFormat = global::ProtoBuf.DataFormat.Default)]
    public global::System.Collections.Generic.List<pokerface.MsgPlayerCardList> playerCardLists
    {
      get { return _playerCardLists; }
    }
  

    private pokerface.MsgHandScore _scores = null;
    [global::ProtoBuf.ProtoMember(3, IsRequired = false, Name=@"scores", DataFormat = global::ProtoBuf.DataFormat.Default)]
    [global::System.ComponentModel.DefaultValue(null)]
    public pokerface.MsgHandScore scores
    {
      get { return _scores; }
      set { _scores = value; }
    }

    private bool _continueAble = default(bool);
    [global::ProtoBuf.ProtoMember(4, IsRequired = false, Name=@"continueAble", DataFormat = global::ProtoBuf.DataFormat.Default)]
    [global::System.ComponentModel.DefaultValue(default(bool))]
    public bool continueAble
    {
      get { return _continueAble; }
      set { _continueAble = value; }
    }
    private global::ProtoBuf.IExtension extensionObject;
    global::ProtoBuf.IExtension global::ProtoBuf.IExtensible.GetExtensionObject(bool createIfMissing)
      { return global::ProtoBuf.Extensible.GetExtensionObject(ref extensionObject, createIfMissing); }
  }
  
    [global::ProtoBuf.ProtoContract(Name=@"CardID")]
    public enum CardID
    {
            
      [global::ProtoBuf.ProtoEnum(Name=@"R2H", Value=0)]
      R2H = 0,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R2D", Value=1)]
      R2D = 1,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R2C", Value=2)]
      R2C = 2,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R2S", Value=3)]
      R2S = 3,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R3H", Value=4)]
      R3H = 4,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R3D", Value=5)]
      R3D = 5,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R3C", Value=6)]
      R3C = 6,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R3S", Value=7)]
      R3S = 7,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R4H", Value=8)]
      R4H = 8,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R4D", Value=9)]
      R4D = 9,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R4C", Value=10)]
      R4C = 10,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R4S", Value=11)]
      R4S = 11,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R5H", Value=12)]
      R5H = 12,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R5D", Value=13)]
      R5D = 13,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R5C", Value=14)]
      R5C = 14,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R5S", Value=15)]
      R5S = 15,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R6H", Value=16)]
      R6H = 16,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R6D", Value=17)]
      R6D = 17,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R6C", Value=18)]
      R6C = 18,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R6S", Value=19)]
      R6S = 19,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R7H", Value=20)]
      R7H = 20,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R7D", Value=21)]
      R7D = 21,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R7C", Value=22)]
      R7C = 22,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R7S", Value=23)]
      R7S = 23,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R8H", Value=24)]
      R8H = 24,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R8D", Value=25)]
      R8D = 25,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R8C", Value=26)]
      R8C = 26,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R8S", Value=27)]
      R8S = 27,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R9H", Value=28)]
      R9H = 28,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R9D", Value=29)]
      R9D = 29,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R9C", Value=30)]
      R9C = 30,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R9S", Value=31)]
      R9S = 31,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R10H", Value=32)]
      R10H = 32,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R10D", Value=33)]
      R10D = 33,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R10C", Value=34)]
      R10C = 34,
            
      [global::ProtoBuf.ProtoEnum(Name=@"R10S", Value=35)]
      R10S = 35,
            
      [global::ProtoBuf.ProtoEnum(Name=@"JH", Value=36)]
      JH = 36,
            
      [global::ProtoBuf.ProtoEnum(Name=@"JD", Value=37)]
      JD = 37,
            
      [global::ProtoBuf.ProtoEnum(Name=@"JC", Value=38)]
      JC = 38,
            
      [global::ProtoBuf.ProtoEnum(Name=@"JS", Value=39)]
      JS = 39,
            
      [global::ProtoBuf.ProtoEnum(Name=@"QH", Value=40)]
      QH = 40,
            
      [global::ProtoBuf.ProtoEnum(Name=@"QD", Value=41)]
      QD = 41,
            
      [global::ProtoBuf.ProtoEnum(Name=@"QC", Value=42)]
      QC = 42,
            
      [global::ProtoBuf.ProtoEnum(Name=@"QS", Value=43)]
      QS = 43,
            
      [global::ProtoBuf.ProtoEnum(Name=@"KH", Value=44)]
      KH = 44,
            
      [global::ProtoBuf.ProtoEnum(Name=@"KD", Value=45)]
      KD = 45,
            
      [global::ProtoBuf.ProtoEnum(Name=@"KC", Value=46)]
      KC = 46,
            
      [global::ProtoBuf.ProtoEnum(Name=@"KS", Value=47)]
      KS = 47,
            
      [global::ProtoBuf.ProtoEnum(Name=@"AH", Value=48)]
      AH = 48,
            
      [global::ProtoBuf.ProtoEnum(Name=@"AD", Value=49)]
      AD = 49,
            
      [global::ProtoBuf.ProtoEnum(Name=@"AC", Value=50)]
      AC = 50,
            
      [global::ProtoBuf.ProtoEnum(Name=@"AS", Value=51)]
      AS = 51,
            
      [global::ProtoBuf.ProtoEnum(Name=@"JOB", Value=52)]
      JOB = 52,
            
      [global::ProtoBuf.ProtoEnum(Name=@"JOR", Value=53)]
      JOR = 53,
            
      [global::ProtoBuf.ProtoEnum(Name=@"CARDMAX", Value=54)]
      CARDMAX = 54
    }
  
    [global::ProtoBuf.ProtoContract(Name=@"MessageCode")]
    public enum MessageCode
    {
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPInvalid", Value=0)]
      OPInvalid = 0,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPAction", Value=1)]
      OPAction = 1,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPActionResultNotify", Value=2)]
      OPActionResultNotify = 2,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPActionAllowed", Value=3)]
      OPActionAllowed = 3,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPReActionAllowed", Value=5)]
      OPReActionAllowed = 5,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPDeal", Value=6)]
      OPDeal = 6,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPHandOver", Value=7)]
      OPHandOver = 7,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPRestore", Value=8)]
      OPRestore = 8,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPPlayerLeaveRoom", Value=9)]
      OPPlayerLeaveRoom = 9,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPPlayerEnterRoom", Value=10)]
      OPPlayerEnterRoom = 10,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPDisbandRequest", Value=11)]
      OPDisbandRequest = 11,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPDisbandNotify", Value=12)]
      OPDisbandNotify = 12,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPDisbandAnswer", Value=13)]
      OPDisbandAnswer = 13,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPPlayerReady", Value=14)]
      OPPlayerReady = 14,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPRoomDeleted", Value=15)]
      OPRoomDeleted = 15,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPRoomUpdate", Value=16)]
      OPRoomUpdate = 16,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPRoomShowTips", Value=17)]
      OPRoomShowTips = 17,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPGameOver", Value=18)]
      OPGameOver = 18,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPKickout", Value=19)]
      OPKickout = 19,
            
      [global::ProtoBuf.ProtoEnum(Name=@"OPDonate", Value=20)]
      OPDonate = 20
    }
  
}