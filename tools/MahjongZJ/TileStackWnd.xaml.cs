using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Text.RegularExpressions;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Media.Imaging;
using mahjong;

namespace MahjongTest
{
    /// <summary>
    /// TileStackWnd.xaml 的交互逻辑
    /// </summary>
    public partial class TileStackWnd : UserControl
    {
        private int expectedReadyHandFlags;

        public TileStackWnd()
        {
            InitializeComponent();

            InitButtonArray();

            HideAllButtons();
        }


        public  Button [] ButtonsSp1 { get; } = new Button[14];
        public Button[] ButtonsSp2 { get; } = new Button[14];
        public Button[] ButtonsSp3 { get; } = new Button[14];
        public Button[] ButtonsSp4 { get; } = new Button[14];
        public Button[] ButtonsAct { get; } = new Button[5];

        public Dictionary<int, BitmapImage> ImagesSrc { get; private set; }

        public List<int> TilesHandList { get; }= new List<int>();
        public List<int> TilesFlowerList { get; } = new List<int>();
        public List<MsgMeldTile> MeldList { get; } = new List<MsgMeldTile>();
        public int BankerChairId { get; private set; }

        public bool IsBandker => BankerChairId == MyPlayer.ChairId;

        public MainWindow MyOwner { get; private set; }
        public MsgAllowPlayerAction CurrentAllowPlayerAction { get; private set; }

        public MsgAllowPlayerReAction CurrentAllowPlayerReAction { get; private set; }

        public IEnumerator<Button> GetEnumerator()
        {
            foreach (var b in ButtonsSp1)
            {
                yield return b;
            }
            foreach (var b in ButtonsSp2)
            {
                yield return b;
            }
            foreach (var b in ButtonsSp3)
            {
                yield return b;
            }
            foreach (var b in ButtonsSp4)
            {
                yield return b;
            }
            foreach (var b in ButtonsAct)
            {
                yield return b;
            }
        }

        private void InitButtonArray()
        {
            var i = 0;
            foreach (var child in Sp1.Children)
            {
                ButtonsSp1[i++] = child as Button;
            }

            i = 0;
            foreach (var child in Sp2.Children)
            {
                ButtonsSp2[i++] = child as Button;
            }

            i = 0;
            foreach (var child in Sp3.Children)
            {
                ButtonsSp3[i++] = child as Button;
            }

            i = 0;
            foreach (var child in Sp4.Children)
            {
                ButtonsSp4[i++] = child as Button;
            }

            ButtonsAct[0] = BtnAction1;
            ButtonsAct[1] = BtnAction2;
            ButtonsAct[2] = BtnAction3;
            ButtonsAct[3] = BtnAction4;
            ButtonsAct[4] = BtnAction5;
        }

        public void SetImageSrc(Dictionary<int, BitmapImage> imageDict, MainWindow owner)
        {
            MyOwner = owner;
            ImagesSrc = imageDict;

            var i = 0;
            foreach (var button in ButtonsSp1)
            {
                button.Content = new Image() {Source = ImagesSrc[i++]};
            }

            i = 0;
            foreach (var button in ButtonsSp2)
            {
                button.Content = new Image() {Source = ImagesSrc[i++]};
            }
            i = 0;
            foreach (var button in ButtonsSp3)
            {
                button.Content = new Image() { Source = ImagesSrc[i++] };
            }
            i = 0;
            foreach (var button in ButtonsSp4)
            {
                button.Content = new Image() { Source = ImagesSrc[i++] };
            }
        }

        private void OnAction1_Button_Click(object sender, RoutedEventArgs e)
        {
            //enumActionType_KONG_Concealed
            // enumActionType_KONG_Exposed

            var button = sender as Button;
            if (button == null)
            {
                return;
            }

            var action = (int)button.Tag;
            var tile1 = -1;
            List<MsgMeldTile> meldList;
            switch (action)
            {
                case (int)ActionType.enumActionType_KONG_Concealed:
                    meldList = CurrentAllowPlayerAction.meldsForAction.Select(x => x).Where(x => x.meldType == (int)MeldType.enumMeldTypeConcealedKong).ToList();
                    if (!ChowPongKongWnd.ShowDialog(meldList, out tile1, this))
                    {
                        return;
                    }
                    OnTakeActionKongConcealedTile(tile1);
                    break;
                case (int)ActionType.enumActionType_KONG_Exposed:
                    meldList = CurrentAllowPlayerReAction.meldsForAction.Select(x => x).Where(x => x.meldType == (int)MeldType.enumMeldTypeExposedKong).ToList();
                    if (!ChowPongKongWnd.ShowDialog(meldList, out tile1, this))
                    {
                        return;
                    }
                    OnTakeActionKongExposedTile(tile1);
                    break;
            }
            HideAllActionButtons();
        }
        private void OnAction2_Button_Click(object sender, RoutedEventArgs e)
        {
            //enumActionType_KONG_Triplet2
            //enumActionType_CHOW
            //throw new NotImplementedException();

            var button = sender as Button;
            if (button == null)
            {
                return;
            }

            var action = (int)button.Tag;
            var tile1 = -1;
            List<MsgMeldTile> meldList;
            switch (action)
            {
                case (int)ActionType.enumActionType_KONG_Triplet2:
                    meldList = CurrentAllowPlayerAction.meldsForAction.Select(x => x).Where(x => x.meldType == (int)MeldType.enumMeldTypeTriplet2Kong).ToList();
                    if (!ChowPongKongWnd.ShowDialog(meldList, out tile1, this))
                    {
                        return;
                    }
                    OnTakeActionKong2TripletTile(tile1);
                    break;
                case (int)ActionType.enumActionType_CHOW:
                    meldList = CurrentAllowPlayerReAction.meldsForAction.Select(x => x).Where(x => x.meldType == (int)MeldType.enumMeldTypeSequence).ToList();
                    if (!ChowPongKongWnd.ShowDialog(meldList, out tile1, this))
                    {
                        return;
                    }
                    OnTakeActionChowTile(tile1);
                    break;
            }

            HideAllActionButtons();
        }

        private void OnAction3_Button_Click(object sender, RoutedEventArgs e)
        {
            //enumActionType_WIN_SelfDrawn
            //enumActionType_PONG
            //throw new NotImplementedException();
            var button = sender as Button;
            if (button == null)
            {
                return;
            }

            var action = (int)button.Tag;
            switch (action)
            {
                case (int)ActionType.enumActionType_WIN_SelfDrawn:
                    OnTakeActionWinSelfDraw();
                    break;
                case (int)ActionType.enumActionType_PONG:
                    var tile1 = -1;
                    var meldList = CurrentAllowPlayerReAction.meldsForAction.Select(x => x).Where(x => x.meldType == (int)MeldType.enumMeldTypeTriplet).ToList();
                    if (!ChowPongKongWnd.ShowDialog(meldList, out tile1, this))
                    {
                        return;
                    }
                    OnTakeActionPongTile(tile1);
                    break;
            }

            HideAllActionButtons();
        }

        private void OnAction4_Button_Click(object sender, RoutedEventArgs e)
        {
            //enumActionType_WIN_FirstReadyHand
            //enumActionType_SKIP
            //throw new NotImplementedException();
            var button = sender as Button;
            if (button == null)
            {
                return;
            }

            var action = (int)button.Tag;
            switch (action)
            {
                case (int)ActionType.enumActionType_FirstReadyHand:
                    OnNonBankerTakeActionFirstHand(0);
                    break;
                case (int)ActionType.enumActionType_SKIP:
                    OnTakeActionSkip();
                    break;
            }

            HideAllActionButtons();
        }

        private void OnAction5_Button_Click(object sender, RoutedEventArgs e)
        {
            //enumActionType_DISCARD
            //enumActionType_WIN_Chuck
            //throw new NotImplementedException();
            var button = sender as Button;
            if (button == null)
            {
                return;
            }

            var action = (int)button.Tag;
            switch (action)
            {
                case (int)ActionType.enumActionType_FirstReadyHand:
                    if (!RichiWnd.ShowDialog(CurrentAllowPlayerAction.tipsForAction[0], this))
                    {
                        return;
                    }
                    OnNonBankerTakeActionFirstHand(1);
                    break;
                case (int)ActionType.enumActionType_DISCARD:
                    var tile1 = -1;
                   
                    if (IsBandkerReadyHand || expectedReadyHandFlags != 0)
                    {
                        int readyHandFlags;
                        if (!DiscardWnd.ShowDialog(CurrentAllowPlayerAction.tipsForAction, out tile1, out readyHandFlags, expectedReadyHandFlags, this))
                        {
                            return;
                        }

                        IsBandkerReadyHand = false;
                        OnTakeActionDiscardTile(tile1, readyHandFlags);
                    }
                    else
                    {
                        if (!DiscardWnd.ShowDialog(CurrentAllowPlayerAction.tipsForAction, out tile1, this))
                        {
                            return;
                        }

                        OnTakeActionDiscardTile(tile1);
                    }

                    break;
                case (int)ActionType.enumActionType_WIN_Chuck:
                    OnTakeActionWinChuck();
                    break;
            }

            HideAllActionButtons();
        }

        public bool IsBandkerReadyHand { get; set; }

        private void OnTakeActionWinSelfDraw()
        {
            var msgAction = new MsgPlayerAction
            {
                qaIndex = CurrentAllowPlayerAction.qaIndex,
                action = (int)ActionType.enumActionType_WIN_SelfDrawn,
            };

            MyPlayer.SendMessage((int)MessageCode.OPAction, msgAction.ToBytes());
            //MyOwner.AppendActionLog($"[winselfdraw]({MyPlayer.Name}),({CurrentAllowPlayerAction.qaIndex})");
        }

        private void OnTakeActionFinalDraw()
        {
            var msgAction = new MsgPlayerAction
            {
                qaIndex = CurrentAllowPlayerAction.qaIndex,
                //action = (int)ActionType.enumActionType_AccumulateWin,
            };

            MyPlayer.SendMessage((int)MessageCode.OPAction, msgAction.ToBytes());
            //MyOwner.AppendActionLog($"[winselfdraw]({MyPlayer.Name}),({CurrentAllowPlayerAction.qaIndex})");
        }

        private void OnTakeActionKong2TripletTile(int tile1)
        {
            var msgAction = new MsgPlayerAction
            {
                qaIndex = CurrentAllowPlayerAction.qaIndex,
                action = (int)ActionType.enumActionType_KONG_Triplet2,
                tile = tile1,
            };

            var sb = new StringBuilder();
            for (var i = 0; i < 4; i++)
            {
                var tileId = tile1;
                sb.Append($"{MyOwner.TileId2Name(tileId)},");
            }

            //MyOwner.AppendActionLog($"[triplet2kong]({MyPlayer.Name}):{sb}({CurrentAllowPlayerAction.qaIndex})");

            MyPlayer.SendMessage((int)MessageCode.OPAction, msgAction.ToBytes());
        }

        private void OnTakeActionChowTile(int tile1)
        {
            var msgAction = new MsgPlayerAction
            {
                qaIndex = CurrentAllowPlayerReAction.qaIndex,
                action = (int)ActionType.enumActionType_CHOW,
                tile = CurrentAllowPlayerReAction.victimTileID,
            };

            msgAction.meldType = (int) MeldType.enumMeldTypeSequence;
            msgAction.meldTile1 = tile1;

            var sb = new StringBuilder();
            for (var i = 0; i < 3; i++)
            {
                var tileId = tile1 + i;
                sb.Append($"{MyOwner.TileId2Name(tileId)},");
            }

            //MyOwner.AppendActionLog($"[chow]({MyPlayer.Name}):{sb}({CurrentAllowPlayerReAction.qaIndex})");

            MyPlayer.SendMessage((int)MessageCode.OPAction, msgAction.ToBytes());
        }

        private void OnTakeActionPongTile(int tile1)
        {
            var msgAction = new MsgPlayerAction
            {
                qaIndex = CurrentAllowPlayerReAction.qaIndex,
                action = (int)ActionType.enumActionType_PONG,
                tile = tile1,
            };

            msgAction.meldType = (int) MeldType.enumMeldTypeTriplet;
            msgAction.meldTile1 = tile1;

            var sb = new StringBuilder();
            for (var i = 0; i < 3; i++)
            {
                var tileId = tile1;
                sb.Append($"{MyOwner.TileId2Name(tileId)},");
            }

            //MyOwner.AppendActionLog($"[pong]({MyPlayer.Name}):{sb}({CurrentAllowPlayerReAction.qaIndex})");
            MyPlayer.SendMessage((int)MessageCode.OPAction, msgAction.ToBytes());
        }

        private void OnTakeActionKongExposedTile(int tile1)
        {
            var msgAction = new MsgPlayerAction
            {
                qaIndex = CurrentAllowPlayerReAction.qaIndex,
                action = (int)ActionType.enumActionType_KONG_Exposed,
                tile = tile1,
            };

            msgAction.meldType = (int) MeldType.enumMeldTypeExposedKong;
            msgAction.meldTile1 = tile1;

            var sb = new StringBuilder();
            for (var i = 0; i < 4; i++)
            {
                var tileId = tile1;
                sb.Append($"{MyOwner.TileId2Name(tileId)},");
            }

            //MyOwner.AppendActionLog($"[kongExposed]({MyPlayer.Name}):{sb}({CurrentAllowPlayerReAction.qaIndex})");
            MyPlayer.SendMessage((int)MessageCode.OPAction, msgAction.ToBytes());
        }

        private void OnTakeActionKongConcealedTile(int tile1)
        {
            var msgAction = new MsgPlayerAction
            {
                qaIndex = CurrentAllowPlayerAction.qaIndex,
                action = (int)ActionType.enumActionType_KONG_Concealed,
                tile = tile1,
            };

            //msgAction.actionMeld = new MsgMeldTile() { meldType = (int)MeldType.enumMeldTypeConcealedKong, tile1 = tile1 };
            var sb = new StringBuilder();
            for (var i = 0; i < 4; i++)
            {
                var tileId = tile1;
                sb.Append($"{MyOwner.TileId2Name(tileId)},");
            }

            //MyOwner.AppendActionLog($"[kongConcealed]({MyPlayer.Name}):{sb}({CurrentAllowPlayerAction.qaIndex})");
            MyPlayer.SendMessage((int)MessageCode.OPAction, msgAction.ToBytes());
        }

        private void OnTakeActionSkip()
        {
            var qaIndex = 0;
            if (CurrentAllowPlayerAction != null)
            {
                qaIndex = CurrentAllowPlayerAction.qaIndex;
            } else if (CurrentAllowPlayerReAction != null)
            {
                qaIndex = CurrentAllowPlayerReAction.qaIndex;
            }
            var msgAction = new MsgPlayerAction
            {
                qaIndex = qaIndex,
                action = (int)ActionType.enumActionType_SKIP,
            };

            MyPlayer.SendMessage((int)MessageCode.OPAction, msgAction.ToBytes());

            //MyOwner.AppendActionLog($"[skip]({MyPlayer.Name}),({CurrentAllowPlayerReAction.qaIndex})");
        }

        private void OnNonBankerTakeActionFirstHand(int tid)
        {
            var msgAction = new MsgPlayerAction
            {
                qaIndex = CurrentAllowPlayerAction.qaIndex,
                action = (int)ActionType.enumActionType_FirstReadyHand,
                tile = tid, // 1表示听牌，0表示不听牌
                flags = tid,
            };

            MyPlayer.SendMessage((int)MessageCode.OPAction, msgAction.ToBytes());
            //MyOwner.AppendActionLog($"[richi]({MyPlayer.Name}),({CurrentAllowPlayerAction.qaIndex}),({tid})");
        }

        private void OnTakeActionWinChuck()
        {
            var msgAction = new MsgPlayerAction
            {
                qaIndex = CurrentAllowPlayerReAction.qaIndex,
                action = (int)ActionType.enumActionType_WIN_Chuck,
                tile = CurrentAllowPlayerReAction.victimTileID
            };

            MyPlayer.SendMessage((int)MessageCode.OPAction, msgAction.ToBytes());
            //MyOwner.AppendActionLog($"[winchuck]({MyPlayer.Name}),({CurrentAllowPlayerReAction.qaIndex})");
        }

        private void OnTakeActionDiscardTile(int tile2Discarded)
        {
            var msgAction = new MsgPlayerAction
            {
                qaIndex = CurrentAllowPlayerAction.qaIndex,
                action = (int) ActionType.enumActionType_DISCARD,
                tile = tile2Discarded
            };

            MyPlayer.SendMessage((int) MessageCode.OPAction, msgAction.ToBytes());
            //MyOwner.AppendActionLog($"[discard]({MyPlayer.Name}):{MyOwner.TileId2Name(tile2Discarded)},({CurrentAllowPlayerAction.qaIndex})");
        }

        private void OnTakeActionDiscardTile(int tile2Discarded, int readyHandFlags)
        {
            var msgAction = new MsgPlayerAction
            {
                qaIndex = CurrentAllowPlayerAction.qaIndex,
                action = (int)ActionType.enumActionType_DISCARD,
                tile = tile2Discarded,
                flags = readyHandFlags,
            };
            // var ix = readyHand ? 1 : 0;
            MyPlayer.SendMessage((int)MessageCode.OPAction, msgAction.ToBytes());
            //MyOwner.AppendActionLog($"[discard]({MyPlayer.Name}):{MyOwner.TileId2Name(tile2Discarded)},({CurrentAllowPlayerAction.qaIndex}),({ix})");
        }

        private void HideAllButtons()
        {
            foreach (var button in this)
            {
                button.Visibility = Visibility.Hidden;
            }
        }

        public void ResetPlayStatus()
        {
            //throw new NotImplementedException();
            // hide all
            Reset2New();
        }

        public void OnDeal(MsgDeal msg)
        {
            Reset2New();

            MsgPlayerTileList myPlayList = null;
            foreach (var ptl in msg.playerTileLists)
            {
                if (ptl.chairID == MyPlayer.ChairId)
                {
                    myPlayList = ptl;
                    break;
                }
            }

            if (myPlayList == null)
                return;

            if (myPlayList.tilesHand.Count < 1)
                return;

            // 庄家标记
            TbName.Text = msg.bankerChairID == MyPlayer.ChairId ? $"{MyPlayer.Name}(庄)" : $"{MyPlayer.Name}";
            BankerChairId = msg.bankerChairID;
            TbScore.Text = "";

            if (!MyOwner.IsPlaying)
            {
                MyOwner.IsPlaying = true;
                MyOwner.ClearLog();
                MyOwner.AppendLog("[begin]\r\n");
                MyOwner.ResetActionListWndIndex();
            }

            if (IsBandker)
            {
                TbPseudoFlower.Text = MyOwner.TileId2Name(msg.windFlowerID);
                MyOwner.AppendLog($"[bank]:{MyPlayer.Name}\r\n");
                MyOwner.AppendLog($"[wind]:{MyOwner.TileId2Name(msg.windFlowerID)}\r\n");
                MyOwner.TbTileInWallRemain.Text = msg.tilesInWall.ToString();
                MyOwner.ResetScoreWnd();
            }

            // 手牌列表
            TilesHandList.AddRange(myPlayList.tilesHand);
            Hand2Buttons();

            // 花牌列表
            TilesFlowerList.AddRange(myPlayList.tilesFlower);
            Flower2Buttons();

            var sb = new StringBuilder();
            sb.Append($"[deal]({MyPlayer.Name})(hand):");
            foreach (var i in TilesHandList)
            {
                sb.Append(MyOwner.TileId2Name(i));
                sb.Append(",");
            }
            sb.AppendLine();
            sb.Append("\t(flower):");
            foreach (var i in TilesFlowerList)
            {
                sb.Append(MyOwner.TileId2Name(i));
                sb.Append(",");
            }
            sb.AppendLine();
            MyOwner.AppendLog(sb.ToString());
        }

        private void Flower2Buttons()
        {
            SortFlowerTiles();
            var i = 0;
            foreach (var t in TilesFlowerList)
            {
                if (i >= ButtonsSp3.Length)
                {
                    return;
                }

                ButtonsSp3[i].Content = new Image() {Source = ImagesSrc[t]};
                ButtonsSp3[i].Visibility = Visibility.Visible;
                ++i;
            }
        }

        private void Hand2Buttons()
        {
            HideSp4Buttons();
            SortHandTiles();
            var i = 0;
            foreach (var t in TilesHandList)
            {
                ButtonsSp4[i].Content = new Image() {Source = ImagesSrc[t]};
                ButtonsSp4[i].Visibility = Visibility.Visible;
                ++i;
            }
        }

        private void SortHandTiles()
        {
            TilesHandList.Sort((x, y) => x - y);
        }

        private void SortFlowerTiles()
        {
            TilesFlowerList.Sort((x, y) => x - y);
        }
        private void HideSp4Buttons()
        {
            foreach (var button in ButtonsSp4)
            {
                button.Visibility = Visibility.Hidden;
            }
        }

        private void HideSp2Buttons()
        {
            foreach (var button in ButtonsSp2)
            {
                button.Visibility = Visibility.Hidden;
            }
        }
        public void Reset2New()
        {
            HideAllButtons();
            TilesFlowerList.Clear();
            TilesHandList.Clear();
            MeldList.Clear();
            BankerChairId = 0;
            TbRichi.Text = "";
            TbPseudoFlower.Text = "";
        }

        internal void SetPlayer(Player player)
        {
            MyPlayer = player;
            TbName.Text = player.Name;

            Visibility = Visibility.Visible;
        }

        public Player MyPlayer { get; set; }

        public void OnAllowedReActions(MsgAllowPlayerReAction msg)
        {
            HideAllActionButtons();
            CurrentAllowPlayerReAction = msg;
            CurrentAllowPlayerAction = null;

            var actions = msg.allowedActions;


            if ((actions & (int)ActionType.enumActionType_KONG_Exposed) != 0)
            {
                BtnAction1.Visibility = Visibility.Visible;
                BtnAction1.Content = "明杠";
                BtnAction1.Tag = (int)ActionType.enumActionType_KONG_Exposed;
            }

            if ((actions & (int)ActionType.enumActionType_CHOW) != 0)
            {
                BtnAction2.Visibility = Visibility.Visible;
                BtnAction2.Content = "吃";
                BtnAction2.Tag = (int)ActionType.enumActionType_CHOW;
            }

            if ((actions & (int)ActionType.enumActionType_PONG) != 0)
            {
                BtnAction3.Visibility = Visibility.Visible;
                BtnAction3.Content = "碰";
                BtnAction3.Tag = (int)ActionType.enumActionType_PONG;
            }

            if ((actions & (int)ActionType.enumActionType_SKIP) != 0)
            {
                BtnAction4.Visibility = Visibility.Visible;
                BtnAction4.Content = "过";
                BtnAction4.Tag = (int)ActionType.enumActionType_SKIP;
            }

            if ((actions & (int)ActionType.enumActionType_WIN_Chuck) != 0)
            {
                BtnAction5.Visibility = Visibility.Visible;
                BtnAction5.Content = "胡";
                BtnAction5.Tag = (int)ActionType.enumActionType_WIN_Chuck;
            }

            if (MyOwner.CheckBoxAutoAction.IsChecked == false && IsAutoX)
            {
                // 自动打牌
                if ((actions & (int)ActionType.enumActionType_SKIP) != 0)
                {
                    OnTakeActionSkip();

                    HideAllActionButtons();
                }
            }
        }

        public void OnAllowedActions(MsgAllowPlayerAction msg)
        {
            HideAllActionButtons();
            CurrentAllowPlayerAction = msg;
            CurrentAllowPlayerReAction = null;


            var actions = msg.allowedActions;

            if ((actions & (int)ActionType.enumActionType_KONG_Concealed) != 0)
            {
                BtnAction1.Visibility = Visibility.Visible;
                BtnAction1.Content = "暗杠";
                BtnAction1.Tag = (int)ActionType.enumActionType_KONG_Concealed;
            }
            if ((actions & (int)ActionType.enumActionType_KONG_Triplet2) != 0)
            {
                BtnAction2.Visibility = Visibility.Visible;
                BtnAction2.Content = "加杠";
                BtnAction2.Tag = (int)ActionType.enumActionType_KONG_Triplet2;
            }
            if ((actions & (int)ActionType.enumActionType_WIN_SelfDrawn) != 0)
            {
                BtnAction3.Visibility = Visibility.Visible;
                BtnAction3.Content = "胡";
                BtnAction3.Tag = (int)ActionType.enumActionType_WIN_SelfDrawn;
            }
            if ((actions & (int)ActionType.enumActionType_FirstReadyHand) != 0)
            {
                if (!IsBandker)
                {
                    BtnAction5.Visibility = Visibility.Visible;
                    BtnAction5.Content = "听牌";
                    BtnAction5.Tag = (int) ActionType.enumActionType_FirstReadyHand;

                    BtnAction4.Visibility = Visibility.Visible;
                    BtnAction4.Content = "过";
                    BtnAction4.Tag = (int)ActionType.enumActionType_FirstReadyHand;
                    return;
                }
                else
                {
                    // 庄家起手听
                    IsBandkerReadyHand = true;
                }
            }

            if ((actions & (int)ActionType.enumActionType_SKIP) != 0)
            {
                BtnAction4.Visibility = Visibility.Visible;
                BtnAction4.Content = "过";
                BtnAction4.Tag = (int)ActionType.enumActionType_SKIP;
            }

            //if ((actions & (int)ActionType.enumActionType_AccumulateWin) != 0)
            //{
            //    BtnAction5.Visibility = Visibility.Visible;
            //    BtnAction5.Content = "抓";
            //    BtnAction5.Tag = (int)ActionType.enumActionType_AccumulateWin;
            //}

            expectedReadyHandFlags = 0;
            if ((actions & (int)ActionType.enumActionType_ReadyHand) != 0 || (actions & (int)ActionType.enumActionType_FirstReadyHand) != 0)
            {
                expectedReadyHandFlags = 1;
            }

            //if ((actions & (int)ActionType.enumActionType_FlyReadyHand) != 0)
            //{
            //    expectedReadyHandFlags |= 2;
            //}

            if ((actions & (int)ActionType.enumActionType_DISCARD) != 0)
            {
                BtnAction5.Visibility = Visibility.Visible;
                BtnAction5.Content = "出牌";
                BtnAction5.Tag = (int)ActionType.enumActionType_DISCARD;
                
                if (MyOwner.CheckBoxAutoAction.IsChecked == false && IsAutoX)
                {
                    var handTips = msg.tipsForAction;
                    // 自动打牌
                    if ((actions & (int)ActionType.enumActionType_FirstReadyHand) != 0)
                    {
                        if (!IsBandker)
                        {
                            // 绝对不听牌
                            OnNonBankerTakeActionFirstHand(0);
                        }
                        else
                        {
                            OnTakeActionDiscardTile(handTips[0].targetTile, 0);
                        }
                    }
                    else if ((actions & (int)ActionType.enumActionType_DISCARD) != 0)
                    {
                        OnTakeActionDiscardTile(handTips[0].targetTile);
                    }
                    else if ((actions & (int)ActionType.enumActionType_SKIP) != 0)
                    {
                        OnTakeActionSkip();
                    }

                    HideAllActionButtons();
                }
                
            }
        }

        private void HideAllActionButtons()
        {
            foreach (var button in ButtonsAct)
            {
                button.Visibility = Visibility.Hidden;
            }
        }

        public void OnActionResult(MsgActionResultNotify msg)
        {
            // 抽牌
            if (msg.action == (int) ActionType.enumActionType_DRAW)
            {
                if (msg.actionTile != (int) TileID.enumTid_MAX + 1)
                {
                    TilesHandList.Add(msg.actionTile);

                    MyOwner.BtnDraw.Content = new Image() { Source = ImagesSrc[msg.actionTile] };
                    MyOwner.TbDraw.Text = $"{MyPlayer.Name}<-{MyOwner.TileId2Name(msg.actionTile)}";
                    // log action
                    //MyOwner.AppendLog($"[draw]({MyPlayer.Name}):{MyOwner.TileId2Name(msg.actionTile)}\r\n");
                }

                TilesFlowerList.AddRange(msg.newFlowers);
                Hand2Buttons();
                Flower2Buttons();

                //foreach (var flower in msg.newFlowers)
                //{
                    //MyOwner.AppendLog($"[draw]({MyPlayer.Name}):{MyOwner.TileId2Name(flower)}\r\n");
                //}
                MyOwner.TbTileInWallRemain.Text = msg.tilesInWall.ToString();
                return;
            }

            // 出牌
            if (msg.action == (int) ActionType.enumActionType_DISCARD)
            {
                TilesHandList.Remove(msg.actionTile);
                Hand2Buttons();

                MyOwner.BtnDiscard.Content = new Image() {Source = ImagesSrc[msg.actionTile]};
                MyOwner.TbDiscard.Text = $"{MyPlayer.Name}->{MyOwner.TileId2Name(msg.actionTile)}";
                //MyOwner.AppendActionLog($"[discard]({MyPlayer.Name}):{MyOwner.TileId2Name(msg.actionTile)}");
                return;
            }

            // 吃
            if (msg.action == (int)ActionType.enumActionType_CHOW)
            {
                for (var i = 0; i < 3; i++)
                {
                    var tileId = msg.actionMeld.tile1 + i;
                    if (tileId != msg.actionTile)
                        TilesHandList.Remove(tileId);
                }

                MeldList.Add(msg.actionMeld);

                Hand2Buttons();
                MeldList2Buttons();

                //var sb = new StringBuilder();
                //for (var i = 0; i < 3; i++)
                //{
                //    var tileId = msg.actionMeld.tile1 + i;
                //    sb.Append($"{MyOwner.TileId2Name(tileId)},");
                //}

                //MyOwner.AppendActionLog($"[chow]({MyPlayer.Name}):{sb}");
                return;
            }

            // 碰
            if (msg.action == (int)ActionType.enumActionType_PONG)
            {
                for (var i = 0; i < 2; i++)
                {
                    TilesHandList.Remove(msg.actionMeld.tile1);
                }
                MeldList.Add(msg.actionMeld);

                Hand2Buttons();
                MeldList2Buttons();

                //var sb = new StringBuilder();
                //for (var i = 0; i < 3; i++)
                //{
                //    var tileId = msg.actionMeld.tile1;
                //    sb.Append($"{MyOwner.TileId2Name(tileId)},");
                //}

                //MyOwner.AppendActionLog($"[pong]({MyPlayer.Name}):{sb}");
                return;
            }

            // 明杠
            if (msg.action == (int)ActionType.enumActionType_KONG_Exposed)
            {
                for (var i = 0; i < 3; i++)
                {
                    TilesHandList.Remove(msg.actionMeld.tile1);
                }

                MeldList.Add(msg.actionMeld);

                Hand2Buttons();
                MeldList2Buttons();

                //var sb = new StringBuilder();
                //for (var i = 0; i < 4; i++)
                //{
                //    var tileId = msg.actionMeld.tile1;
                //    sb.Append($"{MyOwner.TileId2Name(tileId)},");
                //}

                //MyOwner.AppendActionLog($"[kongExposed]({MyPlayer.Name}):{sb}");
                return;
            }

            // 暗杠
            if (msg.action == (int)ActionType.enumActionType_KONG_Concealed)
            {
                for (var i = 0; i < 4; i++)
                {
                    TilesHandList.Remove(msg.actionTile);
                }

                MeldList.Add(new MsgMeldTile() {meldType = (int)MeldType.enumMeldTypeConcealedKong, tile1 = msg.actionTile, contributor = MyPlayer.ChairId});

                Hand2Buttons();
                MeldList2Buttons();

                //var sb = new StringBuilder();
                //for (var i = 0; i < 4; i++)
                //{
                //    var tileId = msg.actionTile;
                //    sb.Append($"{MyOwner.TileId2Name(tileId)},");
                //}

                //MyOwner.AppendActionLog($"[kongConcealed]({MyPlayer.Name}):{sb}");
                return;
            }

            // 加杠
            if (msg.action == (int)ActionType.enumActionType_KONG_Triplet2)
            {
                for (var i = 0; i < 1; i++)
                {
                    TilesHandList.Remove(msg.actionTile);
                }

                var meld = MeldList.Find((x)=>x.meldType==(int)MeldType.enumMeldTypeTriplet&&x.tile1==msg.actionTile);
                meld.meldType = (int) MeldType.enumMeldTypeTriplet2Kong;
                //MeldList.Add(meld);

                Hand2Buttons();
                MeldList2Buttons();

                //var sb = new StringBuilder();
                //for (var i = 0; i < 4; i++)
                //{
                //    var tileId = meld.tile1;
                //    sb.Append($"{MyOwner.TileId2Name(tileId)},");
                //}

                //MyOwner.AppendActionLog($"[triplet2kong]({MyPlayer.Name}):{sb}");
                return;
            }

            // 起手听
            TbRichi.Text = "听";
            //MyOwner.AppendActionLog($"[richi]({MyPlayer.Name}),idx:{CurrentAllowPlayerAction.qaIndex}");
        }

        private void MeldList2Buttons()
        {
            HideSp2Buttons();
            HideSp1Buttons();

            var i = 0;
            foreach (var meld in MeldList)
            {
                if (i > 13)
                {
                    break;
                }

                if (meld.meldType == (int)MeldType.enumMeldTypeTriplet2Kong
                    || meld.meldType == (int)MeldType.enumMeldTypeConcealedKong
                    || meld.meldType == (int)MeldType.enumMeldTypeExposedKong)
                {
                    for (var j = 0; j < 4; j++)
                    {
                        var btn = ButtonsSp2[i];
                        btn.Content = new Image() { Source = ImagesSrc[meld.tile1] };
                        btn.Tag = meld.tile1;

                        btn.Visibility = Visibility.Visible;
                        if (j == 0)
                        {
                            SetContributor(ButtonsSp1[i+1], meld.contributor);
                            SetMeldFlag(ButtonsSp1[i], meld.meldType);
                        }
                        i++;

                        if (i > 13)
                        {
                            break;
                        }
                    }

                }
                else if (meld.meldType == (int)MeldType.enumMeldTypeTriplet)
                {
                    for (var j = 0; j < 3; j++)
                    {
                        var btn = ButtonsSp2[i];
                        btn.Content = new Image() { Source = ImagesSrc[meld.tile1] };
                        btn.Tag = meld.tile1;

                        btn.Visibility = Visibility.Visible;
                        if (j == 0)
                        {
                            SetContributor(ButtonsSp1[i + 1], meld.contributor);
                            SetMeldFlag(ButtonsSp1[i], meld.meldType);
                        }
                        i++;

                        if (i > 13)
                        {
                            break;
                        }
                    }
                }
                else if (meld.meldType == (int)MeldType.enumMeldTypeSequence)
                {
                    for (var j = 0; j < 3; j++)
                    {
                        var btn = ButtonsSp2[i];
                        btn.Content = new Image() { Source = ImagesSrc[meld.tile1 + j] };
                        btn.Tag = meld.tile1;

                        btn.Visibility = Visibility.Visible;

                        if (j == 0)
                        {
                            SetContributor(ButtonsSp1[i + 1], meld.contributor);
                            SetMeldFlag(ButtonsSp1[i], meld.meldType);
                        }
                        
                        i++;

                        if (i > 13)
                        {
                            break;
                        }
                    }

                }
            }
        }

        private void SetMeldFlag(Button button, int meldMeldType)
        {
            if (meldMeldType == (int) MeldType.enumMeldTypeTriplet2Kong)
            {
                button.Content = "加";
            }
            else if (meldMeldType == (int) MeldType.enumMeldTypeConcealedKong)
            {
                button.Content = "暗";
            }
            else if (meldMeldType == (int) MeldType.enumMeldTypeExposedKong)
            {
                button.Content = "明";
            }
            else if (meldMeldType == (int)MeldType.enumMeldTypeTriplet)
            {
                button.Content = "碰";
            }
            else if (meldMeldType == (int)MeldType.enumMeldTypeSequence)
            {
                button.Content = "吃";
            }

            button.Visibility = Visibility.Visible;
        }

        private void HideSp1Buttons()
        {
            foreach (var button in ButtonsSp1)
            {
                button.Visibility = Visibility.Hidden;
            }
        }

        public void SetContributor(Button btn, int contr)
        {
            var chair2Name = new []
            {
                "A",
                "B",
                "C",
                "D"
            };
            btn.Content = chair2Name[contr];
            btn.Visibility = Visibility.Visible;
        }

        public void CancelAllowedAction()
        {
            HideAllActionButtons();
        }

        public void OnHandScore(MsgHandOver msg)
        {
            if (MyOwner.IsPlaying)
            {
                MyOwner.IsPlaying = false;
                MyOwner.AppendLog("[end]\r\n");
                var handoverType = msg.endType;
                switch (handoverType)
                {
                    case (int)HandOverType.enumHandOverType_None:
                        MyOwner.AppendLog("流局\r\n");
                        break;
                    case (int)HandOverType.enumHandOverType_Win_SelfDrawn:
                        MyOwner.AppendLog("自摸胡牌\r\n");
                        break;
                    case (int)HandOverType.enumHandOverType_Win_Chuck:
                        MyOwner.AppendLog("放铳胡牌\r\n");
                        break;
                }
            }

            var handScore = msg.scores;
            if (handScore == null)
            {
                return;
            }

            var myScore = handScore.playerScores.FirstOrDefault(playerScore => playerScore.targetChairID == MyPlayer.ChairId);
            if (myScore == null)
                return;


            TbScore.Text = myScore.score.ToString();

            if (MyOwner.IsFirstPlayer(MyPlayer))
            {
                var scoreMsg = FormatScore(msg);
                MyOwner.ShowScoreWnd(scoreMsg);
            }
        }

        private string FormatScore(MsgHandOver msg)
        {
            var sb = new StringBuilder();

            // 结束类型
            sb.AppendLine($"{Enum2StrHelper.EndType2String(msg.endType)}");

            var handScore = msg.scores;
            // 每个玩家的得分和赢牌类型
            foreach (var playerScore in handScore.playerScores)
            {
                sb.Append(Enum2StrHelper.ChairId2Name(playerScore.targetChairID));
                sb.Append(":");
                sb.Append(playerScore.score);
                if (playerScore.winType != 0)
                {
                    sb.Append(",");
                    sb.Append(Enum2StrHelper.WinType2String(playerScore.winType));
                }
                
                sb.AppendLine();
            }

            foreach (var playerScore in handScore.playerScores)
            {
                if (playerScore.fakeList.Count > 0)
                {
                    sb.AppendLine();
                    sb.AppendLine("马牌：");
                    foreach(var t in playerScore.fakeList)
                    {
                        sb.Append($"{MyOwner.IdNames[t]},");
                    }
                    sb.AppendLine();
                    break;
                }
            }

            // 每个玩家得分详细信息
            sb.AppendLine("------------------Details:---------------------");
            foreach (var playerScore in handScore.playerScores)
            {
                sb.AppendLine($"名字：{Enum2StrHelper.ChairId2Name(playerScore.targetChairID)}");
                sb.AppendLine($"得分：{playerScore.score}");
                sb.AppendLine($"得分类型：{Enum2StrHelper.WinType2String(playerScore.winType)}");

                if (playerScore.winType != (int) HandOverType.enumHandOverType_None
                    && playerScore.winType != (int)HandOverType.enumHandOverType_Chucker)
                {
                    if (playerScore.greatWin != null)
                    {
                        var greatWin = playerScore.greatWin;
                        sb.AppendLine($"应得分数：{greatWin.baseWinScore}");
                        sb.AppendLine($"牌型倍数：{greatWin.greatWinPoints}");
                        sb.AppendLine($"中马个数：{playerScore.specialScore}");
                    }
                }

                sb.AppendLine();
            }

            return sb.ToString();
        }

        public void OnEnterRoom(MsgEnterRoomResult msg)
        {
            if (msg.status != (int) EnterRoomStatus.Success)
            {
                var x = (mahjong.EnterRoomStatus)msg.status;
                MessageBox.Show($"enter room failed:{x.ToString()}");
                return;
            }

            MyPlayer.SendMessage((int)MessageCode.OPPlayerReady, null);
        }

        public void OnDisbandNotify(MsgDisbandNotify msg)
        {
            if (msg.waits != null)
            {
                var me = msg.waits.Any((x) => x == MyPlayer.ChairId);
                if (me)
                {
                    var result = MessageBox.Show(MyOwner, "有人请求解散房间，是否同意？", "解散房间询问", MessageBoxButton.YesNo);
                    var agree = result == MessageBoxResult.Yes;


                    var msgAnswer = new MsgDisbandAnswer();
                    msgAnswer.agree = agree;

                    MyPlayer.SendMessage((int)MessageCode.OPDisbandAnswer, msgAnswer.ToBytes());
                }
            }
        }

        public void SendReady2Server()
        {
            MyPlayer.SendMessage((int)MessageCode.OPPlayerReady, null);
        }

        public void OnShowRoomTips(MsgRoomShowTips msg)
        {
            MyOwner.AppendLog($"{MyPlayer.UserId}:  {msg.tips}\r\n");

            if (string.IsNullOrWhiteSpace(msg.tips))
            {
                return;
            }

            if (MyOwner.CheckBoxAutoAction.IsChecked == true)
            {
                DoAutoAction(msg.tips);
            }
        }

        public struct ActionParams
        {
            public string Action;
            public string TileString;
            public bool HasRichi;
            public int RichiFlags;
        }

        private void DoAutoAction(string msgTips)
        {
            ActionParams actionParams = ParseActionMsgTips(msgTips);
            var tileId = 0;
            var handled = true;
            switch (actionParams.Action)
            {
                case "richi":
                    bool isRichi = bool.Parse(actionParams.TileString);
                    OnNonBankerTakeActionFirstHand(isRichi ? 1:0);
                    break;
                case "discard":
                    tileId = MyOwner.NameIds[actionParams.TileString];
                    if (actionParams.HasRichi)
                    {
                        OnTakeActionDiscardTile(tileId, actionParams.RichiFlags);
                    }
                    else
                    {
                        OnTakeActionDiscardTile(tileId);
                    }
                    break;
                case "chow":
                    tileId = MyOwner.NameIds[actionParams.TileString];
                    OnTakeActionChowTile(tileId);
                    break;
                case "pong":
                    tileId = MyOwner.NameIds[actionParams.TileString];
                    OnTakeActionPongTile(tileId);
                    break;
                case "c-kong":
                    tileId = MyOwner.NameIds[actionParams.TileString];
                    OnTakeActionKongConcealedTile(tileId);
                    break;
                case "skip":
                    OnTakeActionSkip();
                    break;
                case "e-kong":
                    tileId = MyOwner.NameIds[actionParams.TileString];
                    OnTakeActionKongExposedTile(tileId);
                    break;
                case "t-kong":
                    tileId = MyOwner.NameIds[actionParams.TileString];
                    OnTakeActionKong2TripletTile(tileId);
                    break;
                case "winChuck":
                    OnTakeActionWinChuck();
                    break;
                case "winSelf":
                    OnTakeActionWinSelfDraw();
                    break;
                case "finalDraw":
                    OnTakeActionFinalDraw();
                    break;
                default:
                    handled = false;
                    break;
            }

            if (handled)
            {
                HideAllActionButtons();
            }
        }

        private ActionParams ParseActionMsgTips(string msgTips)
        {
            ActionParams actionParams = new ActionParams {Action = ""};

            var pattern = @"^\[(?<action>[^\s]+).*\]$";
            var rgx = new Regex(pattern, RegexOptions.IgnoreCase);
            var matches = rgx.Matches(msgTips);
            if (matches.Count > 0)
            {
                var match = matches[0];
                var action = match.Groups["action"].Value;
                actionParams.Action = action;

                var needTileN = true;
                switch (action)
                {
                    case "chow":
                        break;
                    case "pong":
                        break;
                    case "c-kong":
                        break;
                    case "e-kong":
                        break;
                    case "t-kong":
                        break;
                    case "richi":
                        break;
                    case "discard":
                        break;
                    default:
                        needTileN = false;
                        break;
                }

                if (needTileN)
                {
                    var secondPattern = @"^\[[^\s]+\s(?<paramA>[^\s]+).*\]$";
                    rgx = new Regex(secondPattern, RegexOptions.IgnoreCase);
                    matches = rgx.Matches(msgTips);
                    match = matches[0];
                    var paramA = match.Groups["paramA"].Value;
                    //Console.WriteLine(paramA);
                    actionParams.TileString = paramA;
                }

                if (action == "discard")
                {
                    var secondPattern = @"^\[.*richi\s(?<paramB>[^\s]+).*\]$";
                    rgx = new Regex(secondPattern, RegexOptions.IgnoreCase);
                    matches = rgx.Matches(msgTips);
                    if (matches.Count > 0)
                    {
                        match = matches[0];
                        var paramB = match.Groups["paramB"].Value;
                        //Console.WriteLine(paramB);
                        actionParams.HasRichi = true;
                        actionParams.RichiFlags = int.Parse(paramB);
                    }
                }
            }

            return actionParams;
        }

        private void OnAuoX_Button_Click(object sender, RoutedEventArgs e)
        {
            IsAutoX = !IsAutoX;
            if (!IsAutoX)
            {
                AutoX.Content = "AutoX";
            }
            else
            {
                AutoX.Content = "C-AutoX";
            }
        }

        public bool IsAutoX { get; set; }
    }
}
