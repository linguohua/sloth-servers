using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MahjongTest
{
    public class MyCsvHeaders
    {
        private static string[] _HeadersForYC = new[] {
            "名称", "类型", "庄家userID", "庄家手牌",  "庄家动作提示", "userID2", "手牌", "动作提示", "userID3", "手牌",
             "动作提示","userID4", "手牌", "动作提示", "抽牌序列", "强制一致", "房间配置ID", "是否连庄"
        };

        public static string[] GetHeaders()
        {
            return _HeadersForYC;
        }

        public static string GetRoomTypeName()
        {
            return "湛江";
        }
    }
}
