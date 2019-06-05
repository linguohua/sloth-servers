using System.IO;
using System.Reflection;
using Newtonsoft.Json.Linq;
using System.Collections.Generic;
using Newtonsoft.Json;

namespace PokerTest
{
    class ProgramConfig
    {
        public static string ServerUrl = "http://localhost:3001";
        public static string RecentUsedRoomNumber = "";

        public const string ConfigFileName = "config.json";
        //public static List<string> OptionalUrls = new List<string>();

        public static ProgramConfigJSON configJSON;
        public static string Account = "linux";
        public static string Password = "MMzz----88~bbabcdefgqqb";

        public class ProgramConfigJSON
        {
            [JsonProperty("serverUrl")]
            public string serverURL;
            [JsonProperty("account")]
            public string account;
            [JsonProperty("password")]
            public string password;
            [JsonProperty("optionalUrls")]
            public string[] optionalURLs;
        }

        public static void LoadConfigFromFile()
        {
            var dir = System.IO.Path.GetDirectoryName(Assembly.GetExecutingAssembly().Location);
            if (string.IsNullOrWhiteSpace(dir))
            {
                return;
            }

            var filePath = Path.Combine(dir, ConfigFileName);
            if (!File.Exists(filePath))
            {
                return;
            }

            var str = HttpHandlers.WriteSafeReadAllLines(filePath);
            if (!string.IsNullOrWhiteSpace(str))
            {
                //var a = JObject.Parse(str);
                //var url = (string) a["serverUrl"];
                //if (!string.IsNullOrWhiteSpace(url))
                //{
                //    ServerUrl = url;
                //}

                var des = (ProgramConfigJSON)JsonConvert.DeserializeObject(str, typeof(ProgramConfigJSON));
                configJSON = des;

                ServerUrl = configJSON.serverURL;

                if (!string.IsNullOrWhiteSpace(configJSON.account))
                {
                    Account = configJSON.account;
                }

                if (!string.IsNullOrWhiteSpace(configJSON.password))
                {
                    Password = configJSON.password;
                }

            }
        }

        public static void SaveConfig2File()
        {
            if (configJSON == null)
            {
                configJSON = new ProgramConfigJSON();
            }

            configJSON.serverURL = ServerUrl;
            var str = JsonConvert.SerializeObject(configJSON, Formatting.Indented);
            var dir = System.IO.Path.GetDirectoryName(Assembly.GetExecutingAssembly().Location);
            if (string.IsNullOrWhiteSpace(dir))
            {
                return;
            }

            var filePath = Path.Combine(dir, ConfigFileName);
            File.WriteAllText(filePath, str);
        }
    }
}
