using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MahjongTest
{
    public class MyCsvHeaders
    {
        private static string[] _HeadersForDF = new[] {
            "名称","类型", "庄家userID", "庄家手牌", "庄家花牌", "庄家动作提示", "userID2", "手牌", "花牌", "动作提示", "userID3", "手牌",
            "花牌", "动作提示",
            "userID4", "手牌", "花牌", "动作提示", "抽牌序列", "风牌", "强制一致", "房间配置ID","是否连庄","家家庄"
        };
        private static string[] _HeadersForDT = new[] {
            "名称","类型", "庄家userID", "庄家手牌", "庄家花牌", "庄家动作提示", "userID2", "手牌", "花牌", "动作提示", "userID3", "手牌",
            "花牌", "动作提示",
            "userID4", "手牌", "花牌", "动作提示", "抽牌序列", "风牌", "强制一致", "房间配置ID","是否连庄"
        };
        private static string[] _HeadersForYC = new[] {
            "名称", "类型", "庄家userID", "庄家手牌", "庄家花牌", "庄家动作提示", "userID2", "手牌", "花牌", "动作提示", "userID3", "手牌",
            "花牌", "动作提示",
            "userID4", "手牌", "花牌", "动作提示", "抽牌序列", "加价局", "强制一致", "房间配置ID"
        };

        public static string[] GetHeaders(RoomType rt)
        {
            switch (rt)
            {
                case RoomType.DafengMJ:
                    return _HeadersForDF;
                case RoomType.DongTaiMJ:
                    return _HeadersForDT;
                case RoomType.YanChengMJ:
                    return _HeadersForYC;
                default:
                    return _HeadersForDF;
            }
        }

        public static string GetRoomTypeName(RoomType rt)
        {
            switch (rt)
            {
                case RoomType.DafengMJ:
                    return "大丰";
                case RoomType.DongTaiMJ:
                    return "东台";
                case RoomType.YanChengMJ:
                    return "盐城";
                default:
                    return "大丰";
            }
        }
    }
}
