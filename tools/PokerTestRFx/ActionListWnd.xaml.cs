using System;
using System.Collections.Generic;
using System.Windows;

namespace PokerTest
{
    /// <summary>
    /// ActionListWnd.xaml 的交互逻辑
    /// </summary>
    public partial class ActionListWnd : Window
    {
        public ActionListWnd()
        {
            InitializeComponent();
        }

        private List<string> _Lines;
        private void OnClose_Button_Clicked(object sender, RoutedEventArgs e)
        {
            Hide();
        }

        public void ShowWithNewActionList(List<string> lines)
        {
            LbActionList.Items.Clear();
            lines.ForEach(x=> LbActionList.Items.Add((x)));
            _Lines = lines;
            Show();
        }

        public void SetOwner(MainWindow mainWindow)
        {
            Owner = mainWindow;
        }

        public void SelectIfSame(string y)
        {
            var i = LbActionList.SelectedIndex;
            if (i < 0)
                i = 0;

            for (; i < _Lines.Count; ++i)
            {
                if (string.Compare(_Lines[i], y, StringComparison.Ordinal) == 0)
                    break;
            }

            LbActionList.SelectedIndex = i;
        }

        public void ResetSelectedIndex()
        {
            LbActionList.UnselectAll();
        }
    }
}
