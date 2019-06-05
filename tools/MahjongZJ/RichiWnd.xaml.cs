using System.Collections.Generic;
using System.Windows;
using System.Windows.Controls;
using mahjong;

namespace MahjongTest
{
    /// <summary>
    /// RichiWnd.xaml 的交互逻辑
    /// </summary>
    public partial class RichiWnd : Window
    {
        public RichiWnd()
        {
            InitializeComponent();

            InitButtonArray();

            HideAllButtons();
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

        }
        public static bool ShowDialog(MsgReadyHandTips readyHandTips, TileStackWnd owner)
        {
            var x = new RichiWnd();
            x.SetOwner(owner);
            x.SetReadyHandTips(readyHandTips);

            var result = x.ShowDialog();
            if (result == null || !result.Value)
            {
                // snip
                return false;
            }

            return true;
        }

        private void SetReadyHandTips(MsgReadyHandTips readyHandTips)
        {
            var readyHandList = readyHandTips.readyHandList;
            var i = 0;
            for (var j = 0; j < readyHandList.Count - 1; j += 2)
            {
                var x = ButtonsSp1[i];
                x.Visibility = Visibility.Visible;

                var tid = readyHandList[j];
                x.Content = new Image() { Source = MyOwner.ImagesSrc[tid] };

                var y = ButtonsSp0[i];
                y.Visibility = Visibility.Visible;
                y.Content = readyHandList[j + 1];

                i++;
            }
        }

        private void SetOwner(TileStackWnd owner)
        {
            MyOwner = owner;
            this.Owner = owner.MyOwner;
        }

        public TileStackWnd MyOwner { get; private set; }
        private void OnRichi_Button_Clicked(object sender, RoutedEventArgs e)
        {
            this.DialogResult = true;
        }

        private void OnCancel_Button_Clicked(object sender, RoutedEventArgs e)
        {
            this.DialogResult = false;
        }
    }
}
