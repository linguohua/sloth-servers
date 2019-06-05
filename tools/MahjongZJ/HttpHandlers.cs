using System;
using System.IO;
using System.Net;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using System.Windows;

namespace MahjongTest
{
    internal class HttpHandlers
    {
        //public const string HttpServerUrl = @"http://localhost:3001";
        public const string PathUploadCfgFile = @"/support/uploadCfgs";
        public const string PathExportRoomOps = @"/support/exportRoomOps";
        public const string PathExportRoomCfg = @"/support/exportRoomCfg";
        public const string PathExportRoomSIDss = @"/support/exportRoomSIDs";

        public const string PathAttachDealCfgFile = @"/support/attachDealCfg";
        public const string PathAttachRoomCfgFile = @"/support/attachRoomCfg";

        public static async void SendFileContent(string filePath, MainWindow wnd)
        {
            try
            {
                var content = WriteSafeReadAllLines(filePath);
                using (var httpClient = new HttpClient())
                {
                    // Add a new Request Message
                    var url = $"{ProgramConfig.ServerUrl}{PathUploadCfgFile}?account={ProgramConfig.Account}&password={ProgramConfig.Password}";
                    HttpRequestMessage requestMessage = new HttpRequestMessage(HttpMethod.Post, url)
                    {
                        Content = new StringContent(content, Encoding.UTF8, "application/text")
                    };


                    var response = await httpClient.SendAsync(requestMessage);

                    if (response.StatusCode == HttpStatusCode.OK)
                    {
                        MessageBox.Show(wnd, "Upload OK");
                        var body = await response.Content.ReadAsStringAsync();
                        wnd.OnUploaded(body);
                    }
                    else
                    {

                        var body = await response.Content.ReadAsStringAsync();
                        MessageBox.Show(wnd, body);
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }

        public static async void SendFileContent2(string filePath, string roomNumber, string path, MainWindow wnd)
        {
            try
            {
                var content = WriteSafeReadAllLines(filePath);
                using (var httpClient = new HttpClient())
                {

                    // Add a new Request Message
                    var url = $"{ProgramConfig.ServerUrl}{path}?account={ProgramConfig.Account}&password={ProgramConfig.Password}&roomNumber={roomNumber}";
                    HttpRequestMessage requestMessage = new HttpRequestMessage(HttpMethod.Post, url)
                    {
                        Content = new StringContent(content, Encoding.UTF8, "application/text")
                    };

                    var response = await httpClient.SendAsync(requestMessage);

                    if (response.StatusCode == HttpStatusCode.OK)
                    {
                        MessageBox.Show(wnd, "Upload OK");
                        // var body = await response.Content.ReadAsStringAsync();
                        // wnd.OnUploaded(body);
                    }
                    else
                    {

                        var body = await response.Content.ReadAsStringAsync();
                        MessageBox.Show(wnd, body);
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }

        public static async Task<byte[]> ExportRoomOps(string xID, Window owner)
        {
            try
            {
                using (var httpClient = new HttpClient())
                {
                    // Add a new Request Message
                    var url = $"{ProgramConfig.ServerUrl}{PathExportRoomOps}?{xID}&account={ProgramConfig.Account}&password={ProgramConfig.Password}";
                    HttpRequestMessage requestMessage = new HttpRequestMessage(HttpMethod.Get, url);

                    var response = await httpClient.SendAsync(requestMessage);

                    if (response.StatusCode == HttpStatusCode.OK)
                    {
                        var body = await response.Content.ReadAsByteArrayAsync();
                        return body;
                    }
                    else
                    {

                        var body = await response.Content.ReadAsStringAsync();
                        MessageBox.Show(owner, body);
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }

            return null;
        }

       public static async Task<byte[]> ExportRoomOpsAcc(string url, Window owner)
        {
            try
            {
                using (var httpClient = new HttpClient())
                {
                    // Add a new Request Message
                    HttpRequestMessage requestMessage = new HttpRequestMessage(HttpMethod.Get, url);

                    var response = await httpClient.SendAsync(requestMessage);

                    if (response.StatusCode == HttpStatusCode.OK)
                    {
                        var body = await response.Content.ReadAsByteArrayAsync();
                        return body;
                    }
                    else
                    {

                        var body = await response.Content.ReadAsStringAsync();
                        MessageBox.Show(owner, body);
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }

            return null;
        }

        public static async Task<string> ExportRoomShareIDs(string xID, Window owner)
        {
            try
            {
                using (var httpClient = new HttpClient())
                {
                    // Add a new Request Message
                    var url = $"{ProgramConfig.ServerUrl}{PathExportRoomSIDss}?{xID}&account={ProgramConfig.Account}&password={ProgramConfig.Password}";
                    HttpRequestMessage requestMessage = new HttpRequestMessage(HttpMethod.Get, url);

                    var response = await httpClient.SendAsync(requestMessage);

                    if (response.StatusCode == HttpStatusCode.OK)
                    {
                        var body = await response.Content.ReadAsByteArrayAsync();
                        return Encoding.UTF8.GetString(body);
                    }
                    else
                    {

                        var body = await response.Content.ReadAsStringAsync();
                        MessageBox.Show(owner, body);
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }

            return null;
        }

        public static async Task<byte[]> ExportRoomCfg(string roomConfigId, Window owner)
        {
            try
            {
                using (var httpClient = new HttpClient())
                {
                    // Add a new Request Message
                    var url = $"{ProgramConfig.ServerUrl}{PathExportRoomCfg}?roomConfigID={roomConfigId}&account={ProgramConfig.Account}&password={ProgramConfig.Password}";
                    HttpRequestMessage requestMessage = new HttpRequestMessage(HttpMethod.Get, url);

                    var response = await httpClient.SendAsync(requestMessage);

                    if (response.StatusCode == HttpStatusCode.OK)
                    {
                        var body = await response.Content.ReadAsByteArrayAsync();
                        return body;
                    }
                    else
                    {

                        var body = await response.Content.ReadAsStringAsync();
                        MessageBox.Show(owner, body);
                    }
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }

            return null;
        }

        public static string WriteSafeReadAllLines(string path)
        {
            using (var csv = new FileStream(path, FileMode.Open, FileAccess.Read, FileShare.ReadWrite))
            using (var sr = new StreamReader(csv, Encoding.Default))
            {
                return sr.ReadToEnd();
            }
        }
        public static async void SendPostMethod(string path, string content, string extraQueryString)
        {
            try
            {
                using (var httpClient = new HttpClient())
                {
                    // Add a new Request Message
                    var url = $"{ProgramConfig.ServerUrl}{path}?account={ProgramConfig.Account}&password={ProgramConfig.Password}";
                    if (!string.IsNullOrWhiteSpace(extraQueryString))
                    {
                        url = url + extraQueryString;
                    }

                    var requestMessage = new HttpRequestMessage(HttpMethod.Post, url)
                    {
                        Content = new StringContent(content, Encoding.UTF8, "application/text")
                    };

                    var response = await httpClient.SendAsync(requestMessage);

                    if (response.StatusCode == HttpStatusCode.OK)
                    {
                        //MessageBox.Show("OK");
                        return;
                    }

                    // Get the response
                    var body = await response.Content.ReadAsStringAsync();
                    MessageBox.Show(body);
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }

        public static async void SendGetMethod(string path, string extraQueryString)
        {
            try
            {
                using (var httpClient = new HttpClient())
                {
                    // Add a new Request Message
                    var url = $"{ProgramConfig.ServerUrl}{path}?account={ProgramConfig.Account}&password={ProgramConfig.Password}";
                    if (!string.IsNullOrWhiteSpace(extraQueryString))
                    {
                        url = url + extraQueryString;
                    }

                    var requestMessage = new HttpRequestMessage(HttpMethod.Get, url);

                    var response = await httpClient.SendAsync(requestMessage);

                    var str = "OK";
                    if (response.StatusCode != HttpStatusCode.OK)
                    {
                        str = "Error";
                    }

                    // Get the response
                    var body = await response.Content.ReadAsStringAsync();
                    MessageBox.Show(body, str);
                }
            }
            catch (Exception ex)
            {
                MessageBox.Show(ex.Message);
            }
        }
    }
}
