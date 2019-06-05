using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Shapes;

namespace PokerTest
{
    /// <summary>
    /// Host.xaml 的交互逻辑
    /// </summary>
    public partial class HostWnd : Window
    {
        public HostWnd()
        {
            InitializeComponent();

            combobHost.Text = ProgramConfig.ServerUrl;

            if (ProgramConfig.configJSON.optionalURLs != null)
            {
                foreach (var str in ProgramConfig.configJSON.optionalURLs)
                {
                    combobHost.Items.Add(str);
                }

            }

        }

        private string _hostUrl;

        public string HostUrl
        {
            get
            {
                return _hostUrl;
            }
        }

        private void OnOK_Button_Clicked(object sender, RoutedEventArgs e)
        {
            var hostUrl = combobHost.Text;
            try
            {
                if (string.IsNullOrWhiteSpace(hostUrl))
                {
                    MessageBox.Show("please input host url");
                    return;
                }

                Uri uriResult;
                bool result = Uri.TryCreate(hostUrl, UriKind.Absolute, out uriResult)
                    && uriResult.Scheme == Uri.UriSchemeHttp;

                if (!result)
                {
                    MessageBox.Show("please input a valid host url");
                    return;
                }

                _hostUrl = hostUrl;
                DialogResult = true;
            }
            catch (System.Exception exx)
            {
                MessageBox.Show(exx.Message);
            }

        }

        private void OnCancel_Button_Clicked(object sender, RoutedEventArgs e)
        {
            DialogResult = false;
        }
    }
}
