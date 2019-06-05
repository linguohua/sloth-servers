using System.Collections.Generic;
using System.Windows;

namespace PokerTest
{
    /// <summary>
    /// CfgSelectWnd.xaml 的交互逻辑
    /// </summary>
    public partial class CfgSelectWnd : Window
    {
        public CfgSelectWnd()
        {
            InitializeComponent();
        }

        private void OnSelect_Button_Clicked(object sender, RoutedEventArgs e)
        {
            SelectedIdx = LbCfgs.SelectedIndex;
            DialogResult = true;
        }

        private void OnCancel_Button_Clicked(object sender, RoutedEventArgs e)
        {
            DialogResult = false;
        }

        public static bool ShowDialog(List<DealCfg> cfgs, int currentSelect, out int newSelected, MainWindow owner)
        {
            newSelected = 0;

            var x = new CfgSelectWnd();
            x.SetOwner(owner);
            x.SetCfs(cfgs, currentSelect);

            var result = x.ShowDialog();
            if (result == null || !result.Value)
            {
                // snip
                return false;
            }

            newSelected  = x.SelectedIdx;
            return true;
        }

        private void SetCfs(List<DealCfg> cfgs, int currentSelect)
        {
            foreach (var cfg in cfgs)
            {
                LbCfgs.Items.Add(cfg);
            }

            LbCfgs.SelectedIndex = currentSelect;
        }

        public int SelectedIdx { get; set; }

        private void SetOwner(MainWindow owner)
        {
            Owner = owner;
        }
    }
}
