using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Text;
using System.Text.RegularExpressions;
using System.Windows;
using System.Windows.Media.Imaging;
using CsvHelper;
using pokerface;
using Newtonsoft.Json.Linq;

namespace PokerTest
{
    /// <summary>
    /// MainWindow.xaml 的交互逻辑
    /// </summary>
    public partial class MainWindow : Window
    {
        public const string Version = "1.13";
        public const string ToolName = "RunFast";
        public MainWindow()
        {
            InitIdNames();
            DataContext = this;
            InitializeComponent();

            LoadMahjongPics();

            HideTileStackWnds();

            ProgramConfig.LoadConfigFromFile();

            UpdateTile();
        }

        private void UpdateTile()
        {
            Title = $"PokerTestTool-[{ToolName}](ver:{Version})[{ProgramConfig.ServerUrl}]";
        }

        private readonly List<Player> _players = new List<Player>();
        public int TileCfgIndex { get; private set; }

        public Dictionary<int, BitmapImage> ImageDict { get;  } = new Dictionary<int, BitmapImage>();

        public List<DealCfgSimple> DealCfgs { get; } = new List<DealCfgSimple>();

        //private DispatcherTimer _dispatcherTimer;
        private void OnUploadCfgFile_Button_Click(object sender, RoutedEventArgs e)
        {
            // Configure open file dialog box
            Microsoft.Win32.OpenFileDialog dlg = new Microsoft.Win32.OpenFileDialog();
            dlg.FileName = "origin"; // Default file name
            dlg.DefaultExt = ".csv"; // Default file extension
            dlg.Filter = "CSV documents (.csv)|*.csv"; // Filter files by extension

            // Show open file dialog box
            bool? result = dlg.ShowDialog();

            // Process open file dialog box results
            if (result == true)
            {
                try
                {
                    // Open document
                    var filePath = dlg.FileName;
                    HttpHandlers.SendFileContent(filePath, this);
                }
                catch (Exception ex)
                {
                    MessageBox.Show(ex.Message);
                }
            }
        }

        private static DealCfgSimple FromString(JObject x)
        {
            var a = x;
            return new DealCfgSimple() {Name = (string)a["name"], PlayerCount = (int)a["playerRequired"] };
        }

        public void OnUploaded(string body)
        {
            DealCfgs.Clear();
            var a = JObject.Parse(body);
            JArray configNames = (JArray) a["cfgs"];
            DealCfgs.AddRange(configNames.Select(c => FromString((JObject)c)).ToArray());

            if (DealCfgs.Count > 0)
            {
                TileCfgIndex = 0;
                TbCurrentCfg.Text = CurrentDealCfg.Name;
            }
        }

        private void BuildPlayers()
        {
            var names = new [] {"A","B", "C", "D"};
            HideTileStackWnds();
            var wnds = new TileStackWnd[4] { Auc, Buc, Cuc, Duc };
            var userIds = new string[4] {"1", "2", "3", "4"};
            // TODO: 为了和unity配合测试，少启动一个
            for (var i = 0; i < CurrentDealCfg.PlayerCount -1 ; ++i)
            {
                var player = new Player(names[i], userIds[i], "monkey-room", wnds[i], this);
                player.Connect();

                _players.Add(player);
            }
        }
        private void OnSinglePlayer_Button_Click(object sender, RoutedEventArgs e)
        {
            if (_players.Count == 4)
            {
                MessageBox.Show("already 4 players");
                return;
            }

            var wnds = new TileStackWnd[4] { Auc, Buc, Cuc, Duc };
            TileStackWnd freeWnd = null;
            foreach (var wnd in wnds)
            {
                var p = _players.Find((x) => x.MyWnd == wnd);
                if (p == null)
                {
                    freeWnd = wnd;
                    break;
                }
            }

            if (freeWnd == null)
            {
                MessageBox.Show("no free player-view to use");
                return;
            }

            var inputWnd = new InputWnd {Owner = this};
            var result = inputWnd.ShowDialog();
            if (result == false)
            {
                return;
            }

            var userId = inputWnd.TextBoxUserId.Text;
            var roomNumber = inputWnd.TextBoxRoomId.Text;
            var player = new Player(userId, userId, roomNumber, freeWnd, this);
            player.Connect();

            _players.Add(player);
        }

        public void OnPlayerDisconnected(Player player)
        {
            _players.Remove(player);
        }

        private void HideTileStackWnds()
        {
            var wnds = new TileStackWnd[4] { Auc, Buc, Cuc, Duc };
            foreach (var wnd in wnds)
            {
                wnd.Visibility = Visibility.Collapsed;
            }
        }

        public DealCfgSimple CurrentDealCfg
        {
            get
            {
                if (DealCfgs.Count < 1)
                    return DealCfgSimple.Empty;

                if (TileCfgIndex < 0)
                    TileCfgIndex = 0;

                if (TileCfgIndex >= DealCfgs.Count)
                    TileCfgIndex = DealCfgs.Count - 1;

                return DealCfgs[TileCfgIndex];
            }
        }

        private void OnStartGame_Button_Click(object sender, RoutedEventArgs e)
        {
            if (_players.Count == 0)
            {
                BuildPlayers();
            }
            else
            {
                foreach (var player in _players)
                {
                    player.SendReady2Server();
                }
            }
            //else if (_players.Count != CurrentDealCfg.PlayerCount)
            //{
            //    _players.ForEach(KillPlayer);
            //    _players.Clear();
            //    BuildPlayers();
            //}

            //if (string.IsNullOrWhiteSpace(CurrentDealCfg.Name))
            //{
            //    MessageBox.Show("no config to start game");
            //    return;
            //}

            //HttpHandlers.SendPostMethod(@"/start", CurrentDealCfg.Name);
        }

        private void KillPlayer(Player player)
        {
            player.Dispose();
        }

        private string _saveCfgPrefix;
        private int _saveCfgIndex;
        private void OnExportCfg_Button_Click(object sender, RoutedEventArgs e)
        {
            var log = TbLogger.Text;
            if (string.IsNullOrWhiteSpace(log))
                return;
            //var log = @"";
            var log2X = new Log2X();
            try
            {
                if (!log2X.Parse(log))
                {
                    MessageBox.Show(this, "log parse error");
                    return;
                }

                // write csv file
                // 名称	庄家手牌	庄家花牌	闲家1手牌	闲家1花牌	闲家2牌	闲家2花牌	闲家3手牌	闲家3花牌	抽牌序列	强制庄家	强制风牌
                Microsoft.Win32.SaveFileDialog dlg = new Microsoft.Win32.SaveFileDialog();
                if (_saveCfgPrefix != null)
                {
                    dlg.FileName = _saveCfgPrefix + _saveCfgIndex;
                }
                else
                {
                    dlg.FileName = "dealcfg"; // Default file name
                }

                dlg.DefaultExt = ".csv"; // Default file extension
                dlg.Filter = "CSV documents (.csv)|*.csv"; // Filter files by extension

                // Show save file dialog box
                Nullable<bool> result = dlg.ShowDialog();
                // Process save file dialog box results
                if (result == true)
                {
                    // Save document
                    string filename = dlg.FileName;
                    ParsePrefixAndIndex(filename, out _saveCfgPrefix, out _saveCfgIndex);
                    using (var textWriter = new StreamWriter(new FileStream(filename, FileMode.Create, FileAccess.ReadWrite), Encoding.Default))
                    {
                        var headers = new [] {"名称", "庄家手牌", "庄家花牌",
                            "闲家1手牌", "闲家1花牌", "闲家2牌", "闲家2花牌", "闲家3手牌", "闲家3花牌", "抽牌序列", "强制庄家", "强制风牌" };
                        var csv = new CsvWriter(textWriter);
                        foreach (var header in headers)
                        {
                            csv.WriteField(header);
                        }
                        csv.NextRecord();

                        csv.WriteField(string.IsNullOrWhiteSpace(CurrentDealCfg.Name) ? "bug" : CurrentDealCfg.Name);

                        var writedDeal = 0;
                        foreach (var item in log2X.Deals)
                        {
                            csv.WriteField(item.HandTiles);
                            csv.WriteField(item.FlowerTiles);
                            writedDeal+=2;
                        }

                        for (; writedDeal < 8; writedDeal++)
                        {
                            csv.WriteField("");
                        }

                        var sb = new StringBuilder();
                        foreach (var draw in log2X.Draws)
                        {
                            sb.Append(draw);
                            sb.Append(",");
                        }
                        csv.WriteField(sb.ToString());

                        if (null != log2X.Banker)
                        {
                            csv.WriteField(log2X.Banker);
                        }
                        if (null != log2X.Wind)
                        {
                            csv.WriteField(log2X.Wind);
                        }

                        csv.NextRecord();
                    }

                    // write action file
                    var actionFileName = filename.Substring(0, filename.LastIndexOf(".", StringComparison.Ordinal)) + ".txt";
                    using (var textWriter = new StreamWriter(new FileStream(actionFileName, FileMode.Create, FileAccess.ReadWrite), Encoding.Default))
                    {
                        foreach (var actionLine in log2X.ActionLines)
                        {
                            textWriter.WriteLine(actionLine);
                        }
                    }

                    // write result file
                    var resultMsg = GetScoreWndResult();
                    if (!string.IsNullOrWhiteSpace(resultMsg))
                    {
                        var resultFileName = filename.Substring(0, filename.LastIndexOf(".", StringComparison.Ordinal)) + ".result";
                        using (var textWriter = new StreamWriter(new FileStream(resultFileName, FileMode.Create, FileAccess.ReadWrite), Encoding.Default))
                        {
                            textWriter.Write(resultMsg);
                        }
                    }

                }
            }
            catch (Exception ex)
            {
                MessageBox.Show(this, ex.Message);
            }

        }

        private void OnExportRoomCfg_Button_Click(object sender, RoutedEventArgs e)
        {
            try
            {
                var inputWnd = new InputWnd { Owner = this, IsNeedUserId = false };
                var result = inputWnd.ShowDialog();
                if (result == false)
                {
                    return;
                }

                var roomNumber = inputWnd.TextBoxRoomId.Text;

                HttpHandlers.SendPostMethod("/resetRoom", roomNumber, "&roomNumber=" + roomNumber);
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }

        private void OnExportRoomOps_Button_Click(object sender, RoutedEventArgs e)
        {
            ExportRoomWnd.ShowExportDialog(ExportRoomWnd.ExportRoomType.Operations, this);
        }

        private void OnLoadAction_Button_Click(object sender, RoutedEventArgs e)
        {
            //Microsoft.Win32.OpenFileDialog dlg = new Microsoft.Win32.OpenFileDialog();
            //dlg.FileName = "dealcfg"; // Default file name
            //dlg.DefaultExt = ".txt"; // Default file extension
            //dlg.Filter = "txt documents (.txt)|*.txt"; // Filter files by extension

            //// Show open file dialog box
            //var result = dlg.ShowDialog();

            //// Process open file dialog box results
            //if (result == true)
            //{
            //    // Open document
            //    var x = HttpHandlers.WriteSafeReadAllLines(dlg.FileName);
            //    if (!string.IsNullOrWhiteSpace(x))
            //    {
            //        ShowActionListWnd(x);
            //    }
            //}
            TestAutoActionRegex();
        }

        private void TestAutoActionRegex()
        {
            var actionMsgs = new string[]
            {
                "[c-kong 8筒 8筒 8筒 8筒]",
                "[discard 9万]",
                "[discard 9万 richi true]",
                "[chow 1条 2条 3条]",
                "[richi False]",
                "[skip]",
                "[e-kong 7条 7条 7条 7条]",
                "[t-kong 7条 7条 7条 7条]",
                "[pong 6筒 6筒 6筒]",
                "[winChuck]",
                "[winSelf]",

            };

            string pattern;
            Regex rgx;
            MatchCollection matches;
            Match match;

            pattern = @"^\[(?<action>[^\s]+).*\]$";


            foreach (var msg in actionMsgs)
            {
                rgx = new Regex(pattern, RegexOptions.IgnoreCase);
                matches = rgx.Matches(msg);
                if (matches.Count > 0)
                {
                    match = matches[0];
                    var action = match.Groups["action"].Value;
                    Console.WriteLine(action);

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
                        matches = rgx.Matches(msg);
                        match = matches[0];
                        var paramA = match.Groups["paramA"].Value;
                        Console.WriteLine(paramA);
                    }

                    if (action == "discard")
                    {
                        var secondPattern = @"^\[.*richi\s(?<paramB>[^\s]+).*\]$";
                        rgx = new Regex(secondPattern, RegexOptions.IgnoreCase);
                        matches = rgx.Matches(msg);
                        if (matches.Count > 0)
                        {
                            match = matches[0];
                            var paramB = match.Groups["paramB"].Value;
                            Console.WriteLine(paramB);
                        }
                    }
                }

            }

        }

        public void SyncActionWnd(string x)
        {
            AlWnd?.SelectIfSame(x);
        }
        public ActionListWnd AlWnd { get; private set; }
        private void ShowActionListWnd(string x)
        {
            if (AlWnd == null)
            {
                AlWnd = new ActionListWnd();
                AlWnd.SetOwner(this);
            }

            var lines = new List<string>();
            using (StringReader reader = new StringReader(x))
            {

                string rline;
                while ((rline = reader.ReadLine()) != null)
                {
                    // Do something with the line
                    lines.Add(rline);
                }
            }

            AlWnd.ShowWithNewActionList(lines);
            //throw new NotImplementedException();
        }

        public void ResetActionListWndIndex()
        {
            AlWnd?.ResetSelectedIndex();
        }
        private static void ParsePrefixAndIndex(string filename, out string savePrefix, out int saveNextIndex)
        {
            savePrefix = null;
            saveNextIndex = 0;
            var fname = Path.GetFileNameWithoutExtension(filename);
            if (fname == null)
                return;
            for (var i = 0; i < fname.Length; ++i)
            {
                if (char.IsDigit(fname[i]))
                {
                    savePrefix = fname.Substring(0, i);
                    if (int.TryParse(fname.Substring(i), out saveNextIndex))
                    {
                        saveNextIndex++;
                    }

                    break;
                }
            }
        }
        private void OnSelectCfg_Button_Click(object sender, RoutedEventArgs e)
        {
            //  选择配置
            //if (DealCfgs.Count < 1)
            //    return;
            //int newSelected;
            //if (CfgSelectWnd.ShowDialog(DealCfgs, TileCfgIndex,out newSelected, this))
            //{
            //    TileCfgIndex = newSelected;
            //    TbCurrentCfg.Text = CurrentDealCfg.Name;
            //}

            var hostWnd = new HostWnd { Owner = this };
            var result = hostWnd.ShowDialog();
            if (result == false)
            {
                return;
            }

            ProgramConfig.ServerUrl = hostWnd.HostUrl;
            ProgramConfig.SaveConfig2File();

            UpdateTile();
        }

        private void OnCreateRoom_Button_Click(object sender, RoutedEventArgs e)
        {
            HttpHandlers.SendPostMethod(@"/support/createMonkeyRoom", "monkey", null);
        }

        private void OnDestroyRoom_Button_Click(object sender, RoutedEventArgs e)
        {
            HttpHandlers.SendPostMethod(@"/support/destroyMonkeyRoom", "",null);
            IsPlaying = false;
        }

        private void LoadMahjongPics()
        {
            var dir = Environment.CurrentDirectory;
            ImageDict.Clear();
            var suitCh = new string[] {"红桃", "方片", "草花", "黑桃", };
            var rankNames = new string[] { "2", "3", "4", "5", "6","7", "8", "9", "10", "J","Q","K","A"};
            // 4张牌
            var index = 0;
            for (var i = 0; i < 13; i++)
            {
                for (var suit = 0; suit < 4; suit++)
                {
                    var suitC = suitCh[suit];
                    var rank = rankNames[i];
                    var fileName = $"{suitC}{rank}.png";
                    var path = System.IO.Path.Combine(dir, "images", fileName);

                    var bitmap = new BitmapImage(new Uri(path));
                    ImageDict.Add(index++, bitmap);

                }
            }

            // 大小wang
            {
                var jokersName = new string[] { "小王", "大王" };
                for (var suit = 0; suit < 2; suit++)
                {
                    var joker = jokersName[suit];
                    
                    var fileName = $"{joker}.png";
                    var path = System.IO.Path.Combine(dir, "images", fileName);

                    var bitmap = new BitmapImage(new Uri(path));
                    ImageDict.Add(index++, bitmap);

                }
            }

            Auc.SetImageSrc(ImageDict, this);
            Buc.SetImageSrc(ImageDict, this);
            Cuc.SetImageSrc(ImageDict, this);
            Duc.SetImageSrc(ImageDict, this);
        }


        public void AppendLog(string logMsg)
        {
            TbLogger.AppendText(logMsg);
            TbLogger.ScrollToEnd();
        }
        public void AppendActionLog(string logMsg)
        {
            TbLogger.AppendText(logMsg);
            TbLogger.AppendText("\r\n");
            TbLogger.ScrollToEnd();

            SyncActionWnd(logMsg);
        }

        public void ClearLog()
        {
            TbLogger.Clear();
        }

        public string TileId2Name(int msgWindFlowerId)
        {
            return IdNames[msgWindFlowerId];
        }

        public int TileId2Name(string name)
        {
            return NameIds[name];
        }

        public Dictionary<int, string> IdNames { get; } = new Dictionary<int, string>();
        public Dictionary<string, int> NameIds { get; } = new Dictionary<string, int>();
        public bool IsPlaying { get; set; }

        public void InitIdNames()
        {
            NameIds["红桃2"] = (int)CardID.R2H;
            NameIds["方块2"] = (int)CardID.R2D;
            NameIds["梅花2"] = (int)CardID.R2C;
            NameIds["黑桃2"] = (int)CardID.R2S;

            NameIds["红桃3"] = (int)CardID.R3H;
            NameIds["方块3"] = (int)CardID.R3D;
            NameIds["梅花3"] = (int)CardID.R3C;
            NameIds["黑桃3"] = (int)CardID.R3S;

            NameIds["红桃4"] = (int)CardID.R4H;
            NameIds["方块4"] = (int)CardID.R4D;
            NameIds["梅花4"] = (int)CardID.R4C;
            NameIds["黑桃4"] = (int)CardID.R4S;

            NameIds["红桃5"] = (int)CardID.R5H;
            NameIds["方块5"] = (int)CardID.R5D;
            NameIds["梅花5"] = (int)CardID.R5C;
            NameIds["黑桃5"] = (int)CardID.R5S;

            NameIds["红桃6"] = (int)CardID.R6H;
            NameIds["方块6"] = (int)CardID.R6D;
            NameIds["梅花6"] = (int)CardID.R6C;
            NameIds["黑桃6"] = (int)CardID.R6S;

            NameIds["红桃7"] = (int)CardID.R7H;
            NameIds["方块7"] = (int)CardID.R7D;
            NameIds["梅花7"] = (int)CardID.R7C;
            NameIds["黑桃7"] = (int)CardID.R7S;

            NameIds["红桃8"] = (int)CardID.R8H;
            NameIds["方块8"] = (int)CardID.R8D;
            NameIds["梅花8"] = (int)CardID.R8C;
            NameIds["黑桃8"] = (int)CardID.R8S;

            NameIds["红桃9"] = (int)CardID.R9H;
            NameIds["方块9"] = (int)CardID.R9D;
            NameIds["梅花9"] = (int)CardID.R9C;
            NameIds["黑桃9"] = (int)CardID.R9S;

            NameIds["红桃10"] = (int)CardID.R10H;
            NameIds["方块10"] = (int)CardID.R10D;
            NameIds["梅花10"] = (int)CardID.R10C;
            NameIds["黑桃10"] = (int)CardID.R10S;

            NameIds["红桃J"] = (int)CardID.JH;
            NameIds["方块J"] = (int)CardID.JD;
            NameIds["梅花J"] = (int)CardID.JC;
            NameIds["黑桃J"] = (int)CardID.JS;

            NameIds["红桃Q"] = (int)CardID.QH;
            NameIds["方块Q"] = (int)CardID.QD;
            NameIds["梅花Q"] = (int)CardID.QC;
            NameIds["黑桃Q"] = (int)CardID.QS;

            NameIds["红桃K"] = (int)CardID.KH;
            NameIds["方块K"] = (int)CardID.KD;
            NameIds["梅花K"] = (int)CardID.KC;
            NameIds["黑桃K"] = (int)CardID.KS;

            NameIds["红桃A"] = (int)CardID.AH; ;
            NameIds["方块A"] = (int)CardID.AD; ;
            NameIds["梅花A"] = (int)CardID.AC; ;
            NameIds["黑桃A"] = (int)CardID.AS; ;

            NameIds["黑小丑"] = (int)CardID.JOB;
            NameIds["红小丑"] = (int)CardID.JOR;

            IdNames[(int)CardID.R2H] = "红桃2";
            IdNames[(int)CardID.R2D] = "方块2";
            IdNames[(int)CardID.R2C] = "梅花2";
            IdNames[(int)CardID.R2S] = "黑桃2";

            IdNames[(int)CardID.R3H] = "红桃3";
            IdNames[(int)CardID.R3D] = "方块3";
            IdNames[(int)CardID.R3C] = "梅花3";
            IdNames[(int)CardID.R3S] = "黑桃3";

            IdNames[(int)CardID.R4H] = "红桃4";
            IdNames[(int)CardID.R4D] = "方块4";
            IdNames[(int)CardID.R4C] = "梅花4";
            IdNames[(int)CardID.R4S] = "黑桃4";

            IdNames[(int)CardID.R5H] = "红桃5";
            IdNames[(int)CardID.R5D] = "方块5";
            IdNames[(int)CardID.R5C] = "梅花5";
            IdNames[(int)CardID.R5S] = "黑桃5";

            IdNames[(int)CardID.R6H] = "红桃6";
            IdNames[(int)CardID.R6D] = "方块6";
            IdNames[(int)CardID.R6C] = "梅花6";
            IdNames[(int)CardID.R6S] = "黑桃6";

            IdNames[(int)CardID.R7H] = "红桃7";
            IdNames[(int)CardID.R7D] = "方块7";
            IdNames[(int)CardID.R7C] = "梅花7";
            IdNames[(int)CardID.R7S] = "黑桃7";

            IdNames[(int)CardID.R8H] = "红桃8";
            IdNames[(int)CardID.R8D] = "方块8";
            IdNames[(int)CardID.R8C] = "梅花8";
            IdNames[(int)CardID.R8S] = "黑桃8";

            IdNames[(int)CardID.R9H] = "红桃9";
            IdNames[(int)CardID.R9D] = "方块9";
            IdNames[(int)CardID.R9C] = "梅花9";
            IdNames[(int)CardID.R9S] = "黑桃9";

            IdNames[(int)CardID.R10H] = "红桃10";
            IdNames[(int)CardID.R10D] = "方块10";
            IdNames[(int)CardID.R10C] = "梅花10";
            IdNames[(int)CardID.R10S] = "黑桃10";

            IdNames[(int)CardID.JH] = "红桃J";
            IdNames[(int)CardID.JD] = "方块J";
            IdNames[(int)CardID.JC] = "梅花J";
            IdNames[(int)CardID.JS] = "黑桃J";

            IdNames[(int)CardID.QH] = "红桃Q";
            IdNames[(int)CardID.QD] = "方块Q";
            IdNames[(int)CardID.QC] = "梅花Q";
            IdNames[(int)CardID.QS] = "黑桃Q";

            IdNames[(int)CardID.KH] = "红桃K";
            IdNames[(int)CardID.KD] = "方块K";
            IdNames[(int)CardID.KC] = "梅花K";
            IdNames[(int)CardID.KS] = "黑桃K";

            IdNames[(int)CardID.AH] = "红桃A";
            IdNames[(int)CardID.AD] = "方块A";
            IdNames[(int)CardID.AC] = "梅花A";
            IdNames[(int)CardID.AS] = "黑桃A";

            IdNames[(int)CardID.JOB] = "黑小丑";
            IdNames[(int)CardID.JOR] = "红小丑";

        }

        private ScoreWnd _scoreWnd;
        internal int rb111;

        public void ShowScoreWnd(string msg)
        {
            if (_scoreWnd == null)
                _scoreWnd = new ScoreWnd();

            _scoreWnd.ShowWithMsg(msg, this);
        }

        public void ResetScoreWnd()
        {
            _scoreWnd?.Clear();
        }

        private void OnShowScoreWnd_Button_Click(object sender, RoutedEventArgs e)
        {
            _scoreWnd?.Show();
        }

        public string GetScoreWndResult()
        {
            return _scoreWnd == null ? "" : _scoreWnd.GetResultMsg();
        }

        private void OnAttachDealCfg_Button_Click(object sender, RoutedEventArgs e)
        {
            //TestManualProto();
            // Configure open file dialog box
            Microsoft.Win32.OpenFileDialog dlg = new Microsoft.Win32.OpenFileDialog();
            dlg.FileName = "origin"; // Default file name
            dlg.DefaultExt = ".csv"; // Default file extension
            dlg.Filter = "CSV documents (.csv)|*.csv"; // Filter files by extension

            // Show open file dialog box
            bool? result = dlg.ShowDialog();

            // Process open file dialog box results
            if (result == true)
            {
                try
                {
                    var inputWnd = new InputWnd { Owner = this, IsNeedUserId = false};
                    result = inputWnd.ShowDialog();
                    if (result == false)
                    {
                        return;
                    }

                    var roomNumber = inputWnd.TextBoxRoomId.Text;
                    // Open document
                    var filePath = dlg.FileName;
                    HttpHandlers.SendFileContent2(filePath, roomNumber, HttpHandlers.PathAttachDealCfgFile, this);
                }
                catch (Exception ex)
                {
                    MessageBox.Show(ex.Message);
                }
            }
        }

        private void TestManualProto()
        {
            // 编码
            using (MemoryStream ms = new MemoryStream(4096))
            {
                using (BinaryWriter bw = new BinaryWriter(ms))
                {
                    ManualProtoCoder.EncodeInt32(bw, (uint)MessageCode.OPDeal, 1);
                    var abc = new byte[4000];
                    abc[0] = 1;
                    abc[1] = 2;
                    abc[2] = 3;

                    ManualProtoCoder.EncodeBytes(bw, abc, 2);

                    var encoded = ms.ToArray();
                    var gmsg = encoded.ToProto<GameMessage>();
                    Debug.Assert(gmsg.Ops == (int)MessageCode.OPDeal);
                    Debug.Assert(abc.SequenceEqual(gmsg.Data));

                    // 解码
                    using (MemoryStream ms2 = new MemoryStream(encoded))
                    {
                        using (BinaryReader br = new BinaryReader(ms2))
                        {
                            uint ops = 0;
                            ManualProtoCoder.DecodeInt32(br, 1, out ops);
                            Debug.Assert(ops == (uint)MessageCode.OPDeal);

                            byte[] bytes;
                            ManualProtoCoder.DecodeBytes(br, 2, out bytes);

                            Debug.Assert(bytes.SequenceEqual(abc));
                        }
                    }
                }
            }
        }

        private void OnAttachRoomCfg_Button_Click(object sender, RoutedEventArgs e)
        {
            // Configure open file dialog box
            Microsoft.Win32.OpenFileDialog dlg = new Microsoft.Win32.OpenFileDialog();
            dlg.FileName = "room"; // Default file name
            dlg.DefaultExt = ".json"; // Default file extension
            dlg.Filter = "JSON documents (.json)|*.json"; // Filter files by extension

            // Show open file dialog box
            bool? result = dlg.ShowDialog();

            // Process open file dialog box results
            if (result == true)
            {
                try
                {
                    var inputWnd = new InputWnd { Owner = this, IsNeedUserId = false };
                    result = inputWnd.ShowDialog();
                    if (result == false)
                    {
                        return;
                    }

                    var roomNumber = inputWnd.TextBoxRoomId.Text;
                    // Open document
                    var filePath = dlg.FileName;
                    HttpHandlers.SendFileContent2(filePath, roomNumber, HttpHandlers.PathAttachRoomCfgFile, this);
                }
                catch (Exception ex)
                {
                    MessageBox.Show(ex.Message);
                }
            }
        }

        private void OnKickAllInRoom_Button_Click(object sender, RoutedEventArgs e)
        {
            try
            {
                var inputWnd = new InputWnd { Owner = this, IsNeedUserId = false };
                var result = inputWnd.ShowDialog();
                if (result == false)
                {
                    return;
                }

                var roomNumber = inputWnd.TextBoxRoomId.Text;

                HttpHandlers.SendPostMethod("/support/kickAll", roomNumber, "&roomNumber="+roomNumber);
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }

        private void OnDealCfg_Button_Click(object sender, RoutedEventArgs e)
        {
            var dealCfgWnd = new DealCfgWnd(this);
            dealCfgWnd.ShowDialog();
        }

        public bool IsFirstPlayer(Player myPlayer)
        {
            if (_players.Count < 1)
            {
                return false;
            }

            return _players[0] == myPlayer;
        }

        private void OnDisbandRoom_Button_Click(object sender, RoutedEventArgs e)
        {
            try
            {
                //var inputWnd = new InputWnd { Owner = this, IsNeedUserId = false };
                //var result = inputWnd.ShowDialog();
                //if (result == false)
                //{
                //    return;
                //}

                //var roomNumber = inputWnd.TextBoxRoomId.Text;

                //HttpHandlers.SendPostMethod("/support/disbandRoom", roomNumber, "&roomNumber=" + roomNumber);
                if (_players.Count < 1)
                {
                    return;
                }
                var player = _players[0];
                player.SendMessage((int)MessageCode.OPDisbandRequest, null);
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }

        private void OnRoomCount_Button_Click(object sender, RoutedEventArgs e)
        {
            try
            {
                HttpHandlers.SendGetMethod("/support/roomCount", null);
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }

        private void OnUserCount_Button_Click(object sender, RoutedEventArgs e)
        {
            try
            {
                HttpHandlers.SendGetMethod("/support/userCount", null);
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }

        private void OnExceptionCount_Button_Click(object sender, RoutedEventArgs e)
        {
            try
            {
                HttpHandlers.SendGetMethod("/support/roomException", null);
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }

        private void OnClearExceptionCount_Button_Click(object sender, RoutedEventArgs e)
        {
            try
            {
                HttpHandlers.SendGetMethod("/support/clearRoomException", null);
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }
    }
}
