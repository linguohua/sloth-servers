using System.Collections.Generic;
using System.Windows;
using System.Windows.Controls;
using pokerface;

namespace PokerTest
{
    /// <summary>
    /// DiscardWnd.xaml 的交互逻辑
    /// </summary>
    public partial class DiscardWnd : Window
    {
        //lic int SelectedTile { get; private set; }
        public DiscardWnd()
        {
            InitializeComponent();

            InitButtonArray();

            HideAllButtons();

            //SelectedTile = -1;
        }
        private void HideAllButtons()
        {
            foreach (var button in this)
            {
                button.Visibility = Visibility.Hidden;
            }
        }

        public Button[] ButtonsSp0 { get; } = new Button[DealCfg.handMax];
        public Button[] ButtonsSp1 { get; } = new Button[DealCfg.handMax];
        public Button[] ButtonsSp2 { get; } = new Button[DealCfg.handMax];

        public List<int> SelectedTiles = new List<int>();
        public List<int> HandTiles = new List<int>();

        private pokerface.MsgCardHand prevCardHand;

        private int specialCardID;
        private List<MsgCardHand> discardAbleTips;
        private int discardAbleTipsIndex;

        //public List<MsgReadyHandTips> ReadyHandTips { get; private set; }
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
                var btn = child as Button;
                if (btn != null)
                {
                    btn.Click += OnTileDeSelected;
                    ButtonsSp0[i++] = btn;
                }
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
        private void OnTileDeSelected(object sender, RoutedEventArgs e)
        {
            var btn = sender as Button;
            if (btn == null)
                return;

            var tileId = (int)btn.Tag;
            SelectedTiles.Remove(tileId);
            HandTiles.Add(tileId);

            Hand2UI();
            Selected2UI();
        }

        private void OnTileSelected(object sender, RoutedEventArgs e)
        {
            // HideSp1Buttons();
            // HideSp0Buttons();

            var btn = sender as Button;
            if (btn == null)
                return;

            var tileId = (int)btn.Tag;
            HandTiles.Remove(tileId);
            SelectedTiles.Add(tileId);

            Hand2UI();
            Selected2UI();

            //SelectedTile = tileId;
            //var readyHandTip = FindReadyHandTip(tileId);
            //if (readyHandTip == null)
            //    return;

            //var readyHandList = readyHandTip.readyHandList;
            //var i = 0;
            //for (var j = 0; j < readyHandList.Count - 1; j += 2)
            //{
            //    var x = ButtonsSp1[i];
            //    x.Visibility = Visibility.Visible;

            //    var tid = readyHandList[j];
            //    x.Content = new Image() { Source = MyOwner.ImagesSrc[tid] };

            //    var y = ButtonsSp0[i];
            //    y.Visibility = Visibility.Visible;
            //    y.Content = readyHandList[j + 1];

            //    i++;
            //}
        }

        //private MsgReadyHandTips FindReadyHandTip(int tileId)
        //{
        //    return ReadyHandTips.Find((x) => x.targetTile == tileId);
        //}

        private void HideSp1Buttons()
        {
            foreach (var button in ButtonsSp1)
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
            if (SelectedTiles.Count < 0)
            {

                MessageBox.Show(this, "please select a tile to discard");
                return;
            }

            if (specialCardID >= 0)
            {
                var found = false;
                foreach(var c in SelectedTiles)
                {
                    if (c == specialCardID)
                    {
                        found = true;
                        break;
                    }
                }

                if (!found)
                {
                    var cardName =  this.MyOwner.MyOwner.IdNames[specialCardID];
                    MessageBox.Show(this, $"please select a hand that include {cardName} to discard");
                    return;
                }
            }

            var current = AgariIndex.agariConvertMsgCardHand(SelectedTiles.ToArray());
            if (current == null)
            {
                MessageBox.Show(this, "please select a valid hand");
                return;
            }

            if (prevCardHand != null && !AgariIndex.agariGreatThan(prevCardHand, current))
            {
                MessageBox.Show(this, "please select a hand great than prev");
                return;
            }

            DialogResult = true;
        }

        private void OnDiscardRichi_Button_Clicked(object sender, RoutedEventArgs e)
        {
            //if (SelectedTiles.Count < 0)
            //{

            //    MessageBox.Show(this, "please select a tile to discard");
            //    return;
            //}
            ////var readyHandList = FindReadyHandTip(SelectedTile);
            ////if (readyHandList == null || readyHandList.readyHandList.Count < 1)
            ////    return;

            //ReadyHandFlags = (int)1;
            //DialogResult = true;

            if (discardAbleTips == null || discardAbleTips.Count < 1)
            {
                return;
            }

            discardAbleTipsIndex++;
            if (discardAbleTipsIndex >= discardAbleTips.Count)
            {
                discardAbleTipsIndex = 0;
            }

            var current = discardAbleTips[discardAbleTipsIndex];
            HandTiles.AddRange(SelectedTiles);
            SelectedTiles.Clear();

            SelectedTiles.AddRange(current.cards);
            foreach (var c in current.cards)
            {
                HandTiles.Remove(c);
            }

            Hand2UI();
            Selected2UI();

        }

        private void OnDiscardFlyRichi_Button_Clicked(object sender, RoutedEventArgs e)
        {
            if (SelectedTiles.Count < 0)
            {

                MessageBox.Show(this, "please select a tile to discard");
                return;
            }
            //var readyHandList = FindReadyHandTip(SelectedTile);
            //if (readyHandList == null || readyHandList.readyHandList.Count < 1)
            //    return;

            ReadyHandFlags = (int)2;
            DialogResult = true;
        }

        public int ReadyHandFlags { get; set; }

        private void OnCancel_Button_Clicked(object sender, RoutedEventArgs e)
        {
            //throw new NotImplementedException();
            DialogResult = false;
        }

        private void SetReadyHandTips(List<int> readyHandTips)
        {
            HandTiles.Clear();
            HandTiles.AddRange(readyHandTips);

            Hand2UI();

        }

        private void Hand2UI()
        {
            HideSp2Buttons();
            HandTiles.Sort((x, y) =>
            {
                return x - y;
            });

            var j = 0;
            foreach (var ri in HandTiles)
            {
                var btn = ButtonsSp2[j];
                btn.Content = new Image() { Source = MyOwner.ImagesSrc[ri] };
                btn.Tag = ri;
                btn.Visibility = Visibility.Visible;

                ++j;
            }
        }

        private void Selected2UI()
        {
            HideSp0Buttons();
            SelectedTiles.Sort((x, y) =>
            {
                return x - y;
            });

            var j = 0;
            foreach (var ri in SelectedTiles)
            {
                var btn = ButtonsSp0[j];
                btn.Content = new Image() { Source = MyOwner.ImagesSrc[ri] };
                btn.Tag = ri;
                btn.Visibility = Visibility.Visible;

                ++j;
            }
        }

        public static bool ShowDialog(List<int> tileDiscarded, pokerface.MsgCardHand prevCardHand, int specialCardID, TileStackWnd owner)
        {
            var tiles2Discarded = owner.TilesHandList;
            tileDiscarded.Clear();

            var x = new DiscardWnd();
            x.SetOwner(owner);
            x.BtnExtra.Visibility = Visibility.Hidden;
            x.BtnExtraXX.Visibility = Visibility.Hidden;
            x.SetReadyHandTips(tiles2Discarded);
            x.prevCardHand = prevCardHand;
            x.specialCardID = specialCardID;

            if (prevCardHand != null)
            {
                var currents = AgariIndex.FindGreatThanCardHand(prevCardHand, tiles2Discarded,specialCardID);
                if (null == currents || currents.Count == 0)
                {
                    MessageBox.Show("oh shit, a huge bug");
                    throw new System.Exception("huge bug");
                }

                if (currents.Count > 1)
                {
                    x.BtnExtra.Visibility = Visibility.Visible;
                    x.BtnExtra.Content = "下一个提示";
                }
 
                var current = currents[0];
                x.SelectedTiles.AddRange(current.cards);
                foreach(var c in current.cards)
                {
                    x.HandTiles.Remove(c);
                }

                x.discardAbleTips = currents;
                x.discardAbleTipsIndex = 0;

                x.Hand2UI();
                x.Selected2UI();
            }
            else
            {
                var current = AgariIndex.SearchLongestDiscardCardHand(tiles2Discarded, specialCardID);
                x.SelectedTiles.AddRange(current.cards);
                foreach (var c in current.cards)
                {
                    x.HandTiles.Remove(c);
                }

                x.Hand2UI();
                x.Selected2UI();
            }

            var result = x.ShowDialog();
            if (result == null || !result.Value)
            {
                // snip
                return false;
            }

            tileDiscarded.AddRange(x.SelectedTiles);
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
