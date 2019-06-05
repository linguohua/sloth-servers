using System.Windows;

namespace MahjongTest
{
    /// <summary>
    /// InputWnd.xaml 的交互逻辑
    /// </summary>
    public partial class InputWnd : Window
    {
        public InputWnd()
        {
            IsNeedUserId = true;
            IsNeedRoomId = true;
            InitializeComponent();

            TextBoxRoomId.Text = ProgramConfig.RecentUsedRoomNumber;
        }

        private bool _isNeedUserId;
        private bool _isNeedRoomId;

        public bool IsNeedUserId
        {
            get { return _isNeedUserId; }
            set
            {
                _isNeedUserId = value;
                if (!_isNeedUserId)
                {
                    TextBoxUserId.IsEnabled = false;
                }
            }
        }

        public bool IsNeedRoomId
        {
            get { return _isNeedRoomId; }
            set
            {
                _isNeedRoomId = value;
                if (!_isNeedRoomId)
                {
                    TextBoxRoomId.IsEnabled = false;
                }
            }
        }

        private void OnOK_Button_Clicked(object sender, RoutedEventArgs e)
        {
            if (IsNeedUserId && string.IsNullOrWhiteSpace(TextBoxUserId.Text))
            {
                MessageBox.Show("please input a valid userID");
                return;
            }
            if (IsNeedRoomId && string.IsNullOrWhiteSpace(TextBoxRoomId.Text))
            {
                MessageBox.Show("please input a valid roomID");
                return;
            }

            if (IsNeedRoomId && !string.IsNullOrWhiteSpace(TextBoxRoomId.Text))
            {
                ProgramConfig.RecentUsedRoomNumber = TextBoxRoomId.Text;
            }
            
            DialogResult = true;
        }

        private void OnCancel_Button_Clicked(object sender, RoutedEventArgs e)
        {
            DialogResult = false;
        }
    }
}
