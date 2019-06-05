using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using CsvHelper;
using pokerface;
using Path = System.IO.Path;

namespace PokerTest
{
    /// <summary>
    /// ExportRoomWnd.xaml 的交互逻辑
    /// </summary>
    public partial class ExportRoomWnd : Window
    {
        public ExportRoomWnd()
        {
            InitializeComponent();
        }

        public enum ExportRoomType
        {
            RoomCfg,
            Operations,
        }

        private ExportRoomType _exprotType;
        private MainWindow _myWindow;

        public static void ShowExportDialog(ExportRoomType type, MainWindow owner)
        {
            var wnd = new ExportRoomWnd
            {
                Owner = owner,
                _exprotType = type,
                _myWindow = owner
            };
            wnd.ShowDialog();
        }

        private async void OnExportButton_Clicked(object sender, RoutedEventArgs e)
        {
            var text = TextBoxUserId.Text;
            if (string.IsNullOrWhiteSpace(text))
            {
                MessageBox.Show("please input a valid user id");
                return;
            }

            if (cbAll.IsChecked == true)
            {
                if (true == RbUserId.IsChecked)
                {
                    MessageBox.Show("need record ID type");
                    return;
                }

                try
                {
                    await ExportAllRoomOps(text);
                    return;
                }
                catch (Exception ex)
                {
                    MessageBox.Show(ex.Message);
                }

            }

            if (_exprotType == ExportRoomType.Operations)
            {
                try
                {
                    var filename = "user-" + text + "-ops";
                    await ExportOps(text, filename);
                }
                catch (Exception ex)
                {
                    MessageBox.Show(ex.Message);
                }
            }
            else if (_exprotType == ExportRoomType.RoomCfg)
            {
                //await ExportCfg(text);
            }
        }

        private Task ExportCfg(string text)
        {
            throw new NotImplementedException();
        }


        private async Task ExportAllRoomOps(string text)
        {
            string xid = "recordSID=" + text;
            var result = await HttpHandlers.ExportRoomShareIDs(xid, this);
            if (result == null)
            {
                return;
            }

            var strArray = ParseRoomShareIDs(result);

            for (var i = 0; i < strArray.Length; i++)
            {
                string recordID = strArray[i];
                var filename = $"room-{i}-{recordID}";
                await ExportOps(recordID, filename);
            }
        }

        private string[] ParseRoomShareIDs(string strx)
        {
            var strArray = new List<string>();
            using (var strReader = new StringReader(strx))
            {
                string str;
                do
                {
                    str = strReader.ReadLine();
                    if (!string.IsNullOrWhiteSpace(str))
                    {
                        strArray.Add(str);
                    }

                } while (str != null);

            }

            return strArray.ToArray();
        }


        private async Task ExportOps(string text, string fileName)
        {
            string xid;
            if (RbUserId.IsChecked == true)
            {
                xid = "userID=" + text;
            }
            else
            {
                xid = "recordSID=" + text;
            }

            var result = await HttpHandlers.ExportRoomOps(xid, this);

            if (result == null)
            {
                return;
            }

            var recorder = result.ToProto<SRMsgHandRecorder>();

            List<PlayerData> players = FromRecorder(recorder);
            //string drawSequence = ExtractDrawSequence(recorder);
            //string kongSequence = ExtractKongSequence(recorder);

            var windId = recorder.windFlowerID;
            var bankerPlayer = players.Find((p) => p.ChairId == recorder.bankerChairID);
            var bankderUserId = bankerPlayer.UserId;
            players = sortPlayers(players, bankderUserId);
            var extra = recorder.extra;
            var markup = extra.markup;

            //名称	userID1	手牌	花牌	动作提示	userID2	手牌	花牌	动作提示	userID3	手牌	花牌	动作提示	userID4	手牌	花牌	动作提示	抽牌序列	庄家	风牌
            Microsoft.Win32.SaveFileDialog dlg = new Microsoft.Win32.SaveFileDialog();
            dlg.FileName = fileName; // Default file name

            dlg.DefaultExt = ".csv"; // Default file extension
            dlg.Filter = "CSV documents (.csv)|*.csv"; // Filter files by extension

            // Show save file dialog box
            bool? dlgResult = dlg.ShowDialog();
            // Process save file dialog box results
            if (dlgResult == true)
            {
                // Save document
                string filename = dlg.FileName;
                using (var textWriter =
                    new StreamWriter(new FileStream(filename, FileMode.Create, FileAccess.ReadWrite), Encoding.Default))
                {
                    var headers = DealCfgWnd.Headers;

                    // 第一行
                    var csv = new CsvWriter(textWriter);
                    foreach (var header in headers)
                    {
                        csv.WriteField(header);
                    }
                    csv.NextRecord();

                    // 第二行
                    csv.WriteField("bug"); // 名字
                    csv.WriteField("大丰关张"); // 名字

                    foreach (var playerData in players)
                    {
                        playerData.WriteCsv(csv);
                    }

                    int pad = 4 - players.Count();
                    for (int i = 0; i < pad; ++i)
                    {
                        for (int j = 0; j < 3; ++j)
                            csv.WriteField("");
                    }

                    //csv.WriteField(drawSequence);
                    //csv.WriteField(kongSequence);//杠后牌
                    //csv.WriteField(markup); // 上楼计数
                    csv.WriteField(1); // 强制一致
                    csv.WriteField(recorder.roomConfigID);
                    csv.WriteField(recorder.isContinuousBanker ? "1" : "0");
                    csv.NextRecord();
                }

                if (!string.IsNullOrWhiteSpace(recorder.roomConfigID))
                {
                    var dir = Path.GetDirectoryName(filename);
                    if (dir == null)
                    {
                        return;
                    }

                    var jsonFileName = Path.Combine(dir, Path.GetFileNameWithoutExtension(filename) + ".json");
                    await LoadAndSaveRoomConfig(recorder.roomConfigID, jsonFileName);
                }
            }
        }

        private async Task LoadAndSaveRoomConfig(string recorderRoomConfigId, string jsonFileName)
        {
            var result = await HttpHandlers.ExportRoomCfg(recorderRoomConfigId, this);

            if (result == null)
            {
                return;
            }

            using (var textWriter =
                new StreamWriter(new FileStream(jsonFileName, FileMode.Create, FileAccess.ReadWrite), Encoding.Default))
            {
                textWriter.Write(Encoding.UTF8.GetString(result));
            }
        }

        private List<PlayerData> sortPlayers(List<PlayerData> players, string bankerUserId)
        {
            var ps = players.ToArray();
            Array.Sort(ps, (x, y) => x.ChairId - y.ChairId);

            int bankerIdx = 0;
            foreach (var pd in ps)
            {
                if (bankerUserId == pd.UserId)
                {
                    break;
                }
                bankerIdx++;
            }

            var result = new List<PlayerData>();
            for (int i = 0; i < ps.Length; ++i)
            {
                result.Add(ps[bankerIdx % ps.Length]);
                bankerIdx++;
            }

            return result;
        }

        //private string ExtractDrawSequence(SRMsgHandRecorder recorder)
        //{
        //    var drawActions = recorder.actions.Where((x) => x.action == (int)ActionType.enumActionType_DRAW);
        //    var sb = new StringBuilder();
        //    foreach (var drawAction in drawActions)
        //    {
        //        if (drawAction.cards.Count > 1 || 0 != (drawAction.flags & (int)pokerface.SRFlags.SRFlyRichi))
        //        {
        //            // 杠后牌切换
        //            continue;
        //        }

        //        foreach (var tile in drawAction.cards)
        //        {
        //            if (tile < (int)CardID.CARDMAX)
        //            {
        //                sb.Append(_myWindow.IdNames[tile]);
        //                sb.Append(",");
        //            }

        //        }
        //    }

        //    return sb.ToString();
        //}

        //private string ExtractKongSequence(SRMsgHandRecorder recorder)
        //{
        //    var drawActions = recorder.actions.Where((x) => x.action == (int)ActionType.enumActionType_DRAW);
        //    var sb = new StringBuilder();

        //    // 第一个风牌也是杠后牌
        //    var windID = recorder.windFlowerID;
        //    sb.Append(_myWindow.IdNames[windID]);
        //    sb.Append(",");

        //    foreach (var drawAction in drawActions)
        //    {
        //        if (drawAction.cards.Count > 1)
        //        {
        //            var tile = drawAction.cards[1];
        //            sb.Append(_myWindow.IdNames[tile]);
        //            sb.Append(",");
        //        }
        //    }

        //    return sb.ToString();
        //}

        private List<PlayerData> FromRecorder(SRMsgHandRecorder recorder)
        {
            var players = recorder.players;
            var playerDatas = new List<PlayerData>();

            foreach (var player in players)
            {
                var pd = new PlayerData() { UserId = player.userID, ChairId = player.chairID, MyWindow = _myWindow };
                playerDatas.Add(pd);
            }

            foreach (var dealDetail in recorder.deals)
            {
                var pd = playerDatas.Find((x) => x.ChairId == dealDetail.chairID);
                pd.Deal = dealDetail;
            }

            foreach (var pd in playerDatas)
            {
                var actions = recorder.actions.Where((x) => x.chairID == pd.ChairId).ToList();
                pd.Actions = actions;
            }

            return playerDatas;
        }

        internal class PlayerData
        {
            public string UserId { get; set; }

            public MainWindow MyWindow { get; set; }

            public int ChairId { get; set; }
            public SRDealDetail Deal { get; set; }

            public List<SRAction> Actions { get; set; }

            public void WriteCsv(CsvWriter csv)
            {
                csv.WriteField(UserId);

                var handStr = ToHandString();
                csv.WriteField(handStr);

                //var flowerStr = ToFlowerString();
                //csv.WriteField(flowerStr);

                var actionTips = ActionTips();
                csv.WriteField(actionTips);
            }

            public string ToHandString()
            {
                var hands = Deal.cardsHand;
                var sb = new StringBuilder();
                foreach (var hand in hands)
                {
                    sb.Append(MyWindow.IdNames[hand]);
                    sb.Append(",");
                }
                return sb.ToString();
            }

            //public string ToFlowerString()
            //{
            //    var flowers = Deal.tilesFlower;
            //    var sb = new StringBuilder();
            //    foreach (var flower in flowers)
            //    {
            //        sb.Append(MyWindow.IdNames[flower]);
            //        sb.Append(",");
            //    }
            //    return sb.ToString();
            //}

            public string ActionTips()
            {
                var qaIndex = 0;
                var sb = new StringBuilder();
                foreach (var action in Actions)
                {
                    // 忽略抽牌动作
                    if (action.action == (int)ActionType.enumActionType_DRAW)
                    {
                        continue;
                    }

                    // 由于扑克牌没有吃椪杠同时发生因此没有同样的qaIndex
                    // 但是扑克牌最后一个win-selfdrawn的qaIndex前一个动作是一致的
                    // 因此需要注释掉
                    // 同样的qaIndex，忽略
                    //if (action.qaIndex == qaIndex)
                    //{
                    //    continue;
                    //}

                    qaIndex = action.qaIndex;
                    var act = (ActionType)action.action;
                    switch (act)
                    {
                        case ActionType.enumActionType_DISCARD:
                            sb.Append(
                                "[discard ");
                            for(var i = 1; i < action.cards.Count; i++)
                            {
                                var tileId = action.cards[i];
                                sb.Append($"{MyWindow.IdNames[tileId]} ");
                            }
                            sb.Append("],");
                            break;
                        case ActionType.enumActionType_Win_SelfDrawn:
                            sb.Append(
                                $"[winSelf],");
                            break;
                        default:
                            sb.Append(
                                $"[skip],");
                            break;
                    }
                }

                return sb.ToString();
            }
        }
    }

}
