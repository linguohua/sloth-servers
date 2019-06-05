using System.Collections.Generic;
using System.Windows;
using System.Windows.Controls;
using mahjong;

namespace MahjongTest
{
    /// <summary>
    /// ChowPongKongWnd.xaml 的交互逻辑
    /// </summary>
    public partial class ChowPongKongWnd : Window
    {
        public ChowPongKongWnd()
        {
            InitializeComponent();

            InitButtonArray();
            HideAllButtons();

            SelectedTile = -1;
        }

        private void OnOK_Button_Clicked(object sender, RoutedEventArgs e)
        {
            //throw new NotImplementedException();
            if (SelectedTile < 0)
            {
                MessageBox.Show(this, "please select a meld to chow/pong/kong");
                return;
            }
            //throw new NotImplementedException();
            this.DialogResult = true;
        }

        private void OnCancel_Button_Clicked(object sender, RoutedEventArgs e)
        {
            //throw new NotImplementedException();
            this.DialogResult = false;
        }


        private void SetOwner(TileStackWnd owner)
        {
            MyOwner = owner;
            this.Owner = owner.MyOwner;
        }

        private void HideAllButtons()
        {
            foreach (var button in this)
            {
                button.Visibility = Visibility.Hidden;
            }
        }
        public Button[] ButtonsSp1 { get; } = new Button[14];
        public Button[] ButtonsSp2 { get; } = new Button[14];

        public TileStackWnd MyOwner { get; private set; }

        public static bool ShowDialog(List<MsgMeldTile> meldList, out int tileDiscarded, TileStackWnd owner)
        {
            tileDiscarded = 0;
            var x = new ChowPongKongWnd();
            x.SetOwner(owner);
            x.SetMeldList(meldList);

            var result = x.ShowDialog();
            if (result == null || !result.Value)
            {
                // snip
                return false;
            }

            tileDiscarded = x.SelectedTile;
            return true;
        }

        private void SetMeldList(List<MsgMeldTile> meldList)
        {
            var i = 0;
            foreach (var meld in meldList)
            {
                if (meld.meldType == (int) MeldType.enumMeldTypeTriplet2Kong
                    || meld.meldType == (int) MeldType.enumMeldTypeConcealedKong
                    || meld.meldType == (int) MeldType.enumMeldTypeExposedKong)
                {
                    for (int j = 0; j < 4; j++)
                    {
                        var btn = ButtonsSp2[i++];
                        btn.Content = new Image() { Source = MyOwner.ImagesSrc[meld.tile1] };
                        btn.Tag = meld.tile1;

                        btn.Visibility = Visibility.Visible;
                    }

                    i++;

                }
                else if (meld.meldType == (int)MeldType.enumMeldTypeTriplet)
                {
                    for (int j = 0; j < 3; j++)
                    {
                        var btn = ButtonsSp2[i++];
                        btn.Content = new Image() { Source = MyOwner.ImagesSrc[meld.tile1] };
                        btn.Tag = meld.tile1;

                        btn.Visibility = Visibility.Visible;
                    }

                    i++;
                }
                else if (meld.meldType == (int)MeldType.enumMeldTypeSequence)
                {
                    for (int j = 0; j < 3; j++)
                    {
                        var btn = ButtonsSp2[i++];
                        btn.Content = new Image() { Source = MyOwner.ImagesSrc[meld.tile1 + j] };
                        btn.Tag = meld.tile1;

                        btn.Visibility = Visibility.Visible;
                    }

                    i++;
                }
            }
        }

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
            HideSp1Button();

            var btn = sender as Button;
            if (btn == null)
                return;

            var tileId = (int)btn.Tag;
            SelectedTile = tileId;
        }

        public int SelectedTile { get; private set; }

        private void HideSp1Button()
        {
            foreach (var button in ButtonsSp1)
            {
                button.Visibility = Visibility.Hidden;
            }
        }

    }
}
