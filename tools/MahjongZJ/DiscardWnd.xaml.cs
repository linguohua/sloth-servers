using System.Collections.Generic;
using System.Windows;
using System.Windows.Controls;
using mahjong;

namespace MahjongTest
{
    /// <summary>
    /// DiscardWnd.xaml 的交互逻辑
    /// </summary>
    public partial class DiscardWnd : Window
    {
        public int SelectedTile { get; private set; }
        public DiscardWnd()
        {
            InitializeComponent();

            InitButtonArray();

            HideAllButtons();

            SelectedTile = -1;
        }
        private void HideAllButtons()
        {
            foreach (var button in this)
            {
                button.Visibility = Visibility.Hidden;
            }
        }

        public Button[] ButtonsSp0 { get; } = new Button[14];
        public Button[] ButtonsSp1 { get; } = new Button[14];
        public Button[] ButtonsSp2 { get; } = new Button[14];

        public List<MsgReadyHandTips> ReadyHandTips { get; private set; }
        public IEnumerator<Button> GetEnumerator()
        {

            foreach (var b in ButtonsSp0)
            {
                yield return b;
            }

            foreach (var b in ButtonsSp1)
            {
                yield return b;
            }
            foreach (var b in ButtonsSp2)
            {
                yield return b;
            }

        }

        private void InitButtonArray()
        {
            var i = 0;
            foreach (var child in Sp0.Children)
            {
                ButtonsSp0[i++] = child as Button;
            }
            i = 0;
            foreach (var child in Sp1.Children)
            {
                ButtonsSp1[i++] = child as Button;
            }
            i = 0;
            foreach (var child in Sp2.Children)
            {
                var btn = child as Button;
                if (btn != null)
                {
                    btn.Click += OnTileSelected;
                    ButtonsSp2[i++] = btn;
                }
            }
        }

        private void OnTileSelected(object sender, RoutedEventArgs e)
        {
            HideSp1Buttons();
            HideSp0Buttons();

            var btn = sender as Button;
            if (btn == null)
                return;

            var tileId = (int)btn.Tag;
            SelectedTile = tileId;
            var readyHandTip = FindReadyHandTip(tileId);
            if (readyHandTip == null)
                return;

            var readyHandList = readyHandTip.readyHandList;
            var i = 0;
            for (var j =0; j < readyHandList.Count-1; j +=2)
            {
                var x = ButtonsSp1[i];
                x.Visibility = Visibility.Visible;

                var tid = readyHandList[j];
                x.Content = new Image() { Source = MyOwner.ImagesSrc[tid] };

                var y = ButtonsSp0[i];
                y.Visibility = Visibility.Visible;
                y.Content = readyHandList[j+1];

                i++;
            }
        }

        private MsgReadyHandTips FindReadyHandTip(int tileId)
        {
            return ReadyHandTips.Find((x) => x.targetTile == tileId);
        }

        private void HideSp1Buttons()
        {
            foreach (var button in ButtonsSp1)
            {
                button.Visibility = Visibility.Hidden;
            }
        }

        private void HideSp0Buttons()
        {
            foreach (var button in ButtonsSp0)
            {
                button.Visibility = Visibility.Hidden;
            }
        }
        private void OnDiscard_Button_Clicked(object sender, RoutedEventArgs e)
        {
            //throw new NotImplementedException();
            if (SelectedTile < 0)
            {

                MessageBox.Show(this, "please select a tile to discard");
                return;
            }
            DialogResult = true;
        }
        private void OnDiscardRichi_Button_Clicked(object sender, RoutedEventArgs e)
        {
            if (SelectedTile < 0)
            {

                MessageBox.Show(this, "please select a tile to discard");
                return;
            }
            var readyHandList = FindReadyHandTip(SelectedTile);
            if (readyHandList == null || readyHandList.readyHandList.Count < 1)
                return;

            ReadyHandFlags = (int)1;
            DialogResult = true;
        }

        private void OnDiscardFlyRichi_Button_Clicked(object sender, RoutedEventArgs e)
        {
            if (SelectedTile < 0)
            {

                MessageBox.Show(this, "please select a tile to discard");
                return;
            }
            var readyHandList = FindReadyHandTip(SelectedTile);
            if (readyHandList == null || readyHandList.readyHandList.Count < 1)
                return;

            ReadyHandFlags = (int)2;
            DialogResult = true;
        }

        public int ReadyHandFlags { get; set; }

        private void OnCancel_Button_Clicked(object sender, RoutedEventArgs e)
        {
            //throw new NotImplementedException();
            DialogResult = false;
        }

        private void SetReadyHandTips(List<MsgReadyHandTips> readyHandTips)
        {
            readyHandTips.Sort((x, y) =>
            {
                if (y.readyHandList.Count == x.readyHandList.Count)
                    return x.targetTile - y.targetTile;

                var i = 0;
                var sumx = 0;
                var sumy = 0;
                foreach (var xy in x.readyHandList)
                {
                    if ((i % 2) == 1)
                        sumx += xy;
                    ++i;
                }

                i = 0;
                foreach (var xy in y.readyHandList)
                {
                    if ((i % 2) == 1)
                        sumy += xy;
                    ++i;
                }

                if (sumy != sumx)
                    return sumy - sumx;

                return y.readyHandList.Count - x.readyHandList.Count;
            });


            ReadyHandTips = readyHandTips;
            var j = 0;
            foreach (var ri in readyHandTips)
            {
                var btn = ButtonsSp2[j];
                btn.Content = new Image() { Source = MyOwner.ImagesSrc[ri.targetTile] };
                btn.Tag = ri.targetTile;
                btn.Visibility = Visibility.Visible;

                ++j;
            }
            
        }

        public static bool ShowDialog(List<MsgReadyHandTips>tiles2Discarded, out int tileDiscarded, TileStackWnd owner)
        {
            tileDiscarded = 0;
            var x = new DiscardWnd();
            x.SetOwner(owner);
            x.BtnExtra.Visibility = Visibility.Hidden;
            x.BtnExtraXX.Visibility = Visibility.Hidden;
            x.SetReadyHandTips(tiles2Discarded);

            var result = x.ShowDialog();
            if (result == null|| !result.Value)
            {
                // snip
                return false;
            }

            tileDiscarded = x.SelectedTile;
            return true;
        }

        public static bool ShowDialog(List<MsgReadyHandTips> tiles2Discarded, out int tileDiscarded, out int readyHandFlags, int expectedReadyHandFlags, TileStackWnd owner)
        {
            tileDiscarded = 0;
            readyHandFlags = 0;
            var x = new DiscardWnd();
            x.SetOwner(owner);
            x.SetReadyHandTips(tiles2Discarded);
            x.BtnExtra.Visibility = Visibility.Hidden;
            x.BtnExtraXX.Visibility = Visibility.Hidden;
            if ((expectedReadyHandFlags&1) != 0) {
                x.BtnExtra.Visibility = Visibility.Visible;
            }
            if ((expectedReadyHandFlags & 2) != 0)
            {
                x.BtnExtraXX.Visibility = Visibility.Visible;
            }
            var result = x.ShowDialog();
            if (result == null || !result.Value)
            {
                // snip
                return false;
            }

            tileDiscarded = x.SelectedTile;
            readyHandFlags = x.ReadyHandFlags;
            return true;
        }

        private void SetOwner(TileStackWnd owner)
        {
            MyOwner = owner;
            this.Owner = owner.MyOwner;
        }

        public TileStackWnd MyOwner { get; private set; }

    }
}
