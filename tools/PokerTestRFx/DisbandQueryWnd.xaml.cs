using System.Windows;

namespace PokerTest
{
    /// <summary>
    /// DisbandQueryWnd.xaml 的交互逻辑
    /// </summary>
    public partial class DisbandQueryWnd : Window
    {
        public DisbandQueryWnd()
        {
            InitializeComponent();
        }

        private void OnAgreeButton_Clicked(object sender, RoutedEventArgs e)
        {
            DialogResult = true;
        }

        private void OnRefuseButton_Clicked(object sender, RoutedEventArgs e)
        {
            DialogResult = false;
        }
    }
}
