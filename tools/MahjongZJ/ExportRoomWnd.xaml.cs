using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using CsvHelper;
using mahjong;
using Path = System.IO.Path;

namespace MahjongTest
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
        private async void OnAccExportButton_Clicked(object sender, RoutedEventArgs e)
        {
            var text = AccTextBox.Text;
            if (string.IsNullOrWhiteSpace(text))
            {
                MessageBox.Show("please input a valid acc-url");
                return;
            }

            try
            {
                var filename = "user-acc-export-room-ops";
                await ExportOpsFromAcc(text, filename);
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
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

        private  Task ExportCfg(string text)
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

        private async Task ExportOpsFromAcc(string text, string fileName)
        {
            var result = await HttpHandlers.ExportRoomOpsAcc(text, this);

            if (result == null)
            {
                return;
            }

            var accRR = result.ToProto<accessory.MsgAccLoadReplayRecord>();
            await parseAndSaveRecord(accRR.replayRecordBytes, fileName, false);
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
            
            await parseAndSaveRecord(result, fileName, true);
        }

        private async Task parseAndSaveRecord(byte [] result, string fileName, bool loadConfig)
        {

            var recorder = result.ToProto<SRMsgHandRecorder>();

            List<PlayerData> players = FromRecorder(recorder);
            string drawSequnece = ExtractDrawSequence(recorder);
            var bankerPlayer = players.Find((p) => p.ChairId == recorder.bankerChairID);
            var bankderUserId = bankerPlayer.UserId;
            players = sortPlayers(players, bankderUserId);
            var extra = recorder.extra;

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
                    var headers = MyCsvHeaders.GetHeaders();

                    // 第一行
                    var csv = new CsvWriter(textWriter);
                    foreach (var header in headers)
                    {
                        csv.WriteField(header);
                    }
                    csv.NextRecord();

                    // 第二行
                    csv.WriteField("bug"); // 名字
                    csv.WriteField(MyCsvHeaders.GetRoomTypeName()); // 类型
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
                    csv.WriteField(drawSequnece);

                    csv.WriteField(1); // 强制一致
                    csv.WriteField(recorder.roomConfigID);
                    csv.WriteField(recorder.isContinuousBanker ? "1" : "0");

                    csv.NextRecord();
                }

                if (loadConfig && !string.IsNullOrWhiteSpace(recorder.roomConfigID))
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
            Array.Sort(ps, (x,y)=> x.ChairId - y.ChairId);

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

        private string ExtractDrawSequence(SRMsgHandRecorder recorder)
        {
            var drawActions = recorder.actions.Where((x) => x.action == (int) ActionType.enumActionType_DRAW);
            var sb = new StringBuilder();
            foreach (var drawAction in drawActions)
            {
                foreach (var tile in drawAction.tiles)
                {
                    if (tile < (int)TileID.enumTid_MAX)
                    {
                        sb.Append(_myWindow.IdNames[tile]);
                        sb.Append(",");
                    }

                }
            }

            drawActions = recorder.actions.Where((x) => x.action == (int)ActionType.enumActionType_CustomA);
            foreach (var drawAction in drawActions)
            {
                foreach (var tile in drawAction.tiles)
                {
                    if (tile < (int)TileID.enumTid_MAX)
                    {
                        sb.Append(_myWindow.IdNames[tile]);
                        sb.Append(",");
                    }
                }
            }

            return sb.ToString();
        }

        private List<PlayerData> FromRecorder(SRMsgHandRecorder recorder)
        {
            var players = recorder.players;
            var playerDatas = new List<PlayerData>();

            foreach (var player in players)
            {
                var pd = new PlayerData() {UserId = player.userID, ChairId = player.chairID, MyWindow =_myWindow};
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
                
                var actionTips = ActionTips();
                csv.WriteField(actionTips);
            }

            public string ToHandString()
            {
                var hands = Deal.tilesHand;
                var sb = new StringBuilder();
                foreach (var hand in hands)
                {
                    sb.Append(MyWindow.IdNames[hand]);
                    sb.Append(",");
                }
                return sb.ToString();
            }

            public string ActionTips()
            {
                var qaIndex = 0;
                var sb = new StringBuilder();
                foreach (var action in Actions)
                {
                    // 忽略抽牌动作
                    if (action.action == (int) ActionType.enumActionType_DRAW)
                    {
                        continue;
                    }

                    // 同样的qaIndex，忽略
                    if (action.qaIndex == qaIndex)
                    {
                        continue;
                    }

                    qaIndex = action.qaIndex;
                    var act = (ActionType) action.action;
                    var tileId = 0;
                    switch (act)
                    {
                            case ActionType.enumActionType_CHOW:
                                tileId = action.tiles[0];
                                sb.Append(
                                    $"[chow {MyWindow.IdNames[tileId]} {MyWindow.IdNames[tileId + 1]} {MyWindow.IdNames[tileId + 2]}],");
                            break;
                            case ActionType.enumActionType_PONG:
                                tileId = action.tiles[0];
                                sb.Append(
                                    $"[pong {MyWindow.IdNames[tileId]} {MyWindow.IdNames[tileId]} {MyWindow.IdNames[tileId]}],");
                            break;
                            case ActionType.enumActionType_KONG_Exposed:
                                tileId = action.tiles[0];
                                sb.Append(
                                    $"[e-kong {MyWindow.IdNames[tileId]} {MyWindow.IdNames[tileId]} {MyWindow.IdNames[tileId]} {MyWindow.IdNames[tileId]}],");
                            break;
                            case ActionType.enumActionType_KONG_Triplet2:
                                tileId = action.tiles[0];
                                sb.Append(
                                    $"[t-kong {MyWindow.IdNames[tileId]} {MyWindow.IdNames[tileId]} {MyWindow.IdNames[tileId]} {MyWindow.IdNames[tileId]}],");
                            break;
                            case ActionType.enumActionType_SKIP:
                                sb.Append(
                                    $"[skip],");
                            break;
                            case ActionType.enumActionType_KONG_Concealed:
                                tileId = action.tiles[0];
                                sb.Append(
                                    $"[c-kong {MyWindow.IdNames[tileId]} {MyWindow.IdNames[tileId]} {MyWindow.IdNames[tileId]} {MyWindow.IdNames[tileId]}],");
                            break;
                            case ActionType.enumActionType_DISCARD:
                                tileId = action.tiles[0];
                                if (action.flags == (int) SRFlags.SRRichi)
                                {
                                    sb.Append(
                                        $"[discard {MyWindow.IdNames[tileId]} richi true],");
                                }
                                else
                                {
                                    sb.Append(
                                        $"[discard {MyWindow.IdNames[tileId]}],");
                                }

                            break;
                            case ActionType.enumActionType_FirstReadyHand:
                                sb.Append(
                                    $"[richi {action.flags == (int)SRFlags.SRRichi}],");
                            break;
                            case ActionType.enumActionType_WIN_Chuck:
                                sb.Append(
                                    $"[winChuck],");
                            break;
                            case ActionType.enumActionType_WIN_SelfDrawn:
                                sb.Append(
                                    $"[winSelf],");
                            break;
                            //case ActionType.enumActionType_AccumulateWin:
                            //    sb.Append(
                            //        $"[finalDraw],");
                            //break;
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
