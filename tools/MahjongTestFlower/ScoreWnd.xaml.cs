using System.ComponentModel;
using System.Globalization;
using System.Windows;
using System.Windows.Media;

namespace MahjongTest
{
    /// <summary>
    /// ScoreWnd.xaml 的交互逻辑
    /// </summary>
    public partial class ScoreWnd : Window
    {
        public ScoreWnd()
        {
            InitializeComponent();
        }

        private void OnOK_Button_Clicked(object sender, RoutedEventArgs e)
        {
            Hide();
        }

        public void ShowWithMsg(string msg, MainWindow owner)
        {
            Owner = owner;
            TbMsg.Text = msg;
            Size2Content();
            
            Show();
        }

        private void Size2Content()
        {
            if (string.IsNullOrWhiteSpace(TbMsg.Text))
                return;

            var formattedText = new FormattedText(
                TbMsg.Text,
                CultureInfo.CurrentUICulture,
                FlowDirection.LeftToRight,
                new Typeface(TbMsg.FontFamily, TbMsg.FontStyle, TbMsg.FontWeight, TbMsg.FontStretch),
                TbMsg.FontSize,
                Brushes.Black);

            var width = Width;
            var height = Height;

            if (formattedText.Width > width)
                width = formattedText.Width;
            if ((formattedText.Height+120) > height)
                height = 120 + formattedText.Height;

            Height = height;
            Width = width;
        }

        public void Clear()
        {
            TbMsg.Clear();
        }

        public string GetResultMsg()
        {
            return TbMsg.Text;
        }

        private void ScoreWnd_OnClosing(object sender, CancelEventArgs e)
        {
            Hide();
            e.Cancel = true;
        }
    }
}
