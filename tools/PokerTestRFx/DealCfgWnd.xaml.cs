using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Windows;
using System.Windows.Controls;
using CsvHelper;
using pokerface;

namespace PokerTest
{
    /// <summary>
    /// DealCfgWnd.xaml 的交互逻辑
    /// </summary>
    public partial class DealCfgWnd : Window
    {
        public class XTiles
        {
            public List<int> Tiles = new List<int>();
            public Button[] ButtonsHand;
            public GroupBox GroupBox;
            private readonly DealCfgWnd _owner;
            private readonly string _groupBoxName;
            public const int MaxXCount = 16;
            public readonly int _fieldIndex;
            public XTiles(DealCfgWnd owner)
            {
                _owner = owner;

                _groupBoxName = "杠后牌";
                _fieldIndex = 15;

                ButtonsHand = new Button[MaxXCount];
            }

            public string ToTilesString()
            {
                //var hands = Tiles.ToArray();
                //Array.Sort(hands);

                var sb = new StringBuilder();
                foreach (var tile in Tiles)
                {
                    sb.Append(_owner._owner.IdNames[tile]);
                    sb.Append(",");
                }
                return sb.ToString();
            }

            public void HideAllButtons()
            {
                foreach (var button in ButtonsHand)
                {
                    button.Visibility = Visibility.Collapsed;
                }
            }

            public void Tiles2Ui()
            {
                HideAllButtons();

                int i = 0;
                foreach (var tile in Tiles)
                {
                    var btn = ButtonsHand[i];
                    btn.Tag = tile;
                    btn.Content = new Image() { Source = _owner._owner.ImageDict[tile] };
                    btn.Visibility = Visibility.Visible;
                    i++;
                }


                GroupBox.Header = $"{_groupBoxName}:{Tiles.Count}";
            }

            public void ReadCsv(CsvReader csvReader)
            {
                var drawSeqStrs = csvReader.GetField(_fieldIndex);
                var drawSeqStrArray = drawSeqStrs.Split(',', '，', ' ', '\t');

                foreach (var s in drawSeqStrArray)
                {
                    if (!string.IsNullOrWhiteSpace(s))
                    {
                        var tid = _owner._owner.NameIds[s];
                        {
                            Tiles.Add(tid);
                        }
                    }
                }
            }
        }

        public DealCfgWnd(MainWindow owner)
        {
            _owner = owner;
            Owner = owner;

            InitializeComponent();

            InitWallTiles();

            InitWallTilesUi();

            InitDealCfgs();

            InitDealCfgUi();

            _drawCfg = new DrawCfg(this);
            InitDrawUi();
            _endXTile = new XTiles(this);
            InitXTilesUi();
            //switch (_owner.rb111)
            //{
            //    case 0:
            //        RB136.IsChecked = true;
            //        break;
            //    case 1:
            //        RB108.IsChecked = true;
            //        break;
            //}
        }

        private readonly Random _random = new Random();
        public readonly MainWindow _owner;
        public readonly int[] _wallTiles = new int[(int)CardID.CARDMAX];

        private readonly Button[] _wallTilesButtons = new Button[(int)CardID.CARDMAX];
        private readonly Label[] _wallTilesLabels = new Label[(int)CardID.CARDMAX];

        private readonly DealCfg[] _dealCfgs = new DealCfg[4];
        private DrawCfg _drawCfg ;
        private readonly XTiles _endXTile;

        //public readonly int WindId = (int)CardID.CARDMAX;

        // "抽牌序列", "杠后牌", "上楼计数"
        public static string[] Headers = new[] {
            "名称", "类型", "庄家userID", "庄家手牌", "庄家动作提示", "userID2", "手牌", "动作提示", "userID3", "手牌",
            "动作提示", "userID4", "手牌", "动作提示",  "强制一致", "房间配置ID", "是否连庄"
        };

        private void InitDealCfgs()
        {
            for (int i = 0; i < 4; ++i)
            {
                _dealCfgs[i] = new DealCfg(i, this);
            }
        }

        private void InitWallTiles()
        {
            for (int i = 0; i < (int) CardID.CARDMAX; i++)
            {
                _wallTiles[i] = 1;
            }

            for (var i = (int)CardID.R2D; i <= (int)CardID.R2S; i++)
            {
                _wallTiles[i] = 0;
            }

            for (var i = (int)CardID.JOB; i <= (int)CardID.JOR; i++)
            {
                _wallTiles[i] = 0;
            }

            for (var i = (int)CardID.AS; i <= (int)CardID.AS; i++)
            {
                _wallTiles[i] = 0;
            }

            //for (int i = (int) CardID.CARDMAX; i < (int) CardID.CARDMAX; i++)
            //{
            //    _wallTiles[i] = 1;
            //}
        }

        private void InitDrawUi()
        {
            _drawCfg.GroupBox = EGroupBox;

            //for (int j = 0; j < DrawCfg.MaxDrawCount; j++)
            //{
            //    Button btn = new Button();
            //    btn.Width = 45;
            //    btn.Height = 53;
            //    btn.Visibility = Visibility.Collapsed;

            //    ETiles.Children.Add(btn);
            //    _drawCfg.ButtonsHand[j] = btn;

            //    btn.Click += OnDrawCfgButtonClicked;
            //}
        }

        public Button AddDrawUIButton()
        {
            Button btn = new Button();
            btn.Width = 45;
            btn.Height = 53;
            btn.Visibility = Visibility.Collapsed;

            ETiles.Children.Add(btn);
            btn.Click += OnDrawCfgButtonClicked;

            return btn;
        }

        private void InitXTilesUi()
        {
            _endXTile.GroupBox = KongXGroupBox;

            for (int j = 0; j < XTiles.MaxXCount; j++)
            {
                Button btn = new Button();
                btn.Width = 45;
                btn.Height = 53;
                btn.Visibility = Visibility.Collapsed;

                KongXTileNonFlower.Children.Add(btn);
                _endXTile.ButtonsHand[j] = btn;

                btn.Click += OnKongXButtonClicked;
            }
        }

        private void InitDealCfgUi()
        {
            var groupBoies = new GroupBox[]
            {
                AGroupBox,
                BGroupBox,
                CGroupBox,
                DGroupBox,
            };

            var nonFlowersPanels = new WrapPanel[]
            {
                ATileNonFlower,
                BTileNonFlower,
                CTileNonFlower,
                DTileNonFlower,
            };

            //var flowersPanels = new WrapPanel[]
            //{
            //    ATileFlower,
            //    BTileFlower,
            //    CTileFlower,
            //    DTileFlower,
            //};

            for (int i = 0; i < 4; ++i)
            {
                var dealCfg = _dealCfgs[i];
                dealCfg.GroupBox = groupBoies[i];

                var nonFlowerPanel = nonFlowersPanels[i];
                //var flowerPanel = flowersPanels[i];

                for (int j = 0; j < DealCfg.handMax; j++)
                {
                    Button btn = new Button();
                    btn.Width = 45;
                    btn.Height = 53;
                    btn.Visibility = Visibility.Collapsed;
                    
                    nonFlowerPanel.Children.Add(btn);
                    dealCfg.ButtonsHand[j] = btn;

                    btn.Click += OnDealCfgButtonClicked;
                }

                //for (int j = 0; j < 20; j++)
                //{
                //    Button btn = new Button();
                //    btn.Width = 45;
                //    btn.Height = 53;
                //    btn.Visibility = Visibility.Collapsed;

                //    flowerPanel.Children.Add(btn);
                //    dealCfg.ButtonsFlower[j] = btn;

                //    btn.Click += OnDealCfgButtonClicked;
                //}
            }
        }

        private void WallTiles2Ui()
        {
            int sum = 0;
            for (int i = 0; i < (int)CardID.CARDMAX; i++)
            {
                sum += _wallTiles[i];
                _wallTilesLabels[i].Content = _wallTiles[i];

                _wallTilesButtons[i].IsEnabled = _wallTiles[i] > 0;
            }

            WallTilesGroup.Header = $"牌墙({sum})";
        }

        private void InitWallTilesUi()
        {
            for (int i = 0; i < (int) CardID.CARDMAX; i++)
            {
                var stackPanel = CreateTileStackPanel(i);

                WallTileNonFlower.Children.Add(stackPanel);
            }

            //for (int i = (int)TileID.enumTid_HAK; i < (int)CardID.CARDMAX; i++)
            //{
            //    var stackPanel = CreateTileStackPanel(i);

            //    WallTileFlower.Children.Add(stackPanel);
            //}

            WallTiles2Ui();
        }

        private StackPanel CreateTileStackPanel(int i)
        {
            var stackPanel = new StackPanel();
            stackPanel.Orientation = Orientation.Vertical;
            stackPanel.Margin = new Thickness(0, 0, 10, 0);
            var btn = new Button();
            btn.Width = 45;
            btn.Height = 53;
            btn.Content = new Image() {Source = _owner.ImageDict[i]};

            btn.Click += OnWallTileButtonClicked;
            btn.Tag = i;

            var label = new Label();
            label.Content = _wallTiles[i];
            label.HorizontalAlignment = HorizontalAlignment.Center;
            label.VerticalAlignment = VerticalAlignment.Center;

            stackPanel.Children.Add(btn);
            stackPanel.Children.Add(label);

            _wallTilesButtons[i] = btn;
            _wallTilesLabels[i] = label;
            return stackPanel;
        }

        private void OnDealCfgButtonClicked(object sender, RoutedEventArgs e)
        {
            Button btn = sender as Button;
            if (btn == null)
            {
                return;
            }

            DealCfg.DealCfgTag select = (DealCfg.DealCfgTag)btn.Tag;

            bool found = false;
            if (select.Tile >= (int)CardID.CARDMAX /*|| select.Tile == WindId*/)
            {
                // found = select.DealCfg.TilesFlower.Remove(select.Tile);
            }
            else
            {
                found = select.DealCfg.TilesHand.Remove(select.Tile);
            }

            if (found)
            {
                _wallTiles[select.Tile]++;
            }

            select.DealCfg.Tiles2Ui();
            WallTiles2Ui();
        }

        private void OnDrawCfgButtonClicked(object sender, RoutedEventArgs e)
        {
            Button btn = sender as Button;
            if (btn == null)
            {
                return;
            }

            var select = (int)btn.Tag;

            bool found = _drawCfg.Tiles.Remove(select);

            if (found)
            {
                _wallTiles[select]++;
            }

            _drawCfg.Tiles2Ui();
            WallTiles2Ui();
        }
        private void OnKongXButtonClicked(object sender, RoutedEventArgs e)
        {
            Button btn = sender as Button;
            if (btn == null)
            {
                return;
            }

            var select = (int)btn.Tag;

            bool found = _endXTile.Tiles.Remove(select);

            if (found)
            {
                _wallTiles[select]++;
                WallTiles2Ui();
            }

            _endXTile.Tiles2Ui();
        }

        private void OnWallTileButtonClicked(object sender, RoutedEventArgs e)
        {
            Button btn = sender as Button;
            if (btn == null)
            {
                return;
            }

            int select = (int)btn.Tag;
            if (RadioButtonKongX.IsChecked == true)
            {
                if (_endXTile.Tiles.Count == XTiles.MaxXCount)
                {
                    return;
                }

                _endXTile.Tiles.Add(select);
                _endXTile.Tiles2Ui();
            }
            else if (RadioButtonE.IsChecked == true)
            {
                //if (_drawCfg.Tiles.Count == DrawCfg.MaxDrawCount)
                //{
                //    return;
                //}

                _drawCfg.Tiles.Add(select);
                _drawCfg.Tiles2Ui();
            }
            else
            {
                var dealCfg = GetCurrentSelectDealCfg();
                if (select >= (int)CardID.CARDMAX /*|| select == WindId*/)
                {
                    //dealCfg.TilesFlower.Add(select);
                }
                else
                {
                    int total = DealCfg.handMax;
                    if (dealCfg.IsBanker)
                    {
                        total = DealCfg.handMax;
                    }

                    if (dealCfg.TilesHand.Count == total)
                    {
                        return;
                    }

                    dealCfg.TilesHand.Add(select);
                }

                dealCfg.Tiles2Ui();
            }

            _wallTiles[select]--;
            WallTiles2Ui();
        }

        private DealCfg GetCurrentSelectDealCfg()
        {
            int i = 0;
            if (RadioButtonB.IsChecked == true)
            {
                i = 1;
            } else if (RadioButtonC.IsChecked == true)
            {
                i = 2;
            }
            else if (RadioButtonD.IsChecked == true)
            {
                i = 3;
            }

            return _dealCfgs[i];
        }

        private bool DrawNonFlower(out int tile, List<int>flowers)
        {
            tile = 0;

            while (true)
            {
                var remains = _wallTiles.Select((v, i) => new { value = v, index = i }).Where((p) => p.value > 0)
                    .Select((p) => p.index).ToArray();
                int remainCount = remains.Length;

                if (remainCount < 1)
                    return false;

                var select = remains[_random.Next(0, remainCount)];
                _wallTiles[select]--;

                if (select >= (int) CardID.CARDMAX /*|| select == WindId*/)
                {
                    flowers.Add(select);
                }
                else
                {
                    tile = select;
                    return true;
                }
            }
        }
        private bool DrawOne(out int tile)
        {
            tile = 0;

            var remains = _wallTiles.Select((v, i) => new { value = v, index = i }).Where((p) => p.value > 0)
                .Select((p) => p.index).ToArray();
            int remainCount = remains.Length;

            if (remainCount < 1)
                return false;

            var select = remains[_random.Next(0, remainCount)];
            _wallTiles[select]--;

            tile = select;

            return true;
        }

        private void DrawForDealCfg(DealCfg dealCfg)
        {
            int total = DealCfg.handMax;
            if (dealCfg.IsBanker)
            {
                total = DealCfg.handMax;
            }

            int current = dealCfg.TilesHand.Count;

            if (current == total)
            {
                return;
            }

            while (current < total)
            {
                int tile;
                List<int> flowers = new List<int>();
                var ok = DrawNonFlower(out tile, flowers);

                //dealCfg.TilesFlower.AddRange(flowers);

                if (!ok)
                {
                    break;
                }

                dealCfg.TilesHand.Add(tile);
                current++;
            }

            WallTiles2Ui();
            dealCfg.Tiles2Ui();
        }

        private void ClearDealCfgTiles(DealCfg dealCfg)
        {
            foreach (var t in dealCfg.TilesHand)
            {
                _wallTiles[t]++;
            }

            //foreach (var t in dealCfg.TilesFlower)
            //{
            //    _wallTiles[t]++;
            //}

            dealCfg.TilesHand.Clear();
            //dealCfg.TilesFlower.Clear();

            WallTiles2Ui();
            dealCfg.Tiles2Ui();
        }

        private void OnX0_Btn_Gernerate_Clicked(object sender, RoutedEventArgs e)
        {
            DrawForDealCfg(_dealCfgs[0]);
        }

        private void OnX0_Btn_Clear_Clicked(object sender, RoutedEventArgs e)
        {
            ClearDealCfgTiles(_dealCfgs[0]);
        }

        private void OnX1_Btn_Gernerate_Clicked(object sender, RoutedEventArgs e)
        {
            DrawForDealCfg(_dealCfgs[1]);
        }

        private void OnX1_Btn_Clear_Clicked(object sender, RoutedEventArgs e)
        {
            ClearDealCfgTiles(_dealCfgs[1]);
        }

        private void OnX2_Btn_Gernerate_Clicked(object sender, RoutedEventArgs e)
        {
            DrawForDealCfg(_dealCfgs[2]);
        }

        private void OnX2_Btn_Clear_Clicked(object sender, RoutedEventArgs e)
        {
            ClearDealCfgTiles(_dealCfgs[2]);
        }

        private void OnX3_Btn_Gernerate_Clicked(object sender, RoutedEventArgs e)
        {
            DrawForDealCfg(_dealCfgs[3]);
        }

        private void OnX3_Btn_Clear_Clicked(object sender, RoutedEventArgs e)
        {
            ClearDealCfgTiles(_dealCfgs[3]);
        }

        private void OnGenerate_Button_Clicked(object sender, RoutedEventArgs e)
        {
            foreach (var dealCfg in _dealCfgs)
            {
                DrawForDealCfg(dealCfg);
            }
        }

        private void OnSave_Button_Clicked(object sender, RoutedEventArgs e)
        {
            foreach (var dealCfg in _dealCfgs)
            {
                int total = DealCfg.handMax;
                if (dealCfg.IsBanker)
                {
                    total = DealCfg.handMax;
                }

                if (dealCfg.TilesHand.Count != 0 && dealCfg.TilesHand.Count != total)
                {
                    DrawForDealCfg(dealCfg);
                }
            }

            foreach (var dealCfg in _dealCfgs)
            {
                int total = DealCfg.handMax;
                if (dealCfg.IsBanker)
                {
                    total = DealCfg.handMax;
                }

                if (dealCfg.TilesHand.Count != 0 && dealCfg.TilesHand.Count != total)
                {
                    MessageBox.Show($"The {dealCfg.Index} set config hand tiles must equal to {total}");
                    return;
                }
            }

            for (int i = 3; i > 0; i--)
            {
                if (_dealCfgs[i].TilesHand.Count > 0 && _dealCfgs[i - 1].TilesHand.Count == 0)
                {
                    MessageBox.Show("config must continuous");
                    return;
                }
            }

            if (_dealCfgs[0].TilesHand.Count < 1 || _dealCfgs[1].TilesHand.Count < 1)
            {
                MessageBox.Show("at least have 2 player config");
                return;
            }

            //名称	userID1	手牌	花牌	动作提示	userID2	手牌	花牌	动作提示	userID3	手牌	花牌	动作提示	userID4	手牌	花牌	动作提示	抽牌序列	庄家	风牌
            Microsoft.Win32.SaveFileDialog dlg = new Microsoft.Win32.SaveFileDialog();
            var cfgName = "xyz";
            if (string.IsNullOrWhiteSpace(tbCfgName.Text) == false)
            {
                cfgName = tbCfgName.Text;
            }

            dlg.FileName = "user-new-" + cfgName; // Default file name

            dlg.DefaultExt = ".csv"; // Default file extension
            dlg.Filter = "CSV documents (.csv)|*.csv"; // Filter files by extension

            // Show save file dialog box
            bool? dlgResult = dlg.ShowDialog();
            // Process save file dialog box results
            if (dlgResult == true)
            {
                try
                {
                    // Save document
                    string filename = dlg.FileName;
                    using (var textWriter =
                        new StreamWriter(new FileStream(filename, FileMode.Create, FileAccess.ReadWrite),
                            Encoding.Default))
                    {

                        // 第一行
                        var csv = new CsvWriter(textWriter);
                        foreach (var header in Headers)
                        {
                            csv.WriteField(header);
                        }
                        csv.NextRecord();

                        // 第二行
                        csv.WriteField(cfgName); // 名字
                        csv.WriteField("大丰关张"); // 名字

                        foreach (var dealCfg in _dealCfgs)
                        {
                            dealCfg.WriteCsv(csv);
                        }

                        //csv.WriteField(_drawCfg.ToTilesString()); // 抽牌序列
                        //csv.WriteField(_endXTile.ToTilesString());
                        //csv.WriteField(0); // 上楼计数
                        csv.WriteField(0); // 强制一致
                        csv.WriteField(""); // 房间配置ID
                        csv.WriteField("0"); // 是否连庄
                        csv.NextRecord();
                    }
                }
                catch (Exception ex)
                {
                    MessageBox.Show(ex.Message);
                }
            }
        }

        private void OnX4_Btn_Gernerate_Clicked(object sender, RoutedEventArgs e)
        {
            //throw new NotImplementedException();
            int count = 0;
            //if (_drawCfg.Tiles.Count >= DrawCfg.MaxDrawCount)
            //{
            //    return;
            //}

            while (true)
            {
                int tile;
                var result = DrawOne(out tile);
                if (!result)
                {
                    break;
                }

                _drawCfg.Tiles.Add(tile);
                count++;
            }

            _drawCfg.Tiles2Ui();

            WallTiles2Ui();
        }

        private void OnX4_Btn_Clear_Clicked(object sender, RoutedEventArgs e)
        {
            foreach (var t in _drawCfg.Tiles)
            {
                _wallTiles[t]++;
            }

            _drawCfg.Tiles.Clear();

            WallTiles2Ui();
            _drawCfg.Tiles2Ui();
        }
        
        private void OnLoad_Button_Clicked(object sender, RoutedEventArgs e)
        {
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
                    var x = HttpHandlers.WriteSafeReadAllLines(dlg.FileName);
                    if (!string.IsNullOrWhiteSpace(x))
                    {
                        using (var stringReader = new StringReader(x))
                        {
                            // csv reader
                            var csvReader = new CsvReader(stringReader);
                            
                            // 读取头部
                            if (!csvReader.ReadHeader())
                            {
                                MessageBox.Show("Invalid input csv file, no record found");
                                return;
                            }


                            if (Headers.Length != csvReader.FieldHeaders.Length)
                            {
                                MessageBox.Show("Invalid input csv file, header count not match");
                                return;
                            }

                            for (var i = 0; i < Headers.Length; i++)
                            {
                                if (string.Compare(Headers[i], csvReader.FieldHeaders[i], StringComparison.Ordinal) != 0)
                                {
                                    MessageBox.Show("Invalid input csv file, header not match");
                                    return;
                                }
                            }

                            while (csvReader.Read())
                            {
                                for (var i = 0; i < 4; i++)
                                {
                                    var dealCfg = _dealCfgs[i];

                                    dealCfg.ReadCsv(csvReader);

                                    dealCfg.Tiles2Ui();
                                    
                                }

                                _drawCfg.ReadCsv(csvReader);
                                _drawCfg.Tiles2Ui();

                                //_endXTile.ReadCsv(csvReader);
                                //_endXTile.Tiles2Ui();

                                WallTiles2Ui();

                                // 仅读取第一个记录
                                break;
                            }
                        }
                    }
                }
                catch (Exception ex)
                {
                    MessageBox.Show(ex.Message);
                }
            }
        }

        private void clearXXX()
        {
            OnX0_Btn_Clear_Clicked(null, null);
            OnX1_Btn_Clear_Clicked(null, null);
            OnX2_Btn_Clear_Clicked(null, null);
            OnX3_Btn_Clear_Clicked(null, null);
            OnX4_Btn_Clear_Clicked(null, null);
        }

        //private void RadioButton136_Checked(object sender, RoutedEventArgs e)
        //{
        //    _owner.rb111 = 0;
        //    clearXXX();

        //    for (int i = 0; i < (int)CardID.CARDMAX; i++)
        //    {
        //        _wallTiles[i] = 1;
        //    }

        //    WallTiles2Ui();
        //}

        //private void RadioButton108_Checked(object sender, RoutedEventArgs e)
        //{
        //    _owner.rb111 = 1;
        //    clearXXX();
        //    for (int i = 0; i < (int)CardID.CARDMAX; i++)
        //    {
        //        _wallTiles[i] = 1;
        //    }

        //    for (int i = (int)CardID.CARDMAX; i < (int)CardID.CARDMAX; i++)
        //    {
        //        _wallTiles[i] =1;
        //    }

        //    WallTiles2Ui();
        //}

        private void OnXKongX_Btn_Gernerate_Clicked(object sender, RoutedEventArgs e)
        {
            int count = 0;
            if (_endXTile.Tiles.Count >= XTiles.MaxXCount)
            {
                return;
            }

            while (count < XTiles.MaxXCount)
            {
                int tile;
                var result = DrawOne(out tile);
                if (!result)
                {
                    break;
                }

                _endXTile.Tiles.Add(tile);
                count++;
            }

            _endXTile.Tiles2Ui();
            WallTiles2Ui();
        }

        private void OnXKongX_Btn_Clear_Clicked(object sender, RoutedEventArgs e)
        {
            foreach (var t in _endXTile.Tiles)
            {
                _wallTiles[t]++;
            }

            _endXTile.Tiles.Clear();

            WallTiles2Ui();
            _endXTile.Tiles2Ui();
        }

    }
}
