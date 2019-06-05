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
using mahjong;
using Newtonsoft.Json.Linq;

namespace MahjongTest
{
    /// <summary>
    /// MainWindow.xaml 的交互逻辑
    /// </summary>
    public partial class MainWindow : Window
    {
        public const string Version = "1.17";
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
            Title = $"MJTestTool-[ZJ](ver:{Version})[{ProgramConfig.ServerUrl}]";
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
                wnd.Visibility = Visibility.Hidden;
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

            //HttpHandlers.SendPostMethod(@"/support/start", CurrentDealCfg.Name);
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

                HttpHandlers.SendPostMethod("/support/resetRoom", roomNumber, "&roomNumber=" + roomNumber);
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

        private void OnUnlimit_Button_Click(object sender, RoutedEventArgs e)
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

                HttpHandlers.SendPostMethod("/support/ulimitRound", roomNumber, "&roomNumber=" + roomNumber);
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
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

            // 条子 MJs7.png
            for (var i = 1; i < 10; i++)
            {
                var fileName = $"MJs{i}.png";
                var path = System.IO.Path.Combine(dir, "images", fileName);

                var bitmap = new BitmapImage(new Uri(path));
                ImageDict.Add(AgariIndex.SOU1+i-1, bitmap);
            }
            // 筒子 MJt1.png
            for (var i = 1; i < 10; i++)
            {
                var fileName = $"MJt{i}.png";
                var path = System.IO.Path.Combine(dir, "images", fileName);

                var bitmap = new BitmapImage(new Uri(path));
                ImageDict.Add(AgariIndex.PIN1 + i - 1, bitmap);
            }
            // 万子 MJw7.png
            for (var i = 1; i < 10; i++)
            {
                var fileName = $"MJw{i}.png";
                var path = System.IO.Path.Combine(dir, "images", fileName);

                var bitmap = new BitmapImage(new Uri(path));
                ImageDict.Add(AgariIndex.MAN1 + i - 1, bitmap);
            }

            // 风牌 MJf1.png
            for (var i = 1; i < 5; i++)
            {
                var fileName = $"MJf{i}.png";
                var path = System.IO.Path.Combine(dir, "images", fileName);

                var bitmap = new BitmapImage(new Uri(path));
                ImageDict.Add(AgariIndex.TON + i - 1, bitmap);
            }

            // 箭牌 中：MJd1.png， 发：MJd2.png，白：MJd3.png
            for (var i = 1; i < 4; i++)
            {
                var fileName = $"MJd{i}.png";
                var path = System.IO.Path.Combine(dir, "images", fileName);

                var bitmap = new BitmapImage(new Uri(path));
                ImageDict.Add(AgariIndex.HAK + i - 1, bitmap);
            }

            // 花牌 MJh8.png
            for (var i = 1; i < 9; i++)
            {
                var fileName = $"MJh{i}.png";
                var path = System.IO.Path.Combine(dir, "images", fileName);

                var bitmap = new BitmapImage(new Uri(path));
                ImageDict.Add(AgariIndex.FlowerBegin + i - 1, bitmap);
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
        public int CurGameTypeSelectedIndex { get; internal set; }

        public void InitIdNames()
        {
            NameIds.Add("1万", AgariIndex.MAN1);
            NameIds.Add("2万", AgariIndex.MAN2);
            NameIds.Add("3万", AgariIndex.MAN3);
            NameIds.Add("4万", AgariIndex.MAN4);
            NameIds.Add("5万", AgariIndex.MAN5);
            NameIds.Add("6万", AgariIndex.MAN6);
            NameIds.Add("7万", AgariIndex.MAN7);
            NameIds.Add("8万", AgariIndex.MAN8);
            NameIds.Add("9万", AgariIndex.MAN9);

            NameIds.Add("1筒", AgariIndex.PIN1);
            NameIds.Add("2筒", AgariIndex.PIN2);
            NameIds.Add("3筒", AgariIndex.PIN3);
            NameIds.Add("4筒", AgariIndex.PIN4);
            NameIds.Add("5筒", AgariIndex.PIN5);
            NameIds.Add("6筒", AgariIndex.PIN6);
            NameIds.Add("7筒", AgariIndex.PIN7);
            NameIds.Add("8筒", AgariIndex.PIN8);
            NameIds.Add("9筒", AgariIndex.PIN9);

            NameIds.Add("1条", AgariIndex.SOU1);
            NameIds.Add("2条", AgariIndex.SOU2);
            NameIds.Add("3条", AgariIndex.SOU3);
            NameIds.Add("4条", AgariIndex.SOU4);
            NameIds.Add("5条", AgariIndex.SOU5);
            NameIds.Add("6条", AgariIndex.SOU6);
            NameIds.Add("7条", AgariIndex.SOU7);
            NameIds.Add("8条", AgariIndex.SOU8);
            NameIds.Add("9条", AgariIndex.SOU9);

            NameIds.Add("东", AgariIndex.TON);
            NameIds.Add("南", AgariIndex.NAN);
            NameIds.Add("西", AgariIndex.SHA);
            NameIds.Add("北", AgariIndex.PEI);
            NameIds.Add("白", AgariIndex.HAK);
            NameIds.Add("发", AgariIndex.HAT);
            NameIds.Add("中", AgariIndex.CHU);
            NameIds.Add("梅", AgariIndex.PLUM);
            NameIds.Add("兰", AgariIndex.ORCHID);
            NameIds.Add("竹", AgariIndex.BAMBOO);
            NameIds.Add("菊", AgariIndex.CHRYSANTHEMUM);
            NameIds.Add("春", AgariIndex.SPRING);
            NameIds.Add("夏", AgariIndex.SUMMER);
            NameIds.Add("秋", AgariIndex.AUTUMN);
            NameIds.Add("冬", AgariIndex.WINTER);


            IdNames.Add(AgariIndex.MAN1, "1万" );
            IdNames.Add(AgariIndex.MAN2,"2万" );
            IdNames.Add(AgariIndex.MAN3,"3万" );
            IdNames.Add(AgariIndex.MAN4,"4万" );
            IdNames.Add(AgariIndex.MAN5,"5万" );
            IdNames.Add(AgariIndex.MAN6,"6万" );
            IdNames.Add(AgariIndex.MAN7,"7万" );
            IdNames.Add(AgariIndex.MAN8,"8万" );
            IdNames.Add(AgariIndex.MAN9,"9万" );

            IdNames.Add(AgariIndex.PIN1,"1筒" );
            IdNames.Add(AgariIndex.PIN2,"2筒" );
            IdNames.Add(AgariIndex.PIN3,"3筒" );
            IdNames.Add(AgariIndex.PIN4,"4筒" );
            IdNames.Add(AgariIndex.PIN5,"5筒" );
            IdNames.Add(AgariIndex.PIN6,"6筒" );
            IdNames.Add(AgariIndex.PIN7,"7筒" );
            IdNames.Add(AgariIndex.PIN8,"8筒" );
            IdNames.Add(AgariIndex.PIN9,"9筒" );
              
            IdNames.Add(AgariIndex.SOU1,"1条" );
            IdNames.Add(AgariIndex.SOU2,"2条" );
            IdNames.Add(AgariIndex.SOU3,"3条" );
            IdNames.Add(AgariIndex.SOU4,"4条" );
            IdNames.Add(AgariIndex.SOU5,"5条" );
            IdNames.Add(AgariIndex.SOU6,"6条" );
            IdNames.Add(AgariIndex.SOU7,"7条" );
            IdNames.Add(AgariIndex.SOU8,"8条" );
            IdNames.Add(AgariIndex.SOU9,"9条" );

            IdNames.Add(AgariIndex.TON,"东");
            IdNames.Add(AgariIndex.NAN,"南");
            IdNames.Add(AgariIndex.SHA,"西");
            IdNames.Add(AgariIndex.PEI,"北");
            IdNames.Add(AgariIndex.HAK,"白");
            IdNames.Add(AgariIndex.HAT,"发");
            IdNames.Add(AgariIndex.CHU,"中");
            IdNames.Add(AgariIndex.PLUM, "梅");
            IdNames.Add(AgariIndex.ORCHID,"兰");
            IdNames.Add(AgariIndex.BAMBOO,"竹");
            IdNames.Add(AgariIndex.CHRYSANTHEMUM, "菊");
            IdNames.Add(AgariIndex.SPRING,"春");
            IdNames.Add(AgariIndex.SUMMER, "夏");
            IdNames.Add(AgariIndex.AUTUMN,"秋");
            IdNames.Add(AgariIndex.WINTER,"冬");
        }

        private ScoreWnd _scoreWnd;
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
                HttpHandlers.SendGetMethod("/roomCount", null);
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
                HttpHandlers.SendGetMethod("/userCount", null);
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
                HttpHandlers.SendGetMethod("/roomException", null);
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
                HttpHandlers.SendGetMethod("/clearRoomException", null);
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }

        private void OnSetExceptionCount_Button_Click(object sender, RoutedEventArgs e)
        {
            try
            {
                HttpHandlers.SendGetMethod("/incrRoomException", null);
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }
    }
}
